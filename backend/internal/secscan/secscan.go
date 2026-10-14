// Package secscan checks devices for services that are risky to leave open on a home
// network, and alerts when one appears.
package secscan

import (
	"fmt"
	"log/slog"
	"net"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/identify"
	"github.com/sarat03/lanlookout/internal/models"
	"github.com/sarat03/lanlookout/internal/notify"
)

// Risk - why an open port matters
type Risk struct {
	Service  string
	Severity string // high, medium, low
	Advice   string
}

// Risky - ports worth knowing about. Ordinary web UIs (80/443) are listed as low.
var Risky = map[int]Risk{
	21:    {"FTP", "high", "Sends passwords in plain text. Turn it off or use SFTP."},
	23:    {"Telnet", "high", "Unencrypted remote login, a favourite of IoT malware. Turn it off."},
	135:   {"Windows RPC", "medium", "Should not be reachable from other devices unless you need it."},
	139:   {"NetBIOS / SMB", "medium", "Old file sharing. Turn off SMBv1; share only what you need."},
	445:   {"SMB file sharing", "medium", "Keep it patched and password protected; never forward it to the internet."},
	554:   {"RTSP video", "medium", "Camera streams are often viewable without a password. Set one."},
	1883:  {"MQTT (unencrypted)", "medium", "Require a username and password on the broker."},
	2323:  {"Telnet (alt)", "high", "Alternate telnet port used by IoT devices. Turn it off."},
	2375:  {"Docker API (no TLS)", "high", "Anyone on the network can control Docker and the host. Close it."},
	3306:  {"MySQL", "medium", "Databases should only listen on localhost or a private network."},
	3389:  {"Remote Desktop", "medium", "Use a strong password and network level authentication; never forward it."},
	5432:  {"PostgreSQL", "medium", "Databases should only listen on localhost or a private network."},
	5555:  {"Android debug (ADB)", "high", "Gives full control of the device. Turn off wireless debugging."},
	5900:  {"VNC", "medium", "Often weakly protected. Use a strong password or tunnel it over SSH."},
	6379:  {"Redis", "high", "Usually has no password. Bind it to localhost."},
	8080:  {"HTTP admin (alt)", "low", "Check this web page has a password and isn't the default one."},
	9200:  {"Elasticsearch", "high", "Usually has no password. Bind it to localhost."},
	11211: {"Memcached", "medium", "No authentication. Bind it to localhost."},
	27017: {"MongoDB", "high", "Often has no password. Enable auth and bind to localhost."},
	80:    {"HTTP admin", "low", "Check this web page has a password and isn't the default one."},
	22:    {"SSH", "low", "Fine if you use it. Prefer keys over passwords."},
}

var (
	mu       sync.Mutex
	running  bool
	lastScan time.Time
)

// Finding - one risky thing on one device
type Finding struct {
	HostID   int
	Name     string
	IP       string
	Port     int
	Service  string
	Severity string
	Advice   string
}

// Start - scan on a schedule
func Start() {
	go func() {
		time.Sleep(2 * time.Minute) // let the first ARP scan finish
		for {
			if conf.AppConfig.SecScanHours > 0 {
				Scan()
			}
			time.Sleep(time.Duration(max(conf.AppConfig.SecScanHours, 1)) * time.Hour)
		}
	}()
}

// LastScan - when the last scan finished, and whether one is running
func LastScan() (time.Time, bool) {
	mu.Lock()
	defer mu.Unlock()
	return lastScan, running
}

// Scan - check every online host. Returns false if a scan is already running.
func Scan() bool {

	mu.Lock()
	if running {
		mu.Unlock()
		return false
	}
	running = true
	mu.Unlock()
	defer func() {
		mu.Lock()
		running, lastScan = false, time.Now()
		mu.Unlock()
	}()

	hosts, _ := gdb.Select("now")
	ports := make([]int, 0, len(Risky))
	for p := range Risky {
		ports = append(ports, p)
	}

	slog.Info("Security scan started", "hosts", len(hosts))
	sem := make(chan struct{}, 16) // ponytail: 16 hosts at once x ~20 ports; lower if the router complains
	var wg sync.WaitGroup
	for _, h := range hosts {
		if h.Now != 1 || h.IP == "" {
			continue
		}
		wg.Add(1)
		go func(h models.Host) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			open := OpenPorts(h.IP, ports, 1500*time.Millisecond)
			newRisky := Added(h.Ports, open)
			gdb.SetPorts(h.ID, join(open))

			for _, p := range newRisky {
				r := Risky[p]
				if r.Severity == "low" {
					continue
				}
				name := h.Name
				if name == "" {
					name = h.IP
				}
				notify.Event("port", h, fmt.Sprintf("Risky service on '%s' (%s): %s, port %d. %s", name, h.IP, r.Service, p, r.Advice), conf.AppConfig.AlertSec)
			}
		}(h)
	}
	wg.Wait()
	slog.Info("Security scan finished")
	return true
}

// OpenPorts - TCP connect check, all ports in parallel
func OpenPorts(ip string, ports []int, timeout time.Duration) []int {
	var (
		m    sync.Mutex
		open []int
		wg   sync.WaitGroup
	)
	for _, p := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			c, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(p)), timeout)
			if err == nil {
				c.Close()
				m.Lock()
				open = append(open, p)
				m.Unlock()
			}
		}(p)
	}
	wg.Wait()
	sort.Ints(open)
	return open
}

// Added - ports open now that weren't in the previous scan ("" = never scanned: no alerts)
func Added(prev string, now []int) (added []int) {
	if prev == "" {
		return nil
	}
	old := split(prev)
	for _, p := range now {
		if !slices.Contains(old, p) {
			added = append(added, p)
		}
	}
	return added
}

// Findings - risky open ports of all hosts plus routers with UPnP, worst first
func Findings() (list []Finding) {

	hosts, _ := gdb.Select("now")
	gateways := identify.UPnPGateways()
	for _, h := range hosts {
		for _, p := range split(h.Ports) {
			if r, ok := Risky[p]; ok {
				list = append(list, Finding{h.ID, h.Name, h.IP, p, r.Service, r.Severity, r.Advice})
			}
		}
		if slices.Contains(gateways, h.IP) {
			list = append(list, Finding{h.ID, h.Name, h.IP, 0, "UPnP port forwarding", "medium",
				"Any device on your network can open ports to the internet. Turn UPnP off in the router unless a game console or app needs it."})
		}
	}
	rank := map[string]int{"high": 0, "medium": 1, "low": 2}
	sort.SliceStable(list, func(i, j int) bool { return rank[list[i].Severity] < rank[list[j].Severity] })
	return list
}

func split(s string) (ports []int) {
	for _, p := range strings.Split(s, ",") {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			ports = append(ports, n)
		}
	}
	return ports
}

// join - "none" marks a scanned host with nothing open, so the next scan can alert
func join(ports []int) string {
	if len(ports) == 0 {
		return "none"
	}
	s := make([]string, len(ports))
	for i, p := range ports {
		s[i] = strconv.Itoa(p)
	}
	return strings.Join(s, ",")
}
