import { j as json } from './index-CJiV4IBS.js';
import { S3Client, HeadBucketCommand } from '@aws-sdk/client-s3';

async function GET() {
  const checks = {
    status: "ok",
    timestamp: (/* @__PURE__ */ new Date()).toISOString(),
    s3: "unknown"
  };
  try {
    await new S3Client({ region: process.env.AWS_REGION }).send(new HeadBucketCommand({ Bucket: process.env.S3_CONTENT_BUCKET }));
    checks.s3 = "ok";
  } catch {
    checks.s3 = "error";
    checks.status = "degraded";
  }
  return json(checks, { status: checks.status === "ok" ? 200 : 503 });
}

export { GET };
//# sourceMappingURL=_server.ts-C39N6Xy_.js.map
