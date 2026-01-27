import * as server from '../entries/pages/_page.server.ts.js';

export const index = 2;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/_page.svelte.js')).default;
export { server };
export const server_id = "src/routes/+page.server.ts";
export const imports = ["_app/immutable/nodes/2.ellRxtD2.js","_app/immutable/chunks/CWj6FrbW.js","_app/immutable/chunks/69_IOA4Y.js","_app/immutable/chunks/DIeogL5L.js","_app/immutable/chunks/DiUamYel.js","_app/immutable/chunks/CDPdhYCU.js","_app/immutable/chunks/BagUfPUu.js","_app/immutable/chunks/CsoU4Bzd.js"];
export const stylesheets = [];
export const fonts = [];
