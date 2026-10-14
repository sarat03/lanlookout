package check

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/sarat03/lanlookout/internal/models"
)

// DNS - returns DNS names of a host.
// Tries the system resolver (e.g. OpenWrt's dnsmasq knows DHCP hostnames),
// then asks the device itself over mDNS (Apple, Avahi/Linux, printers, Chromecasts).
func DNS(host models.Host) (name, dns string) {

	dnsNames, _ := net.LookupAddr(host.IP)

	if len(dnsNames) == 0 {
		dnsNames = mdnsLookup(host.IP)
	}

	return joinNames(dnsNames)
}

func joinNames(dnsNames []string) (name, dns string) {

	for i := range dnsNames {
		dnsNames[i] = strings.TrimSuffix(dnsNames[i], ".")
	}

	if len(dnsNames) > 0 {
		name = dnsNames[0]
		dns = strings.Join(dnsNames, " ")
	}

	return name, dns
}

// mdnsLookup sends a reverse (PTR) query straight to the device on port 5353.
// mDNS responders answer such "legacy unicast" queries with a normal DNS reply
// (RFC 6762 section 6.7), so Go's own resolver can read it: no extra dependency.
// ponytail: one device per query with a 1s timeout; Android and most IoT gear
// don't answer reverse mDNS, those stay unnamed until DNS knows them.
func mdnsLookup(ip string) []string {

	if net.ParseIP(ip) == nil {
		return nil
	}

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "udp", net.JoinHostPort(ip, "5353"))
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	names, _ := r.LookupAddr(ctx, ip)

	return names
}
