import * as server from '../entries/pages/services/_id_/_page.server.ts.js';

export const index = 11;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/services/_id_/_page.svelte.js')).default;
export { server };
export const server_id = "src/routes/services/[id]/+page.server.ts";
export const imports = ["_app/immutable/nodes/11.DTsU1vyc.js","_app/immutable/chunks/CWj6FrbW.js","_app/immutable/chunks/69_IOA4Y.js","_app/immutable/chunks/DIeogL5L.js","_app/immutable/chunks/DiUamYel.js","_app/immutable/chunks/CDPdhYCU.js","_app/immutable/chunks/oi9Nr6Ux.js","_app/immutable/chunks/B-Yn3XEC.js","_app/immutable/chunks/BagUfPUu.js","_app/immutable/chunks/CsoU4Bzd.js"];
export const stylesheets = [];
export const fonts = [];
