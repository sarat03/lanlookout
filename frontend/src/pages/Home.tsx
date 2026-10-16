import { createEffect, createMemo, createResource, createSignal, For, Show } from "solid-js";
import { apiGetEvents, apiNetStatus, apiSaveMemberOrder } from "../functions/api";
import { appConfig, bkpHosts, Event, Host, hostTypes, isShared, setAppConfig } from "../functions/exports";

// Types that tell us whether a person is home.
const presenceTypes = ["phone", "watch"];

const byName = (a: Host, b: Host) => (a.Name || a.IP).localeCompare(b.Name || b.IP, undefined, { numeric: true });
const label = (h: Host) => h.Name || h.DNS || h.IP;
const online = (hs: Host[]) => hs.filter(h => h.Now == 1).length;
const today = () => new Date().toISOString().slice(0, 10);

// "2026-10-16 08:12:03" -> "08:12" today, else "Oct 15 08:12".
function when(date?: string) {
  if (!date) return "";
  const t = date.slice(11, 16);
  if (date.startsWith(today())) return t;
  const d = new Date(date.slice(0, 10) + "T00:00:00");
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric" }) + " " + t;
}

function DeviceRow(_props: { host: Host }) {
  const h = _props.host;
  return (
    <li>
      <span class={"wy-dot" + (h.Now == 1 ? " on" : "")} title={h.Now == 1 ? "Online" : "Offline"}></span>
      <Show when={hostTypes[h.Type]}><i class={"bi " + hostTypes[h.Type] + " wy-muted"} aria-hidden="true"></i></Show>
      <a class="wy-link" href={"/host/" + h.ID}>{label(h)}</a>
    </li>
  );
}

