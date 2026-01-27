import { error, type Handle } from '@sveltejs/kit';
import { validateApiToken } from '$lib/server/auth';
import { logger } from '$lib/logger';

const PROTECTED_ROUTES = [
  '/orders',
  '/marketing/*',
  '/services/scheduling'
].map(route => new RegExp(`/api(/v\\d+)?${route.replace(/\*/g, '(/.*)?')}(/|$)`));

const ADMIN_ROUTES = [
  '/admin/manage/content/*',
].map(route => new RegExp(`/api(/v\\d+)?${route.replace(/\*/g, '(/.*)?')}(/|$)`));

export const handle: Handle = async ({ event, resolve }) => {
  try {
    const start = Date.now();
    const path = event.url.pathname;
    const token = event.request.headers.get('Authorization')?.replace('Bearer ', '');
    
    if (token) {
      event.locals.user = await validateApiToken(token);
    }
    
    if (isValidRoute(PROTECTED_ROUTES, path) && !event.locals.user) {
      throw error(401, 'Unauthorized - Valid token required');
    }
    if (isValidRoute(ADMIN_ROUTES, path) && !event.locals.user?.isAdmin) {
      throw error(403, 'Forbidden - Admin access required');
    }
    
    const response = await resolve(event);
    
    logger.info('Request completed', {
      method: event.request.method,
      path,
      status: response.status,
      duration: Date.now() - start,
      userId: event.locals.user?.id
    });
    
    return response;
  } catch (err) {
    logger.error('Request failed', err as Error);
    throw err;
  }
};

const isValidRoute = (routes: RegExp[], path: string) => routes.some(pattern => pattern.test(path));