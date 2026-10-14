import { createResource, createSignal, onMount, Show } from "solid-js";
import { Host, setShow } from "../../functions/exports";
import { apiGetHistory } from "../../functions/api";
import MacHistory from "../MacHistory"

function HistCard(_props: any) {

  const [today, setToday] = createSignal('');

  onMount(() => {
    setShow(15000);
    setToday(new Date().toLocaleDateString("en-CA"));
  });

  // Uptime over the latest scans the API returns (up to 210).
  const [uptime] = createResource(() => _props.mac || undefined, async (mac: string) => {
    const hist: Host[] = await apiGetHistory(mac);
    if (!hist?.length) return null;
    const on = hist.filter(h => h.Now == 1).length;
    return { pct: Math.round(100 * on / hist.length), n: hist.length, since: hist[hist.length - 1].Date };
  });

  const handleDate = (date: string) => {
    setToday("");
    setToday(date);
  };

  return (
    <section class="wy-card" aria-label="History">
      <div class="wy-row-wrap">
        <div style={{ "flex-grow": 1 }}>
          <h2>Online history</h2>
          <Show when={uptime()}>{(u) =>
            <div class="wy-sub">Online <b>{u().pct}%</b> of the last {u().n} scans (since {u().since})</div>
          }</Show>
        </div>
        <label class="wy-inline-label">Day
          <input type="date" class="wy-input" value={today()}
            onInput={(e) => handleDate(e.currentTarget.value)}/>
        </label>
      </div>
      <Show when={_props.mac !== "" && today() !== ""} fallback={<div class="wy-muted">Loading…</div>}>
        <div class="wy-hist">
          <MacHistory mac={_props.mac} date={today()}></MacHistory>
        </div>
      </Show>
      <div class="wy-legend">
        <span><i style={{ background: "var(--wy-on)" }}></i>Online</span>
        <span><i style={{ background: "var(--wy-off)" }}></i>Offline</span>
        <span class="wy-muted">Newest first · hover a bar for details</span>
      </div>
    </section>
  )
}

export default HistCard
