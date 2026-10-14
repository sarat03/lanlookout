import { apiPath } from "../../functions/api"
import { appConfig } from "../../functions/exports"
import Field from "./Field"

function Prometheus() {

  return (
    <section class="wy-card" aria-label="Prometheus">
      <div>
        <h2>Prometheus</h2>
        <div class="wy-sub">Expose host status for scraping.</div>
      </div>
      <form action={apiPath + '/api/config_prometheus/'} method="post">
        <Field id="penable" label="Prometheus metrics" help={<>Served at <a class="wy-link wy-mono" href="/metrics" target="_blank">/metrics</a>.</>} tight>
          <input id="penable" type="checkbox" role="switch" class="wy-switch" name="enable" checked={appConfig().PrometheusEnable}></input>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save Prometheus settings</button>
        </div>
      </form>
    </section>
  )
}

export default Prometheus
