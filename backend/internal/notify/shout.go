package notify

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/nicholas-fedor/shoutrrr"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
)

// Unknown - send message to log and shoutrrr
func Unknown(host models.Host) {

	name := host.DNS
	if name == "" {
		name = host.Name // e.g. from a DHCP lease
	}
	msg := fmt.Sprintf("Unknown host found. Name: '%s', IP: '%s', MAC: '%s', Hw: '%s', Iface: '%s'", name, host.IP, host.Mac, host.Hw, host.Iface)

	event("new", host, msg, conf.AppConfig.AlertNew)
}

// Status - a watched host went offline or came back
func Status(host models.Host) {

	kind, state := "offline", "is offline"
	if host.Now == 1 {
		kind, state = "online", "is back online"
	}
	msg := fmt.Sprintf("Host '%s' %s. IP: '%s', MAC: '%s'", label(host), state, host.IP, host.Mac)

	event(kind, host, msg, conf.AppConfig.AlertOffline)
}

// IPChange - a known host moved to a new IP
func IPChange(host models.Host, oldIP string) {

	msg := fmt.Sprintf("Host '%s' changed IP: '%s' -> '%s'. MAC: '%s'", label(host), oldIP, host.IP, host.Mac)

	event("ip", host, msg, conf.AppConfig.AlertIP)
}

// Event - an alert about one host or the whole network (empty Mac).
// kind: outage, wan_ip, port, threat... enabled: that alert type is on.
func Event(kind string, host models.Host, msg string, enabled bool) {
	event(kind, host, msg, enabled)
}

func label(host models.Host) string {
	if host.Name != "" {
		return host.Name
	}
	return host.Mac
}

// event - always logged and kept in Recent events; sent only if that alert is on
// and it's not quiet hours.
func event(kind string, host models.Host, msg string, enabled bool) {

	slog.Warn(msg)
	gdb.AddEvent(models.Event{
		Date: time.Now().Format("2006-01-02 15:04:05"),
		Kind: kind,
		Mac:  host.Mac,
		Msg:  msg,
	})

	if enabled && !Quiet(time.Now(), conf.AppConfig.QuietFrom, conf.AppConfig.QuietTo) {
		shout(msg)
	}
}

// Quiet - is now inside quiet hours "HH:MM"-"HH:MM" (may wrap past midnight)
func Quiet(now time.Time, from, to string) bool {

	f, err1 := time.Parse("15:04", from)
	t, err2 := time.Parse("15:04", to)
	if err1 != nil || err2 != nil || from == to {
		return false
	}
	m := now.Hour()*60 + now.Minute()
	fm := f.Hour()*60 + f.Minute()
	tm := t.Hour()*60 + t.Minute()

	if fm < tm {
		return m >= fm && m < tm
	}
	return m >= fm || m < tm
}

// Test Shoutrrr notification
func Test() {

	msg := "test notification"
	slog.Info("Sending " + msg)
	shout(msg)
}

// shout - send msg to Shoutrrr
func shout(msg string) {

	hostname, _ := os.Hostname()
	prefix := "LanLookout on '" + hostname + "': "

	if conf.AppConfig.ShoutURL != "" {
		err := shoutrrr.Send(conf.AppConfig.ShoutURL, prefix+msg)
		if err != nil {
			slog.Error("Notification failed (shoutrrr): ", "", err)
		}
	}
}
