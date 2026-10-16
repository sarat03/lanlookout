package conf

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/sarat03/lanlookout/internal/check"
	"github.com/sarat03/lanlookout/internal/models"
)

func read(path string) (config models.Conf) {

	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8840")
	viper.SetDefault("THEME", "cosmo")
	viper.SetDefault("COLOR", "dark")
	viper.SetDefault("NODEPATH", "")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ARP_ARGS", "")
	viper.SetDefault("ARP_STRS_JOINED", "")
	viper.SetDefault("IFACES", "")
	viper.SetDefault("TIMEOUT", 120)
	viper.SetDefault("TRIM_HIST", 48)
	viper.SetDefault("SHOUTRRR_URL", "")

	viper.SetDefault("DHCP_LEASES", "")
	viper.SetDefault("MQTT_URL", "")
	viper.SetDefault("MQTT_USER", "")
	viper.SetDefault("MQTT_PASS", "")

	viper.SetDefault("OFFLINE_AFTER", 2)
	viper.SetDefault("ALERT_NEW", true)
	viper.SetDefault("ALERT_OFFLINE", true)
	viper.SetDefault("ALERT_IP", true)
	viper.SetDefault("ALERT_NET", true)
	viper.SetDefault("ALERT_SEC", true)
	viper.SetDefault("ALERT_DNS", true)
	viper.SetDefault("QUIET_FROM", "")
	viper.SetDefault("QUIET_TO", "")
	viper.SetDefault("MEMBER_ORDER", "")

	viper.SetDefault("IDENTIFY_ENABLE", true)
	viper.SetDefault("NET_CHECK_MIN", 5)
	viper.SetDefault("PUBLIC_IP_URL", "https://api.ipify.org")
	viper.SetDefault("SPEED_HOURS", 0)
	viper.SetDefault("SPEED_URL", "https://speed.cloudflare.com")
	viper.SetDefault("SEC_SCAN_HOURS", 24)
	viper.SetDefault("DNS_TYPE", "")
	viper.SetDefault("DNS_URL", "")
	viper.SetDefault("DNS_USER", "")
	viper.SetDefault("DNS_PASS", "")
	viper.SetDefault("THREAT_LIST_URL", "")

	viper.SetDefault("USE_DB", "sqlite")
	viper.SetDefault("PG_CONNECT", "")

	viper.SetDefault("INFLUX_ENABLE", false)

	viper.SetDefault("PROMETHEUS_ENABLE", false)

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	check.IfError(err)

	viper.AllowEmptyEnv(true) // PUBLIC_IP_URL="" etc. means "off", not "use the default"
	viper.AutomaticEnv()      // Get ENVIRONMENT variables

	config.Host = viper.Get("HOST").(string)
	config.Port = viper.Get("PORT").(string)
	config.Theme = viper.Get("THEME").(string)
	config.Color = viper.Get("COLOR").(string)
	config.NodePath = viper.Get("NODEPATH").(string)
	config.LogLevel = viper.Get("LOG_LEVEL").(string)
	config.ArpArgs = viper.Get("ARP_ARGS").(string)
	config.ArpStrs = viper.GetStringSlice("ARP_STRS")
	config.Ifaces = viper.Get("IFACES").(string)
	config.Timeout = viper.GetInt("TIMEOUT")
	config.TrimHist = viper.GetInt("TRIM_HIST")
	config.ShoutURL = viper.Get("SHOUTRRR_URL").(string)

	config.DHCPLeases = viper.GetString("DHCP_LEASES")
	config.MqttURL = viper.GetString("MQTT_URL")
	config.MqttUser = viper.GetString("MQTT_USER")
	config.MqttPass = viper.GetString("MQTT_PASS")

	config.OfflineAfter = viper.GetInt("OFFLINE_AFTER")
	config.AlertNew = viper.GetBool("ALERT_NEW")
	config.AlertOffline = viper.GetBool("ALERT_OFFLINE")
	config.AlertIP = viper.GetBool("ALERT_IP")
	config.AlertNet = viper.GetBool("ALERT_NET")
	config.AlertSec = viper.GetBool("ALERT_SEC")
	config.AlertDNS = viper.GetBool("ALERT_DNS")
	config.QuietFrom = viper.GetString("QUIET_FROM")
	config.QuietTo = viper.GetString("QUIET_TO")
	config.MemberOrder = viper.GetString("MEMBER_ORDER")

	config.IdentifyEnable = viper.GetBool("IDENTIFY_ENABLE")
	config.NetCheckMin = viper.GetInt("NET_CHECK_MIN")
	config.PublicIPURL = viper.GetString("PUBLIC_IP_URL")
	config.SpeedHours = viper.GetInt("SPEED_HOURS")
	config.SpeedURL = viper.GetString("SPEED_URL")
	config.SecScanHours = viper.GetInt("SEC_SCAN_HOURS")
	config.DNSType = viper.GetString("DNS_TYPE")
	config.DNSURL = viper.GetString("DNS_URL")
	config.DNSUser = viper.GetString("DNS_USER")
	config.DNSPass = viper.GetString("DNS_PASS")
	config.ThreatListURL = viper.GetString("THREAT_LIST_URL")

	config.UseDB = viper.Get("USE_DB").(string)
	config.PGConnect = viper.Get("PG_CONNECT").(string)

	config.InfluxEnable = viper.GetBool("INFLUX_ENABLE")
	config.InfluxSkipTLS = viper.GetBool("INFLUX_SKIP_TLS")
	config.InfluxAddr, _ = viper.Get("INFLUX_ADDR").(string)
	config.InfluxToken, _ = viper.Get("INFLUX_TOKEN").(string)
	config.InfluxOrg, _ = viper.Get("INFLUX_ORG").(string)
	config.InfluxBucket, _ = viper.Get("INFLUX_BUCKET").(string)

	config.PrometheusEnable = viper.GetBool("PROMETHEUS_ENABLE")

	joined := viper.Get("ARP_STRS_JOINED").(string)
	// slog.Info("ARP_STRS_JOINED: " + joined)

	if joined != "" {
		config.ArpStrs = strings.Split(joined, ",")
	}

	return config
}
