import { createMemo, createResource, createSignal, For, Show } from "solid-js";
import { apiGetEvents, apiNetStats, apiNetStatus, apiSpeedTest } from "../functions/api";
import { appConfig, Event, NetStat } from "../functions/exports";

const W = 800, H = 160, PAD = 28;

// Latency line with outage checks shaded red. X is check order (checks are evenly spaced).
function LatencyChart(_props: { stats: NetStat[] }) {
  const pts = createMemo(() => {
    const s = _props.stats;
    const max = Math.max(10, ...s.filter(x => x.Up).map(x => x.LatencyMs));
    const x = (i: number) => PAD + (s.length < 2 ? 0 : i * (W - PAD) / (s.length - 1));
    const y = (ms: number) => H - 18 - (ms / max) * (H - 30);
    const segs: string[] = [];
    let cur = "";
    s.forEach((p, i) => {
      if (p.Up) cur += `${x(i).toFixed(1)},${y(p.LatencyMs).toFixed(1)} `;
      else if (cur) { segs.push(cur); cur = ""; }
    });
    if (cur) segs.push(cur);
    const step = s.length < 2 ? W : (W - PAD) / (s.length - 1);
    return { max, segs, down: s.map((p, i) => p.Up ? null : x(i) - step / 2).filter(v => v !== null) as number[], step };
  });

  return (
    <svg class="wy-chart" viewBox={`0 0 ${W} ${H}`} preserveAspectRatio="none" role="img"
      aria-label={`Latency over time, peak ${Math.round(pts().max)} ms`}>
      <line class="grid" x1={PAD} x2={W} y1={12} y2={12}/>
      <line class="grid" x1={PAD} x2={W} y1={H - 18} y2={H - 18}/>
      <text x="0" y="16">{Math.round(pts().max)}</text>
      <text x="0" y={H - 14}>0 ms</text>
      <For each={pts().down}>{(x) => <rect class="down" x={x} y={12} width={Math.max(pts().step, 2)} height={H - 30}/>}</For>
      <For each={pts().segs}>{(p) => <polyline class="line" points={p} vector-effect="non-scaling-stroke"/>}</For>
      <Show when={_props.stats.length}>
        <text x={PAD} y={H - 2}>{_props.stats[0].Date.slice(5, 16)}</text>
        <text x={W} y={H - 2} text-anchor="end">{_props.stats[_props.stats.length - 1].Date.slice(5, 16)}</text>
      </Show>
    </svg>
  );
}

// Status strip: one bar per check, green = online, red = down. Hover a bar for its time.
function UptimeChart(_props: { stats: NetStat[] }) {
  const UH = 40;
  const step = () => (W - PAD) / Math.max(_props.stats.length, 1);
  return (
    <svg class="wy-chart" viewBox={`0 0 ${W} ${UH + 16}`} preserveAspectRatio="none" role="img"
      aria-label="Online status over time" style={{ height: "72px" }}>
      <For each={_props.stats}>{(p, i) =>
        <rect x={PAD + i() * step()} y={0} width={Math.max(step() - 1, 1)} height={UH}
          fill={p.Up ? "var(--wy-on)" : "var(--wy-danger)"}>
          <title>{p.Date.slice(5, 16)} · {p.Up ? "online" : "down"}</title>
        </rect>
      }</For>
      <text x={PAD} y={UH + 14}>{_props.stats[0].Date.slice(5, 16)}</text>
      <text x={W} y={UH + 14} text-anchor="end">{_props.stats[_props.stats.length - 1].Date.slice(5, 16)}</text>
    </svg>
  );
}

