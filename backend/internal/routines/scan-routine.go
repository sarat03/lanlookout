package routines

import (
	"sync"
	"time"

	"github.com/sarat03/lanlookout/internal/arp"
	"github.com/sarat03/lanlookout/internal/check"
	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/identify"
	"github.com/sarat03/lanlookout/internal/influx"
	"github.com/sarat03/lanlookout/internal/leases"
	"github.com/sarat03/lanlookout/internal/models"
	"github.com/sarat03/lanlookout/internal/mqtt"
	"github.com/sarat03/lanlookout/internal/notify"
	"github.com/sarat03/lanlookout/internal/prometheus"
)

func startScan(quit chan bool) {

	for {
		foundHosts := arp.Scan(conf.AppConfig.Ifaces, conf.AppConfig.ArpArgs, conf.AppConfig.ArpStrs)

		// Make map of found hosts
		foundHostsMap := make(map[string]models.Host)
		for _, fHost := range foundHosts {
			foundHostsMap[fHost.Mac] = fHost
		}

		// Leases only add hosts arp-scan missed; arp data wins.
		if conf.AppConfig.DHCPLeases != "" {
			for _, lHost := range leases.Read(conf.AppConfig.DHCPLeases) {
				if _, ok := foundHostsMap[lHost.Mac]; !ok {
					foundHostsMap[lHost.Mac] = lHost
				}
			}
		}

		select {
		case <-quit: // restarted while scanning: the new routine does the compare
			return
		default:
		}
		compareHosts(foundHostsMap)

		select {
		case <-quit:
			return
		case <-time.After(time.Duration(max(conf.AppConfig.Timeout, 1)) * time.Second):
		}
	}
}

// offlineAfter - missed scans before a watched host counts as offline.
// arp-scan sometimes misses a reply; raise OFFLINE_AFTER if alerts flap.
func offlineAfter() int {
	return max(conf.AppConfig.OfflineAfter, 1)
}

// lastScan - when compareHosts last finished, for the "Scan now" button
var lastScan time.Time

// LastScan - time the last scan finished (zero before the first one)
func LastScan() time.Time {
	compareMu.Lock()
	defer compareMu.Unlock()
	return lastScan
}

// misses - consecutive missed scans per MAC.
// ponytail: in memory, so a restart can repeat one offline alert.
var misses = map[string]int{}

// compareMu - a rescan can start while the old scan goroutine is still comparing.
var compareMu sync.Mutex

func compareHosts(foundHostsMap map[string]models.Host) {

	compareMu.Lock()
	defer compareMu.Unlock()

	allHosts, ok := gdb.Select("now")
	if !ok {
		return
	}

	for _, aHost := range allHosts {

		fHost, exists := foundHostsMap[aHost.Mac]
		oldIP := aHost.IP
		hadName, hadType := aHost.Name != "", aHost.Type != ""
		if exists {

			aHost.Iface = fHost.Iface
			aHost.IP = fHost.IP
			aHost.Date = fHost.Date
			aHost.Now = 1

			// Retry names the first lookup missed (DNS may learn them later).
			// A name set by the user is never overwritten.
			if aHost.Name == "" {
				aHost.Name, aHost.DNS = check.DNS(aHost)
				if aHost.Name == "" {
					aHost.Name = fHost.Name // from DHCP lease, if any
				}
			}

			delete(foundHostsMap, aHost.Mac)

			if conf.AppConfig.IdentifyEnable {
				identify.Apply(&aHost)
			}

			if aHost.Watch == 1 && misses[aHost.Mac] >= offlineAfter() {
				notify.Status(aHost)
			}
			misses[aHost.Mac] = 0

			if aHost.Known == 1 && oldIP != "" && oldIP != aHost.IP {
				notify.IPChange(aHost, oldIP)
			}

		} else {
			aHost.Now = 0
			misses[aHost.Mac]++
			if aHost.Watch == 1 && misses[aHost.Mac] == offlineAfter() {
				notify.Status(aHost)
			}
		}
		gdb.UpdateScan(aHost, !hadName && aHost.Name != "", !hadType && aHost.Type != "")

		aHost.ID = 0
		aHost.Date = time.Now().Format("2006-01-02 15:04:05")
		gdb.Update("history", aHost)

		if conf.AppConfig.InfluxEnable {
			influx.Add(conf.AppConfig, aHost)
		}
		if conf.AppConfig.PrometheusEnable {
			prometheus.Add(aHost)
		}
		if conf.AppConfig.MqttURL != "" {
			mqtt.Add(aHost)
		}
	}

	for _, fHost := range foundHostsMap {

		if name, dns := check.DNS(fHost); name != "" {
			fHost.Name, fHost.DNS = name, dns
		}
		if conf.AppConfig.IdentifyEnable {
			identify.Apply(&fHost)
		}
		notify.Unknown(fHost) // Log and Shoutrrr

		gdb.Update("now", fHost)
	}

	lastScan = time.Now()
}
