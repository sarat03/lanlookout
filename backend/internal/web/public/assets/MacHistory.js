import{P as y,X as f,Y as v,o as m,Z as H,g as w,w as d,t as D,b as I,s as b,e as g,D as h,F as M}from"./index.js";async function u(a,r){let e=[];return r===""?e=await y(a):e=await f(a,r),e!=null?(e.sort((s,t)=>s.Date<t.Date?1:-1),e):[]}var x=D("<i>");function P(a){const[r,e]=v([]);let s;return m(async()=>{const t=await u(a.mac,a.date);e(t),s=setInterval(async()=>{const o=await u(a.mac,a.date);e(o)},6e4)}),H(()=>{clearInterval(s)}),w(M,{each:r,children:(t,o)=>w(d,{get when(){return o()<h()},get children(){var i=x();return I(n=>{var c="Date:"+t.Date+`
Iface:`+t.Iface+`
IP:`+t.IP+`
Known:`+t.Known,l=t.Now===0?"my-box-off":"my-box-on";return c!==n.e&&b(i,"title",n.e=c),l!==n.t&&g(i,n.t=l),n},{e:void 0,t:void 0}),i}})})}export{P as M};
