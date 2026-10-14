import { apiGetConfig } from "./api";
import { appConfig, setAppConfig } from "./exports";

// Settings forms post here instead of reloading the page. Submit events bubble,
// so one listener on the Settings page covers every form.
export async function saveInPlace(e: SubmitEvent) {

  const form = e.target as HTMLFormElement;
  if (form.method.toLowerCase() !== "post") return;
  e.preventDefault();

  const btn = form.querySelector<HTMLButtonElement>("button[type=submit]");
  const label = btn?.textContent ?? "";
  const show = (t: string) => { if (btn) btn.textContent = t; };

  const old = appConfig();
  show("Saving…");
  try {
    // The handlers answer with a redirect meant for plain form posts; don't follow it.
    const res = await fetch(form.action, { method: "POST", body: new FormData(form), redirect: "manual" });
    if (!res.ok && res.type !== "opaqueredirect") throw new Error(res.statusText);
    const cfg = await apiGetConfig();
    // Theme and colour are applied once at page load.
    if (cfg.Theme !== old.Theme || cfg.Color !== old.Color) {
      location.reload();
      return;
    }
    setAppConfig(cfg);
    show("Saved ✓");
  } catch {
    show("Save failed, try again");
  }
  setTimeout(() => show(label), 2000);
}
