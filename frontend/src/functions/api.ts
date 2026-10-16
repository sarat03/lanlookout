export const apiPath = 'http://0.0.0.0:8840';

export const apiGetAllHosts = async () => {
  const url = apiPath+'/api/all';
  const hosts = await (await fetch(url)).json();

  return hosts;
};

export const apiGetIfaces = async (): Promise<string[]> => {
  return await (await fetch(apiPath+'/api/ifaces')).json();
};

export const apiGetConfig = async () => {

  const url = apiPath+'/api/config';
  const res = await (await fetch(url)).json();

  return res;
};

export const apiGetVersion = async () => {

  const url = apiPath+'/api/version';
  const res = await (await fetch(url)).json();

  return res;
};

export const apiTestNotify = async () => {

  const url = apiPath+'/api/notify_test';
  await fetch(url);
};

export const apiEditHost = async (id:number, name:string, known:string) => {

  const url = apiPath+'/api/edit/'+id+'/'+name+'/'+known;
  const res = await (await fetch(url)).json();

  return res;
};

export const apiGetHost = async (id:string) => {

  const url = apiPath+'/api/host/'+id;
  const res = await (await fetch(url)).json();

  return res;
};

export const apiDelHost = async (id:number) => {

  const url = apiPath+'/api/host/del/'+id;
  const res = await (await fetch(url)).json();

  return res;
};

export const apiPortScan = async (ip:string, port:number) => {

  const url = apiPath+'/api/port/'+ip+'/'+port;
  const res = await (await fetch(url)).json();

  return res;
};

export const apiGetHistory = async (mac:string) => {
  const url = apiPath+'/api/history/'+mac+'/?num=210';
  const hosts = await (await fetch(url)).json();

  return hosts;
};

export const apiGetHistoryByDate = async (mac:string, date: string) => {
  const url = apiPath+'/api/history/'+mac+'/'+date;
  const hosts = await (await fetch(url)).json();

  return hosts;
};

export const apiWOL = async (mac:string) => {

  const url = apiPath+'/api/wol/'+mac;
  const res = await (await fetch(url)).json();

  return res;
};
export const apiAddHost = async (mac:string, name:string, ip:string) => {

  const q = new URLSearchParams({ name, ip });
  const url = apiPath+'/api/host/add/'+encodeURIComponent(mac)+'?'+q;
  const res = await (await fetch(url)).json();

  return res;
};

export const apiRescan = async () => {

  const url = apiPath+'/api/rescan';
  await fetch(url);
};

export const apiScanStatus = async () => {

  const url = apiPath+'/api/scan_status';
  const res = await (await fetch(url)).json();

  return res.LastScan as string;
};

export const apiGetEvents = async (num: number) => {

  const url = apiPath+'/api/events?num='+num;
  const res = await (await fetch(url)).json();

  return res;
};

export const apiHostMeta = async (id:number, meta: { watch?: number, type?: string, tags?: string, owner?: string }) => {

  const q = new URLSearchParams(Object.entries(meta).map(([k, v]) => [k, String(v)]));
  const url = apiPath+'/api/host/meta/'+id+'?'+q;
  const res = await (await fetch(url)).json();

  return res;
};

export const apiKnownAll = async () => {

  const url = apiPath+'/api/known/all';
  const res = await (await fetch(url)).json();

  return res;
};

const getJSON = async (path: string) => (await fetch(apiPath + path)).json();

export const apiNetStatus = () => getJSON('/api/net/status');
export const apiNetStats = (hours: number) => getJSON('/api/net/stats?hours=' + hours);
export const apiSpeedTest = async () => {
  const res = await fetch(apiPath + '/api/net/speedtest');
  const body = await res.json();
  if (!res.ok) throw new Error(body);
  return body;
};
export const apiSecurity = () => getJSON('/api/security');
export const apiSecurityScan = () => getJSON('/api/security/scan');
export const apiDNSTop = (ip: string, days: number) => getJSON('/api/dns/' + encodeURIComponent(ip) + '?days=' + days);
