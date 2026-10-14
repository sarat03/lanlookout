import { createMemo, createResource, createSignal, For, Show } from "solid-js";
import { apiSecurity, apiSecurityScan } from "../functions/api";
import { appConfig, Finding } from "../functions/exports";

function Security() {

  const [data, { refetch }] = createResource(apiSecurity);
  const [starting, setStarting] = createSignal(false);
  const [showLow, setShowLow] = createSignal(false);

  const findings = createMemo(() => ((data()?.Findings ?? []) as Finding[]));
  const shown = createMemo(() => findings().filter(f => showLow() || f.Severity !== "low"));
  const count = (s: string) => findings().filter(f => f.Severity === s).length;
  const running = () => starting() || data()?.Running;

  const scan = async () => {
    setStarting(true);
    await apiSecurityScan();
    await new Promise(r => setTimeout(r, 1000));
    // Poll until the background scan finishes.
    for (let i = 0; i < 300; i++) {
      await refetch();
      if (!data()?.Running) break;
      await new Promise(r => setTimeout(r, 2000));
    }
    setStarting(false);
  };

  return (
    <>
    <header class="wy-head">
      <div>
        <h1>Security</h1>
        <div class="wy-sub">
          Services on your devices that are risky to leave open.
          {data()?.LastScan ? " Last scan " + data().LastScan + "." : findings().length ? " Results from an earlier scan." : " Not scanned since the app started."}
          {appConfig().SecScanHours > 0 ? ` Runs every ${appConfig().SecScanHours} h.` : ""}
        </div>
      </div>
      <button type="button" class="wy-btn primary" onClick={scan} disabled={running()}>
        <i class="bi bi-shield-check" aria-hidden="true"></i>{running() ? "Scanning…" : "Scan now"}
      </button>
    </header>

    <section class="wy-stats" aria-label="Summary">
      <div class="wy-stat"><span class="k"><span class="wy-sev high">High</span></span><span class="n">{count("high")}</span><span class="h">Fix these first</span></div>
      <div class="wy-stat"><span class="k"><span class="wy-sev medium">Medium</span></span><span class="n">{count("medium")}</span><span class="h">Check they need to be open</span></div>
      <div class="wy-stat"><span class="k"><span class="wy-sev low">Low</span></span><span class="n">{count("low")}</span><span class="h">Usually fine (SSH, web UIs)</span></div>
    </section>

    <section class="wy-card" aria-label="Findings">
      <div class="wy-row-wrap">
        <h2 style={{ "flex-grow": 1 }}>Findings</h2>
        <label class="wy-inline-label">
          <input type="checkbox" role="switch" class="wy-switch" checked={showLow()} onChange={e => setShowLow(e.currentTarget.checked)}/>Show low
        </label>
      </div>
      <Show when={shown().length} fallback={<div class="wy-muted">
        {data.loading ? "Loading…" : data()?.LastScan ? "Nothing risky found." : "Nothing risky found so far. Press Scan now to check again."}
      </div>}>
        <ul class="wy-list" style={{ "list-style": "none", padding: 0, margin: 0 }}>
          <For each={shown()}>{(f) =>
            <li class="wy-host" style={{ "align-items": "flex-start" }}>
              <span class={"wy-sev " + f.Severity}>{f.Severity}</span>
              <div style={{ flex: "1 1 260px", "min-width": 0 }}>
                <div style={{ "font-weight": 600 }}>{f.Service}{f.Port ? ` · port ${f.Port}` : ""}</div>
                <div class="wy-muted" style={{ "font-size": "14px" }}>{f.Advice}</div>
              </div>
              <a class="wy-link" href={"/host/" + f.HostID} style={{ "white-space": "nowrap" }}>{f.Name || f.IP}<span class="wy-muted wy-mono"> {f.IP}</span></a>
            </li>
          }</For>
        </ul>
      </Show>
      <div class="wy-muted" style={{ "font-size": "13px" }}>
        Checks 21 well-known ports on online devices from inside your network, and whether the router allows UPnP port forwarding.
        It can't see what's exposed to the internet; check your router's port-forwarding page for that.
      </div>
    </section>
    </>
  )
}

export default Security
