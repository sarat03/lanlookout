import { createSignal, For, Show } from "solid-js";
import { apiPortScan } from "../../functions/api";

// Ports worth checking first: ssh, web, file shares, printers, media, home automation.
const commonPorts = [21, 22, 23, 53, 80, 139, 443, 445, 554, 631, 1883, 3000, 3389, 5000, 5353, 5900, 8000, 8080, 8123, 8443, 8840, 9000, 9100, 32400];

function Ping(_props: any) {

  let stop = false;

  const [beginStr, setBegin] = createSignal("");
  const [endStr, setEnd] = createSignal("");
  const [curPort, setCurPort] = createSignal("");
  const [running, setRunning] = createSignal(false);
  const [foundPorts, setFoundPorts] = createSignal<number[]>([]);

  const handleScan = async () => {
    stop = false;
    setRunning(true);

    let begin = Number(beginStr());
    if (Number.isNaN(begin) || begin < 1 || begin > 65535) {
      begin = 1;
    }
    let end = Number(endStr());
    if (Number.isNaN(end) || end < 1 || end > 65535) {
      end = 65535;
    }

    let portOpened:boolean;
    for (let i = begin ; i <= end; i++) {

      if (stop) {
          break;
      }
      setCurPort(i.toString());
      portOpened = await apiPortScan(_props.IP, i);
      if (portOpened && !foundPorts().includes(i)) {
        setFoundPorts([...foundPorts(), i]);
      }
    }
    setRunning(false);
  };

  const handleQuick = async () => {
    setRunning(true);
    setCurPort("common");
    const open = await Promise.all(commonPorts.map(p => apiPortScan(_props.IP, p).then(ok => ok ? p : 0)));
    setFoundPorts([...new Set([...foundPorts(), ...open.filter(Boolean)])].sort((a, b) => a - b));
    setCurPort("");
    setRunning(false);
  };

  const handleStop = () => {
    if (stop) {
      setBegin(curPort());
      handleScan();
    } else {
      stop = true;
    }
  }

  return (
    <section class="wy-card" aria-label="Port scan">
      <h2>Open ports</h2>
      <form class="wy-row-wrap" onSubmit={e => { e.preventDefault(); handleScan(); }}>
        <label class="wy-label" style={{ flex: "1 1 120px" }}>From port
          <input type="number" min="1" max="65535" class="wy-input wy-mono" placeholder="1"
            value={beginStr()} onInput={e => setBegin(e.currentTarget.value)}></input></label>
        <label class="wy-label" style={{ flex: "1 1 120px" }}>To port
          <input type="number" min="1" max="65535" class="wy-input wy-mono" placeholder="65535"
            value={endStr()} onInput={e => setEnd(e.currentTarget.value)}></input></label>
        <button type="button" class="wy-btn primary" style={{ "align-self": "flex-end" }} disabled={running()} onClick={handleQuick}
          title={"Checks " + commonPorts.length + " common ports at once"}>Quick scan</button>
        <button type="submit" class="wy-btn" style={{ "align-self": "flex-end" }} disabled={running()}>Scan range</button>
        <Show when={curPort() != "" && curPort() != "common"}>
          <button type="button" onClick={handleStop} class="wy-btn" style={{ "align-self": "flex-end" }}>
            {running() ? "Pause" : "Continue"}
          </button>
        </Show>
      </form>
      <div role="status" class="wy-muted" style={{ "font-size": "14px" }}>
        {curPort() == "common" ? "Checking common ports…"
          : curPort() == "" ? (foundPorts().length ? "Found " + foundPorts().length + " open." : "Checks which TCP ports answer on " + (_props.IP || "this host") + ".")
          : (running() ? "Checking port " + curPort() + "…" : "Stopped at port " + curPort() + ".")
            + " Found " + foundPorts().length + " open."}
      </div>
      <Show when={foundPorts().length}>
        <div class="wy-row-wrap" style={{ gap: "8px" }}>
          <For each={foundPorts()}>{(port) =>
            <a class="wy-btn wy-chip wy-mono" href={"http://" + _props.IP + ":" + port} target="_blank">{port}</a>
          }</For>
        </div>
      </Show>
    </section>
  )
}

export default Ping
