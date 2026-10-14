import { createResource, createSignal, For, Show } from "solid-js";
import { apiGetEvents } from "../functions/api";
import { allHosts, Event } from "../functions/exports";

const look: Record<Event["Kind"], { icon: string, color: string, label: string }> = {
  new: { icon: "bi-plus-circle-fill", color: "var(--wy-warn)", label: "New device" },
  offline: { icon: "bi-arrow-down-circle-fill", color: "var(--wy-off)", label: "Offline" },
  online: { icon: "bi-arrow-up-circle-fill", color: "var(--wy-on)", label: "Back online" },
  ip: { icon: "bi-arrow-left-right", color: "var(--wy-accent)", label: "IP changed" },
  outage: { icon: "bi-wifi-off", color: "var(--wy-danger)", label: "Internet outage" },
  wan_ip: { icon: "bi-globe2", color: "var(--wy-accent)", label: "Public IP changed" },
  port: { icon: "bi-shield-exclamation", color: "var(--wy-warn)", label: "Risky port" },
  threat: { icon: "bi-exclamation-octagon-fill", color: "var(--wy-danger)", label: "Threat domain" },
};

function Events() {

  const [num, setNum] = createSignal(20);
  const [events] = createResource(num, (n) => apiGetEvents(n) as Promise<Event[]>);
  const hostID = (mac: string) => allHosts.find(h => h.Mac === mac)?.ID;

  return (
    <section class="wy-card" aria-label="Recent events">
      <h2>Recent events</h2>
      <Show when={events()?.length} fallback={<div class="wy-muted">{events.loading ? "Loading…" : "Nothing yet. New devices, watched devices going offline and IP changes show up here."}</div>}>
        <ul class="wy-events">
          <For each={events()}>{(e) =>
            <li>
              <i class={"bi " + (look[e.Kind]?.icon ?? "bi-dot")} style={{ color: look[e.Kind]?.color }} title={look[e.Kind]?.label} aria-hidden="true"></i>
              <span class="wy-muted wy-mono">{e.Date}</span>
              <Show when={hostID(e.Mac)} fallback={<span>{e.Msg}</span>}>
                <a href={"/host/" + hostID(e.Mac)}>{e.Msg}</a>
              </Show>
            </li>
          }</For>
        </ul>
        <Show when={events()!.length >= num()}>
          <button type="button" class="wy-btn" style={{ "align-self": "flex-start" }} onClick={() => setNum(num() + 50)}>Show more</button>
        </Show>
      </Show>
    </section>
  )
}

export default Events
