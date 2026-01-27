import { e as ensure_array_like, b as bind_props } from './index2-DTE7gUAU.js';
import { e as escape_html } from './context-D64tQzLr.js';

function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let data = $$props["data"];
    $$renderer2.push(`<h1>${escape_html(data.content.title)}</h1> <!--[-->`);
    const each_array = ensure_array_like(data.content.sections);
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let section = each_array[$$index];
      $$renderer2.push(`<section><h2>${escape_html(section.heading)}</h2> <p>${escape_html(section.content)}</p></section>`);
    }
    $$renderer2.push(`<!--]-->`);
    bind_props($$props, { data });
  });
}

export { _page as default };
//# sourceMappingURL=_page.svelte-ZI5QsAMS.js.map
