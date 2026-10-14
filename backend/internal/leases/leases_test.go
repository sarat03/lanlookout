package leases

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	now := time.Unix(1000, 0)
	in := "2000 AA:BB:CC:DD:EE:01 192.168.1.5 laptop 01:aa\n" +
		"500 aa:bb:cc:dd:ee:02 192.168.1.6 old *\n" + // expired
		"0 aa:bb:cc:dd:ee:03 192.168.1.7 * *\n" + // infinite, no name
		"garbage\n"

	h := parse(bufio.NewScanner(strings.NewReader(in)), now)

	if len(h) != 2 || h[0].Mac != "aa:bb:cc:dd:ee:01" || h[0].Name != "laptop" || h[1].Name != "" || h[1].IP != "192.168.1.7" || h[0].Now != 1 {
		t.Fatalf("got %+v", h)
	}
}
