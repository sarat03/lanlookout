import { apiPath } from "../../functions/api"
import { appConfig } from "../../functions/exports"
import Field from "./Field"

function DNSLog() {

  return (
    <section class="wy-card" aria-label="DNS logs">
      <div>
        <h2>DNS logs (Pi-hole / AdGuard Home)</h2>
        <div class="wy-sub">See which sites each device looks up, and get alerts for known-bad domains. Read every 5 minutes.</div>
      </div>
      <form action={apiPath + '/api/config_dns/'} method="post">
        <Field id="dtype" label="DNS server">
          <select id="dtype" name="type" class="wy-input">
            <option value="" selected={!appConfig().DNSType}>Off</option>
            <option value="pihole" selected={appConfig().DNSType == "pihole"}>Pi-hole (v6 or newer)</option>
            <option value="adguard" selected={appConfig().DNSType == "adguard"}>AdGuard Home</option>
          </select>
        </Field>
        <Field id="durl" label="Address" help={<>For example <code>http://192.168.1.2</code> (Pi-hole) or <code>http://192.168.1.2:3000</code> (AdGuard Home).</>}>
          <input id="durl" name="url" type="text" class="wy-input wy-mono" value={appConfig().DNSURL}></input>
        </Field>
        <Field id="duser" label="Username" help="AdGuard Home only.">
          <input id="duser" name="user" type="text" class="wy-input" autocomplete="off" value={appConfig().DNSUser}></input>
        </Field>
        <Field id="dpass" label="Password" help="Pi-hole: an app password from Settings → Web interface / API.">
          <input id="dpass" name="pass" type="password" class="wy-input" autocomplete="new-password" value={appConfig().DNSPass}></input>
        </Field>
        <Field id="dthreat" label="Threat list URL" help={<>Optional hosts-format list of malicious domains, downloaded daily. For example URLhaus: <code>https://urlhaus.abuse.ch/downloads/hostfile/</code></>}>
          <input id="dthreat" name="threats" type="text" class="wy-input wy-mono" value={appConfig().ThreatListURL}></input>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save DNS settings</button>
        </div>
      </form>
    </section>
  )
}

export default DNSLog
