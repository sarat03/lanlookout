package identify

import (
	"log/slog"
	"net"
	"strconv"
	"strings"
)

// dhcpListen - read DHCP requests devices broadcast when they join.
// Needs UDP port 67 and host networking; skipped if the port is taken.
// ponytail: only sees broadcasts; renewals sent straight to the router are missed.
func dhcpListen() {

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 67})
	if err != nil {
		slog.Warn("DHCP fingerprinting off (port 67 unavailable)", "err", err)
		return
	}

	buf := make([]byte, 1500)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		parseDHCP(buf[:n])
	}
}

// parseDHCP - BOOTP request: client MAC at 28, magic cookie at 236, options after.
func parseDHCP(p []byte) {

	if len(p) < 240 || p[0] != 1 || p[1] != 1 || p[2] != 6 ||
		p[236] != 99 || p[237] != 130 || p[238] != 83 || p[239] != 99 {
		return
	}
	mac := net.HardwareAddr(p[28:34]).String()

	var hostname, vendor, params string
	for i := 240; i < len(p); {
		code := p[i]
		if code == 255 {
			break
		}
		if code == 0 || i+1 >= len(p) {
			i++
			continue
		}
		l := int(p[i+1])
		if i+2+l > len(p) {
			break
		}
		v := p[i+2 : i+2+l]
		switch code {
		case 12:
			hostname = string(v)
		case 60:
			vendor = string(v)
		case 55:
			s := make([]string, len(v))
			for j, b := range v {
				s[j] = strconv.Itoa(int(b))
			}
			params = strings.Join(s, ",")
		}
		i += 2 + l
	}

	mu.Lock()
	defer mu.Unlock()
	info := get(byMAC, mac)
	info.Hostname = first(hostname, info.Hostname)
	info.DHCP = first(vendor, info.DHCP)
	info.DHCPOpts = first(params, info.DHCPOpts)
}
