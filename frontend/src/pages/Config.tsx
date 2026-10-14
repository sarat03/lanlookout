import { createSignal, For, Match, Switch } from "solid-js"
import About from "../components/Config/About"
import Basic from "../components/Config/Basic"
import Influx from "../components/Config/Influx"
import Prometheus from "../components/Config/Prometheus"
import Mqtt from "../components/Config/Mqtt"
import Scan from "../components/Config/Scan"
import Alerts from "../components/Config/Alerts"
import Monitor from "../components/Config/Monitor"
import DNSLog from "../components/Config/DNSLog"
import { saveInPlace } from "../functions/save"

const tabs = [
  { k: "general", label: "General" },
  { k: "scan", label: "Scanning" },
  { k: "integ", label: "Integrations" },
  { k: "about", label: "About" },
];

function Config() {

  // Tab lives in the query string so a reload (theme change) keeps it.
  const [tab, setTabRaw] = createSignal(new URLSearchParams(window.location.search).get("tab") || "general");
  const setTab = (k: string) => { setTabRaw(k); history.replaceState(null, "", "?tab=" + k); };

  return (
    <>
    <header class="wy-head">
      <div>
        <h1>Settings</h1>
        <div class="wy-sub">Each section saves on its own. Saving scan settings restarts the scanner.</div>
      </div>
    </header>

    <div class="wy-row-wrap" role="group" aria-label="Sections" style={{ gap: "8px" }}>
      <For each={tabs}>{(t) =>
        <button type="button" class="wy-btn wy-chip" aria-pressed={tab() === t.k} onClick={() => setTab(t.k)}>{t.label}</button>
      }</For>
    </div>

    <div style={{ display: "flex", "flex-direction": "column", gap: "20px", "max-width": "980px" }} onSubmit={saveInPlace}>
      <Switch fallback={<><Basic></Basic><Alerts></Alerts></>}>
        <Match when={tab() === "scan"}><Scan></Scan><Monitor></Monitor></Match>
        <Match when={tab() === "integ"}><DNSLog></DNSLog><Mqtt></Mqtt><Prometheus></Prometheus><Influx></Influx></Match>
        <Match when={tab() === "about"}><About></About></Match>
      </Switch>
    </div>
    </>
  )
}

export default Config
