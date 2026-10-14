package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/dnslog"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/models"
	"github.com/sarat03/lanlookout/internal/netcheck"
	"github.com/sarat03/lanlookout/internal/secscan"
)

// netStatus godoc
// @Summary      Internet status
// @Description  Up/down, outage start, public IP and the latest check
// @Tags         internet
// @Produce      json
// @Success      200  {object}  netcheck.Status
// @Router       /net/status [get]
func netStatus(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, netcheck.Current())
}

// netStats godoc
// @Summary      Internet checks
// @Description  Latency, up/down and speed test results, oldest first
// @Tags         internet
// @Produce      json
// @Param        hours  query  int  false  "How far back (default 24)"
// @Success      200  {array}  models.NetStat
// @Router       /net/stats [get]
func netStats(c *gin.Context) {
	hours, err := strconv.Atoi(c.Query("hours"))
	if err != nil || hours < 1 {
		hours = 24
	}
	since := time.Now().Add(-time.Duration(hours) * time.Hour).Format("2006-01-02 15:04:05")
	c.IndentedJSON(http.StatusOK, gdb.SelectNetStats(since))
}

// speedTest godoc
// @Summary      Run a speed test now
// @Description  Takes 10-60 seconds. Result is stored with the other checks
// @Tags         internet
// @Produce      json
// @Success      200  {object}  models.NetStat
// @Router       /net/speedtest [get]
func speedTest(c *gin.Context) {
	stat := models.NetStat{Date: time.Now().Format("2006-01-02 15:04:05"), Up: 1}
	stat.LatencyMs, _ = netcheck.Latency(netcheck.Targets, 3*time.Second)
	down, up, err := netcheck.Speed(conf.AppConfig.SpeedURL)
	if err != nil {
		c.IndentedJSON(http.StatusBadGateway, err.Error())
		return
	}
	stat.DownMbps, stat.UpMbps = down, up
	stat.PublicIP = netcheck.Current().PublicIP
	gdb.AddNetStat(stat)
	c.IndentedJSON(http.StatusOK, stat)
}

// securityFindings godoc
// @Summary      Security findings
// @Description  Risky open ports and routers with UPnP, worst first
// @Tags         security
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       /security [get]
func securityFindings(c *gin.Context) {
	last, running := secscan.LastScan()
	lastStr := ""
	if !last.IsZero() {
		lastStr = last.Format("2006-01-02 15:04:05")
	}
	findings := secscan.Findings()
	if findings == nil {
		findings = []secscan.Finding{}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"Findings": findings, "LastScan": lastStr, "Running": running})
}

// securityScan godoc
// @Summary      Start a security scan
// @Description  Runs in the background; poll /security until Running is false
// @Tags         security
// @Produce      json
// @Success      200  {string}  string  "OK"
// @Router       /security/scan [get]
func securityScan(c *gin.Context) {
	go secscan.Scan()
	c.IndentedJSON(http.StatusOK, "OK")
}

// dnsTop godoc
// @Summary      Top domains of a device
// @Description  From Pi-hole / AdGuard Home query logs
// @Tags         dns
// @Produce      json
// @Param        ip    path   string  true   "Device IP"
// @Param        days  query  int     false  "Days back, 1 = today (default 1)"
// @Success      200  {object}  map[string]any
// @Router       /dns/{ip} [get]
func dnsTop(c *gin.Context) {
	days, err := strconv.Atoi(c.Query("days"))
	if err != nil || days < 1 {
		days = 1
	}
	since := time.Now().AddDate(0, 0, 1-days).Format("2006-01-02")
	c.IndentedJSON(http.StatusOK, gin.H{
		"Enabled": conf.AppConfig.DNSType != "" && conf.AppConfig.DNSURL != "",
		"Error":   dnslog.LastError(),
		"Top":     gdb.SelectDNSTop(c.Param("ip"), since, 50),
	})
}
