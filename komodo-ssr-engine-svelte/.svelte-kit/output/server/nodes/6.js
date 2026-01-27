import * as server from '../entries/pages/landing/_page.server.ts.js';

export const index = 6;
let component_cache;
export const component = async () => component_cache ??= (await import('../entries/pages/landing/_page.svelte.js')).default;
export { server };
export const server_id = "src/routes/landing/+page.server.ts";
export const imports = ["_app/immutable/nodes/6.CWG1ehzT.js","_app/immutable/chunks/CWj6FrbW.js","_app/immutable/chunks/69_IOA4Y.js","_app/immutable/chunks/DIeogL5L.js"];
export const stylesheets = [];
export const fonts = [];
