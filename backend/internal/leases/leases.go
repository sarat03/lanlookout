// Package leases reads a dnsmasq-style DHCP leases file (dnsmasq, Pi-hole, OpenWrt).
package leases

import (
	"bufio"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sarat03/lanlookout/internal/models"
)

// Read - hosts with an active lease. Line format: "expiry mac ip hostname clientid".
// ponytail: an active lease means "was on the network recently", not "online now";
// a device that left stays online until its lease expires (often hours).
func Read(path string) (hosts []models.Host) {

	f, err := os.Open(path)
	if err != nil {
		slog.Error("Can't read DHCP leases", "path", path, "err", err)
		return nil
	}
	defer f.Close()

	return parse(bufio.NewScanner(f), time.Now())
}

func parse(sc *bufio.Scanner, now time.Time) (hosts []models.Host) {

	for sc.Scan() {
		p := strings.Fields(sc.Text())
		if len(p) < 4 {
			continue
		}
		exp, err := strconv.ParseInt(p[0], 10, 64)
		if err != nil || (exp != 0 && exp < now.Unix()) { // 0 = infinite lease
			continue
		}

		var h models.Host
		h.Mac = strings.ToLower(p[1])
		h.IP = p[2]
		if p[3] != "*" {
			h.Name = p[3]
		}
		h.Iface = "dhcp"
		h.Now = 1 // like arp results: an active lease counts as online
		h.Date = now.Format("2006-01-02 15:04:05")
		hosts = append(hosts, h)
	}
	return hosts
}
