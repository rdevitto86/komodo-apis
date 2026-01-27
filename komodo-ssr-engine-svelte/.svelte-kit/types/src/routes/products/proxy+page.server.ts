// @ts-nocheck
import { getPageContentFromS3 } from '$lib/server/content';
import type { PageServerLoad } from './$types';

export const load = async ({ setHeaders }: Parameters<PageServerLoad>[0]) => {
  setHeaders({
    'cache-control': 'public, max-age=300, s-maxage=600, stale-while-revalidate=86400'
  });

  const content = await getPageContentFromS3('products');
  
  return { content };
};
