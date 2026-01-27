import { b as bind_props } from './index2-DTE7gUAU.js';
import { e as escape_html } from './context-D64tQzLr.js';

function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let data = $$props["data"];
    $$renderer2.push(`<h1>Landing Page</h1> <p>Seed: ${escape_html(data.seed)}</p>`);
    bind_props($$props, { data });
  });
}

export { _page as default };
//# sourceMappingURL=_page.svelte-D9ZlaW-G.js.map
