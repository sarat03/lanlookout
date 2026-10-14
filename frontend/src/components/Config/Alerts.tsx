import { apiPath } from "../../functions/api"
import { appConfig } from "../../functions/exports"
import Field from "./Field"

function Alerts() {

  return (
    <section class="wy-card" aria-label="Alerts">
      <div>
        <h2>Alerts</h2>
        <div class="wy-sub">Which events send a notification. Every event is still listed under History → Recent events.</div>
      </div>
      <form action={apiPath + '/api/config_alerts/'} method="post">
        <Field id="anew" label="New device" help="A device not marked known joins the network." tight>
          <input id="anew" type="checkbox" role="switch" class="wy-switch" name="new" checked={appConfig().AlertNew}></input>
        </Field>
        <Field id="aoff" label="Watched device offline / back online" help="Only devices with Watch turned on, on their host page." tight>
          <input id="aoff" type="checkbox" role="switch" class="wy-switch" name="offline" checked={appConfig().AlertOffline}></input>
        </Field>
        <Field id="aip" label="IP address change" help="A known device got a different IP." tight>
          <input id="aip" type="checkbox" role="switch" class="wy-switch" name="ip" checked={appConfig().AlertIP}></input>
        </Field>
        <Field id="anet" label="Internet outage / public IP change" help="Sent when the connection comes back." tight>
          <input id="anet" type="checkbox" role="switch" class="wy-switch" name="net" checked={appConfig().AlertNet}></input>
        </Field>
        <Field id="asec" label="New risky open port" help="Found by the security scan." tight>
          <input id="asec" type="checkbox" role="switch" class="wy-switch" name="sec" checked={appConfig().AlertSec}></input>
        </Field>
        <Field id="adns" label="Threat-list domain" help="A device looked up a known-bad domain (needs DNS logs)." tight>
          <input id="adns" type="checkbox" role="switch" class="wy-switch" name="dns" checked={appConfig().AlertDNS}></input>
        </Field>
        <Field id="aafter" label="Offline after (missed scans)" help="Raise it if phones going to sleep cause false offline alerts.">
          <input id="aafter" name="offlineafter" type="number" min="1" class="wy-input wy-mono" value={appConfig().OfflineAfter}></input>
        </Field>
        <Field label="Quiet hours" help="No notifications between these times (server time). Leave empty for none.">
          <div class="wy-row-wrap">
            <input name="quietfrom" type="time" class="wy-input" aria-label="Quiet from" value={appConfig().QuietFrom}></input>
            <span class="wy-muted">to</span>
            <input name="quietto" type="time" class="wy-input" aria-label="Quiet to" value={appConfig().QuietTo}></input>
          </div>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save alert settings</button>
        </div>
      </form>
    </section>
  )
}

export default Alerts
