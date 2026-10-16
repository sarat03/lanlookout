import { useParams } from "@solidjs/router";
import { createSignal, onMount, Show } from "solid-js";

import { apiDelHost, apiGetHost, apiWOL } from "../functions/api";

import HostCard from "../components/HostPage/HostCard";
import Ping from "../components/HostPage/Ping";
import HistCard from "../components/HostPage/HistCard";
import DNSCard from "../components/HostPage/DNSCard";
import { hostStatus } from "../components/Body/TableRow";
import { emptyHost, Host } from "../functions/exports";

function HostPage() {

  const [currentHost, setCurrentHost] = createSignal<Host>(emptyHost);
  const [woke, setWoke] = createSignal(false);
  const [confirm, setConfirm] = createSignal(false);

  onMount(async () => {
    const params = useParams();
    const host = await apiGetHost(params.id);

    setCurrentHost(host);
  });

  const st = () => hostStatus(currentHost());
  const title = () => currentHost().Name || currentHost().IP || "Host";

  const handleWOL = async () => {
    await apiWOL(currentHost().Mac);
    setWoke(true);
    setTimeout(() => setWoke(false), 2000);
  };

  const handleDel = async () => {
    await apiDelHost(currentHost().ID);
    window.location.href = '/';
  };

  return (
    <>
    <a class="wy-link" href="/hosts/" style={{ display: "inline-flex", "align-items": "center", gap: "6px", "min-height": "44px", "font-weight": 500 }}>
      <i class="bi bi-chevron-left" aria-hidden="true"></i>All hosts
    </a>

    <header class="wy-head">
      <div style={{ display: "flex", "align-items": "center", gap: "14px" }}>
        <span class={"wy-dot " + st().cls} style={{ width: "16px", height: "16px" }}></span>
        <div>
          <h1>{title()}</h1>
          <div class={"wy-status " + st().cls} style={{ "font-size": "14px" }}>{st().label} · last seen {currentHost().Date}</div>
        </div>
      </div>
      <button type="button" class="wy-btn" onClick={handleWOL}
        title="Sends a magic packet that powers on a sleeping or shut-down device. Needs Wake-on-LAN enabled in its BIOS/network settings and a wired connection on the same network.">
        <i class={woke() ? "bi bi-check-lg" : "bi bi-lightning-charge"} aria-hidden="true"></i>{woke() ? "Magic packet sent" : "Wake-on-LAN"}
      </button>
      <button type="button" class="wy-btn danger" onClick={() => setConfirm(true)}>Delete host</button>
    </header>

    <p class="wy-muted" style={{ margin: "0 0 16px", "font-size": "13px" }}>
      <i class="bi bi-lightning-charge" aria-hidden="true"></i> <b>Wake-on-LAN</b> sends a "magic packet" that turns this device on remotely.
      It only works if Wake-on-LAN is enabled in the device's BIOS/network settings, and usually only over Ethernet on the same network.
    </p>

    <Show when={confirm()}>
      <div class="wy-alert" role="alert">
        <span style={{ "flex-grow": 1 }}>Delete {title()} and its history? This cannot be undone.</span>
        <button type="button" class="wy-btn" onClick={() => setConfirm(false)}>Keep it</button>
        <button type="button" class="wy-btn danger solid" onClick={handleDel}>Delete</button>
      </div>
    </Show>

    <div class="wy-cols">
      <HostCard host={currentHost()}></HostCard>
      <div class="wide" style={{ display: "flex", "flex-direction": "column", gap: "20px", "min-width": 0 }}>
        <Ping IP={currentHost().IP}></Ping>
        <HistCard mac={currentHost().Mac}></HistCard>
        <DNSCard ip={currentHost().IP}></DNSCard>
      </div>
    </div>
    </>
  )
}

export default HostPage
