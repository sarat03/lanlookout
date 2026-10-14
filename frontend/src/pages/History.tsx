import { createEffect, createSignal, For, Show } from "solid-js"
import Filter from "../components/Filter"
import { allHosts, histUpdOnFilter, Host, setHistUpdOnFilter, setShow, show } from "../functions/exports"
import MacHistory from "../components/MacHistory"
import HistShow from "../components/HistShow"
import Events from "../components/Events"
import { hostStatus } from "../components/Body/TableRow"

function History() {

  let hosts: Host[] = [];
  hosts.push(...allHosts);

  const showStr = localStorage.getItem("histShow") as string;
  setShow(+showStr);
  (show() === 0 || isNaN(show())) ? setShow(200) : '';

  let saved = "hosts";
  try { saved = localStorage.getItem("histTab") || "hosts"; } catch {}
  const [tab, setTab] = createSignal(saved);
  const pick = (t: string) => { setTab(t); try { localStorage.setItem("histTab", t); } catch {} };

  createEffect(() => {
    if (histUpdOnFilter()) {
      hosts = [];
      hosts.push(...allHosts);
      setHistUpdOnFilter(false);
    }
  });

  return (
    <>
    <header class="wy-head">
      <div>
        <h1>History</h1>
        <div class="wy-sub">Each bar is one scan, newest on the left. Hover a bar for details.</div>
      </div>
    </header>

    <div class="wy-seg" role="tablist" aria-label="History view">
      <For each={[["hosts", "Hosts"], ["events", "Recent events"]]}>{([id, label]) =>
        <button type="button" role="tab" aria-selected={tab() === id} onClick={() => pick(id)}>{label}</button>
      }</For>
    </div>

    <Show when={tab() === "hosts"} fallback={<Events></Events>}>
    <section class="wy-row-wrap" aria-label="Filters">
      <Filter></Filter>
      <HistShow name="histShow"></HistShow>
    </section>

    <section class="wy-list" aria-label="History by host">
      <Show when={!histUpdOnFilter()}>
        <For each={hosts} fallback={<div class="wy-empty">No hosts to show.</div>}>{(host) =>
          <div class="wy-host">
            <div class="who" style={{ flex: "1 1 200px" }}>
              <span class={"wy-dot " + hostStatus(host).cls} title={hostStatus(host).label}></span>
              <div style={{ "min-width": 0 }}>
                <a class="name" href={"/host/" + host.ID}>{host.Name || "Unnamed device"}</a>
                <div class="wy-mono wy-muted" style={{ "font-size": "13px" }}>{host.IP}</div>
              </div>
            </div>
            <div class="wy-hist" style={{ flex: "4 1 300px", "min-width": 0 }}>
              <MacHistory mac={host.Mac} date=""></MacHistory>
            </div>
          </div>
        }</For>
      </Show>
    </section>
    </Show>
    </>
  )
}

export default History
