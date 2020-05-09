define(["react","@grafana/ui","@grafana/data","@grafana/runtime"],(function(t,n,e,r){return function(t){var n={};function e(r){if(n[r])return n[r].exports;var o=n[r]={i:r,l:!1,exports:{}};return t[r].call(o.exports,o,o.exports,e),o.l=!0,o.exports}return e.m=t,e.c=n,e.d=function(t,n,r){e.o(t,n)||Object.defineProperty(t,n,{enumerable:!0,get:r})},e.r=function(t){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(t,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(t,"__esModule",{value:!0})},e.t=function(t,n){if(1&n&&(t=e(t)),8&n)return t;if(4&n&&"object"==typeof t&&t&&t.__esModule)return t;var r=Object.create(null);if(e.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:t}),2&n&&"string"!=typeof t)for(var o in t)e.d(r,o,function(n){return t[n]}.bind(null,o));return r},e.n=function(t){var n=t&&t.__esModule?function(){return t.default}:function(){return t};return e.d(n,"a",n),n},e.o=function(t,n){return Object.prototype.hasOwnProperty.call(t,n)},e.p="/",e(e.s=51)}([function(n,e){n.exports=t},function(t,e){t.exports=n},function(t,n){t.exports=function(t){var n=typeof t;return null!=t&&("object"==n||"function"==n)}},function(t,n,e){var r=e(8),o=e(28),i=e(29),a=r?r.toStringTag:void 0;t.exports=function(t){return null==t?void 0===t?"[object Undefined]":"[object Null]":a&&a in Object(t)?o(t):i(t)}},function(t,n,e){var r=e(9),o="object"==typeof self&&self&&self.Object===Object&&self,i=r||o||Function("return this")();t.exports=i},function(t,n){t.exports=function(t){return null!=t&&"object"==typeof t}},function(t,n){t.exports=function(t){return t}},function(t,n,e){var r=e(3),o=e(2);t.exports=function(t){if(!o(t))return!1;var n=r(t);return"[object Function]"==n||"[object GeneratorFunction]"==n||"[object AsyncFunction]"==n||"[object Proxy]"==n}},function(t,n,e){var r=e(4).Symbol;t.exports=r},function(t,n,e){(function(n){var e="object"==typeof n&&n&&n.Object===Object&&n;t.exports=e}).call(this,e(27))},function(t,n){t.exports=function(t,n){return t===n||t!=t&&n!=n}},function(t,n,e){var r=e(7),o=e(12);t.exports=function(t){return null!=t&&o(t.length)&&!r(t)}},function(t,n){t.exports=function(t){return"number"==typeof t&&t>-1&&t%1==0&&t<=9007199254740991}},function(t,n){var e=/^(?:0|[1-9]\d*)$/;t.exports=function(t,n){var r=typeof t;return!!(n=null==n?9007199254740991:n)&&("number"==r||"symbol"!=r&&e.test(t))&&t>-1&&t%1==0&&t<n}},function(t,n){t.exports=function(t){return t.webpackPolyfill||(t.deprecate=function(){},t.paths=[],t.children||(t.children=[]),Object.defineProperty(t,"loaded",{enumerable:!0,get:function(){return t.l}}),Object.defineProperty(t,"id",{enumerable:!0,get:function(){return t.i}}),t.webpackPolyfill=1),t}},function(t,n){t.exports=e},function(t,n){t.exports=r},function(t,n,e){var r=e(18),o=e(10),i=e(35),a=e(36),u=Object.prototype,c=u.hasOwnProperty,p=r((function(t,n){t=Object(t);var e=-1,r=n.length,p=r>2?n[2]:void 0;for(p&&i(n[0],n[1],p)&&(r=1);++e<r;)for(var f=n[e],s=a(f),l=-1,y=s.length;++l<y;){var v=s[l],b=t[v];(void 0===b||o(b,u[v])&&!c.call(t,v))&&(t[v]=f[v])}return t}));t.exports=p},function(t,n,e){var r=e(6),o=e(19),i=e(21);t.exports=function(t,n){return i(o(t,n,r),t+"")}},function(t,n,e){var r=e(20),o=Math.max;t.exports=function(t,n,e){return n=o(void 0===n?t.length-1:n,0),function(){for(var i=arguments,a=-1,u=o(i.length-n,0),c=Array(u);++a<u;)c[a]=i[n+a];a=-1;for(var p=Array(n+1);++a<n;)p[a]=i[a];return p[n]=e(c),r(t,this,p)}}},function(t,n){t.exports=function(t,n,e){switch(e.length){case 0:return t.call(n);case 1:return t.call(n,e[0]);case 2:return t.call(n,e[0],e[1]);case 3:return t.call(n,e[0],e[1],e[2])}return t.apply(n,e)}},function(t,n,e){var r=e(22),o=e(34)(r);t.exports=o},function(t,n,e){var r=e(23),o=e(24),i=e(6),a=o?function(t,n){return o(t,"toString",{configurable:!0,enumerable:!1,value:r(n),writable:!0})}:i;t.exports=a},function(t,n){t.exports=function(t){return function(){return t}}},function(t,n,e){var r=e(25),o=function(){try{var t=r(Object,"defineProperty");return t({},"",{}),t}catch(t){}}();t.exports=o},function(t,n,e){var r=e(26),o=e(33);t.exports=function(t,n){var e=o(t,n);return r(e)?e:void 0}},function(t,n,e){var r=e(7),o=e(30),i=e(2),a=e(32),u=/^\[object .+?Constructor\]$/,c=Function.prototype,p=Object.prototype,f=c.toString,s=p.hasOwnProperty,l=RegExp("^"+f.call(s).replace(/[\\^$.*+?()[\]{}|]/g,"\\$&").replace(/hasOwnProperty|(function).*?(?=\\\()| for .+?(?=\\\])/g,"$1.*?")+"$");t.exports=function(t){return!(!i(t)||o(t))&&(r(t)?l:u).test(a(t))}},function(t,n){var e;e=function(){return this}();try{e=e||new Function("return this")()}catch(t){"object"==typeof window&&(e=window)}t.exports=e},function(t,n,e){var r=e(8),o=Object.prototype,i=o.hasOwnProperty,a=o.toString,u=r?r.toStringTag:void 0;t.exports=function(t){var n=i.call(t,u),e=t[u];try{t[u]=void 0;var r=!0}catch(t){}var o=a.call(t);return r&&(n?t[u]=e:delete t[u]),o}},function(t,n){var e=Object.prototype.toString;t.exports=function(t){return e.call(t)}},function(t,n,e){var r,o=e(31),i=(r=/[^.]+$/.exec(o&&o.keys&&o.keys.IE_PROTO||""))?"Symbol(src)_1."+r:"";t.exports=function(t){return!!i&&i in t}},function(t,n,e){var r=e(4)["__core-js_shared__"];t.exports=r},function(t,n){var e=Function.prototype.toString;t.exports=function(t){if(null!=t){try{return e.call(t)}catch(t){}try{return t+""}catch(t){}}return""}},function(t,n){t.exports=function(t,n){return null==t?void 0:t[n]}},function(t,n){var e=Date.now;t.exports=function(t){var n=0,r=0;return function(){var o=e(),i=16-(o-r);if(r=o,i>0){if(++n>=800)return arguments[0]}else n=0;return t.apply(void 0,arguments)}}},function(t,n,e){var r=e(10),o=e(11),i=e(13),a=e(2);t.exports=function(t,n,e){if(!a(e))return!1;var u=typeof n;return!!("number"==u?o(e)&&i(n,e.length):"string"==u&&n in e)&&r(e[n],t)}},function(t,n,e){var r=e(37),o=e(48),i=e(11);t.exports=function(t){return i(t)?r(t,!0):o(t)}},function(t,n,e){var r=e(38),o=e(39),i=e(41),a=e(42),u=e(13),c=e(44),p=Object.prototype.hasOwnProperty;t.exports=function(t,n){var e=i(t),f=!e&&o(t),s=!e&&!f&&a(t),l=!e&&!f&&!s&&c(t),y=e||f||s||l,v=y?r(t.length,String):[],b=v.length;for(var d in t)!n&&!p.call(t,d)||y&&("length"==d||s&&("offset"==d||"parent"==d)||l&&("buffer"==d||"byteLength"==d||"byteOffset"==d)||u(d,b))||v.push(d);return v}},function(t,n){t.exports=function(t,n){for(var e=-1,r=Array(t);++e<t;)r[e]=n(e);return r}},function(t,n,e){var r=e(40),o=e(5),i=Object.prototype,a=i.hasOwnProperty,u=i.propertyIsEnumerable,c=r(function(){return arguments}())?r:function(t){return o(t)&&a.call(t,"callee")&&!u.call(t,"callee")};t.exports=c},function(t,n,e){var r=e(3),o=e(5);t.exports=function(t){return o(t)&&"[object Arguments]"==r(t)}},function(t,n){var e=Array.isArray;t.exports=e},function(t,n,e){(function(t){var r=e(4),o=e(43),i=n&&!n.nodeType&&n,a=i&&"object"==typeof t&&t&&!t.nodeType&&t,u=a&&a.exports===i?r.Buffer:void 0,c=(u?u.isBuffer:void 0)||o;t.exports=c}).call(this,e(14)(t))},function(t,n){t.exports=function(){return!1}},function(t,n,e){var r=e(45),o=e(46),i=e(47),a=i&&i.isTypedArray,u=a?o(a):r;t.exports=u},function(t,n,e){var r=e(3),o=e(12),i=e(5),a={};a["[object Float32Array]"]=a["[object Float64Array]"]=a["[object Int8Array]"]=a["[object Int16Array]"]=a["[object Int32Array]"]=a["[object Uint8Array]"]=a["[object Uint8ClampedArray]"]=a["[object Uint16Array]"]=a["[object Uint32Array]"]=!0,a["[object Arguments]"]=a["[object Array]"]=a["[object ArrayBuffer]"]=a["[object Boolean]"]=a["[object DataView]"]=a["[object Date]"]=a["[object Error]"]=a["[object Function]"]=a["[object Map]"]=a["[object Number]"]=a["[object Object]"]=a["[object RegExp]"]=a["[object Set]"]=a["[object String]"]=a["[object WeakMap]"]=!1,t.exports=function(t){return i(t)&&o(t.length)&&!!a[r(t)]}},function(t,n){t.exports=function(t){return function(n){return t(n)}}},function(t,n,e){(function(t){var r=e(9),o=n&&!n.nodeType&&n,i=o&&"object"==typeof t&&t&&!t.nodeType&&t,a=i&&i.exports===o&&r.process,u=function(){try{var t=i&&i.require&&i.require("util").types;return t||a&&a.binding&&a.binding("util")}catch(t){}}();t.exports=u}).call(this,e(14)(t))},function(t,n,e){var r=e(2),o=e(49),i=e(50),a=Object.prototype.hasOwnProperty;t.exports=function(t){if(!r(t))return i(t);var n=o(t),e=[];for(var u in t)("constructor"!=u||!n&&a.call(t,u))&&e.push(u);return e}},function(t,n){var e=Object.prototype;t.exports=function(t){var n=t&&t.constructor;return t===("function"==typeof n&&n.prototype||e)}},function(t,n){t.exports=function(t){var n=[];if(null!=t)for(var e in Object(t))n.push(e);return n}},function(t,n,e){"use strict";e.r(n);var r=e(15),o=function(t,n){return(o=Object.setPrototypeOf||{__proto__:[]}instanceof Array&&function(t,n){t.__proto__=n}||function(t,n){for(var e in n)n.hasOwnProperty(e)&&(t[e]=n[e])})(t,n)};
/*! *****************************************************************************
Copyright (c) Microsoft Corporation.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
PERFORMANCE OF THIS SOFTWARE.
***************************************************************************** */function i(t,n){function e(){this.constructor=t}o(t,n),t.prototype=null===n?Object.create(n):(e.prototype=n.prototype,new e)}var a=function(){return(a=Object.assign||function(t){for(var n,e=1,r=arguments.length;e<r;e++)for(var o in n=arguments[e])Object.prototype.hasOwnProperty.call(n,o)&&(t[o]=n[o]);return t}).apply(this,arguments)};var u=function(t){function n(n,e){var r=t.call(this,n)||this;return r.templateSrv=e,r}return n.$inject=["instanceSettings","templateSrv"],i(n,t),n.prototype.applyTemplateVariables=function(t){return a(a({},t),{stationCRS:this.templateSrv.replace(t.stationCRS)})},n}(e(16).DataSourceWithBackend),c=e(0),p=e.n(c),f=e(1),s="https://api.departureboard.io/api/v2.0",l={endpoint:"getDeparturesByCRS"},y=f.LegacyForms.SecretFormField,v=f.LegacyForms.FormField,b=function(t){function n(){var n=null!==t&&t.apply(this,arguments)||this;return n.onAPIEndpointChange=function(t){var e=n.props,r=e.onOptionsChange,o=e.options,i=a(a({},o.jsonData),{apiEndpoint:t.target.value});r(a(a({},o),{jsonData:i}))},n.onAPIKeyChange=function(t){var e=n.props,r=e.onOptionsChange,o=e.options;r(a(a({},o),{secureJsonData:{apiKey:t.target.value}}))},n.onResetAPIKey=function(){var t=n.props,e=t.onOptionsChange,r=t.options;e(a(a({},r),{secureJsonFields:a(a({},r.secureJsonFields),{apiKey:!1}),secureJsonData:a(a({},r.secureJsonData),{apiKey:""})}))},n.componentDidMount=function(){var t=n.props,e=t.onOptionsChange,r=t.options,o=a(a({},r.jsonData),{apiEndpoint:s});e(a(a({},r),{jsonData:o}))},n}return i(n,t),n.prototype.render=function(){var t=this.props.options,n=t.jsonData,e=t.secureJsonFields,r=t.secureJsonData||{};return p.a.createElement("div",{className:"gf-form-group"},p.a.createElement("div",{className:"gf-form"},p.a.createElement(v,{label:"API Endpoint",labelWidth:10,inputWidth:20,onChange:this.onAPIEndpointChange,value:n.apiEndpoint||"",placeholder:"departureboard.io API URL"})),p.a.createElement("div",{className:"gf-form-inline"},p.a.createElement("div",{className:"gf-form"},p.a.createElement(y,{isConfigured:e&&e.apiKey,value:r.apiKey||"",label:"API Key",placeholder:"National Rail API key",labelWidth:10,inputWidth:17,onReset:this.onResetAPIKey,onChange:this.onAPIKeyChange}))))},n}(c.PureComponent),d=e(17),h=e.n(d),g=f.LegacyForms.FormField,j=function(t){function n(){var n=null!==t&&t.apply(this,arguments)||this;return n.onQueryTextChange=function(t){var e=n.props,r=e.onChange,o=e.query,i=e.onRunQuery;r(a(a({},o),{stationCRS:t.target.value})),i()},n}return i(n,t),n.prototype.render=function(){var t=h()(this.props.query,l).stationCRS;return p.a.createElement("div",{className:"gf-form"},p.a.createElement(g,{labelWidth:10,value:t||"",onChange:this.onQueryTextChange,placeholder:"PAD",label:"CRS code for station",tooltip:""}))},n}(c.PureComponent);e.d(n,"plugin",(function(){return x}));var x=new r.DataSourcePlugin(u).setConfigEditor(b).setQueryEditor(j)}])}));
//# sourceMappingURL=module.js.map