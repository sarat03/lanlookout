package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
	"github.com/sarat03/lanlookout/internal/notify"
	"github.com/sarat03/lanlookout/internal/routines"
)

// getVersion godoc
// @Summary      Get application version
// @Description  Returns the current running version of the application
// @Tags         system
// @Produce      json
// @Success      200  {string}  string
// @Router       /version [get]
func getVersion(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, conf.AppConfig.Version)
}

// triggerRescan godoc
// @Summary      Rescan all interfaces now
// @Description  Manually trigger rescan
// @Tags         system
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Router       /rescan [get]
func triggerRescan(c *gin.Context) {
	routines.ScanRestart()
	c.Status(http.StatusOK)
}

// getConfig godoc
// @Summary      Get application configuration
// @Description  Returns the current configuration used by the app
// @Tags         system
// @Produce      json
// @Success      200  {object}  models.Conf
// @Router       /config [get]
func getConfig(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, masked(conf.AppConfig)) // security.go
}

// notifyTest godoc
// @Summary      Send test notification
// @Description  Trigger a test notification to verify notification settings
// @Tags         system
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Router       /notify_test [get]
func notifyTest(c *gin.Context) {
	notify.Test()
	c.Status(http.StatusOK)
}

// getStatus godoc
// @Summary      Get network status
// @Description  Retrieve summary statistics of hosts, optionally filtered by interface
// @Tags         system
// @Produce      json
// @Param        iface  path      string  false  "Interface name (omit for all interfaces)"
// @Success      200    {object}  models.Stat
// @Router       /status/{iface} [get]
func getStatus(c *gin.Context) {
	var status models.Stat
	var searchHosts []models.Host

	allHosts, _ := gdb.Select("now")

	iface := c.Param("iface")
	iface = iface[1:]

	if iface != "" && iface != "undefined" {
		for _, host := range allHosts {
			if iface == host.Iface {
				searchHosts = append(searchHosts, host)
			}
		}
	} else {
		searchHosts = allHosts
	}

	for _, host := range searchHosts {
		status.Total = status.Total + 1

		if host.Known > 0 {
			status.Known = status.Known + 1
		} else {
			status.Unknown = status.Unknown + 1
		}
		if host.Now > 0 {
			status.Online = status.Online + 1
		} else {
			status.Offline = status.Offline + 1
		}
	}

	c.IndentedJSON(http.StatusOK, status)
}

// scanStatus godoc
// @Summary      Last scan time
// @Description  When the last scan finished. Poll after /rescan to know when results are in
// @Tags         system
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /scan_status [get]
func scanStatus(c *gin.Context) {
	last := ""
	if t := routines.LastScan(); !t.IsZero() {
		last = t.Format("2006-01-02 15:04:05.000")
	}
	c.IndentedJSON(http.StatusOK, gin.H{"LastScan": last})
}

// getEvents godoc
// @Summary      Recent events
// @Description  Alerts (new device, offline, online, IP change), newest first
// @Tags         system
// @Produce      json
// @Param        num  query  int  false  "Max events (default 100)"
// @Success      200  {array}  models.Event
// @Router       /events [get]
func getEvents(c *gin.Context) {
	num, err := strconv.Atoi(c.Query("num"))
	if err != nil || num < 1 {
		num = 100
	}
	c.IndentedJSON(http.StatusOK, gdb.SelectEvents(num))
}
