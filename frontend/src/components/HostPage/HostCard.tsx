import { For } from "solid-js";
import { apiEditHost, apiHostMeta } from "../../functions/api";
import { hostTypes } from "../../functions/exports";

import { debounce } from "@solid-primitives/scheduled";

function HostCard(_props: any) {

  let name:string = "";

  const debouncedTags = debounce(async (val: string) => {
      await apiHostMeta(_props.host.ID, { tags: val });
    }, 500);

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
        <label class="wy-label" style={{ flex: "2 1 180px" }}>Tags
          <input type="text" class="wy-input" placeholder="family, work" value={_props.host.Tags || ""}
            onInput={e => debouncedTags(e.currentTarget.value)}/>
        </label>
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
