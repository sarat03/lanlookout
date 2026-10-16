import{T as y,_ as f,$ as v,o as m,a0 as H,g as u,z as d,t as I,a as D,s as g,e as h,H as M,F as b}from"./index.js";async function w(a,r){let e=[];return r===""?e=await y(a):e=await f(a,r),e!=null?(e.sort((s,t)=>s.Date<t.Date?1:-1),e):[]}var x=I("<i>");function C(a){const[r,e]=v([]);let s;return m(async()=>{const t=await w(a.mac,a.date);e(t),s=setInterval(async()=>{const o=await w(a.mac,a.date);e(o)},6e4)}),H(()=>{clearInterval(s)}),u(b,{each:r,children:(t,o)=>u(d,{get when(){return o()<M()},get children(){var i=x();return D(n=>{var c="Date:"+t.Date+`
Iface:`+t.Iface+`
IP:`+t.IP+`
Known:`+t.Known,l=t.Now===0?"my-box-off":"my-box-on";return c!==n.e&&g(i,"title",n.e=c),l!==n.t&&h(i,n.t=l),n},{e:void 0,t:void 0}),i}})})}export{C as M};
