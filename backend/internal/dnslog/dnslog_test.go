package dnslog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
)

func setup(t *testing.T) {
	conf.AppConfig = models.Conf{UseDB: "sqlite", DBPath: t.TempDir() + "/t.db", AlertDNS: true}
	gdb.Start()
	threats = ParseList(strings.NewReader("# comment\n0.0.0.0 evil.example\nbad.test\n127.0.0.1 localhost\n"))
	alerted = map[string]bool{}
}

func TestThreatList(t *testing.T) {
	setup(t)
	for d, want := range map[string]bool{"evil.example": true, "cdn.evil.example": true, "example": false, "good.example": false, "BAD.test.": true, "localhost": false} {
		if IsThreat(d) != want {
			t.Errorf("IsThreat(%q) != %v", d, want)
		}
	}
}

func TestPihole(t *testing.T) {
	setup(t)
	now := time.Now()
	logins := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/auth":
			logins++
			json.NewEncoder(w).Encode(map[string]any{"session": map[string]any{"valid": true, "sid": "s1"}})
		case "/api/queries":
			if r.Header.Get("X-FTL-SID") != "s1" {
				w.WriteHeader(401)
				return
			}
			json.NewEncoder(w).Encode(map[string]any{"queries": []map[string]any{
				{"time": float64(now.Unix()), "domain": "youtube.com", "status": "FORWARDED", "client": map[string]string{"ip": "10.0.0.5"}},
				{"time": float64(now.Unix()), "domain": "youtube.com", "status": "CACHE", "client": map[string]string{"ip": "10.0.0.5"}},
				{"time": float64(now.Unix()), "domain": "ads.example", "status": "GRAVITY", "client": map[string]string{"ip": "10.0.0.5"}},
				{"time": float64(now.Unix()), "domain": "x.evil.example", "status": "FORWARDED", "client": map[string]string{"ip": "10.0.0.5"}},
			}})
		}
	}))
	defer srv.Close()

	conf.AppConfig.DNSType, conf.AppConfig.DNSURL, conf.AppConfig.DNSPass = "pihole", srv.URL, "pw"
	sid = ""
	if err := Poll(now.Add(-time.Minute), now); err != nil {
		t.Fatal(err)
	}
	if err := Poll(now.Add(-time.Minute), now); err != nil { // same window again: counts add up
		t.Fatal(err)
	}

	top := gdb.SelectDNSTop("10.0.0.5", now.Format("2006-01-02"), 10)
	if len(top) != 3 || top[0].Domain != "youtube.com" || top[0].Count != 4 {
		t.Fatalf("top %+v", top)
	}
	for _, s := range top {
		if s.Domain == "ads.example" && s.Blocked != 2 {
			t.Errorf("blocked %+v", s)
		}
		if s.Domain == "x.evil.example" && s.Threat != 1 {
			t.Errorf("threat %+v", s)
		}
	}
	if ev := gdb.SelectEvents(10); len(ev) != 1 || ev[0].Kind != "threat" {
		t.Errorf("want one threat alert (deduped), got %+v", ev)
	}
	if logins != 1 {
		t.Errorf("session should be reused, logins=%d", logins)
	}
}

func TestAdGuardPaging(t *testing.T) {
	setup(t)
	now := time.Now().UTC()
	ts := func(m int) string { return now.Add(-time.Duration(m) * time.Minute).Format(time.RFC3339Nano) }
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, p, _ := r.BasicAuth(); u != "admin" || p != "pw" {
			w.WriteHeader(401)
			return
		}
		page1 := r.URL.Query().Get("older_than") == ""
		data := []map[string]any{
			{"time": ts(1), "client": "10.0.0.6", "reason": "NotFilteredNotFound", "question": map[string]string{"name": "netflix.com"}},
			{"time": ts(2), "client": "10.0.0.6", "reason": "FilteredBlackList", "question": map[string]string{"name": "tracker.example"}},
		}
		oldest := ts(2)
		if !page1 {
			data = []map[string]any{
				{"time": ts(3), "client": "10.0.0.6", "reason": "NotFilteredNotFound", "question": map[string]string{"name": "netflix.com"}},
				{"time": ts(30), "client": "10.0.0.6", "reason": "NotFilteredNotFound", "question": map[string]string{"name": "old.example"}},
			}
			oldest = ts(30)
		}
		json.NewEncoder(w).Encode(map[string]any{"data": data, "oldest": oldest})
	}))
	defer srv.Close()

	qs, err := adguardQueries(srv.URL, "admin", "pw", now.Add(-10*time.Minute))
	if err != nil || len(qs) != 3 {
		t.Fatalf("got %d %v", len(qs), err)
	}
	if !qs[1].Blocked || qs[0].Blocked {
		t.Errorf("blocked flags %+v", qs)
	}
	if _, err := adguardQueries(srv.URL, "admin", "wrong", now); err == nil {
		t.Error("bad password should error")
	}
}
