import { createResource, For } from "solid-js"
import { appConfig } from "../../functions/exports"
import { apiGetIfaces, apiPath } from "../../functions/api"
import Field from "./Field"

const docs = "https://github.com/sarat03/lanlookout/blob/main/docs/VLAN_ARP_SCAN.md";

function Scan() {

  // Configured ones stay listed even if missing here (e.g. a VLAN that's down).
  const [found] = createResource(apiGetIfaces);
  const chosen = () => (appConfig().Ifaces || "").split(/\s+/).filter(Boolean);
  const ifaces = () => [...new Set([...(found() || []), ...chosen()])];

  return (
    <section class="wy-card" aria-label="Scanning">
      <div>
        <h2>Scanning</h2>
        <div class="wy-sub">Which networks are scanned, how often, and where results are stored.</div>
      </div>
      <form action={apiPath + '/api/config_settings/'} method="post">
        <Field label="Interfaces" help="Network interfaces of this machine to scan. Pick one or more.">
          <div class="wy-row-wrap" style={{ gap: "8px 16px" }}>
            <For each={ifaces()}>{(i) =>
              <label class="wy-inline-label wy-mono">
                <input type="checkbox" name="ifaces" value={i} checked={chosen().includes(i)}/> {i}
              </label>
            }</For>
          </div>
        </Field>
        <Field id="timeout" label="Scan every (seconds)" help="Shorter means faster alerts and more network traffic.">
          <input id="timeout" name="timeout" type="number" min="1" class="wy-input wy-mono" value={appConfig().Timeout}></input>
        </Field>
        <Field id="arpargs" label="Arguments for arp-scan" help={<>Optional, for example <code>-r 1</code>. Set log level to debug to see the full command. <a class="wy-link" href={docs} target="_blank">Docs</a></>}>
          <input id="arpargs" name="arpargs" type="text" class="wy-input wy-mono" value={appConfig().ArpArgs}></input>
        </Field>
        <Field label="Arp strings" help={<>Extra scans for VLANs, docker0 and so on. One per box; the empty box adds a new one. <a class="wy-link" href={docs} target="_blank">Docs</a></>}>
          <For each={appConfig().ArpStrs}>{arpStr =>
            <input name="arpstrs" type="text" class="wy-input wy-mono" aria-label="Arp string" value={arpStr}></input>
          }</For>
          <input name="arpstrs" type="text" class="wy-input wy-mono" aria-label="New arp string" placeholder="Add another"></input>
        </Field>
        <Field id="leases" label="DHCP leases file" help={<>Optional. Also finds devices from a dnsmasq leases file, useful when arp-scan can't see everything. Pi-hole: <code>/etc/pihole/dhcp.leases</code>, OpenWrt: <code>/tmp/dhcp.leases</code>. Mount it into the container read-only. A device counts as online while its lease is active.</>}>
          <input id="leases" name="leases" type="text" class="wy-input wy-mono" placeholder="/etc/pihole/dhcp.leases" value={appConfig().DHCPLeases}></input>
        </Field>
        <Field id="log" label="Log level" help="Shown in the container log.">
          <select id="log" name="log" class="wy-input">
            <For each={["debug","info","warn","error"]}>{level =>
              <option value={level} selected={level == appConfig().LogLevel}>{level}</option>
            }</For>
          </select>
        </Field>
        <Field id="trim" label="Keep history (hours)" help="Older history is removed.">
          <input id="trim" name="trim" type="number" min="1" class="wy-input wy-mono" value={appConfig().TrimHist}></input>
        </Field>
        <Field id="usedb" label="Database" help="SQLite needs no setup. PostgreSQL suits larger networks.">
          <select id="usedb" name="usedb" class="wy-input">
            <option value="sqlite" selected={appConfig().UseDB == "sqlite"}>sqlite</option>
            <option value="postgres" selected={appConfig().UseDB != "sqlite"}>postgres</option>
          </select>
        </Field>
        <Field id="pgconnect" label="PostgreSQL connection URL" help={<>Only used with postgres. Example: <code>postgres://user:pass@192.168.0.1:5432/db?sslmode=disable</code>. <a class="wy-link" href="https://pkg.go.dev/github.com/lib/pq#hdr-Connection_String_Parameters" target="_blank">Parameters</a></>}>
          <textarea id="pgconnect" name="pgconnect" class="wy-input wy-mono" rows="3" wrap="soft">{appConfig().PGConnect}</textarea>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save and rescan</button>
          <span class="wy-muted" style={{ "font-size": "14px" }}>Saving restarts the scanner.</span>
        </div>
      </form>
    </section>
  )
}

export default Scan
