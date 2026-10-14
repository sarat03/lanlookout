package netcheck

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
)

func TestLatency(t *testing.T) {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	defer l.Close()
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()
	if _, ok := Latency([]string{l.Addr().String(), "127.0.0.1:1"}, time.Second); !ok {
		t.Error("one reachable target should count as up")
	}
	if _, ok := Latency([]string{"127.0.0.1:1"}, time.Second); ok {
		t.Error("nothing reachable should be down")
	}
}

func TestSpeedAndOutage(t *testing.T) {
	conf.AppConfig = models.Conf{UseDB: "sqlite", DBPath: t.TempDir() + "/t.db", AlertNet: true}
	gdb.Start()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/__down" {
			n, _ := strconv.Atoi(r.URL.Query().Get("bytes"))
			w.Write(make([]byte, n))
			return
		}
		io.Copy(io.Discard, r.Body)
	}))
	defer srv.Close()

	down, up, err := Speed(srv.URL)
	if err != nil || down <= 0 || up <= 0 {
		t.Fatalf("speed %v %v %v", down, up, err)
	}

	t0 := time.Now()
	track(false, t0)
	track(false, t0.Add(time.Minute)) // still down: no new start
	if !downSince.Equal(t0) {
		t.Fatal("outage start moved")
	}
	track(true, t0.Add(10*time.Minute))
	time.Sleep(200 * time.Millisecond) // alert is sent in a goroutine
	ev := gdb.SelectEvents(5)
	if len(ev) != 1 || ev[0].Kind != "outage" {
		t.Fatalf("events %+v", ev)
	}

	trackIP("1.2.3.4")
	trackIP("1.2.3.4")
	trackIP("5.6.7.8")
	if ev := gdb.SelectEvents(5); len(ev) != 2 || ev[0].Kind != "wan_ip" {
		t.Fatalf("events %+v", ev)
	}
}
