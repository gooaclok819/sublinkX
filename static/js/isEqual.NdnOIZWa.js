import{U as e,q as t,i as r,m as n,S as a,u as o}from"./_initCloneObject.C-kpAzhc.js";import{cJ as i,bP as c,c1 as u,bS as f,bQ as s}from"./index.CBDW5MZN.js";function v(e){var t=-1,r=null==e?0:e.length;for(this.__data__=new i;++t<r;)this.add(e[t])}function l(e,t){for(var r=-1,n=null==e?0:e.length;++r<n;)if(t(e[r],r,e))return!0;return!1}v.prototype.add=v.prototype.push=function(e){return this.__data__.set(e,"__lodash_hash_undefined__"),this},v.prototype.has=function(e){return this.__data__.has(e)};var b=1,h=2;function p(e,t,r,n,a,o){var i=r&b,c=e.length,u=t.length;if(c!=u&&!(i&&u>c))return!1;var f=o.get(e),s=o.get(t);if(f&&s)return f==t&&s==e;var p=-1,_=!0,d=r&h?new v:void 0;for(o.set(e,t),o.set(t,e);++p<c;){var y=e[p],g=t[p];if(n)var j=i?n(g,y,p,t,e,o):n(y,g,p,e,t,o);if(void 0!==j){if(j)continue;_=!1;break}if(d){if(!l(t,(function(e,t){if(i=t,!d.has(i)&&(y===e||a(y,e,r,n,o)))return d.push(t);var i}))){_=!1;break}}else if(y!==g&&!a(y,g,r,n,o)){_=!1;break}}return o.delete(e),o.delete(t),_}function _(e){var t=-1,r=Array(e.size);return e.forEach((function(e,n){r[++t]=[n,e]})),r}function d(e){var t=-1,r=Array(e.size);return e.forEach((function(e){r[++t]=e})),r}var y=1,g=2,j="[object Boolean]",w="[object Date]",m="[object Error]",O="[object Map]",A="[object Number]",S="[object RegExp]",k="[object Set]",z="[object String]",E="[object Symbol]",L="[object ArrayBuffer]",x="[object DataView]",P=c?c.prototype:void 0,B=P?P.valueOf:void 0;var D=1,N=Object.prototype.hasOwnProperty;var q=1,C="[object Arguments]",J="[object Array]",M="[object Object]",Q=Object.prototype.hasOwnProperty;function R(i,c,s,v,l,b){var h=f(i),P=f(c),R=h?J:r(i),U=P?J:r(c),V=(R=R==C?M:R)==M,F=(U=U==C?M:U)==M,G=R==U;if(G&&n(i)){if(!n(c))return!1;h=!0,V=!1}if(G&&!V)return b||(b=new a),h||o(i)?p(i,c,s,v,l,b):function(t,r,n,a,o,i,c){switch(n){case x:if(t.byteLength!=r.byteLength||t.byteOffset!=r.byteOffset)return!1;t=t.buffer,r=r.buffer;case L:return!(t.byteLength!=r.byteLength||!i(new e(t),new e(r)));case j:case w:case A:return u(+t,+r);case m:return t.name==r.name&&t.message==r.message;case S:case z:return t==r+"";case O:var f=_;case k:var s=a&y;if(f||(f=d),t.size!=r.size&&!s)return!1;var v=c.get(t);if(v)return v==r;a|=g,c.set(t,r);var l=p(f(t),f(r),a,o,i,c);return c.delete(t),l;case E:if(B)return B.call(t)==B.call(r)}return!1}(i,c,R,s,v,l,b);if(!(s&q)){var H=V&&Q.call(i,"__wrapped__"),I=F&&Q.call(c,"__wrapped__");if(H||I){var K=H?i.value():i,T=I?c.value():c;return b||(b=new a),l(K,T,s,v,b)}}return!!G&&(b||(b=new a),function(e,r,n,a,o,i){var c=n&D,u=t(e),f=u.length;if(f!=t(r).length&&!c)return!1;for(var s=f;s--;){var v=u[s];if(!(c?v in r:N.call(r,v)))return!1}var l=i.get(e),b=i.get(r);if(l&&b)return l==r&&b==e;var h=!0;i.set(e,r),i.set(r,e);for(var p=c;++s<f;){var _=e[v=u[s]],d=r[v];if(a)var y=c?a(d,_,v,r,e,i):a(_,d,v,e,r,i);if(!(void 0===y?_===d||o(_,d,n,a,i):y)){h=!1;break}p||(p="constructor"==v)}if(h&&!p){var g=e.constructor,j=r.constructor;g==j||!("constructor"in e)||!("constructor"in r)||"function"==typeof g&&g instanceof g&&"function"==typeof j&&j instanceof j||(h=!1)}return i.delete(e),i.delete(r),h}(i,c,s,v,l,b))}function U(e,t,r,n,a){return e===t||(null==e||null==t||!s(e)&&!s(t)?e!=e&&t!=t:R(e,t,r,n,U,a))}function V(e,t){return U(e,t)}export{U as b,V as i};