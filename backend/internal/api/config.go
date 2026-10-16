package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/gdb"
	"github.com/sarat03/lanlookout/internal/routines"
)

func saveConfigHandler(c *gin.Context) {

	conf.AppConfig.Host = c.PostForm("host")
	conf.AppConfig.Port = c.PostForm("port")
	conf.AppConfig.Theme = c.PostForm("theme")
	conf.AppConfig.Color = c.PostForm("color")
	conf.AppConfig.NodePath = c.PostForm("node")
	conf.AppConfig.ShoutURL = keep(c.PostForm("shout"), conf.AppConfig.ShoutURL, maskURL)

	conf.Write(conf.AppConfig)

	c.Redirect(http.StatusFound, c.Request.Referer())
}

func saveSettingsHandler(c *gin.Context) {

	conf.AppConfig.LogLevel = c.PostForm("log")
	conf.AppConfig.ArpArgs = c.PostForm("arpargs")
	conf.AppConfig.Ifaces = strings.Join(c.PostFormArray("ifaces"), " ") // one checkbox per interface
	conf.AppConfig.DHCPLeases = c.PostForm("leases")

	useDB := c.PostForm("usedb")
	pgConnect := keep(c.PostForm("pgconnect"), conf.AppConfig.PGConnect, maskURL)

	if useDB != conf.AppConfig.UseDB || pgConnect != conf.AppConfig.PGConnect {
		conf.AppConfig.UseDB = c.PostForm("usedb")
		conf.AppConfig.PGConnect = pgConnect
		gdb.Connect()
	}

	timeout := c.PostForm("timeout")
	trimHist := c.PostForm("trim")
	conf.AppConfig.Timeout, _ = strconv.Atoi(timeout)
	conf.AppConfig.TrimHist, _ = strconv.Atoi(trimHist)

	arpStrs := c.PostFormArray("arpstrs")
	conf.AppConfig.ArpStrs = []string{}
	for _, s := range arpStrs {
		if s != "" {
			conf.AppConfig.ArpStrs = append(conf.AppConfig.ArpStrs, s)
		}
	}

	conf.Write(conf.AppConfig)

	routines.ScanRestart()

	c.Redirect(http.StatusFound, c.Request.Referer())
}

func saveInfluxHandler(c *gin.Context) {

	conf.AppConfig.InfluxAddr = c.PostForm("addr")
	conf.AppConfig.InfluxToken = keep(c.PostForm("token"), conf.AppConfig.InfluxToken, maskAll)
	conf.AppConfig.InfluxOrg = c.PostForm("org")
	conf.AppConfig.InfluxBucket = c.PostForm("bucket")

	enable := c.PostForm("enable")
	skip := c.PostForm("skip")
	conf.AppConfig.InfluxEnable = enable == "on"
	conf.AppConfig.InfluxSkipTLS = skip == "on"

	conf.Write(conf.AppConfig)

	c.Redirect(http.StatusFound, c.Request.Referer())
}

func savePrometheusHandler(c *gin.Context) {
	enable := c.PostForm("enable")

	conf.AppConfig.PrometheusEnable = enable == "on"

	conf.Write(conf.AppConfig)

	c.Redirect(http.StatusFound, c.Request.Referer())
}

func saveMqttHandler(c *gin.Context) {

	conf.AppConfig.MqttURL = c.PostForm("url")
	conf.AppConfig.MqttUser = c.PostForm("user")
	conf.AppConfig.MqttPass = keep(c.PostForm("pass"), conf.AppConfig.MqttPass, maskAll)

	conf.Write(conf.AppConfig)

	c.Redirect(http.StatusFound, c.Request.Referer())
}

func saveAlertsHandler(c *gin.Context) {

	conf.AppConfig.AlertNew = c.PostForm("new") == "on"
	conf.AppConfig.AlertOffline = c.PostForm("offline") == "on"
	conf.AppConfig.AlertIP = c.PostForm("ip") == "on"
	conf.AppConfig.AlertNet = c.PostForm("net") == "on"
	conf.AppConfig.AlertSec = c.PostForm("sec") == "on"
	conf.AppConfig.AlertDNS = c.PostForm("dns") == "on"
	conf.AppConfig.QuietFrom = c.PostForm("quietfrom")
	conf.AppConfig.QuietTo = c.PostForm("quietto")

	n, err := strconv.Atoi(c.PostForm("offlineafter"))
	if err == nil && n > 0 {
		conf.AppConfig.OfflineAfter = n
	}

	conf.Write(conf.AppConfig)

	c.Redirect(http.StatusFound, c.Request.Referer())
}

func saveMonitorHandler(c *gin.Context) {

	conf.AppConfig.IdentifyEnable = c.PostForm("identify") == "on"
	conf.AppConfig.PublicIPURL = c.PostForm("publicip")
	conf.AppConfig.SpeedURL = c.PostForm("speedurl")
	conf.AppConfig.NetCheckMin, _ = strconv.Atoi(c.PostForm("netmin"))
	conf.AppConfig.SpeedHours, _ = strconv.Atoi(c.PostForm("speedhours"))
	conf.AppConfig.SecScanHours, _ = strconv.Atoi(c.PostForm("sechours"))

	conf.Write(conf.AppConfig)

	c.Redirect(http.StatusFound, c.Request.Referer())
}

func saveDNSHandler(c *gin.Context) {

	conf.AppConfig.DNSType = c.PostForm("type")
	conf.AppConfig.DNSURL = c.PostForm("url")
	conf.AppConfig.DNSUser = c.PostForm("user")
	conf.AppConfig.DNSPass = keep(c.PostForm("pass"), conf.AppConfig.DNSPass, maskAll)
	conf.AppConfig.ThreatListURL = c.PostForm("threats")

	conf.Write(conf.AppConfig)

	c.Redirect(http.StatusFound, c.Request.Referer())
}
