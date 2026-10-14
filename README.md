<h1><img src="assets/logo.svg" width="32" align="top" alt=""/> LanLookout</h1>

A self-hosted monitor for your home network. LanLookout keeps a list of every device on your LAN,
works out what each one is, watches your internet connection, flags risky open services and,
with Pi-hole or AdGuard Home, shows which sites each device talks to. It runs as one small Go binary
or Docker container with a web UI that works on phones and desktops.

![Hosts](assets/hosts.jpg)

## Features

- **Every device, identified.** Finds devices with `arp-scan` (and optionally your DHCP leases), then
  listens to mDNS/Bonjour, UPnP and DHCP to fill in model, OS, services and device type.
- **Alerts that matter.** New device joined, a watched device went offline or came back, a device
  changed IP, the internet went down, your public IP changed, a risky port opened, a device looked
  up a known-bad domain. Quiet hours, per-alert toggles, and a Recent events list.
  Notifications go through [Shoutrrr](https://shoutrrr.nickfedor.com/services/overview/):
  Discord, Telegram, Slack, email, Gotify, ntfy, Pushover, Matrix, webhooks and more.
- **Internet health.** Latency and outage history, public IP, and on-demand or scheduled speed tests.
- **Security check.** Scans devices for services that are risky to leave open (telnet, ADB, Redis,
  Docker API, SMB...) and checks whether the router allows UPnP port forwarding, with plain advice.
- **Per-device DNS.** Reads Pi-hole (v6+) or AdGuard Home query logs: top sites per device,
  blocked lookups, and alerts against a threat list such as URLhaus.
- **Integrations.** Home Assistant (MQTT discovery, one home/away sensor per device), Prometheus,
  InfluxDB2, a REST API with Swagger docs, CSV export.
- **Per device:** names, types, tags, online history and uptime, port scan, Wake-on-LAN.

| Internet | Security | Device |
|---|---|---|
| ![Internet](assets/internet.jpg) | ![Security](assets/security.jpg) | ![Device](assets/host.jpg) |

## Quick start

LanLookout needs **host networking**: ARP scans and mDNS/UPnP/DHCP listening only work on the host's network.

```sh
docker run --name lanlookout \
  -e "IFACES=eth0" \
  -e "TZ=Europe/London" \
  --network="host" \
  -v ~/.dockerdata/lanlookout:/data/lanlookout \
  ghcr.io/sarat03/lanlookout:latest
```

Then open `http://<server-ip>:8840`. Set `IFACES` to the interface(s) to scan (`ip link show` lists them).
A ready-to-edit [docker-compose.yml](docker-compose.yml) is included.

> [!WARNING]
> LanLookout has **no login**. Run it only on a network you trust, or put it behind a reverse proxy
> with authentication. Requests made by other websites in your browser are refused, and passwords
> and tokens are hidden from the settings API, but anyone who can reach the port can use the app.

### Run without Docker

Build with Go 1.25+ (`cd backend && go build ./cmd/lanlookout`), install `arp-scan`, and run
`sudo ./lanlookout -d /path/to/data`. Releases also ship `.deb`, `.rpm` and `.apk` packages with a service file.

## Config

Settings can be changed in the web UI (Settings), in `config_v2.yaml` in the data directory (same
names in lower case), or with environment variables. An empty environment variable means "off".

### Basic config
| Variable  | Description | Default |
| --------  | ----------- | ------- |
| TZ | Set your timezone for correct time | |
| HOST | Listen address | 0.0.0.0 |
| PORT   | Port for web GUI | 8840 |
| THEME | Any theme name from https://bootswatch.com in lower case | cosmo |
| COLOR | Background color: light or dark | dark |
| SHOUTRRR_URL | LanLookout uses [Shoutrrr](https://github.com/nicholas-fedor/shoutrrr) to send notifications. It is already integrated, just needs a correct URL. Examples for Discord, Email, Gotify, Matrix, Ntfy, Pushover, Slack, Telegram, Generic Webhook and etc are [here](https://nicholas-fedor.github.io/shoutrrr/) | |

### Scan settings
| Variable  | Description | Default |
| --------  | ----------- | ------- |
| IFACES | Interfaces to scan. Could be one or more, separated by space. See [docs/VLAN_ARP_SCAN.md](docs/VLAN_ARP_SCAN.md). | |
| TIMEOUT | Time between scans (seconds) | 120 |
| ARP_ARGS | Arguments for `arp-scan`. Enable `debug` log level to see resulting command. (Example: `-r 1`). See [docs/VLAN_ARP_SCAN.md](docs/VLAN_ARP_SCAN.md). | |
| ARP_STRS ARP_STRS_JOINED | See [docs/VLAN_ARP_SCAN.md](docs/VLAN_ARP_SCAN.md). | |
| LOG_LEVEL | Log level: `debug`, `info`, `warn` or `error` | info |
| TRIM_HIST | Remove history after (hours) | 48 |
| DHCP_LEASES | Optional path to a dnsmasq-format leases file (Pi-hole: `/etc/pihole/dhcp.leases`, OpenWrt: `/tmp/dhcp.leases`). Adds devices `arp-scan` missed; a device counts as online while its lease is active | |
| USE_DB | Either `sqlite` or `postgres` | sqlite |
| PG_CONNECT | Address to connect to PostgreSQL. (Example: `postgres://username:password@192.168.0.1:5432/dbname?sslmode=disable`). Full list of URL parameters [here](https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters) | |

### InfluxDB2 config
This config matches Grafana's config for InfluxDB data source

| Variable  | Description | Default | Example |
| --------  | ----------- | ------- | ------- |
| INFLUX_ENABLE | Enable export to InfluxDB2 | false | true |
| INFLUX_SKIP_TLS | Skip TLS Verify | false | true |
| INFLUX_ADDR | Address:port of InfluxDB2 server | | https://192.168.2.3:8086/ |
| INFLUX_BUCKET | InfluxDB2 bucket | | test |
| INFLUX_ORG | InfluxDB2 org | | home |
| INFLUX_TOKEN | Secret token, generated by InfluxDB2 | | |

### Prometheus config
This config configures the Prometheus data source

| Variable  | Description | Default | Example |
| --------  | ----------- | ------- | ------- |
| PROMETHEUS_ENABLE | Enable the Prometheus `/metrics` endpoint | false | true |

### Home Assistant (MQTT)
Every host is published as a `connectivity` binary sensor via [MQTT discovery](https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery), updated after each scan. Empty URL = off.

| Variable  | Description | Default | Example |
| --------  | ----------- | ------- | ------- |
| MQTT_URL | Broker URL | | tcp://192.168.1.10:1883 |
| MQTT_USER | Broker username | | |
| MQTT_PASS | Broker password | | |

### Alerts
Sent through `SHOUTRRR_URL`. Every alert is also listed under History → Recent events.

| Variable  | Description | Default |
| --------  | ----------- | ------- |
| ALERT_NEW | New device found (not marked known) | true |
| ALERT_OFFLINE | Watched device offline or back online. Turn on "Watch this device" on the host page | true |
| ALERT_IP | Known device changed IP | true |
| OFFLINE_AFTER | Missed scans before a watched device counts as offline | 2 |
| QUIET_FROM QUIET_TO | No notifications between these times, `HH:MM` server time (may wrap past midnight) | |

### Monitoring
| Variable  | Description | Default |
| --------  | ----------- | ------- |
| IDENTIFY_ENABLE | Detect model, OS and type from mDNS/Bonjour, UPnP (SSDP) and DHCP requests. Needs host networking; DHCP needs UDP port 67 free | true |
| NET_CHECK_MIN | Minutes between internet checks (latency to 1.1.1.1 / 8.8.8.8 / 9.9.9.9, outage detection). 0 = off | 5 |
| PUBLIC_IP_URL | URL returning your public IP as text, for IP-change alerts. Empty = off | https://api.ipify.org |
| SPEED_HOURS | Hours between speed tests (~25 MB down, ~10 MB up). 0 = manual only | 0 |
| SPEED_URL | Server with Cloudflare-style `/__down` and `/__up` endpoints | https://speed.cloudflare.com |
| SEC_SCAN_HOURS | Hours between security scans of risky ports. 0 = manual only | 24 |
| ALERT_NET ALERT_SEC ALERT_DNS | Alerts for internet outage / public IP change, new risky port, threat-list domain | true |

### DNS logs (Pi-hole / AdGuard Home)
Shows which sites each device looks up, and alerts when one looks up a domain on a threat list.

| Variable  | Description | Default |
| --------  | ----------- | ------- |
| DNS_TYPE | `pihole` (v6+) or `adguard`. Empty = off | |
| DNS_URL | e.g. `http://192.168.1.2` or `http://192.168.1.2:3000` | |
| DNS_USER | AdGuard Home username | |
| DNS_PASS | Pi-hole app password or AdGuard Home password | |
| THREAT_LIST_URL | Hosts-format list of bad domains, refreshed daily, e.g. `https://urlhaus.abuse.ch/downloads/hostfile/` | |

### Security
The API has no login: run it only on a network you trust, or behind a reverse proxy with authentication. Requests a browser makes from another website are refused, so a page you visit can't change or delete hosts. `/api/config` hides passwords and tokens.


## API

Swagger docs are served at `/swagger/index.html`. See [docs/API.md](docs/API.md) and
[docs/VLAN_ARP_SCAN.md](docs/VLAN_ARP_SCAN.md) for VLAN and custom `arp-scan` setups, and the [FAQ](FAQ.md).

## What it doesn't do

LanLookout watches from inside your network. It isn't a firewall: it can't block traffic, and it
can't see what's exposed to the internet. The DNS view is only as complete as your Pi-hole or
AdGuard Home logs. Device identification is a best guess from what devices announce.

## Credits

The idea comes from [WatchYourLAN](https://github.com/aceberg/WatchYourLAN). LanLookout is released
under the [MIT licence](LICENSE).
