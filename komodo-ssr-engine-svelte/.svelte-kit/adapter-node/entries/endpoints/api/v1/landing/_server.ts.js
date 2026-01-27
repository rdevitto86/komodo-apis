import { json } from "@sveltejs/kit";
async function GET({ params }) {
  return json({ seed: params.id });
}
export {
  GET
};
