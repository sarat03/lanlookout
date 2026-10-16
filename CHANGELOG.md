# Change Log
All notable changes to this project are documented in this file.

## [v0.1.6] - 2026-10-16
### Changed
- Home is now the start page (`/`); the hosts table moved to `/hosts`
- New Home layout: health tiles (devices online, offline, new today, internet), a Needs attention list, a card per person with Home/Away from their phone or watch, and Smart home as compact chips
- Person cards open on click to show their devices; cards with something offline start open

## [v0.1.5] - 2026-10-15
### Added
- Household members: Owner on each device, Home page with a card per member and a shared Smart home card
- Pickers instead of typing: Owner dropdown, tag chips, interface checkboxes (new `/api/ifaces`)
- Online status chart on the Internet page
- Every release is tagged `vX.Y.Z` on GitHub and the Docker image gets the same tag
### Fixed
- `/home` returned 404

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
