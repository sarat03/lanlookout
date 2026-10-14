import { apiPath } from "../../functions/api"
import { appConfig } from "../../functions/exports"
import Field from "./Field"

function Influx() {

  return (
    <section class="wy-card" aria-label="InfluxDB2">
      <div>
        <h2>InfluxDB2</h2>
        <div class="wy-sub">Write host history to an InfluxDB2 bucket.</div>
      </div>
      <form action={apiPath + '/api/config_influx/'} method="post">
        <Field id="ienable" label="Export to InfluxDB2" tight>
          <input id="ienable" type="checkbox" role="switch" class="wy-switch" name="enable" checked={appConfig().InfluxEnable}></input>
        </Field>
        <Field id="addr" label="Address" help="Base URL of the server.">
          <input id="addr" name="addr" type="text" class="wy-input wy-mono" value={appConfig().InfluxAddr} placeholder="http://localhost:8086"></input>
        </Field>
        <Field id="token" label="API token">
          <input id="token" name="token" type="password" autocomplete="off" class="wy-input wy-mono" value={appConfig().InfluxToken}></input>
        </Field>
        <Field id="org" label="Organisation">
          <input id="org" name="org" type="text" class="wy-input" value={appConfig().InfluxOrg}></input>
        </Field>
        <Field id="bucket" label="Bucket">
          <input id="bucket" name="bucket" type="text" class="wy-input" value={appConfig().InfluxBucket}></input>
        </Field>
        <Field id="skip" label="Skip TLS verification" help="Only for self-signed certificates." tight>
          <input id="skip" type="checkbox" role="switch" class="wy-switch" name="skip" checked={appConfig().InfluxSkipTLS}></input>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save InfluxDB2 settings</button>
        </div>
      </form>
    </section>
  )
}

export default Influx
