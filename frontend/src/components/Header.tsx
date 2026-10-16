import { createSignal, For } from "solid-js";
import { appConfig, setAppConfig, themes } from "../functions/exports";
import { apiGetConfig } from "../functions/api";
import { useLocation } from "@solidjs/router";

const links = [
  { href: "/", label: "Hosts", icon: "bi-hdd-stack" },
  { href: "/home/", label: "Home", icon: "bi-house-heart" },
  { href: "/history/", label: "History", icon: "bi-activity" },
  { href: "/internet/", label: "Internet", icon: "bi-globe2" },
  { href: "/security/", label: "Security", icon: "bi-shield-check" },
  { href: "/config/", label: "Settings", icon: "bi-gear" },
];

function Header() {

  const [themePath, setThemePath] = createSignal('');
  const [iconsPath, setIconsPath] = createSignal('');

  const setCurrentTheme = async () => {
    setAppConfig(await apiGetConfig());

    // Unknown theme (e.g. one from an older config) falls back to the default.
    const theme = themes.includes(appConfig().Theme) ? appConfig().Theme : "cosmo";
    const color = appConfig().Color?appConfig().Color:"dark";

    if (appConfig().NodePath == '') {
      setThemePath("https://cdn.jsdelivr.net/npm/bootswatch@5.3.3/dist/"+theme+"/bootstrap.min.css");
      setIconsPath("https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.css");
    } else {
      setThemePath(appConfig().NodePath+"/node_modules/bootswatch/dist/"+theme+"/bootstrap.min.css");
      setIconsPath(appConfig().NodePath+"/node_modules/bootstrap-icons/font/bootstrap-icons.css");
    }

    document.documentElement.setAttribute("data-bs-theme", color);
  }
  setCurrentTheme();

  // Host pages belong to Hosts; everything else matches by first path segment.
  const loc = useLocation();
  const current = (href: string) => {
    const path = loc.pathname;
    return href === "/" ? (path === "/" || path.startsWith("/host/")) : path.startsWith(href.slice(0, -1));
  };

  return (
    <>
    <link rel="stylesheet" href={iconsPath()}></link> {/* icons */}
    <link rel="stylesheet" href={themePath()}></link> {/* theme */}
    <nav class="wy-nav" aria-label="Main">
      <a class="wy-brand" href="/">
        <img src="/fs/public/favicon.svg" alt=""/>
        LanLookout
      </a>
      <For each={links}>{(l) =>
        <a class="wy-nav-link" href={l.href} aria-current={current(l.href) ? "page" : undefined}>
          <i class={"bi " + l.icon} aria-hidden="true"></i>{l.label}
        </a>
      }</For>
      <div class="wy-nav-foot">
        <div>Scans every {appConfig().Timeout} s</div>
        <div>on {appConfig().Ifaces || "—"}</div>
        <a class="wy-link" href="https://github.com/sarat03/lanlookout" target="_blank">
          <i class="bi bi-github" aria-hidden="true"></i> GitHub
        </a>
      </div>
    </nav>
    </>
  )
};

export default Header
