import { createSignal, For } from "solid-js";
import { apiPath, apiTestNotify } from "../../functions/api"
import { appConfig, themes } from "../../functions/exports"
import Field from "./Field";

function Basic() {


  const [sent, setSent] = createSignal(false);

  const handleTestNotify = async () => {
    await apiTestNotify();
    setSent(true);
  };

  return (
    <section class="wy-card" aria-label="General">
      <div>
        <h2>General</h2>
        <div class="wy-sub">How LanLookout is served, how it looks and where alerts go.</div>
      </div>
      <form action={apiPath + '/api/config/'} method="post">
        <Field id="host" label="Listen address" help="Use 0.0.0.0 to reach it from other devices.">
          <input id="host" name="host" type="text" class="wy-input wy-mono" value={appConfig().Host}></input>
        </Field>
        <Field id="port" label="Port" help="Web interface port.">
          <input id="port" name="port" type="text" inputmode="numeric" class="wy-input wy-mono" value={appConfig().Port}></input>
        </Field>
        <Field id="theme" label="Theme" help="Bootswatch theme. Also sets the accent colour.">
          <select id="theme" name="theme" class="wy-input">
            <For each={themes}>{theme =>
              <option value={theme} selected={theme == appConfig().Theme}>{theme}</option>
            }</For>
          </select>
        </Field>
        <Field id="color" label="Colour mode" help="Light or dark.">
          <select id="color" name="color" class="wy-input">
            <option value="dark" selected={appConfig().Color == "dark"}>dark</option>
            <option value="light" selected={appConfig().Color != "dark"}>light</option>
          </select>
        </Field>
        <Field id="node" label="Local node-bootstrap URL" help="Optional. Serve themes and icons locally instead of from the CDN.">
          <input id="node" name="node" type="text" class="wy-input wy-mono" value={appConfig().NodePath}></input>
        </Field>
        <Field id="shout" label="Notification URL" help={<>Shoutrrr link for Discord, Email, Gotify, Telegram and more. <a class="wy-link" href="https://shoutrrr.nickfedor.com/services/overview/" target="_blank">Formats</a>. Leave empty to turn alerts off.</>}>
          <textarea id="shout" name="shout" class="wy-input wy-mono" rows="3" wrap="soft">{appConfig().ShoutURL}</textarea>
          <button onClick={handleTestNotify} type="button" class="wy-btn" style={{ "align-self": "flex-start" }}>
            {sent() ? "Test sent" : "Send test notification"}
          </button>
        </Field>
        <div class="wy-actions">
          <button type="submit" class="wy-btn primary">Save general settings</button>
        </div>
      </form>
    </section>
  )
}

export default Basic
