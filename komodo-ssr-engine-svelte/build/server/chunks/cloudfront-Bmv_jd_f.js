import { CreateInvalidationCommand, CloudFrontClient } from '@aws-sdk/client-cloudfront';

const cf = new CloudFrontClient({
  region: process.env.AWS_REGION || "us-east-1"
});
const DISTRIBUTION_ID = process.env.CLOUDFRONT_DISTRIBUTION_ID;
async function invalidateCloudFrontCache(paths) {
  if (!DISTRIBUTION_ID) {
    console.warn("CloudFront distribution ID not configured");
    return;
  }
  try {
    const command = new CreateInvalidationCommand({
      DistributionId: DISTRIBUTION_ID,
      InvalidationBatch: {
        CallerReference: `invalidation-${Date.now()}`,
        Paths: {
          Quantity: paths.length,
          Items: paths
        }
      }
    });
    const response = await cf.send(command);
    console.log(`[CloudFront] Invalidated paths:`, paths);
    return response.Invalidation;
  } catch (error) {
    console.error("CloudFront invalidation failed:", error);
    throw error;
  }
}
async function invalidatePage(pageKey) {
  return invalidateCloudFrontCache([`/${pageKey}`, `/${pageKey}/`]);
}

export { invalidatePage as i };
//# sourceMappingURL=cloudfront-Bmv_jd_f.js.map