function Home() {
  const [net] = createResource(apiNetStatus);
  const [events] = createResource(async () => ((await apiGetEvents(500)) || []) as Event[]);

  // Latest online/offline event per MAC (events come newest first).
  const lastSeen = createMemo(() => {
    const m = new Map<string, Event>();
    for (const e of events() ?? []) if ((e.Kind === "offline" || e.Kind === "online") && !m.has(e.Mac)) m.set(e.Mac, e);
    return m;
  });

  const all = () => bkpHosts();
  // Shared devices are the house's. Unreviewed ones wait in Needs attention instead.
  const smart = createMemo(() => all().filter(h => isShared(h) && h.Known == 1).sort(byName));
  const newToday = createMemo(() => (events() ?? []).filter(e => e.Kind === "new" && e.Date.startsWith(today())).length);

  // Offline smart-home gear, offline watched devices, and unreviewed devices on the network.
  const attention = createMemo(() => {
    const out: { host: Host, why: string, tag: string }[] = [];
    for (const h of all()) {
      const ev = lastSeen().get(h.Mac);
      const since = ev?.Kind === "offline" ? "offline since " + when(ev.Date) : "offline";
      if (h.Known != 1 && h.Now == 1) out.push({ host: h, why: "new device, not reviewed yet", tag: "new" });
      else if (h.Now != 1 && h.Watch == 1) out.push({ host: h, why: since, tag: "watched" });
      else if (h.Now != 1 && h.Known == 1 && isShared(h)) out.push({ host: h, why: since, tag: "smart home" });
    }
    return out.sort((a, b) => byName(a.host, b.host));
  });

  // Saved card order; members not in it go last, alphabetically.
  const [order, setOrder] = createSignal<string[]>([]);
  createEffect(() => setOrder((appConfig().MemberOrder || "").split(",").filter(Boolean)));
  const rank = (name: string) => { const i = order().indexOf(name); return i < 0 ? Infinity : i; };

  const members = createMemo(() => {
    const m = new Map<string, Host[]>();
    for (const h of all()) if (!isShared(h)) m.set(h.Owner, [...(m.get(h.Owner) || []), h]);
    return [...m.entries()].sort(([a], [b]) => rank(a) - rank(b) || a.localeCompare(b)).map(([name, hs]) => {
      hs.sort(byName);
      const tracked = hs.filter(h => presenceTypes.includes(h.Type));
      // No phone or watch assigned: we can't tell, so say nothing.
      const home = tracked.length ? tracked.some(h => h.Now == 1) : undefined;
      let seen = "";
      if (home === false) {
        const last = tracked.map(h => lastSeen().get(h.Mac)).filter(e => e?.Kind === "offline").map(e => e!.Date).sort().pop();
        if (last) seen = "last seen " + when(last);
      }
      return { name, hosts: hs, home, seen };
    });
  });

  const [showAll, setShowAll] = createSignal(false);

  // Cards start closed; anyone with a device offline starts open.
  const [open, setOpen] = createSignal<Record<string, boolean>>({});
  const isOpen = (m: { name: string, hosts: Host[] }) => open()[m.name] ?? online(m.hosts) < m.hosts.length;
  const toggle = (m: { name: string, hosts: Host[] }) => {
    if (dragged) { dragged = false; return; } // the click that ends a drag
    setOpen({ ...open(), [m.name]: !isOpen(m) });
  };

  // Drag a card by its avatar. Pointer events so it works with touch too.
  const [dragging, setDragging] = createSignal("");
  let dragged = false;
  const dragStart = (e: PointerEvent, name: string) => {
    if (e.button !== 0) return;
    e.preventDefault();
    setDragging(name);
    dragged = false;
    // On the window: reordering re-renders the card, so the handle under the finger goes away.
    window.addEventListener("pointermove", dragMove);
    window.addEventListener("pointerup", dragEnd);
    window.addEventListener("pointercancel", dragEnd);
  };
  const dragMove = (e: PointerEvent) => {
    const name = dragging();
    if (!name) return;
    const over = (document.elementFromPoint(e.clientX, e.clientY)?.closest("[data-member]") as HTMLElement | null)?.dataset.member;
    if (!over || over === name) return;
    const now = members().map(m => m.name);
    const down = now.indexOf(over) > now.indexOf(name);
    const names = now.filter(n => n !== name);
    names.splice(names.indexOf(over) + (down ? 1 : 0), 0, name);
    setOrder(names);
    dragged = true;
  };
  const dragEnd = () => {
    window.removeEventListener("pointermove", dragMove);
    window.removeEventListener("pointerup", dragEnd);
    window.removeEventListener("pointercancel", dragEnd);
    if (!dragging()) return;
    setDragging("");
    if (!dragged) return;
    setTimeout(() => dragged = false); // after the click that ends the drag, if any
    setAppConfig({ ...appConfig(), MemberOrder: order().join(",") });
    apiSaveMemberOrder(order());
  };

  const internet = () => {
    const s = net() as any;
    if (!s) return { value: "…", note: "", cls: "" };
    if (!s.Up) return { value: "Down", note: s.DownSince ? "since " + when(s.DownSince) : "", cls: "bad" };
    const l = s.Last || {};
    const parts = [l.LatencyMs ? Math.round(l.LatencyMs) + " ms" : "", l.DownMbps ? Math.round(l.DownMbps) + " Mbps" : ""].filter(Boolean);
    return { value: "Up", note: parts.join(" · "), cls: "good" };
  };

  return (
    <>
      <header class="wy-head">
        <div>
          <h1>Home</h1>
          <div class="wy-sub">Who's home and what needs attention.</div>
        </div>
      </header>

      <div class="wy-hometiles">
        <div class="wy-card"><div class="wy-muted">Devices online</div><div class="wy-big">{online(all())}/{all().length}</div></div>
        <div class="wy-card"><div class="wy-muted">Offline now</div><div class={"wy-big" + (all().length - online(all()) ? " bad" : "")}>{all().length - online(all())}</div><div class="wy-muted">{smart().length - online(smart())} smart-home</div></div>
        <div class="wy-card"><div class="wy-muted">New today</div><div class={"wy-big" + (newToday() ? " bad" : "")}>{newToday()}</div></div>
        <a class="wy-card" href="/internet/"><div class="wy-muted">Internet</div><div class={"wy-big " + internet().cls}>{internet().value}</div><div class="wy-muted">{internet().note}</div></a>
      </div>

      <Show when={attention().length}>
        <section class="wy-card wy-attn" aria-label="Needs attention">
          <h2><i class="bi bi-exclamation-triangle" aria-hidden="true"></i> Needs attention <span class="wy-muted">{attention().length}</span></h2>
          <ul class="wy-attn-list">
            <For each={showAll() ? attention() : attention().slice(0, 6)}>{(a) =>
              <li>
                <span class={"wy-dot" + (a.host.Now == 1 ? " on" : "")}></span>
                <Show when={hostTypes[a.host.Type]}><i class={"bi " + hostTypes[a.host.Type] + " wy-muted"} aria-hidden="true"></i></Show>
                <a class="wy-link" href={"/host/" + a.host.ID}>{label(a.host)}</a>
                <span class="wy-muted">{a.why}</span>
                <span class="wy-pill">{a.tag}</span>
              </li>
            }</For>
          </ul>
          <Show when={attention().length > 6}>
            <button type="button" class="wy-link wy-more" onClick={() => setShowAll(!showAll())}>
              {showAll() ? "Show less" : "Show " + (attention().length - 6) + " more"}
            </button>
          </Show>
        </section>
      </Show>

      <div class="wy-people">
        <For each={members()} fallback={<section class="wy-card"><p class="wy-muted">No members yet. Set an Owner on a device to add someone.</p></section>}>{(m) =>
          <section class={"wy-card wy-person" + (isOpen(m) ? " open" : "") + (dragging() === m.name ? " dragging" : "")} data-member={m.name}>
            <button type="button" class="wy-person-head" aria-expanded={isOpen(m)} onClick={() => toggle(m)}>
              <i class="bi bi-grip-vertical wy-grip" title="Drag to reorder" aria-hidden="true"
                onPointerDown={e => dragStart(e, m.name)}></i>
              <span class={"wy-avatar" + (m.home ? " home" : "")} title="Drag to reorder"
                onPointerDown={e => dragStart(e, m.name)}>{m.name.charAt(0).toUpperCase()}</span>
              <span class="wy-person-text">
                <strong>{m.name}</strong>
                <span>
                  <Show when={m.home !== undefined}><span class={"wy-pill" + (m.home ? " good" : "")}>{m.home ? "Home" : "Away"}</span> </Show>
                  <span class="wy-muted">{m.seen}</span>
                </span>
                <span class="wy-muted">{online(m.hosts)}/{m.hosts.length} devices online</span>
              </span>
              <i class={"bi wy-muted " + (isOpen(m) ? "bi-chevron-up" : "bi-chevron-down")} aria-hidden="true"></i>
            </button>
            <Show when={isOpen(m)}>
              <ul class="wy-devrows"><For each={m.hosts}>{(h) => <DeviceRow host={h}/>}</For></ul>
            </Show>
          </section>
        }</For>
      </div>

      <section class="wy-card" aria-label="Smart home">
        <h2><i class="bi bi-house-gear" aria-hidden="true"></i> Smart home <span class="wy-muted">{online(smart())}/{smart().length} online</span></h2>
        <Show when={smart().length} fallback={<p class="wy-muted">No shared devices yet. Set a device's Owner to Shared to add it here.</p>}>
          <ul class="wy-devrows wy-devcols"><For each={smart()}>{(h) => <DeviceRow host={h}/>}</For></ul>
        </Show>
      </section>
    </>
  );
}

export default Home
