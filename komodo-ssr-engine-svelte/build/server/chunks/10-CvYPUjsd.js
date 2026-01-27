import { g as getPageContentFromS3 } from './s3-BdyB1uLY.js';
import '@aws-sdk/client-s3';
import './index-DalvyAF_.js';

const load = async ({ setHeaders }) => {
  setHeaders({
    "cache-control": "public, max-age=300, s-maxage=600, stale-while-revalidate=86400"
  });
  const content = await getPageContentFromS3("products");
  return { content };
};

var _page_server_ts = /*#__PURE__*/Object.freeze({
  __proto__: null,
  load: load
});

const index = 10;
let component_cache;
const component = async () => component_cache ??= (await import('./_page.svelte-ZI5QsAMS.js')).default;
const server_id = "src/routes/products/[id]/+page.server.ts";
const imports = ["_app/immutable/nodes/10.BFsSYADE.js","_app/immutable/chunks/CWj6FrbW.js","_app/immutable/chunks/69_IOA4Y.js","_app/immutable/chunks/DIeogL5L.js","_app/immutable/chunks/DiUamYel.js","_app/immutable/chunks/CDPdhYCU.js","_app/immutable/chunks/B-Yn3XEC.js","_app/immutable/chunks/BagUfPUu.js","_app/immutable/chunks/CsoU4Bzd.js"];
const stylesheets = [];
const fonts = [];

export { component, fonts, imports, index, _page_server_ts as server, server_id, stylesheets };
//# sourceMappingURL=10-CvYPUjsd.js.map
