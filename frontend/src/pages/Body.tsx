import { createMemo, createSignal, For, onMount, Show } from "solid-js";

import { bkpHosts, editNames, Host, ifaces, selectedIDs, setEditNames, setSelectedIDs, tagList } from "../functions/exports";
import { apiAddHost, apiDelHost, apiKnownAll, apiRescan, apiScanStatus } from "../functions/api";
import { getHosts } from "../functions/atstart";

import TableRow from "../components/Body/TableRow";

type Status = "all" | "on" | "off" | "new";

const read = (k: string, d: string) => localStorage.getItem(k) ?? d;

function numIP(ip: string) {
  return ip.split(".").reduce((a, b) => a * 256 + Number(b), 0);
}

// ponytail: CSV built in the browser from the host list already loaded
function exportCSV(hosts: Host[]) {
  const cols: (keyof Host)[] = ["Name", "IP", "Mac", "Hw", "Iface", "Type", "Tags", "Known", "Watch", "Now", "Date", "DNS"];
  const esc = (v: unknown) => '"' + String(v ?? "").replace(/"/g, '""') + '"';
  const csv = [cols.join(","), ...hosts.map(h => cols.map(c => esc(h[c])).join(","))].join("\n");
  const a = document.createElement("a");
  a.href = URL.createObjectURL(new Blob([csv], { type: "text/csv" }));
  a.download = "lanlookout-hosts.csv";
  a.click();
  URL.revokeObjectURL(a.href);
}

function Body() {

  const [status, setStatusRaw] = createSignal<Status>(read("wyStatus", "all") as Status);
  const [iface, setIfaceRaw] = createSignal(read("wyIface", ""));
  const [tag, setTagRaw] = createSignal(read("wyTag", ""));
  const [confirmKnown, setConfirmKnown] = createSignal(false);
  const [sort, setSortRaw] = createSignal(read("wySort", "Name"));
  const [desc, setDesc] = createSignal(read("wyDesc", "") === "1");
  const [query, setQuery] = createSignal("");
  const [adding, setAdding] = createSignal(false);
  const [scanning, setScanning] = createSignal(false);
  const [confirmDel, setConfirmDel] = createSignal(false);

  // ponytail: filter prefs live in localStorage, same as the old filter/sort did
  const setStatus = (s: Status) => { setStatusRaw(s); localStorage.setItem("wyStatus", s); };
  const setIface = (s: string) => { setIfaceRaw(s); localStorage.setItem("wyIface", s); };
  const setTag = (s: string) => { setTagRaw(s); localStorage.setItem("wyTag", s); };
  const setSort = (s: string) => { setSortRaw(s); localStorage.setItem("wySort", s); };
  const flipDesc = () => { setDesc(!desc()); localStorage.setItem("wyDesc", desc() ? "1" : ""); };

  onMount(() => {
    getHosts();
  });

  const counts = createMemo(() => {
    const all = bkpHosts();
    const on = all.filter(h => h.Now == 1).length;
    return { all: all.length, on, off: all.length - on, new: all.filter(h => h.Known != 1).length };
  });

  const allTags = createMemo(() => [...new Set(bkpHosts().flatMap(tagList))].sort());

  const visible = createMemo(() => {
    const q = query().trim().toLowerCase();
    const s = status();
    const list = bkpHosts().filter(h =>
      (s === "all" || (s === "on" && h.Now == 1) || (s === "off" && h.Now != 1) || (s === "new" && h.Known != 1))
      && (!iface() || h.Iface === iface())
      && (!allTags().includes(tag()) || tagList(h).includes(tag())) // stale saved tag = no filter
      && (!q || [h.Name, h.IP, h.Mac, h.Hw, h.Iface, h.DNS, h.Type, h.Tags, h.Owner, h.Model, h.OS, h.Services].join(" ").toLowerCase().includes(q))
    );
    const f = sort() as keyof Host;
    const dir = desc() ? -1 : 1;
    return list.slice().sort((a, b) => dir * (f === "IP"
      ? numIP(a.IP) - numIP(b.IP)
      : String(a[f]).localeCompare(String(b[f]), undefined, { numeric: true })));
  });

  const handleKnownAll = async () => {
    await apiKnownAll();
    setConfirmKnown(false);
    getHosts();
  };

  const tiles: { k: Status, label: string, hint: string, dot: string }[] = [
    { k: "all", label: "All hosts", hint: "Seen on your network", dot: "var(--wy-accent)" },
    { k: "on", label: "Online", hint: "Responding right now", dot: "var(--wy-on)" },
    { k: "off", label: "Offline", hint: "Not seen in last scan", dot: "var(--wy-off)" },
    { k: "new", label: "Unknown", hint: "Review and mark as known", dot: "var(--wy-warn)" },
  ];

  const handleRescan = async () => {
    setScanning(true);
    const before = await apiScanStatus();
    await apiRescan();
    // Poll until the scan finishes (arp-scan can take a while on big networks).
    for (let i = 0; i < 120; i++) {
      await new Promise(r => setTimeout(r, 1000));
      if (await apiScanStatus() !== before) break;
    }
    getHosts();
    setScanning(false);
  };

  let macRef!: HTMLInputElement, nameRef!: HTMLInputElement, ipRef!: HTMLInputElement;
  const handleAdd = async (e: SubmitEvent) => {
    e.preventDefault();
    await apiAddHost(macRef.value.trim(), nameRef.value.trim(), ipRef.value.trim());
    (e.target as HTMLFormElement).reset();
    setAdding(false);
    getHosts();
  };

  const toggleEdit = () => {
    if (editNames()) {
      getHosts();
    }
    setSelectedIDs([]);
    setConfirmDel(false);
    setEditNames(!editNames());
  };

  const handleDel = async () => {
    for (let id of selectedIDs()) {
      await apiDelHost(id);
    }
    window.location.href = '/';
  };

  return (
    <>
    <header class="wy-head">
      <div>
        <h1>Hosts</h1>
        <div class="wy-sub">Refreshes every minute</div>
      </div>
      <button type="button" class="wy-btn" onClick={() => exportCSV(visible())} title="Download the hosts shown below as CSV">
        <i class="bi bi-download" aria-hidden="true"></i>Export
      </button>
      <button type="button" class="wy-btn" onClick={() => setAdding(!adding())} aria-expanded={adding()}>
        <i class="bi bi-plus-lg" aria-hidden="true"></i>Add host
      </button>
      <button type="button" class="wy-btn primary" onClick={handleRescan} disabled={scanning()}>
        <i class="bi bi-arrow-clockwise" aria-hidden="true"></i>{scanning() ? "Scanning…" : "Scan now"}
      </button>
    </header>

    <section class="wy-stats" aria-label="Summary">
      <For each={tiles}>{(t) =>
        <button type="button" class="wy-stat" aria-pressed={status() === t.k} onClick={() => setStatus(t.k)}>
          <span class="k"><span class="wy-dot" style={{ background: t.dot, "box-shadow": "none", width: "10px", height: "10px" }}></span>{t.label}</span>
          <span class="n">{counts()[t.k]}</span>
          <span class="h">{t.k === "new" && counts().new === 0 ? "All devices reviewed" : t.hint}</span>
        </button>
      }</For>
    </section>

    <Show when={counts().new > 0}>
      <div class="wy-bar">
        <span style={{ "flex-grow": 1 }}>
          {counts().new} device{counts().new == 1 ? "" : "s"} not yet reviewed.
        </span>
        <Show when={confirmKnown()} fallback={
          <button type="button" class="wy-btn" onClick={() => setConfirmKnown(true)}>Mark all as known</button>
        }>
          <span>Mark all {counts().new} as known? They won't trigger new-device alerts.</span>
          <button type="button" class="wy-btn" onClick={() => setConfirmKnown(false)}>Cancel</button>
          <button type="button" class="wy-btn primary" onClick={handleKnownAll}>Mark known</button>
        </Show>
      </div>
    </Show>

    <Show when={adding()}>
      <form class="wy-card" style={{ "flex-direction": "row", "flex-wrap": "wrap", "align-items": "flex-end", gap: "12px" }} onSubmit={handleAdd} aria-label="Add host">
        <label class="wy-label" style={{ flex: "1 1 180px" }}>MAC address (required)
          <input ref={macRef} class="wy-input wy-mono" required pattern="([0-9A-Fa-f]{2}[:\-]){5}[0-9A-Fa-f]{2}" placeholder="aa:bb:cc:dd:ee:ff"/></label>
        <label class="wy-label" style={{ flex: "1 1 180px" }}>Name
          <input ref={nameRef} class="wy-input" placeholder="Living room TV"/></label>
        <label class="wy-label" style={{ flex: "1 1 160px" }}>IP address
          <input ref={ipRef} class="wy-input wy-mono" placeholder="192.168.1.50"/></label>
        <button type="button" class="wy-btn" onClick={() => setAdding(false)}>Cancel</button>
        <button type="submit" class="wy-btn primary">Save host</button>
      </form>
    </Show>

    <section class="wy-row-wrap" aria-label="Filters">
      <label class="wy-search">
        <i class="bi bi-search wy-muted" aria-hidden="true"></i>
        <input type="search" aria-label="Search hosts" placeholder="Search name, IP, MAC or vendor"
          value={query()} onInput={e => setQuery(e.currentTarget.value)}/>
      </label>
      <Show when={ifaces().length > 1}>
        <label class="wy-inline-label">Interface
          <select class="wy-input" value={iface()} onChange={e => setIface(e.currentTarget.value)}>
            <option value="">All</option>
            <For each={ifaces()}>{(i) => <option value={i}>{i}</option>}</For>
          </select>
        </label>
      </Show>
      <Show when={allTags().length}>
        <label class="wy-inline-label">Tag
          <select class="wy-input" value={tag()} onChange={e => setTag(e.currentTarget.value)}>
            <option value="">All</option>
            <For each={allTags()}>{(t) => <option value={t}>{t}</option>}</For>
          </select>
        </label>
      </Show>
      <label class="wy-inline-label">Sort
        <select class="wy-input" value={sort()} onChange={e => setSort(e.currentTarget.value)}>
          <option value="Name">Name</option>
          <option value="IP">IP address</option>
          <option value="Date">Last seen</option>
          <option value="Mac">MAC</option>
          <option value="Hw">Vendor</option>
          <option value="Iface">Interface</option>
          <option value="Type">Type</option>
        </select>
      </label>
      <button type="button" class="wy-btn icon" onClick={flipDesc} aria-label={desc() ? "Sort descending" : "Sort ascending"} title="Reverse order">
        <i class={desc() ? "bi bi-sort-down" : "bi bi-sort-up"} aria-hidden="true"></i>
      </button>
      <button type="button" class="wy-btn" onClick={toggleEdit} aria-pressed={editNames()}>
        {editNames() ? "Done" : "Edit hosts"}
      </button>
    </section>

    <Show when={editNames()}>
      <div class="wy-bar">
        <span style={{ "flex-grow": 1, color: "var(--wy-accent)", "font-weight": 500 }}>
          {selectedIDs().length ? `${selectedIDs().length} selected` : "Rename hosts inline, or tick them to delete"}
        </span>
        <Show when={confirmDel()} fallback={
          <button type="button" class="wy-btn danger" disabled={!selectedIDs().length} onClick={() => setConfirmDel(true)}>Delete selected</button>
        }>
          <span>Delete {selectedIDs().length} host(s) and their history?</span>
          <button type="button" class="wy-btn" onClick={() => setConfirmDel(false)}>Keep</button>
          <button type="button" class="wy-btn danger solid" onClick={handleDel}>Delete</button>
        </Show>
      </div>
    </Show>

    <section class="wy-list" aria-label="Host list">
      <For each={visible()} fallback={
        <div class="wy-empty">{bkpHosts().length ? "No hosts match. Clear the search or pick All." : "No hosts yet. The first scan will fill this list."}</div>
      }>{(host) =>
        <TableRow host={host}></TableRow>
      }</For>
    </section>
    </>
  )
}

export default Body
