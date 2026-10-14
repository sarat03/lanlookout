package routines

import (
	"testing"
	"time"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
	"github.com/sarat03/lanlookout/internal/notify"
)

// Watched host: miss, miss (offline alert), back with new IP (online + IP alert).
func TestCompareHostsAlerts(t *testing.T) {
	conf.AppConfig = models.Conf{UseDB: "sqlite", DBPath: t.TempDir() + "/t.db", OfflineAfter: 2}
	gdb.Start()

	mac := "aa:bb:cc:00:00:01"
	gdb.Update("now", models.Host{Mac: mac, IP: "10.0.0.5", Name: "phone", Known: 1, Watch: 1, Now: 1})
	seen := func(ip string) map[string]models.Host {
		return map[string]models.Host{mac: {Mac: mac, IP: ip, Date: "d"}}
	}
	none := map[string]models.Host{}

	compareHosts(seen("10.0.0.5"))
	compareHosts(none)
	if n := len(gdb.SelectEvents(10)); n != 0 {
		t.Fatalf("1 miss should not alert, got %d events", n)
	}
	compareHosts(none)
	compareHosts(none) // still offline: no repeat
	compareHosts(seen("10.0.0.9"))
	compareHosts(map[string]models.Host{"aa:bb:cc:00:00:02": {Mac: "aa:bb:cc:00:00:02", IP: "10.0.0.7"}})

	var kinds []string
	ev := gdb.SelectEvents(10)
	for i := len(ev) - 1; i >= 0; i-- {
		kinds = append(kinds, ev[i].Kind)
	}
	want := []string{"offline", "online", "ip", "new"}
	if len(kinds) != len(want) {
		t.Fatalf("events %v, want %v", kinds, want)
	}
	for i := range want {
		if kinds[i] != want[i] {
			t.Fatalf("events %v, want %v", kinds, want)
		}
	}
	if LastScan().IsZero() {
		t.Error("lastScan not set")
	}
}

func TestQuiet(t *testing.T) {
	at := func(h, m int) time.Time { return time.Date(2026, 1, 1, h, m, 0, 0, time.Local) }
	for _, tc := range []struct {
		h, m     int
		from, to string
		want     bool
	}{
		{23, 0, "22:00", "07:00", true},
		{6, 59, "22:00", "07:00", true},
		{7, 0, "22:00", "07:00", false},
		{12, 0, "09:00", "17:00", true},
		{18, 0, "09:00", "17:00", false},
		{12, 0, "", "", false},
		{12, 0, "bad", "07:00", false},
	} {
		if got := notify.Quiet(at(tc.h, tc.m), tc.from, tc.to); got != tc.want {
			t.Errorf("%+v: got %v", tc, got)
		}
	}
}

// A scan holds a snapshot of hosts; its write must not undo edits made meanwhile.
func TestScanKeepsConcurrentEdits(t *testing.T) {
	conf.AppConfig = models.Conf{UseDB: "sqlite", DBPath: t.TempDir() + "/t.db", OfflineAfter: 2}
	gdb.Start()

	gdb.Update("now", models.Host{Mac: "aa:00:00:00:00:01", IP: "10.0.0.2"})
	snap, _ := gdb.Select("now")

	edited := snap[0]
	edited.Tags, edited.Watch, edited.Name = "family", 1, "Dad"
	gdb.Update("now", edited)      // user edit during the scan
	gdb.SetPorts(snap[0].ID, "22") // security scan during the scan

	snap[0].Now, snap[0].IP = 1, "10.0.0.3"
	gdb.UpdateScan(snap[0], snap[0].Name != "", snap[0].Type != "") // scan found no name/type

	h := gdb.SelectByID(snap[0].ID)
	if h.Name != "Dad" || h.Tags != "family" || h.Watch != 1 || h.Ports != "22" || h.IP != "10.0.0.3" || h.Now != 1 {
		t.Errorf("got %+v", h)
	}
}
