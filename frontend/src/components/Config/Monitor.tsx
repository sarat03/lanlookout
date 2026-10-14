import { apiPath } from "../../functions/api"
import { appConfig } from "../../functions/exports"
import Field from "./Field"

function Monitor() {

  return (
    <section class="wy-card" aria-label="Monitoring">
      <div>
        <h2>Monitoring</h2>
        <div class="wy-sub">Device identification, internet checks and the security scan.</div>
      </div>
      <form action={apiPath + '/api/config_monitor/'} method="post">
        <Field id="identify" label="Identify devices" help="Listen to mDNS/Bonjour, UPnP and DHCP to fill in model, OS and type. Needs host networking. Takes effect after a restart." tight>
          <input id="identify" type="checkbox" role="switch" class="wy-switch" name="identify" checked={appConfig().IdentifyEnable}></input>
        </Field>
        <Field id="netmin" label="Internet check (minutes)" help="Connects to 1.1.1.1, 8.8.8.8 and 9.9.9.9 to measure latency and detect outages. 0 = off.">
          <input id="netmin" name="netmin" type="number" min="0" class="wy-input wy-mono" value={appConfig().NetCheckMin}></input>
        </Field>
        <Field id="publicip" label="Public IP lookup URL" help="Must return just the IP as text. Empty = don't look it up (no IP-change alerts).">
          <input id="publicip" name="publicip" type="text" class="wy-input wy-mono" placeholder="https://api.ipify.org" value={appConfig().PublicIPURL}></input>
        </Field>
        <Field id="speedhours" label="Speed test every (hours)" help="Each test downloads ~25 MB and uploads ~10 MB. 0 = only when you press the button.">
          <input id="speedhours" name="speedhours" type="number" min="0" class="wy-input wy-mono" value={appConfig().SpeedHours}></input>
        </Field>
        <Field id="speedurl" label="Speed test server" help="Any server with Cloudflare-style /__down and /__up endpoints.">
          <input id="speedurl" name="speedurl" type="text" class="wy-input wy-mono" placeholder="https://speed.cloudflare.com" value={appConfig().SpeedURL}></input>
        </Field>
        <Field id="sechours" label="Security scan every (hours)" help="Checks risky ports on online devices. 0 = only when you press Scan now.">
          <input id="sechours" name="sechours" type="number" min="0" class="wy-input wy-mono" value={appConfig().SecScanHours}></input>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save monitoring settings</button>
        </div>
      </form>
    </section>
  )
}

export default Monitor
