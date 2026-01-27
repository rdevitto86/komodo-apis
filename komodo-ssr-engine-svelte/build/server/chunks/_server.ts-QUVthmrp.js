import { e as error, j as json, i as isHttpError } from './index-CJiV4IBS.js';
import { g as getPageContentFromS3 } from './s3-BdyB1uLY.js';
import { l as logger } from './index-DalvyAF_.js';
import '@aws-sdk/client-s3';

async function GET({ params }) {
  try {
    const product = await getPageContentFromS3(`product-${params.id}`);
    if (!product) throw error(404, `Product ${params.id} not found`);
    return json(product);
  } catch (err) {
    if (isHttpError(err)) {
      logger.error(err?.body?.message);
      throw err;
    }
    logger.error("Failed to fetch product");
    throw error(500, "Failed to fetch product");
  }
}

export { GET };
//# sourceMappingURL=_server.ts-QUVthmrp.js.map
