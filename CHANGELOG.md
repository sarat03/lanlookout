# Change Log
All notable changes to this project are documented in this file.

## [v0.1.0] - 2026-10-14
First LanLookout release.
### Added
- New responsive web UI (phone and desktop): Hosts, History, Internet, Security, Settings
- Device identification from mDNS/Bonjour, UPnP (SSDP) and DHCP: model, OS, services, type
- Internet health: latency, outage detection, public IP changes, speed tests
- Security check: risky open ports and router UPnP, with advice
- Per-device DNS from Pi-hole v6 / AdGuard Home, threat-list domain alerts
- Watch a device (offline / back online alerts), IP-change alerts, tags, device types
- Alert settings with quiet hours; Recent events list
- Home Assistant via MQTT discovery; DHCP leases as an extra device source
- CSV export, quick port scan, uptime %, Wake-on-LAN
### Security
- Cross-site browser requests to the API are refused; secrets are hidden in `/api/config`
