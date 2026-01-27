import { error, json } from "@sveltejs/kit";
import { i as invalidatePage } from "../../../../../../../../chunks/cloudfront.js";
const POST = async ({ request, locals }) => {
  const { pageKey, cloudfront = true } = await request.json();
  if (!pageKey) {
    throw error(400, "pageKey required");
  }
  if (cloudfront) {
    try {
      await invalidatePage(pageKey);
    } catch (err) {
      console.error("CloudFront invalidation failed:", err);
    }
  }
  return json({
    success: true,
    invalidated: pageKey,
    timestamp: (/* @__PURE__ */ new Date()).toISOString()
  });
};
export {
  POST
};
