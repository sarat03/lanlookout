import { createMemo, For, Show } from "solid-js";
import { bkpHosts, Host, hostTypes, smartTypes } from "../functions/exports";
import { hostStatus } from "../components/Body/TableRow";

function DeviceList(_props: { hosts: Host[] }) {
  return (
    <ul style={{ "list-style": "none", margin: 0, padding: 0, display: "grid", gap: "8px" }}>
      <For each={_props.hosts}>{(h) => {
        const st = hostStatus(h);
        return (
          <li style={{ display: "flex", "align-items": "center", gap: "8px" }}>
            <span class={"wy-dot " + st.cls} title={st.label} aria-label={st.label}></span>
            <Show when={hostTypes[h.Type]}><i class={"bi " + hostTypes[h.Type] + " wy-muted"} aria-hidden="true"></i></Show>
            <a class="wy-link" href={"/host/" + h.ID}>{h.Name || h.DNS || h.IP}</a>
          </li>
        );
      }}</For>
    </ul>
  );
}

function Home() {
  const byName = (a: Host, b: Host) => (a.Name || a.IP).localeCompare(b.Name || b.IP, undefined, { numeric: true });

  // Smart-home gear without an owner lives in the shared home card.
  const smart = createMemo(() => bkpHosts().filter(h => !h.Owner && smartTypes.includes(h.Type)).sort(byName));

  const members = createMemo(() => {
    const m = new Map<string, Host[]>();
    for (const h of bkpHosts()) if (h.Owner) m.set(h.Owner, [...(m.get(h.Owner) || []), h]);
    return [...m.entries()].sort(([a], [b]) => a.localeCompare(b)).map(([name, hs]) => ({ name, hosts: hs.sort(byName) }));
  });

  const online = (hs: Host[]) => hs.filter(h => h.Now == 1).length;

  return (
    <>
      <header class="wy-head">
        <div>
          <h1>Home</h1>
          <div class="wy-sub">Devices by household member. Set the Owner on a device to assign it.</div>
        </div>
      </header>
      <section class="wy-card" aria-label="Smart home">
        <h2><i class="bi bi-house-gear" aria-hidden="true"></i> Smart home <span class="wy-muted">{online(smart())}/{smart().length} online</span></h2>
        <Show when={smart().length} fallback={<p class="wy-muted">No smart-home devices yet. Set a device's type to plug, light, camera, etc.</p>}>
          <DeviceList hosts={smart()}/>
        </Show>
      </section>
      <For each={members()} fallback={<section class="wy-card"><p class="wy-muted">No members yet. Set an Owner on a device to add someone.</p></section>}>{(m) =>
        <section class="wy-card" aria-label={m.name}>
          <h2><i class="bi bi-person" aria-hidden="true"></i> {m.name} <span class="wy-muted">{online(m.hosts)}/{m.hosts.length} online</span></h2>
          <DeviceList hosts={m.hosts}/>
        </section>
      }</For>
    </>
  );
}

export default Home
