package identify

import (
	"testing"

	"golang.org/x/net/dns/dnsmessage"

	"github.com/sarat03/lanlookout/internal/models"
)

func info(model, host string, svc ...string) Info {
	i := Info{Model: model, Hostname: host, Services: map[string]bool{}, SSDP: map[string]string{}}
	for _, s := range svc {
		i.Services[s] = true
	}
	return i
}

func TestClassify(t *testing.T) {
	for _, tc := range []struct {
		i       Info
		vendor  string
		typ, os string
	}{
		{info("iPhone14,2", "", "companion-link", "device-info"), "", "phone", "iOS / macOS"},
		{info("MacBookPro18,1", "", "device-info"), "", "laptop", "iOS / macOS"},
		{info("Chromecast", "", "googlecast"), "", "tv", "Google Cast"},
		{info("Google Nest Mini", "", "googlecast"), "", "speaker", "Google Cast"},
		{info("", "", "ipp"), "", "printer", ""},
		{info("", "", "hap"), "", "iot", ""},
		{info("", "", "sonos"), "", "speaker", ""},
		{info("", "DESKTOP-AB12CD"), "", "desktop", "Windows"},
		{info("", "Galaxy-S23"), "", "phone", ""},
		{info("", ""), "Espressif Inc.", "iot", ""},
		{Info{DHCP: "android-dhcp-14", Services: map[string]bool{}, SSDP: map[string]string{}}, "", "phone", "Android"},
		{Info{Services: map[string]bool{}, SSDP: map[string]string{"deviceType": "urn:schemas-upnp-org:device:InternetGatewayDevice:1"}}, "", "router", ""},
		{info("", "shelly-plug-s-1A2B", "hap"), "", "plug", ""},
		{info("", "ecobee3-lite"), "", "thermostat", ""},
		{info("", "dyson-fan"), "", "fan", ""},
		{info("", "wiz_bulb_3f"), "Espressif Inc.", "light", ""},
		{info("", ""), "Tesla Inc.", "car", ""},
		{info("", ""), "", "", ""},
	} {
		typ, os := Classify(tc.i, tc.vendor)
		if typ != tc.typ || os != tc.os {
			t.Errorf("%+v %q: got %q %q, want %q %q", tc.i, tc.vendor, typ, os, tc.typ, tc.os)
		}
	}
}

func TestParseMDNSAndApply(t *testing.T) {
	name := func(s string) dnsmessage.Name { return dnsmessage.MustNewName(s) }
	b := dnsmessage.NewBuilder(nil, dnsmessage.Header{Response: true})
	b.StartAnswers()
	hdr := func(n string, typ dnsmessage.Type) dnsmessage.ResourceHeader {
		return dnsmessage.ResourceHeader{Name: name(n), Type: typ, Class: dnsmessage.ClassINET, TTL: 120}
	}
	b.PTRResource(hdr("_googlecast._tcp.local.", dnsmessage.TypePTR), dnsmessage.PTRResource{PTR: name("Living Room._googlecast._tcp.local.")})
	b.TXTResource(hdr("Living Room._googlecast._tcp.local.", dnsmessage.TypeTXT), dnsmessage.TXTResource{TXT: []string{"md=Chromecast Ultra", "fn=Living Room TV"}})
	b.AResource(hdr("cc-1234.local.", dnsmessage.TypeA), dnsmessage.AResource{A: [4]byte{10, 0, 0, 9}})
	b.PTRResource(hdr("_services._dns-sd._udp.local.", dnsmessage.TypePTR), dnsmessage.PTRResource{PTR: name("_googlecast._tcp.local.")})
	pkt, err := b.Finish()
	if err != nil {
		t.Fatal(err)
	}

	types := parseMDNS("10.0.0.9", pkt)
	if len(types) != 1 || types[0] != "_googlecast._tcp.local." {
		t.Errorf("types %v", types)
	}

	h := models.Host{IP: "10.0.0.9", Mac: "aa:bb:cc:dd:ee:ff"}
	Apply(&h)
	if h.Model != "Chromecast Ultra" || h.Type != "tv" || h.Name != "Living Room TV" || h.Services != "googlecast" || h.OS != "Google Cast" {
		t.Errorf("applied %+v", h)
	}

	// user-set name and type stay
	h2 := models.Host{IP: "10.0.0.9", Name: "Mine", Type: "other"}
	Apply(&h2)
	if h2.Name != "Mine" || h2.Type != "other" {
		t.Errorf("overwrote user fields %+v", h2)
	}
}

func TestParseDHCP(t *testing.T) {
	p := make([]byte, 240)
	p[0], p[1], p[2] = 1, 1, 6
	copy(p[28:], []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66})
	copy(p[236:], []byte{99, 130, 83, 99})
	p = append(p, 53, 1, 1) // DISCOVER
	p = append(p, 12, 6)
	p = append(p, "iPhone"...)
	p = append(p, 55, 5, 1, 121, 108, 114, 252)
	p = append(p, 255)

	parseDHCP(p)
	i := byMAC["11:22:33:44:55:66"]
	if i == nil || i.Hostname != "iPhone" || i.DHCPOpts != "1,121,108,114,252" {
		t.Fatalf("got %+v", i)
	}
	if typ, os := Classify(*i, ""); typ != "phone" || os != "iOS / macOS" {
		t.Errorf("classify %q %q", typ, os)
	}
	parseDHCP(p[:100]) // short packet: no panic
}
