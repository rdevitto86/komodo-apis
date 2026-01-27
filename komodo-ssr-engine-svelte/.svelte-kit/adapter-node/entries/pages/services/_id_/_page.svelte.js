import { h as head, e as ensure_array_like, a as attr, b as bind_props } from "../../../../chunks/index2.js";
import { e as escape_html } from "../../../../chunks/context.js";
function html(value) {
  var html2 = String(value ?? "");
  var open = "<!---->";
  return open + html2 + "<!---->";
}
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let data = $$props["data"];
    head("1v4eju1", $$renderer2, ($$renderer3) => {
      $$renderer3.title(($$renderer4) => {
        $$renderer4.push(`<title>${escape_html(data.meta.title)}</title>`);
      });
      $$renderer3.push(`<meta name="description"${attr("content", data.meta.description)}/>`);
    });
    $$renderer2.push(`<article><h1>${escape_html(data.content.title)}</h1> <!--[-->`);
    const each_array = ensure_array_like(data.content.sections);
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let section = each_array[$$index];
      $$renderer2.push(`<section><h2>${escape_html(section.heading)}</h2> <div>${html(section.content)}</div> `);
      if (section.image) {
        $$renderer2.push("<!--[-->");
        $$renderer2.push(`<img${attr("src", section.image)}${attr("alt", section.heading)}/>`);
      } else {
        $$renderer2.push("<!--[!-->");
      }
      $$renderer2.push(`<!--]--></section>`);
    }
    $$renderer2.push(`<!--]--></article>`);
    bind_props($$props, { data });
  });
}
export {
  _page as default
};
