import { createResource, createSignal, For, Show } from "solid-js";
import { apiDNSTop } from "../../functions/api";
import { DNSStat } from "../../functions/exports";

// Top domains this device looked up, from Pi-hole / AdGuard Home.
function DNSCard(_props: { ip: string }) {

  const [days, setDays] = createSignal(1);
  const [data] = createResource(() => _props.ip ? [_props.ip, days()] as const : undefined, ([ip, d]) => apiDNSTop(ip, d));
  const top = () => (data()?.Top ?? []) as DNSStat[];

  return (
    <Show when={data()?.Enabled}>
      <section class="wy-card" aria-label="Top sites">
        <div class="wy-row-wrap">
          <h2 style={{ "flex-grow": 1 }}>Top sites</h2>
          <label class="wy-inline-label">Range
            <select class="wy-input" value={days()} onChange={e => setDays(+e.currentTarget.value)}>
              <option value="1">Today</option>
              <option value="7">7 days</option>
            </select>
          </label>
        </div>
        <Show when={data()?.Error}>
          <div class="wy-alert" role="alert">DNS log: {data().Error}</div>
        </Show>
        <Show when={top().length} fallback={<div class="wy-muted">No DNS queries from {_props.ip} yet.</div>}>
          <table class="wy-table">
            <thead><tr><th>Domain</th><th class="num">Lookups</th><th class="num">Blocked</th></tr></thead>
            <tbody>
              <For each={top()}>{(s) =>
                <tr>
                  <td class="wy-mono" style={{ "word-break": "break-all" }}>
                    <Show when={s.Threat}><span class="wy-sev high" style={{ "margin-right": "6px" }}>threat</span></Show>{s.Domain}
                  </td>
                  <td class="num">{s.Count}</td>
                  <td class="num">{s.Blocked || ""}</td>
                </tr>
              }</For>
            </tbody>
          </table>
        </Show>
      </section>
    </Show>
  )
}

export default DNSCard
