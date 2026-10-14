import { createSignal, For, Show } from "solid-js";
import { editNames, Host, hostTypes, selectedIDs, setSelectedIDs, tagList } from "../../functions/exports";
import { apiEditHost, apiWOL } from "../../functions/api";
import { getHosts } from "../../functions/atstart";

import { debounce } from "@solid-primitives/scheduled";

export function hostStatus(h: Host) {
  const on = h.Now == 1;
  if (h.Known != 1) {
    return { cls: on ? "new on" : "new", label: on ? "Online, not yet reviewed" : "Offline, not yet reviewed" };
  }
  return { cls: on ? "on" : "", label: on ? "Online" : "Offline" };
}

function TableRow(_props: { host: Host }) {

  const h = _props.host;
  const [name, setName] = createSignal(h.Name);
  const [woke, setWoke] = createSignal(false);
  const st = hostStatus(h);

  const debouncedApi = debounce(async (val: string) => {
    await apiEditHost(h.ID, val, "");
  }, 300);

  const handleInput = (n: string) => {
    setName(n);
    debouncedApi(n);
  };

  const handleToggle = async () => {
    await apiEditHost(h.ID, name(), "toggle");
    getHosts();
  };

  const handleWOL = async () => {
    await apiWOL(h.Mac);
    setWoke(true);
    setTimeout(() => setWoke(false), 2000);
  };

  const handleCheck = (checked: boolean) => {
    setSelectedIDs(prev => checked
      ? (prev.includes(h.ID) ? prev : [...prev, h.ID])
      : prev.filter(item => item !== h.ID));
  };

  return (
    <div class="wy-host">
      <Show when={editNames()}>
        <input type="checkbox" class="wy-check" aria-label={"Select " + (name() || h.IP)}
          checked={selectedIDs().includes(h.ID)}
          onChange={e => handleCheck(e.currentTarget.checked)}/>
      </Show>
      <div class="who">
        <span class={"wy-dot " + st.cls} title={st.label}></span>
        <div style={{ "min-width": 0, flex: 1 }}>
          <Show when={editNames()} fallback={
            <a class="name" href={"/host/" + h.ID}>
              <Show when={hostTypes[h.Type]}><i class={"bi " + hostTypes[h.Type] + " wy-muted"} title={h.Type} aria-hidden="true" style={{ "margin-right": "6px" }}></i></Show>
              {name() || "Unnamed device"}
              <Show when={h.Watch == 1}><i class="bi bi-bell-fill wy-muted" title="Watched: alerts on offline/online" aria-label="Watched" style={{ "margin-left": "6px", "font-size": "12px" }}></i></Show>
            </a>
          }>
            <input type="text" class="wy-input" style={{ width: "100%" }} value={name()} aria-label="Host name"
              onInput={e => handleInput(e.currentTarget.value)}/>
          </Show>
          <div class={"wy-status " + st.cls}>{st.label}</div>
          <Show when={tagList(h).length}>
            <div class="wy-tags"><For each={tagList(h)}>{(t) => <span>{t}</span>}</For></div>
          </Show>
        </div>
      </div>
      <div class="addr wy-mono">
        <a class="wy-link" href={"http://" + h.IP} target="_blank">{h.IP}</a>
        <div class="wy-muted">{h.Mac}</div>
      </div>
      <div class="hw">
        <div>{h.Model || h.Hw || "Unknown vendor"}</div>
        <div class="wy-muted" style={{ "font-size": "13px" }}>{h.Iface} · seen {h.Date}</div>
      </div>
      <div class="acts">
        <label class="wy-inline-label" title="Known devices don't trigger alerts">
          <input type="checkbox" role="switch" class="wy-switch" checked={h.Known == 1}
            onChange={handleToggle}/>Known
        </label>
        <button type="button" class="wy-btn icon" onClick={handleWOL}
          aria-label={"Wake " + (name() || h.IP)} title={woke() ? "Magic packet sent" : "Wake-on-LAN: send a magic packet to power this device on"}>
          <i class={woke() ? "bi bi-check-lg" : "bi bi-lightning-charge"} aria-hidden="true"></i>
        </button>
      </div>
    </div>
  )
}

export default TableRow
