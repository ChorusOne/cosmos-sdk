(window.webpackJsonp=window.webpackJsonp||[]).push([[215],{812:function(e,t,r){"use strict";r.r(t);var n=r(1),a=Object(n.a)({},(function(){var e=this,t=e.$createElement,r=e._self._c||t;return r("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[r("h1",{attrs:{id:"concepts"}},[r("a",{staticClass:"header-anchor",attrs:{href:"#concepts"}},[e._v("#")]),e._v(" Concepts")]),e._v(" "),r("h2",{attrs:{id:"reference-counting-in-f1-fee-distribution"}},[r("a",{staticClass:"header-anchor",attrs:{href:"#reference-counting-in-f1-fee-distribution"}},[e._v("#")]),e._v(" Reference Counting in F1 Fee Distribution")]),e._v(" "),r("p",[e._v("In F1 fee distribution, in order to calculate the rewards a delegator ought to receive when they\nwithdraw their delegation, we must read the terms of the summation of rewards divided by tokens from\nthe period which they ended when they delegated, and the final period (created when they withdraw).")]),e._v(" "),r("p",[e._v("Additionally, as slashes change the amount of tokens a delegation will have (but we calculate this lazily,\nonly when a delegator un-delegates), we must calculate rewards in separate periods before / after any slashes\nwhich occurred in between when a delegator delegated and when they withdrew their rewards. Thus slashes, like\ndelegations, reference the period which was ended by the slash event.")]),e._v(" "),r("p",[e._v("All stored historical rewards records for periods which are no longer referenced by any delegations\nor any slashes can thus be safely removed, as they will never be read (future delegations and future\nslashes will always reference future periods). This is implemented by tracking a "),r("code",[e._v("ReferenceCount")]),e._v("\nalong with each historical reward storage entry. Each time a new object (delegation or slash)\nis created which might need to reference the historical record, the reference count is incremented.\nEach time one object which previously needed to reference the historical record is deleted, the reference\ncount is decremented. If the reference count hits zero, the historical record is deleted.")])])}),[],!1,null,null,null);t.default=a.exports}}]);