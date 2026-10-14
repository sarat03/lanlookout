import { apiPath } from "../../functions/api"
import { appConfig } from "../../functions/exports"
import Field from "./Field"

function Mqtt() {

  return (
    <section class="wy-card" aria-label="Home Assistant">
      <div>
        <h2>Home Assistant (MQTT)</h2>
        <div class="wy-sub">Each host shows up in Home Assistant as a connectivity sensor (home / away), via MQTT discovery. Leave the URL empty to turn it off.</div>
      </div>
      <form action={apiPath + '/api/config_mqtt/'} method="post">
        <Field id="murl" label="Broker URL" help={<>For example <code>tcp://192.168.1.10:1883</code>. Updates are sent after every scan.</>}>
          <input id="murl" name="url" type="text" class="wy-input wy-mono" placeholder="tcp://homeassistant.local:1883" value={appConfig().MqttURL}></input>
        </Field>
        <Field id="muser" label="Username">
          <input id="muser" name="user" type="text" class="wy-input" autocomplete="off" value={appConfig().MqttUser}></input>
        </Field>
        <Field id="mpass" label="Password">
          <input id="mpass" name="pass" type="password" class="wy-input" autocomplete="new-password" value={appConfig().MqttPass}></input>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save MQTT settings</button>
        </div>
      </form>
    </section>
  )
}

export default Mqtt
