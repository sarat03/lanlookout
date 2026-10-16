package models

// Conf - app config
type Conf struct {
	Host     string
	Port     string
	Theme    string
	Color    string
	DirPath  string
	ConfPath string
	DBPath   string
	NodePath string
	LogLevel string
	Ifaces   string
	ArpArgs  string
	ArpStrs  []string
	Timeout  int
	TrimHist int
	ShoutURL string
	Version  string
	// Extra discovery: dnsmasq / Pi-hole leases file
	DHCPLeases string
	// Home Assistant over MQTT
	MqttURL  string
	MqttUser string
	MqttPass string
	// Alerts
	OfflineAfter int    // missed scans before a watched host is offline
	AlertNew     bool   // new device
	AlertOffline bool   // watched device offline / back online
	AlertIP      bool   // known device changed IP
	AlertNet     bool   // internet outage, public IP change
	AlertSec     bool   // new risky open port
	AlertDNS     bool   // device contacted a threat-list domain
	QuietFrom    string // "22:00", empty = no quiet hours
	QuietTo      string // "07:00"
	// Home page
	MemberOrder string // member names, comma separated, in display order
	// Device identification (mDNS, SSDP, DHCP)
	IdentifyEnable bool
	// Internet health
	NetCheckMin int    // minutes between checks, 0 = off
	PublicIPURL string // returns the public IP as plain text, empty = off
	SpeedHours  int    // hours between speed tests, 0 = only manual
	SpeedURL    string // Cloudflare-compatible __down / __up endpoints
	// Security check
	SecScanHours int // hours between scans, 0 = only manual
	// Per-device DNS
	DNSType       string // "", "pihole" (v6) or "adguard"
	DNSURL        string
	DNSUser       string
	DNSPass       string
	ThreatListURL string // hosts-format list of bad domains
	// PostgreSQL
	UseDB     string
	PGConnect string
	// InfluxDB
	InfluxEnable  bool
	InfluxAddr    string
	InfluxToken   string
	InfluxOrg     string
	InfluxBucket  string
	InfluxSkipTLS bool
	// Prometheus
	PrometheusEnable bool
}

// Host - one host
// SharedOwner - Owner of devices that belong to the whole house
const SharedOwner = "Shared"

type Host struct {
	ID    int    `gorm:"column:ID;primaryKey"`
	Name  string `gorm:"column:NAME"`
	DNS   string `gorm:"column:DNS"`
	Iface string `gorm:"column:IFACE"`
	IP    string `gorm:"column:IP"`
	Mac   string `gorm:"column:MAC"`
	Hw    string `gorm:"column:HW"`
	Date  string `gorm:"column:DATE"`
	Known int    `gorm:"column:KNOWN"`
	Now   int    `gorm:"column:NOW"`
	Watch int    `gorm:"column:WATCH"`                // 1 = notify on offline/online
	Type  string `gorm:"column:TYPE"`                 // phone, laptop, tv...
	Tags  string `gorm:"column:TAGS"`                 // comma separated
	Owner string `gorm:"column:OWNER;default:Shared"` // household member or SharedOwner
	// Filled in automatically
	Model    string `gorm:"column:MODEL"`    // e.g. "Chromecast", "MacBookPro18,1"
	OS       string `gorm:"column:OS"`       // e.g. "iOS / macOS", "Android", "Windows"
	Services string `gorm:"column:SERVICES"` // announced services, comma separated
	Ports    string `gorm:"column:PORTS"`    // open ports from the last security scan
}

// Stat - status
type Stat struct {
	Total   int
	Online  int
	Offline int
	Known   int
	Unknown int
}

// Event - one alert, kept for the Recent events list
type Event struct {
	ID   int    `gorm:"column:ID;primaryKey"`
	Date string `gorm:"column:DATE"`
	Kind string `gorm:"column:KIND"` // new, offline, online, ip
	Mac  string `gorm:"column:MAC"`
	Msg  string `gorm:"column:MSG"`
}

// NetStat - one internet check (and maybe a speed test)
type NetStat struct {
	ID        int     `gorm:"column:ID;primaryKey"`
	Date      string  `gorm:"column:DATE"`
	Up        int     `gorm:"column:UP"`      // 1 = internet reachable
	LatencyMs float64 `gorm:"column:LATENCY"` // TCP connect time, ms
	DownMbps  float64 `gorm:"column:DOWN"`    // 0 = no speed test this time
	UpMbps    float64 `gorm:"column:UPLOAD"`
	PublicIP  string  `gorm:"column:PUBLIC_IP"`
}

// DNSStat - DNS queries of one device to one domain on one day
type DNSStat struct {
	ID      int    `gorm:"column:ID;primaryKey"`
	Day     string `gorm:"column:DAY"` // 2006-01-02
	IP      string `gorm:"column:IP"`
	Domain  string `gorm:"column:DOMAIN"`
	Count   int    `gorm:"column:COUNT"`
	Blocked int    `gorm:"column:BLOCKED"`
	Threat  int    `gorm:"column:THREAT"` // 1 = on the threat list
}
