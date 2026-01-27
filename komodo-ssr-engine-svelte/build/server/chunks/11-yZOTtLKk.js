import { g as getPageContentFromS3 } from './s3-BdyB1uLY.js';
import '@aws-sdk/client-s3';
import './index-DalvyAF_.js';

const load = async ({ setHeaders }) => {
  setHeaders({
    "cache-control": "public, max-age=300, s-maxage=600, stale-while-revalidate=86400"
  });
  const content = await getPageContentFromS3("services");
  return {
    content,
    meta: {
      title: content.title || "Services",
      description: content.description || ""
    }
  };
};

var _page_server_ts = /*#__PURE__*/Object.freeze({
  __proto__: null,
  load: load
});

const index = 11;
let component_cache;
const component = async () => component_cache ??= (await import('./_page.svelte-B-3DbIFz.js')).default;
const server_id = "src/routes/services/[id]/+page.server.ts";
const imports = ["_app/immutable/nodes/11.DTsU1vyc.js","_app/immutable/chunks/CWj6FrbW.js","_app/immutable/chunks/69_IOA4Y.js","_app/immutable/chunks/DIeogL5L.js","_app/immutable/chunks/DiUamYel.js","_app/immutable/chunks/CDPdhYCU.js","_app/immutable/chunks/oi9Nr6Ux.js","_app/immutable/chunks/B-Yn3XEC.js","_app/immutable/chunks/BagUfPUu.js","_app/immutable/chunks/CsoU4Bzd.js"];
const stylesheets = [];
const fonts = [];

export { component, fonts, imports, index, _page_server_ts as server, server_id, stylesheets };
//# sourceMappingURL=11-yZOTtLKk.js.map
