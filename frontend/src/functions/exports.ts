import { createSignal } from "solid-js";
import { createStore } from "solid-js/store";

export interface Host {
	ID:    number;
	Name:  string;
	DNS:   string;
	Iface: string;
	IP:    string;
	Mac:   string;
	Hw:    string;
	Date:  string;
	Known: number;
	Now:   number;
	Watch: number;
	Type:  string;
	Tags:  string;
	Model:    string;
	OS:       string;
	Services: string;
	Ports:    string;
};

// Device types for the host icon. Bootstrap Icons names.
export const hostTypes: Record<string, string> = {
	phone: "bi-phone", tablet: "bi-tablet", laptop: "bi-laptop", desktop: "bi-pc-display",
	tv: "bi-tv", speaker: "bi-speaker", console: "bi-controller", watch: "bi-smartwatch",
	printer: "bi-printer", camera: "bi-camera-video", router: "bi-router", server: "bi-hdd-rack",
	plug: "bi-plug", thermostat: "bi-thermometer-half", fan: "bi-fan", light: "bi-lightbulb",
	car: "bi-car-front", iot: "bi-cpu", other: "bi-question-circle",
};

export const tagList = (h: Host) => (h.Tags || "").split(",").map(t => t.trim()).filter(Boolean);

export interface Conf {
	Host:	   string;
	Port:	   string;
	Theme:	   string;
	Color:     string;
	DirPath:   string;
	Timeout:   number;
	NodePath:  string;
	LogLevel:  string;
	Ifaces:	   string;
	ArpArgs:   string;
	ArpStrs:   string[];
	TrimHist:  number;
	ShoutURL:  string;
	UseDB:     string;
	PGConnect: string;
	DHCPLeases: string;
	MqttURL:   string;
	MqttUser:  string;
	MqttPass:  string;
	OfflineAfter: number;
	AlertNew:     boolean;
	AlertOffline: boolean;
	AlertIP:      boolean;
	QuietFrom:    string;
	QuietTo:      string;
	AlertNet:     boolean;
	AlertSec:     boolean;
	AlertDNS:     boolean;
	IdentifyEnable: boolean;
	NetCheckMin:  number;
	PublicIPURL:  string;
	SpeedHours:   number;
	SpeedURL:     string;
	SecScanHours: number;
	DNSType:      string;
	DNSURL:       string;
	DNSUser:      string;
	DNSPass:      string;
	ThreatListURL: string;
	// InfluxDB
	InfluxEnable:  boolean;
	InfluxAddr:    string;
	InfluxToken:   string;
	InfluxOrg:     string;
	InfluxBucket:  string;
	InfluxSkipTLS: boolean;
	// Prometheus
	PrometheusEnable: boolean;
};

export const emptyHost:Host = {
	ID:    0,
	Name:  "",
	DNS:   "",
	Iface: "",
	IP:    "",
	Mac:   "",
	Hw:    "",
	Date:  "",
	Known: 0,
	Now:   0,
	Watch: 0,
	Type:  "",
	Tags:  "",
	Model: "",
	OS: "",
	Services: "",
	Ports: "",
};

export const emptyConf:Conf = {
	Host:	 "",
	Port:	 "",
	Theme:	 "",
	Color:   "",
	DirPath: "",
	Timeout: 120,
	NodePath: "",
	LogLevel: "",
	Ifaces:	 "",
	ArpArgs: "",
	ArpStrs: [],
	TrimHist: 48,
	ShoutURL: "",
	UseDB: "",
	PGConnect: "",
	DHCPLeases: "",
	MqttURL: "",
	MqttUser: "",
	MqttPass: "",
	OfflineAfter: 2,
	AlertNew: true,
	AlertOffline: true,
	AlertIP: true,
	QuietFrom: "",
	QuietTo: "",
	AlertNet: true,
	AlertSec: true,
	AlertDNS: true,
	IdentifyEnable: true,
	NetCheckMin: 5,
	PublicIPURL: "",
	SpeedHours: 0,
	SpeedURL: "",
	SecScanHours: 24,
	DNSType: "",
	DNSURL: "",
	DNSUser: "",
	DNSPass: "",
	ThreatListURL: "",
	InfluxEnable:  false,
	InfluxAddr:    "",
	InfluxToken:   "",
	InfluxOrg:     "",
	InfluxBucket:  "",
	InfluxSkipTLS: false,
	PrometheusEnable: false,
};

// Official Bootswatch 5.3 themes (https://bootswatch.com)
export const themes = ["cerulean", "cosmo", "cyborg", "darkly", "flatly", "journal", "litera", "lumen", "lux", "materia", "minty", "morph", "pulse", "quartz", "sandstone", "simplex", "sketchy", "slate", "solar", "spacelab", "superhero", "united", "vapor", "yeti", "zephyr"];

export const [allHosts, setAllHosts] = createStore<Host[]>([]);
export const [bkpHosts, setBkpHosts] = createSignal<Host[]>([]);

export const [ifaces, setIfaces] = createSignal<string[]>([]);

export const [appConfig, setAppConfig] = createSignal<Conf>(emptyConf);

export const [editNames, setEditNames] = createSignal(false);

export const [show, setShow] = createSignal<number>(200);

export const [histUpdOnFilter, setHistUpdOnFilter] = createSignal(false);

export const [selectedIDs, setSelectedIDs] = createSignal<number[]>([]);

export interface Event {
	ID:   number;
	Date: string;
	Kind: "new" | "offline" | "online" | "ip" | "outage" | "wan_ip" | "port" | "threat";
	Mac:  string;
	Msg:  string;
};

export interface NetStat {
	ID:        number;
	Date:      string;
	Up:        number;
	LatencyMs: number;
	DownMbps:  number;
	UpMbps:    number;
	PublicIP:  string;
};

export interface Finding {
	HostID:   number;
	Name:     string;
	IP:       string;
	Port:     number;
	Service:  string;
	Severity: "high" | "medium" | "low";
	Advice:   string;
};

export interface DNSStat {
	Domain:  string;
	Count:   number;
	Blocked: number;
	Threat:  number;
};
