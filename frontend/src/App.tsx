import { lazy, onMount } from 'solid-js';
import { Router, Route } from "@solidjs/router";
import './App.css';
import { runAtStart } from './functions/atstart';

import Body from './pages/Body';
import Header from './components/Header';

function App() {

  onMount(() => {
    runAtStart();
  });

  const Config = lazy(() => import("./pages/Config"));
  const Home = lazy(() => import("./pages/Home"));
  const History = lazy(() => import("./pages/History"));
  const HostPage = lazy(() => import("./pages/HostPage"));
  const Internet = lazy(() => import("./pages/Internet"));
  const Security = lazy(() => import("./pages/Security"));

  return (
    <div class="wy-app">
      <Router root={(props) => <><Header></Header><main class="wy-main">{props.children}</main></>}>
          <Route path="/" component={Home}/>
          <Route path="/hosts" component={Body}/>
          <Route path="/config" component={Config}/>
          <Route path="/home" component={Home}/>
          <Route path="/history" component={History}/>
          <Route path="/host/:id" component={HostPage}/>
          <Route path="/internet" component={Internet}/>
          <Route path="/security" component={Security}/>
      </Router>
    </div>
  )
}

export default App
