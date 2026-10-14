// Package netcheck watches the internet connection: reachability and latency every
// few minutes, public IP changes, outages, and optional scheduled speed tests.
package netcheck

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
	"github.com/sarat03/lanlookout/internal/notify"
)

// Targets - TCP endpoints that answer from almost anywhere (Cloudflare, Google DNS over TLS/HTTPS).
var Targets = []string{"1.1.1.1:443", "8.8.8.8:443", "9.9.9.9:443"}

const timeFmt = "2006-01-02 15:04:05"

var (
	mu        sync.Mutex
	downSince time.Time // zero = internet up
	publicIP  string
	lastSpeed time.Time
	speedMu   sync.Mutex // one speed test at a time
)

// Status - for the Internet page
type Status struct {
	Up        bool
	DownSince string
	PublicIP  string
	Last      models.NetStat
}

// Current - latest known state
func Current() Status {
	mu.Lock()
	defer mu.Unlock()
	s := Status{Up: downSince.IsZero(), PublicIP: publicIP}
	if !downSince.IsZero() {
		s.DownSince = downSince.Format(timeFmt)
	}
	if rows := gdb.SelectNetStats(time.Now().Add(-24 * time.Hour).Format(timeFmt)); len(rows) > 0 {
		s.Last = rows[len(rows)-1]
	}
	return s
}

// Start - check loop
func Start() {
	go func() {
		// Remember the last public IP so a restart doesn't alert.
		if rows := gdb.SelectNetStats(time.Now().Add(-30 * 24 * time.Hour).Format(timeFmt)); len(rows) > 0 {
			for i := len(rows) - 1; i >= 0; i-- {
				if rows[i].PublicIP != "" {
					publicIP = rows[i].PublicIP
					break
				}
			}
			for i := len(rows) - 1; i >= 0; i-- {
				if rows[i].DownMbps > 0 {
					lastSpeed, _ = time.ParseInLocation(timeFmt, rows[i].Date, time.Local)
					break
				}
			}
		}
		for {
			if conf.AppConfig.NetCheckMin > 0 {
				Check()
			}
			time.Sleep(time.Duration(max(conf.AppConfig.NetCheckMin, 1)) * time.Minute)
		}
	}()
}

// Check - one reachability check, plus public IP and a due speed test
func Check() models.NetStat {

	stat := models.NetStat{Date: time.Now().Format(timeFmt)}
	lat, ok := Latency(Targets, 3*time.Second)
	if ok {
		stat.Up, stat.LatencyMs = 1, lat
	}
	track(ok, time.Now())

	if ok && conf.AppConfig.PublicIPURL != "" {
		if ip, err := fetchIP(conf.AppConfig.PublicIPURL); err == nil {
			stat.PublicIP = ip
			trackIP(ip)
		} else {
			slog.Debug("Public IP lookup failed", "err", err)
		}
	}

	if ok && conf.AppConfig.SpeedHours > 0 && time.Since(lastSpeed) >= time.Duration(conf.AppConfig.SpeedHours)*time.Hour {
		if down, up, err := Speed(conf.AppConfig.SpeedURL); err == nil {
			stat.DownMbps, stat.UpMbps = down, up
		}
	}

	gdb.AddNetStat(stat)
	return stat
}

// Latency - fastest TCP connect among targets, in ms
func Latency(targets []string, timeout time.Duration) (float64, bool) {
	type res struct {
		ms float64
		ok bool
	}
	ch := make(chan res, len(targets))
	for _, t := range targets {
		go func(t string) {
			start := time.Now()
			c, err := net.DialTimeout("tcp", t, timeout)
			if err != nil {
				ch <- res{}
				return
			}
			c.Close()
			ch <- res{float64(time.Since(start).Microseconds()) / 1000, true}
		}(t)
	}
	best, any := 0.0, false
	for range targets {
		r := <-ch
		if r.ok && (!any || r.ms < best) {
			best, any = r.ms, true
		}
	}
	return best, any
}

// track - outage start/end. You can't be told the internet is down while it's down,
// so the alert goes out when it comes back.
func track(up bool, now time.Time) {
	mu.Lock()
	defer mu.Unlock()
	switch {
	case !up && downSince.IsZero():
		downSince = now
		slog.Warn("Internet is down")
	case up && !downSince.IsZero():
		msg := fmt.Sprintf("Internet was down for %s (from %s)", now.Sub(downSince).Round(time.Minute), downSince.Format(timeFmt))
		downSince = time.Time{}
		go notify.Event("outage", models.Host{}, msg, conf.AppConfig.AlertNet)
	}
}

func trackIP(ip string) {
	mu.Lock()
	old := publicIP
	publicIP = ip
	mu.Unlock()
	if old != "" && old != ip {
		notify.Event("wan_ip", models.Host{}, fmt.Sprintf("Public IP changed: '%s' -> '%s'", old, ip), conf.AppConfig.AlertNet)
	}
}

func fetchIP(url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 100))
	ip := strings.TrimSpace(string(b))
	if net.ParseIP(ip) == nil {
		return "", errors.New("not an IP: " + ip)
	}
	return ip, nil
}

// Speed - download and upload test against Cloudflare-style /__down and /__up endpoints.
// ponytail: one stream each way, so fast links (>500 Mbit) read low; add parallel streams if that matters.
func Speed(base string) (down, up float64, err error) {

	speedMu.Lock()
	defer speedMu.Unlock()

	base = strings.TrimRight(base, "/")
	client := &http.Client{Timeout: 60 * time.Second}

	const downBytes, upBytes = 25_000_000, 10_000_000

	start := time.Now()
	resp, err := client.Get(fmt.Sprintf("%s/__down?bytes=%d", base, downBytes))
	if err != nil {
		return 0, 0, err
	}
	n, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil || n == 0 {
		return 0, 0, fmt.Errorf("download: %v", err)
	}
	down = mbps(n, time.Since(start))

	start = time.Now()
	resp, err = client.Post(base+"/__up", "application/octet-stream", bytes.NewReader(make([]byte, upBytes)))
	if err != nil {
		return down, 0, err
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	up = mbps(upBytes, time.Since(start))

	mu.Lock()
	lastSpeed = time.Now()
	mu.Unlock()
	slog.Info("Speed test", "down_mbps", down, "up_mbps", up)
	return down, up, nil
}

func mbps(n int64, d time.Duration) float64 {
	return float64(int(float64(n)*8/d.Seconds()/1e4)) / 100 // 2 decimals
}
