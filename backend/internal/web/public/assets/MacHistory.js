import{K as u,T as f,U as v,o as m,V as H,g as w,y as d,t as I,b as D,s as b,e as g,w as h,F as M}from"./index.js";async function y(a,r){let e=[];return r===""?e=await u(a):e=await f(a,r),e!=null?(e.sort((s,t)=>s.Date<t.Date?1:-1),e):[]}var x=I("<i>");function K(a){const[r,e]=v([]);let s;return m(async()=>{const t=await y(a.mac,a.date);e(t),s=setInterval(async()=>{const o=await y(a.mac,a.date);e(o)},6e4)}),H(()=>{clearInterval(s)}),w(M,{each:r,children:(t,o)=>w(d,{get when(){return o()<h()},get children(){var i=x();return D(n=>{var c="Date:"+t.Date+`
Iface:`+t.Iface+`
IP:`+t.IP+`
Known:`+t.Known,l=t.Now===0?"my-box-off":"my-box-on";return c!==n.e&&b(i,"title",n.e=c),l!==n.t&&g(i,n.t=l),n},{e:void 0,t:void 0}),i}})})}export{K as M};
