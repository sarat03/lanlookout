import { createEffect, createSignal, For, Show } from "solid-js";
import { apiEditHost, apiHostMeta } from "../../functions/api";
import { bkpHosts, hostTypes, tagList } from "../../functions/exports";

import { debounce } from "@solid-primitives/scheduled";

function HostCard(_props: any) {

  let name:string = "";

  // Choices come from what other devices already use.
  const owners = () => [...new Set([...bkpHosts().map(h => h.Owner), owner()].filter(Boolean))].sort();
  const allTags = () => [...new Set([...bkpHosts().flatMap(tagList), ...tags()])].sort();

  const [owner, setOwner] = createSignal<string>(_props.host.Owner || "");
  const [adding, setAdding] = createSignal(false);
  const saveOwner = (o: string) => {
    setOwner(o.trim());
    setAdding(false);
    apiHostMeta(_props.host.ID, { owner: o.trim() });
  };

  const [tags, setTags] = createSignal<string[]>(tagList(_props.host));

  // The host loads after the card mounts.
  createEffect(() => {
    setOwner(_props.host.Owner || "");
    setTags(tagList(_props.host));
  });
  const saveTags = (t: string[]) => {
    setTags(t);
    apiHostMeta(_props.host.ID, { tags: t.join(",") });
  };
  const toggleTag = (t: string) => saveTags(tags().includes(t) ? tags().filter(x => x !== t) : [...tags(), t]);
  const addTag = (el: HTMLInputElement) => {
    const t = el.value.trim();
    el.value = "";
    if (t && !tags().includes(t)) saveTags([...tags(), t]);
  };

  const debouncedApi = debounce(async (val: string) => {
      await apiEditHost(_props.host.ID, val, "");
    }, 300);

  const handleInput = async (n: string) => {

    name = n;
    debouncedApi(n);
  };

  const handleToggle = async () => {

    if (name == "") {
      name = _props.host.Name;
    }

    await apiEditHost(_props.host.ID, name, 'toggle');
  };

  return (
    <section class="wy-card" aria-label="Details">
      <h2>Details</h2>
      <label class="wy-label">Name
        <input type="text" class="wy-input" value={_props.host.Name}
          onInput={e => handleInput(e.currentTarget.value)}></input>
      </label>
      <div class="wy-row-wrap">
        <label class="wy-label" style={{ flex: "1 1 140px" }}>Type
          <select class="wy-input" value={_props.host.Type || ""}
            onChange={e => apiHostMeta(_props.host.ID, { type: e.currentTarget.value })}>
            <option value="">Not set</option>
            <For each={Object.keys(hostTypes)}>{(t) => <option value={t}>{t}</option>}</For>
          </select>
        </label>
        <label class="wy-label" style={{ flex: "1 1 140px" }}>Owner
          <Show when={!adding()} fallback={
            <input type="text" class="wy-input" placeholder="New member's name" autofocus
              onKeyDown={e => { if (e.key === "Enter") saveOwner(e.currentTarget.value); if (e.key === "Escape") setAdding(false); }}
              onBlur={e => e.currentTarget.value.trim() ? saveOwner(e.currentTarget.value) : setAdding(false)}/>
          }>
            <select class="wy-input" value={owner()}
              onChange={e => e.currentTarget.value === "+" ? setAdding(true) : saveOwner(e.currentTarget.value)}>
              <option value="" selected={!owner()}>Shared</option>
              <For each={owners()}>{(o) => <option value={o} selected={o === owner()}>{o}</option>}</For>
              <option value="+">+ Add member…</option>
            </select>
          </Show>
        </label>
      </div>
      <div class="wy-label">Tags
        <div class="wy-row-wrap" role="group" aria-label="Tags" style={{ gap: "8px" }}>
          <For each={allTags()}>{(t) =>
            <button type="button" class="wy-btn wy-chip" aria-pressed={tags().includes(t)} onClick={() => toggleTag(t)}>{t}</button>
          }</For>
          <input type="text" class="wy-input" placeholder="+ New tag" aria-label="New tag" style={{ width: "120px" }}
            onKeyDown={e => { if (e.key === "Enter") addTag(e.currentTarget); }}
            onBlur={e => addTag(e.currentTarget)}/>
        </div>
      </div>
      <dl class="wy-dl">
        <dt>IP address</dt>
        <dd class="wy-mono"><a class="wy-link" href={"http://" + _props.host.IP} target="_blank">{_props.host.IP}</a></dd>
        <dt>MAC address</dt><dd class="wy-mono">{_props.host.Mac}</dd>
        <dt>DNS name</dt><dd class="wy-mono">{_props.host.DNS || "—"}</dd>
        <dt>Vendor</dt><dd>{_props.host.Hw || "Unknown"}</dd>
        <dt>Model</dt><dd>{_props.host.Model || "—"}</dd>
        <dt>OS</dt><dd>{_props.host.OS || "—"}</dd>
        <dt>Services</dt><dd>{_props.host.Services ? _props.host.Services.split(",").join(", ") : "—"}</dd>
        <dt>Open ports</dt><dd class="wy-mono">{!_props.host.Ports ? "Not scanned" : _props.host.Ports === "none" ? "None risky" : <a class="wy-link" href="/security/">{_props.host.Ports.split(",").join(", ")}</a>}</dd>
        <dt>Interface</dt><dd>{_props.host.Iface}</dd>
        <dt>Last seen</dt><dd>{_props.host.Date}</dd>
        <dt>ID</dt><dd>{_props.host.ID}</dd>
      </dl>
      <div class="wy-actions">
        <div style={{ "flex-grow": 1 }}>
          <div style={{ "font-weight": 500 }}>Known device</div>
          <div class="wy-muted" style={{ "font-size": "13px" }}>Unknown devices trigger a notification</div>
        </div>
        <input type="checkbox" role="switch" class="wy-switch" aria-label="Known device"
          checked={_props.host.Known == 1} onChange={handleToggle}/>
      </div>
      <div class="wy-actions">
        <div style={{ "flex-grow": 1 }}>
          <div style={{ "font-weight": 500 }}>Watch this device</div>
          <div class="wy-muted" style={{ "font-size": "13px" }}>Notify when it goes offline or comes back online</div>
        </div>
        <input type="checkbox" role="switch" class="wy-switch" aria-label="Watch this device"
          checked={_props.host.Watch == 1}
          onChange={e => apiHostMeta(_props.host.ID, { watch: e.currentTarget.checked ? 1 : 0 })}/>
      </div>
    </section>
  )
}

export default HostCard
