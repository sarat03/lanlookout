// Package dnslog pulls DNS query logs from Pi-hole (v6 API) or AdGuard Home, counts
// queries per device and domain, and alerts when a device looks up a domain on a
// threat list. This is how LanLookout sees "which sites each device talks to"
// without sitting in the traffic path.
package dnslog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
	"github.com/sarat03/lanlookout/internal/notify"
)

// Query - one DNS lookup
type Query struct {
	Time    time.Time
	IP      string
	Domain  string
	Blocked bool
}

var (
	client  = &http.Client{Timeout: 20 * time.Second}
	sid     string // Pi-hole session
	threats = map[string]bool{}
	alerted = map[string]bool{} // ip|domain|day
	mu      sync.Mutex
	lastErr string
)

// LastError - for the settings page ("" = working)
func LastError() string {
	mu.Lock()
	defer mu.Unlock()
	return lastErr
}

// Start - poll every 5 minutes, refresh the threat list daily
func Start() {
	go func() {
		since := time.Now().Add(-5 * time.Minute)
		var listAt time.Time
		for {
			if conf.AppConfig.DNSType != "" && conf.AppConfig.DNSURL != "" {
				if conf.AppConfig.ThreatListURL != "" && time.Since(listAt) > 24*time.Hour {
					if err := loadThreats(conf.AppConfig.ThreatListURL); err == nil {
						listAt = time.Now()
					}
				}
				now := time.Now()
				if err := Poll(since, now); err == nil {
					since = now
				}
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

// Poll - fetch queries in [since, until) and store the counts
func Poll(since, until time.Time) error {

	var qs []Query
	var err error
	switch conf.AppConfig.DNSType {
	case "pihole":
		qs, err = piholeQueries(conf.AppConfig.DNSURL, conf.AppConfig.DNSPass, since, until)
	case "adguard":
		qs, err = adguardQueries(conf.AppConfig.DNSURL, conf.AppConfig.DNSUser, conf.AppConfig.DNSPass, since)
	default:
		err = errors.New("unknown DNS_TYPE " + conf.AppConfig.DNSType)
	}

	mu.Lock()
	lastErr = ""
	if err != nil {
		lastErr = err.Error()
	}
	mu.Unlock()
	if err != nil {
		slog.Error("DNS log fetch failed", "type", conf.AppConfig.DNSType, "err", err)
		return err
	}

	store(qs)
	return nil
}

func store(qs []Query) {

	byDay := map[string]map[[2]string]models.DNSStat{}
	for _, q := range qs {
		day := q.Time.Format("2006-01-02")
		if byDay[day] == nil {
			byDay[day] = map[[2]string]models.DNSStat{}
		}
		k := [2]string{q.IP, q.Domain}
		c := byDay[day][k]
		c.Count++
		if q.Blocked {
			c.Blocked++
		}
		if IsThreat(q.Domain) {
			c.Threat = 1
			alertThreat(q, day)
		}
		byDay[day][k] = c
	}
	for day, counts := range byDay {
		gdb.AddDNSCounts(day, counts)
	}
}

func alertThreat(q Query, day string) {
	key := q.IP + "|" + q.Domain + "|" + day
	mu.Lock()
	seen := alerted[key]
	alerted[key] = true
	mu.Unlock()
	if seen {
		return
	}
	host := models.Host{IP: q.IP}
	name := q.IP
	for _, h := range gdb.SelectByIP(q.IP) {
		host = h
		if h.Name != "" {
			name = h.Name
		}
	}
	state := "allowed"
	if q.Blocked {
		state = "blocked by your DNS filter"
	}
	notify.Event("threat", host, fmt.Sprintf("'%s' (%s) looked up a domain on the threat list: %s (%s)", name, q.IP, q.Domain, state), conf.AppConfig.AlertDNS)
}

// IsThreat - domain or any parent domain is on the list
func IsThreat(d string) bool {
	mu.Lock()
	defer mu.Unlock()
	d = strings.ToLower(strings.TrimSuffix(d, "."))
	for d != "" {
		if threats[d] {
			return true
		}
		_, rest, ok := strings.Cut(d, ".")
		if !ok || !strings.Contains(rest, ".") { // stop at the TLD+1
			return false
		}
		d = rest
	}
	return false
}

// loadThreats - hosts-format ("0.0.0.0 bad.com") or one domain per line
func loadThreats(u string) error {
	resp, err := client.Get(u)
	if err != nil {
		slog.Error("Threat list download failed", "err", err)
		return err
	}
	defer resp.Body.Close()
	list := ParseList(io.LimitReader(resp.Body, 100<<20))
	mu.Lock()
	threats = list
	alerted = map[string]bool{}
	mu.Unlock()
	slog.Info("Threat list loaded", "domains", len(list))
	return nil
}

// ParseList - read a hosts-format or plain domain list
func ParseList(r io.Reader) map[string]bool {
	list := map[string]bool{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line, _, _ := strings.Cut(sc.Text(), "#")
		f := strings.Fields(line)
		var d string
		switch {
		case len(f) == 1:
			d = f[0]
		case len(f) >= 2:
			d = f[1]
		default:
			continue
		}
		d = strings.ToLower(d)
		if strings.Contains(d, ".") && d != "localhost" && d != "0.0.0.0" {
			list[d] = true
		}
	}
	return list
}

// ---- Pi-hole v6: POST /api/auth, then GET /api/queries with X-FTL-SID ----

func piholeAuth(base, pass string) error {
	body, _ := json.Marshal(map[string]string{"password": pass})
	resp, err := client.Post(base+"/api/auth", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var r struct {
		Session struct {
			Valid bool   `json:"valid"`
			SID   string `json:"sid"`
		} `json:"session"`
	}
	if json.NewDecoder(resp.Body).Decode(&r) != nil || !r.Session.Valid {
		return errors.New("Pi-hole login failed (check the password / app password)")
	}
	sid = r.Session.SID
	return nil
}

func piholeQueries(base, pass string, since, until time.Time) ([]Query, error) {

	base = strings.TrimRight(base, "/")
	if sid == "" {
		if err := piholeAuth(base, pass); err != nil {
			return nil, err
		}
	}

	get := func() (*http.Response, error) {
		q := url.Values{"from": {fmt.Sprint(since.Unix())}, "until": {fmt.Sprint(until.Unix())}, "length": {"100000"}}
		req, _ := http.NewRequest("GET", base+"/api/queries?"+q.Encode(), nil)
		req.Header.Set("X-FTL-SID", sid)
		return client.Do(req)
	}
	resp, err := get()
	if err == nil && resp.StatusCode == http.StatusUnauthorized { // session expired
		resp.Body.Close()
		if err = piholeAuth(base, pass); err == nil {
			resp, err = get()
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Pi-hole: %s", resp.Status)
	}

	var r struct {
		Queries []struct {
			Time   float64 `json:"time"`
			Domain string  `json:"domain"`
			Status string  `json:"status"`
			Client struct {
				IP string `json:"ip"`
			} `json:"client"`
		} `json:"queries"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	out := make([]Query, 0, len(r.Queries))
	for _, q := range r.Queries {
		s := strings.ToUpper(q.Status)
		blocked := strings.Contains(s, "GRAVITY") || strings.Contains(s, "REGEX") ||
			strings.Contains(s, "DENYLIST") || strings.Contains(s, "BLOCKED") || strings.Contains(s, "BLACKLIST")
		out = append(out, Query{time.Unix(int64(q.Time), 0), q.Client.IP, q.Domain, blocked})
	}
	return out, nil
}

// ---- AdGuard Home: GET /control/querylog (newest first, paged by older_than) ----

func adguardQueries(base, user, pass string, since time.Time) ([]Query, error) {

	base = strings.TrimRight(base, "/")
	var out []Query
	older := ""
	for page := 0; page < 50; page++ { // ponytail: max 50 x 500 queries per poll
		q := url.Values{"limit": {"500"}}
		if older != "" {
			q.Set("older_than", older)
		}
		req, _ := http.NewRequest("GET", base+"/control/querylog?"+q.Encode(), nil)
		if user != "" {
			req.SetBasicAuth(user, pass)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("AdGuard Home: %s", resp.Status)
		}
		var r struct {
			Oldest string `json:"oldest"`
			Data   []struct {
				Time     string `json:"time"`
				Client   string `json:"client"`
				Reason   string `json:"reason"`
				Question struct {
					Name string `json:"name"`
				} `json:"question"`
			} `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&r)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		for _, d := range r.Data {
			t, err := time.Parse(time.RFC3339Nano, d.Time)
			if err != nil {
				continue
			}
			if t.Before(since) {
				return out, nil
			}
			out = append(out, Query{t, d.Client, d.Question.Name, strings.HasPrefix(d.Reason, "Filtered")})
		}
		if len(r.Data) == 0 || r.Oldest == "" {
			break
		}
		older = r.Oldest
	}
	return out, nil
}
