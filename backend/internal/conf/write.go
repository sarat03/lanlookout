package conf

import (
	"log/slog"

	"github.com/spf13/viper"

	"github.com/sarat03/lanlookout/internal/check"
	"github.com/sarat03/lanlookout/internal/models"
)

// Write - write config to file
func Write(config models.Conf) {

	slog.Info("Writing new config to " + config.ConfPath)

	viper.SetConfigFile(config.ConfPath)
	viper.SetConfigType("yaml")

	viper.Set("HOST", config.Host)
	viper.Set("PORT", config.Port)
	viper.Set("THEME", config.Theme)
	viper.Set("COLOR", config.Color)
	viper.Set("NODEPATH", config.NodePath)
	viper.Set("LOG_LEVEL", config.LogLevel)
	viper.Set("ARP_ARGS", config.ArpArgs)
	viper.Set("ARP_STRS", config.ArpStrs)
	viper.Set("ARP_STRS_JOINED", "") // Can be set only with ENV
	viper.Set("IFACES", config.Ifaces)
	viper.Set("TIMEOUT", config.Timeout)
	viper.Set("TRIM_HIST", config.TrimHist)
	viper.Set("SHOUTRRR_URL", config.ShoutURL)

	viper.Set("DHCP_LEASES", config.DHCPLeases)
	viper.Set("MQTT_URL", config.MqttURL)
	viper.Set("MQTT_USER", config.MqttUser)
	viper.Set("MQTT_PASS", config.MqttPass)

	viper.Set("OFFLINE_AFTER", config.OfflineAfter)
	viper.Set("ALERT_NEW", config.AlertNew)
	viper.Set("ALERT_OFFLINE", config.AlertOffline)
	viper.Set("ALERT_IP", config.AlertIP)
	viper.Set("ALERT_NET", config.AlertNet)
	viper.Set("ALERT_SEC", config.AlertSec)
	viper.Set("ALERT_DNS", config.AlertDNS)
	viper.Set("QUIET_FROM", config.QuietFrom)
	viper.Set("QUIET_TO", config.QuietTo)
	viper.Set("MEMBER_ORDER", config.MemberOrder)

	viper.Set("IDENTIFY_ENABLE", config.IdentifyEnable)
	viper.Set("NET_CHECK_MIN", config.NetCheckMin)
	viper.Set("PUBLIC_IP_URL", config.PublicIPURL)
	viper.Set("SPEED_HOURS", config.SpeedHours)
	viper.Set("SPEED_URL", config.SpeedURL)
	viper.Set("SEC_SCAN_HOURS", config.SecScanHours)
	viper.Set("DNS_TYPE", config.DNSType)
	viper.Set("DNS_URL", config.DNSURL)
	viper.Set("DNS_USER", config.DNSUser)
	viper.Set("DNS_PASS", config.DNSPass)
	viper.Set("THREAT_LIST_URL", config.ThreatListURL)

	viper.Set("USE_DB", config.UseDB)
	viper.Set("PG_CONNECT", config.PGConnect)

	viper.Set("influx_enable", config.InfluxEnable)
	viper.Set("influx_skip_tls", config.InfluxSkipTLS)
	viper.Set("influx_addr", config.InfluxAddr)
	viper.Set("influx_token", config.InfluxToken)
	viper.Set("influx_org", config.InfluxOrg)
	viper.Set("influx_bucket", config.InfluxBucket)

	viper.Set("PROMETHEUS_ENABLE", config.PrometheusEnable)

	err := viper.WriteConfig()
	check.IfError(err)
}
