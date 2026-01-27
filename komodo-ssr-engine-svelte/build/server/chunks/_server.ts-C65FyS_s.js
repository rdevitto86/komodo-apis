import { e as error, j as json } from './index-CJiV4IBS.js';
import { i as invalidatePage } from './cloudfront-Bmv_jd_f.js';
import '@aws-sdk/client-cloudfront';

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

export { POST };
//# sourceMappingURL=_server.ts-C65FyS_s.js.map
