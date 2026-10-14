// Package identify works out what a device is from what it announces on the LAN:
// mDNS/Bonjour services and TXT records, UPnP (SSDP) descriptions and the DHCP
// request fingerprint. The rules are our own heuristics, not a vendor database.
package identify

import (
	"log/slog"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/sarat03/lanlookout/internal/models"
)

// Info - everything learned about one device
type Info struct {
	Hostname string            // from mDNS A records or DHCP option 12
	Name     string            // friendly name (UPnP friendlyName, mDNS fn=, instance name)
	Model    string            // e.g. "Chromecast Ultra", "iPhone14,2"
	Maker    string            // UPnP manufacturer
	Services map[string]bool   // e.g. "googlecast", "airplay", "ipp"
	DHCP     string            // vendor class (option 60)
	DHCPOpts string            // parameter request list (option 55), e.g. "1,3,6,15"
	SSDP     map[string]string // deviceType / SERVER header hints
}

var (
	mu    sync.Mutex
	byIP  = map[string]*Info{}
	byMAC = map[string]*Info{}
	// Gateways - IPs that announced a UPnP InternetGatewayDevice (for the security check)
	Gateways = map[string]bool{}
)

func get(m map[string]*Info, key string) *Info {
	i, ok := m[key]
	if !ok {
		i = &Info{Services: map[string]bool{}, SSDP: map[string]string{}}
		m[key] = i
	}
	return i
}

// UPnPGateways - IPs of routers that answer UPnP IGD discovery
func UPnPGateways() (ips []string) {
	mu.Lock()
	defer mu.Unlock()
	for ip := range Gateways {
		ips = append(ips, ip)
	}
	return ips
}

// Start - listeners and periodic probes. Each one logs and gives up on its own
// if it can't bind (e.g. port 67 taken by a DHCP server on the same host).
func Start() {
	go mdnsLoop()
	go ssdpLoop()
	go dhcpListen()
	slog.Info("Device identification started (mDNS, SSDP, DHCP)")
}

// Apply - fill the automatic fields of a host. User-set Name and Type are kept.
func Apply(h *models.Host) {

	mu.Lock()
	var merged Info
	merged.Services = map[string]bool{}
	merged.SSDP = map[string]string{}
	for _, src := range []*Info{byIP[h.IP], byMAC[strings.ToLower(h.Mac)]} {
		if src == nil {
			continue
		}
		merged.Hostname = first(merged.Hostname, src.Hostname)
		merged.Name = first(merged.Name, src.Name)
		merged.Model = first(merged.Model, src.Model)
		merged.Maker = first(merged.Maker, src.Maker)
		merged.DHCP = first(merged.DHCP, src.DHCP)
		merged.DHCPOpts = first(merged.DHCPOpts, src.DHCPOpts)
		for k := range src.Services {
			merged.Services[k] = true
		}
		for k, v := range src.SSDP {
			merged.SSDP[k] = v
		}
	}
	mu.Unlock()

	typ, os := Classify(merged, h.Hw)

	if merged.Model != "" {
		h.Model = merged.Model
	}
	if os != "" {
		h.OS = os
	}
	if len(merged.Services) > 0 {
		svc := make([]string, 0, len(merged.Services))
		for k := range merged.Services {
			svc = append(svc, k)
		}
		sort.Strings(svc)
		h.Services = strings.Join(svc, ",")
	}
	if h.Type == "" {
		h.Type = typ
	}
	if h.Name == "" {
		h.Name = first(merged.Name, merged.Hostname)
	}
}

