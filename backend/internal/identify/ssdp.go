package identify

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ssdpGroup = &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1900}

// ssdpLoop - ask every UPnP device to describe itself, every 10 minutes
func ssdpLoop() {
	for {
		ssdpSearch()
		time.Sleep(10 * time.Minute)
	}
}

func ssdpSearch() {

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		slog.Warn("SSDP search off", "err", err)
		return
	}
	defer conn.Close()

	msg := "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1900\r\nMAN: \"ssdp:discover\"\r\nMX: 2\r\nST: ssdp:all\r\n\r\n"
	conn.WriteToUDP([]byte(msg), ssdpGroup)
	conn.SetReadDeadline(time.Now().Add(4 * time.Second))

	locations := map[string]string{} // ip -> description URL
	buf := make([]byte, 4096)
	for {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		ip := src.IP.String()
		resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf[:n])), nil)
		if err != nil {
			continue
		}
		st, server, loc := resp.Header.Get("ST"), resp.Header.Get("SERVER"), resp.Header.Get("LOCATION")

		mu.Lock()
		info := get(byIP, ip)
		if server != "" {
			info.SSDP["server"] = server
		}
		if strings.Contains(st, "InternetGatewayDevice") {
			Gateways[ip] = true
		}
		mu.Unlock()

		// Only fetch descriptions the device serves itself (no redirects to other hosts).
		if u, err := url.Parse(loc); err == nil && u.Hostname() == ip && (u.Scheme == "http" || u.Scheme == "https") {
			locations[ip] = loc
		}
	}

	client := &http.Client{Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	for ip, loc := range locations {
		describe(client, ip, loc)
	}
}

type upnpRoot struct {
	Device struct {
		DeviceType   string `xml:"deviceType"`
		FriendlyName string `xml:"friendlyName"`
		Manufacturer string `xml:"manufacturer"`
		ModelName    string `xml:"modelName"`
	} `xml:"device"`
}

func describe(client *http.Client, ip, loc string) {

	resp, err := client.Get(loc)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var root upnpRoot
	if xml.NewDecoder(io.LimitReader(resp.Body, 256<<10)).Decode(&root) != nil {
		return
	}
	d := root.Device

	mu.Lock()
	defer mu.Unlock()
	info := get(byIP, ip)
	info.SSDP["deviceType"] = d.DeviceType
	info.Maker = first(info.Maker, d.Manufacturer)
	info.Model = first(info.Model, d.ModelName)
	info.Name = first(info.Name, d.FriendlyName)
	if strings.Contains(d.DeviceType, "InternetGatewayDevice") {
		Gateways[ip] = true
	}
}
