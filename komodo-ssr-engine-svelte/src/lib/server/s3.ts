import { S3Client, GetObjectCommand } from "@aws-sdk/client-s3";
import { logger } from '../logger';

const BUCKET = process.env.S3_CONTENT_BUCKET!;

interface CachedContent {
  data: any;
  expires: number;
}

const contentCache = new Map<string, CachedContent>();
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

const s3 = new S3Client({
  region: process.env.AWS_REGION || "us-east-1"
});

export async function getPageContentFromS3(pageKey: string) {
  const cached = contentCache.get(pageKey);

  if (cached && cached.expires > Date.now()) {
    logger.info(`[Cache HIT] ${pageKey}`);
    return cached.data;
  }

  logger.info(`[Cache MISS] Fetching ${pageKey} from S3`);
  
  try {
    const command = new GetObjectCommand({
      Bucket: BUCKET,
      Key: `pages/${pageKey}.json`
    });

    const response = await s3.send(command);
    const bodyString = await response.Body?.transformToString();
    
    if (!bodyString) {
      throw new Error(`Empty response for ${pageKey}`);
    }

    const content = JSON.parse(bodyString);

    contentCache.set(pageKey, {
      data: content,
      expires: Date.now() + CACHE_TTL
    });

    return content;
  } catch (error) {
    logger.error(`Failed to fetch ${pageKey} from S3:`, error as Error);
    
    if (cached) {
      logger.info(`[Fallback] Using stale cache for ${pageKey}`);
      return cached.data;
    }
    throw error;
  }
}
