package api

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Routes - start API routes
func Routes(router *gin.Engine) {

	r0 := router.Group("/api")
	r0.Use(sameSiteOnly) // security.go
	{
		r0.GET("/all", getAllHosts)                // api-hosts.go
		r0.GET("/edit/:id/:name/*known", editHost) // api-hosts.go
		r0.GET("/host/:id", getHost)               // api-hosts.go
		r0.GET("/host/del/:id", delHost)           // api-hosts.go
		r0.GET("/host/add/:mac", addHost)          // api-hosts.go
		r0.GET("/host/meta/:id", metaHost)         // api-hosts.go
		r0.GET("/known/all", knownAll)             // api-hosts.go

		r0.GET("/config", getConfig)        // api-system.go
		r0.GET("/notify_test", notifyTest)  // api-system.go
		r0.GET("/status/*iface", getStatus) // api-system.go
		r0.GET("/version", getVersion)      // api-system.go
		r0.GET("/ifaces", getIfaces)        // api-network.go
		r0.GET("/rescan", triggerRescan)    // api-system.go
		r0.GET("/scan_status", scanStatus)  // api-system.go
		r0.GET("/events", getEvents)        // api-system.go

		r0.GET("/net/status", netStatus)       // api-monitor.go
		r0.GET("/net/stats", netStats)         // api-monitor.go
		r0.GET("/net/speedtest", speedTest)    // api-monitor.go
		r0.GET("/security", securityFindings)  // api-monitor.go
		r0.GET("/security/scan", securityScan) // api-monitor.go
		r0.GET("/dns/:ip", dnsTop)             // api-monitor.go

		r0.GET("/history", getHistory)                  // api-history.go
		r0.GET("/history/:mac", getHistoryByMAC)        // api-history.go
		r0.GET("/history/:mac/:date", getHistoryByDate) // api-history.go

		r0.GET("/port/:addr/:port", getPortState) // api-network.go
		r0.GET("/wol/:mac", sendWOL)              // api-network.go

		r0.POST("/config/", saveConfigHandler)                // config.go
		r0.POST("/config_settings/", saveSettingsHandler)     // config.go
		r0.POST("/config_influx/", saveInfluxHandler)         // config.go
		r0.POST("/config_prometheus/", savePrometheusHandler) // config.go
		r0.POST("/config_mqtt/", saveMqttHandler)             // config.go
		r0.POST("/config_alerts/", saveAlertsHandler)         // config.go
		r0.POST("/config_monitor/", saveMonitorHandler)       // config.go
		r0.POST("/config_dns/", saveDNSHandler)               // config.go
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
