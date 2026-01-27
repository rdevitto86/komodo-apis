async function load({ setHeaders }) {
  setHeaders({
    "cache-control": "public, max-age=300, s-maxage=600"
  });
  return {
    seed: Math.random()
  };
}

var _page_server_ts = /*#__PURE__*/Object.freeze({
  __proto__: null,
  load: load
});

const index = 2;
let component_cache;
const component = async () => component_cache ??= (await import('./_page.svelte-D9ZlaW-G.js')).default;
const server_id = "src/routes/+page.server.ts";
const imports = ["_app/immutable/nodes/2.ellRxtD2.js","_app/immutable/chunks/CWj6FrbW.js","_app/immutable/chunks/69_IOA4Y.js","_app/immutable/chunks/DIeogL5L.js","_app/immutable/chunks/DiUamYel.js","_app/immutable/chunks/CDPdhYCU.js","_app/immutable/chunks/BagUfPUu.js","_app/immutable/chunks/CsoU4Bzd.js"];
const stylesheets = [];
const fonts = [];

export { component, fonts, imports, index, _page_server_ts as server, server_id, stylesheets };
//# sourceMappingURL=2-DCgFpBuQ.js.map
