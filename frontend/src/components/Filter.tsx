import { createSignal, For } from "solid-js";
import { Host, ifaces, setHistUpdOnFilter } from "../functions/exports";
import { filterFunc } from "../functions/filter";


function Filter() {
  type FilterEvent = Event & {
    currentTarget: HTMLSelectElement;
    target: HTMLSelectElement;
  };

  const [selectValue, setSelectValue] = createSignal("");

  const handleFilter = (field: keyof Host, event: FilterEvent) => {
    const value = event.target ? event.target.value : 0;
    filterFunc(field, value);
    setHistUpdOnFilter(true);
  };

  const handleReset = () => {
    filterFunc("ID", 0);
    setSelectValue("something");
    setSelectValue("");
    setHistUpdOnFilter(true);
  };

  return (
    <div class="wy-row-wrap" role="group" aria-label="Filter hosts">
      <select onChange={(event)=>{handleFilter("Iface", event)}} class="wy-input" aria-label="Filter by interface" value={selectValue()}>
        <option value="" selected disabled>Interface</option>
        <For each={ifaces()}>{(iface) =>
          <option value={iface}>{iface}</option>
        }</For>
      </select>
      <select onChange={(event)=>{handleFilter("Known", event)}} class="wy-input" aria-label="Filter by known" value={selectValue()}>
        <option value="" selected disabled>Known?</option>
        <option value={1}>Known</option>
        <option value={0}>Unknown</option>
      </select>
      <select onChange={(event)=>{handleFilter("Now", event)}} class="wy-input" aria-label="Filter by online" value={selectValue()}>
        <option value="" selected disabled>Status</option>
        <option value={1}>Online</option>
        <option value={0}>Offline</option>
      </select>
      <button type="button" onClick={handleReset} class="wy-btn">Reset filters</button>
    </div>
  )
}

export default Filter
