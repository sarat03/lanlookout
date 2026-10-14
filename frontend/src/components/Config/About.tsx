import { createSignal, onMount } from "solid-js";
import { apiGetVersion } from "../../functions/api"

function About() {

  const [version, setVersion] = createSignal('');
  const [link, setLink] = createSignal('');

  onMount(async () => {
    const v = await apiGetVersion();
    setVersion(v);
    setLink("https://github.com/sarat03/lanlookout/releases/tag/"+v);
  });

  return (
    <section class="wy-card" aria-label="About">
      <div>
        <h2>About</h2>
        <div class="wy-sub">LanLookout <a class="wy-link" href={link()} target="_blank">{version()}</a>, open source.</div>
      </div>
      <dl class="wy-dl" style={{ "grid-template-columns": "150px minmax(0, 1fr)" }}>
        <dt>Settings guide</dt>
        <dd><a class="wy-link" href="https://github.com/sarat03/lanlookout#config" target="_blank">What each setting does and how to set it up</a></dd>
        <dt>Notifications</dt>
        <dd><a class="wy-link" href="https://shoutrrr.nickfedor.com/services/overview/" target="_blank">Shoutrrr URL formats (Discord, Telegram, Email…)</a></dd>
        <dt>API docs</dt>
        <dd><a class="wy-link wy-mono" href="/swagger/index.html" target="_blank">/swagger/index.html</a></dd>
        <dt>Source code</dt>
        <dd><a class="wy-link" href="https://github.com/sarat03/lanlookout" target="_blank">github.com/sarat03/lanlookout</a></dd>
        <dt>VLAN scanning</dt>
        <dd><a class="wy-link" href="https://github.com/sarat03/lanlookout/blob/main/docs/VLAN_ARP_SCAN.md" target="_blank">Arp strings and arp-scan arguments</a></dd>
      </dl>
    </section>
  )
}

export default About
