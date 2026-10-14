package secscan

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestAdded(t *testing.T) {
	if Added("", []int{23}) != nil {
		t.Error("first scan should not alert")
	}
	if a := Added("none", []int{23}); len(a) != 1 || a[0] != 23 {
		t.Errorf("got %v", a)
	}
	if a := Added("22,23", []int{22, 23, 445}); len(a) != 1 || a[0] != 445 {
		t.Errorf("got %v", a)
	}
	if join(nil) != "none" || join([]int{22, 80}) != "22,80" {
		t.Error("join")
	}
}

func TestOpenPorts(t *testing.T) {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	defer l.Close()
	port, _ := strconv.Atoi(l.Addr().(*net.TCPAddr).AddrPort().String()[len("127.0.0.1:"):])
	open := OpenPorts("127.0.0.1", []int{port, 1}, time.Second)
	if len(open) != 1 || open[0] != port {
		t.Errorf("got %v want [%d]", open, port)
	}
}
