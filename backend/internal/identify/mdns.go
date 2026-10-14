package identify

import (
	"log/slog"
	"net"
	"strings"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

var mdnsGroup = &net.UDPAddr{IP: net.IPv4(224, 0, 0, 251), Port: 5353}

// mdnsLoop - listen to all mDNS traffic and ask "what services exist?" every 5 minutes.
// Answers are multicast, so we also learn from queries other devices trigger.
func mdnsLoop() {

	conn, err := net.ListenMulticastUDP("udp4", nil, mdnsGroup)
	if err != nil {
		slog.Warn("mDNS listener off", "err", err)
		return
	}

	asked := map[string]bool{}
	go func() {
		for {
			query(conn, "_services._dns-sd._udp.local.")
			mu.Lock()
			types := []string{}
			for t := range asked {
				types = append(types, t)
			}
			mu.Unlock()
			for _, t := range types {
				query(conn, t)
				time.Sleep(100 * time.Millisecond)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	buf := make([]byte, 9000)
	for {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		for _, t := range parseMDNS(src.IP.String(), buf[:n]) {
			mu.Lock()
			isNew := !asked[t]
			asked[t] = true
			mu.Unlock()
			if isNew {
				query(conn, t)
			}
		}
	}
}

func query(conn *net.UDPConn, name string) {
	b := dnsmessage.NewBuilder(nil, dnsmessage.Header{})
	b.StartQuestions()
	b.Question(dnsmessage.Question{Name: dnsmessage.MustNewName(name), Type: dnsmessage.TypePTR, Class: dnsmessage.ClassINET})
	msg, err := b.Finish()
	if err == nil {
		conn.WriteToUDP(msg, mdnsGroup)
	}
}

// parseMDNS - record what the sender announced; returns service types to ask about
func parseMDNS(ip string, pkt []byte) (types []string) {

	var p dnsmessage.Parser
	h, err := p.Start(pkt)
	if err != nil || !h.Response {
		return nil
	}
	p.SkipAllQuestions()

	mu.Lock()
	defer mu.Unlock()
	info := get(byIP, ip)

	for _, section := range []func() ([]dnsmessage.Resource, error){p.AllAnswers, p.AllAuthorities, p.AllAdditionals} {
		rrs, err := section()
		if err != nil {
			break
		}
		for _, rr := range rrs {
			name := rr.Header.Name.String()
			switch b := rr.Body.(type) {
			case *dnsmessage.PTRResource:
				target := b.PTR.String()
				if name == "_services._dns-sd._udp.local." {
					types = append(types, target)
					continue
				}
				if svc := serviceOf(name); svc != "" {
					info.Services[svc] = true
					// "Living Room TV._googlecast._tcp.local." -> "Living Room TV"
					if inst, _, ok := strings.Cut(target, "._"); ok && info.Name == "" && svc != "device-info" {
						info.Name = unescape(inst)
					}
				}
			case *dnsmessage.TXTResource:
				if svc := serviceOf(name); svc != "" {
					info.Services[svc] = true
				}
				txt(info, b.TXT)
			case *dnsmessage.AResource:
				if strings.HasSuffix(name, ".local.") && info.Hostname == "" {
					info.Hostname = strings.TrimSuffix(name, ".local.")
				}
			}
		}
	}
	return types
}

// serviceOf - "x._airplay._tcp.local." -> "airplay"
func serviceOf(name string) string {
	parts := strings.Split(name, ".")
	for i := len(parts) - 1; i > 0; i-- {
		if (parts[i] == "_tcp" || parts[i] == "_udp") && strings.HasPrefix(parts[i-1], "_") {
			return strings.TrimPrefix(parts[i-1], "_")
		}
	}
	return ""
}

func txt(info *Info, kv []string) {
	for _, s := range kv {
		k, v, ok := strings.Cut(s, "=")
		if !ok || v == "" {
			continue
		}
		switch strings.ToLower(k) {
		case "md", "model", "am", "usb_mdl", "ty":
			if info.Model == "" {
				info.Model = v
			}
		case "fn":
			info.Name = v
		case "usb_mfg", "manufacturer":
			if info.Maker == "" {
				info.Maker = v
			}
		}
	}
}

func unescape(s string) string {
	return strings.NewReplacer(`\ `, " ", `\032`, " ", `\039`, "'", `\.`, ".").Replace(s)
}