func first(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func has(s string, subs ...string) bool {
	s = strings.ToLower(s)
	return slices.ContainsFunc(subs, func(x string) bool { return strings.Contains(s, x) })
}

// Classify - guess device type (one of the frontend's hostTypes) and OS.
// ponytail: a few dozen hand-written rules; extend them as devices show up wrong.
func Classify(i Info, vendor string) (typ, os string) {

	model, host, ssdp := i.Model, i.Hostname+" "+i.Name, i.SSDP["deviceType"]+" "+i.SSDP["server"]+" "+i.Maker
	svc := i.Services

	// OS from the DHCP fingerprint and announcements
	switch {
	case has(i.DHCP, "android"):
		os = "Android"
	case has(i.DHCP, "msft"):
		os = "Windows"
	case has(i.DHCP, "dhcpcd", "udhcp", "linux"):
		os = "Linux"
	case svc["companion-link"] || svc["device-info"] && has(model, "iphone", "ipad", "mac", "watch") || appleDHCP(i.DHCPOpts):
		os = "iOS / macOS"
	case svc["googlecast"]:
		os = "Google Cast"
	}
	if os == "" && has(host, "desktop-", "laptop-") {
		os = "Windows"
	}

	// Type: model first (most precise), then services, then names
	switch {
	case has(model, "iphone"):
		typ = "phone"
	case has(model, "ipad"):
		typ = "tablet"
	case has(model, "macbook"):
		typ = "laptop"
	case has(model, "imac", "macmini", "macpro"): // ponytail: "Mac14,x" ids mix laptops and desktops, left unset
		typ = "desktop"
	case has(model, "watch"):
		typ = "watch"
	case has(model, "appletv") || has(ssdp, "roku", "samsung tv", "webos", "bravia", "tizen"):
		typ = "tv"
	case has(model, "audioaccessory", "homepod", "nest mini", "nest audio", "google home", "sonos", "echo"):
		typ = "speaker"
	case svc["googlecast"]:
		typ = "tv"
	case svc["ipp"] || svc["ipps"] || svc["printer"] || svc["pdl-datastream"] || has(ssdp, "printer"):
		typ = "printer"
	case has(ssdp, "internetgatewaydevice"):
		typ = "router"
	case has(ssdp, "xbox", "playstation") || has(host, "xbox", "playstation", "ps4", "ps5", "nintendo"):
		typ = "console"
	case svc["sonos"] || svc["spotify-connect"] || svc["raop"] && !svc["airplay"]:
		typ = "speaker"
	// Smart-home kinds before the generic iot catch-all; matched on model, hostname and UPnP text.
	// ponytail: name keywords only, misses devices with generic names (they stay iot)
	case has(model+" "+host+" "+ssdp, "plug", "outlet", "socket"):
		typ = "plug"
	case has(model+" "+host+" "+ssdp, "thermostat", "ecobee", "tado", "netatmo"):
		typ = "thermostat"
	case has(model+" "+host+" "+ssdp, "fan", "dyson", "air purifier"):
		typ = "fan"
	case has(model+" "+host+" "+ssdp, "bulb", "light", "lamp", "lifx", "wiz", "nanoleaf"):
		typ = "light"
	case has(vendor, "tesla") || has(model+" "+host, "tesla", "polestar", "rivian"):
		typ = "car"
	case svc["hap"] || svc["matter"] || svc["matterc"] || has(vendor, "espressif", "tuya", "shelly", "sonoff", "philips lighting", "signify"):
		typ = "iot"
	case svc["smb"] || svc["adisk"] || svc["afpovertcp"] || has(vendor, "synology", "qnap"):
		typ = "server"
	case has(ssdp, "mediarenderer"):
		typ = "tv"
	case has(host, "iphone", "android", "galaxy", "pixel", "oneplus", "redmi"):
		typ = "phone"
	case has(host, "ipad", "tablet", "-tab"):
		typ = "tablet"
	case has(host, "macbook", "laptop", "thinkpad", "xps"):
		typ = "laptop"
	case has(host, "desktop", "imac", "-pc"):
		typ = "desktop"
	case has(host, "printer", "epson", "brother", "canon") || strings.HasPrefix(strings.ToLower(i.Hostname), "hp"):
		typ = "printer"
	case has(host, "camera", "cam-", "doorbell") || svc["rtsp"]:
		typ = "camera"
	case os == "Android":
		typ = "phone"
	}
	return typ, os
}

// appleDHCP - Apple devices ask for option 108 (IPv6-only) together with 114 and 252.
func appleDHCP(opts string) bool {
	p := strings.Split(opts, ",")
	return slices.Contains(p, "108") && slices.Contains(p, "114") && slices.Contains(p, "252")
}