function Internet() {

  const [hours, setHours] = createSignal(24);
  const [status, { refetch: refetchStatus }] = createResource(apiNetStatus);
  const [stats, { refetch: refetchStats }] = createResource(hours, apiNetStats);
  const [events] = createResource(async () => (await apiGetEvents(200) as Event[]).filter(e => e.Kind === "outage" || e.Kind === "wan_ip"));
  const [testing, setTesting] = createSignal(false);
  const [testErr, setTestErr] = createSignal("");

  const list = () => (stats() ?? []) as NetStat[];
  const speeds = createMemo(() => list().filter(s => s.DownMbps > 0).reverse());
  const uptime = createMemo(() => list().length ? Math.round(1000 * list().filter(s => s.Up).length / list().length) / 10 : null);
  const avgLatency = createMemo(() => {
    const up = list().filter(s => s.Up);
    return up.length ? Math.round(up.reduce((a, s) => a + s.LatencyMs, 0) / up.length) : null;
  });

  const runTest = async () => {
    setTesting(true);
    setTestErr("");
    try {
      await apiSpeedTest();
      refetchStats();
      refetchStatus();
    } catch (e) {
      setTestErr(String((e as Error).message || e));
    }
    setTesting(false);
  };

  return (
    <>
    <header class="wy-head">
      <div>
        <h1>Internet</h1>
        <div class="wy-sub">
          {appConfig().NetCheckMin > 0 ? `Checked every ${appConfig().NetCheckMin} min` : "Checks are off (Settings → Scanning)"}
          {appConfig().SpeedHours > 0 ? ` · speed test every ${appConfig().SpeedHours} h` : ""}
        </div>
      </div>
      <button type="button" class="wy-btn primary" onClick={runTest} disabled={testing()}>
        <i class="bi bi-speedometer2" aria-hidden="true"></i>{testing() ? "Testing… (up to a minute)" : "Run speed test"}
      </button>
    </header>

    <Show when={testErr()}>
      <div class="wy-alert" role="alert">Speed test failed: {testErr()}</div>
    </Show>

    <section class="wy-stats" aria-label="Summary">
      <div class="wy-stat">
        <span class="k"><span class={"wy-dot " + (status()?.Up ? "on" : "")}></span>Status</span>
        <span class="n">{status() ? (status().Up ? "Online" : "Down") : "…"}</span>
        <span class="h">{status()?.DownSince ? "Down since " + status().DownSince : `Uptime ${uptime() ?? "–"}% in ${hours()} h`}</span>
      </div>
      <div class="wy-stat">
        <span class="k">Latency</span>
        <span class="n">{status()?.Last?.Up ? Math.round(status().Last.LatencyMs) + " ms" : "–"}</span>
        <span class="h">Average {avgLatency() ?? "–"} ms</span>
      </div>
      <div class="wy-stat">
        <span class="k">Speed</span>
        <span class="n">{speeds()[0] ? `${Math.round(speeds()[0].DownMbps)} ↓` : "–"}</span>
        <span class="h">{speeds()[0] ? `${Math.round(speeds()[0].UpMbps)} Mbps up · ${speeds()[0].Date.slice(5, 16)}` : "No test yet"}</span>
      </div>
      <div class="wy-stat">
        <span class="k">Public IP</span>
        <span class="n wy-mono" style={{ "font-size": "20px" }}>{status()?.PublicIP || "–"}</span>
        <span class="h">{appConfig().PublicIPURL ? "Alerts when it changes" : "Lookup is off"}</span>
      </div>
    </section>

    <section class="wy-card" aria-label="Online status">
      <div class="wy-row-wrap">
        <h2 style={{ "flex-grow": 1 }}>Online status <span class="wy-muted" style={{ "font-weight": 400 }}>{uptime() ?? "–"}% uptime</span></h2>
        <label class="wy-inline-label">Range
          <select class="wy-input" value={hours()} onChange={e => setHours(+e.currentTarget.value)}>
            <option value="6">6 hours</option>
            <option value="24">24 hours</option>
            <option value="168">7 days</option>
          </select>
        </label>
      </div>
      <Show when={list().length} fallback={<div class="wy-muted">{stats.loading ? "Loading…" : "No checks yet. The first one runs when the app starts."}</div>}>
        <UptimeChart stats={list()}/>
        <div class="wy-legend">
          <span><i style={{ background: "var(--wy-on)" }}></i>Online</span>
          <span><i style={{ background: "var(--wy-danger)" }}></i>Down</span>
        </div>
      </Show>
    </section>

    <section class="wy-card" aria-label="Latency">
      <h2>Latency</h2>
      <Show when={list().length} fallback={<div class="wy-muted">{stats.loading ? "Loading…" : "No checks yet."}</div>}>
        <LatencyChart stats={list()}/>
        <div class="wy-legend">
          <span><i style={{ background: "var(--wy-accent)" }}></i>Latency to 1.1.1.1 / 8.8.8.8</span>
          <span><i style={{ background: "var(--wy-danger)", opacity: .5 }}></i>Internet down</span>
        </div>
      </Show>
    </section>

    <div class="wy-cols">
      <section class="wy-card" aria-label="Outages">
        <h2>Outages and IP changes</h2>
        <Show when={events()?.length} fallback={<div class="wy-muted">None recorded.</div>}>
          <ul class="wy-events">
            <For each={events()}>{(e) =>
              <li>
                <i class={"bi " + (e.Kind === "outage" ? "bi-wifi-off" : "bi-arrow-left-right")} style={{ color: e.Kind === "outage" ? "var(--wy-danger)" : "var(--wy-accent)" }} aria-hidden="true"></i>
                <span class="wy-muted wy-mono">{e.Date}</span><span>{e.Msg}</span>
              </li>
            }</For>
          </ul>
        </Show>
      </section>
      <section class="wy-card wide" aria-label="Speed tests">
        <h2>Speed tests</h2>
        <Show when={speeds().length} fallback={<div class="wy-muted">No speed tests in this range. Run one, or schedule them in Settings → Scanning.</div>}>
          <table class="wy-table">
            <thead><tr><th>When</th><th class="num">Down (Mbps)</th><th class="num">Up (Mbps)</th><th class="num">Latency</th></tr></thead>
            <tbody>
              <For each={speeds().slice(0, 20)}>{(s) =>
                <tr><td class="wy-mono">{s.Date.slice(5, 16)}</td><td class="num">{s.DownMbps}</td><td class="num">{s.UpMbps}</td><td class="num">{Math.round(s.LatencyMs)} ms</td></tr>
              }</For>
            </tbody>
          </table>
        </Show>
      </section>
    </div>
    </>
  )
}

export default Internet
