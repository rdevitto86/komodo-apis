import { e as error } from './index-CJiV4IBS.js';
import { l as logger } from './index-DalvyAF_.js';

const AUTH_API_URL = process.env.AUTH_API_URL;
if (!AUTH_API_URL) throw new Error("AUTH_API_URL is not defined");
async function validateApiToken(token) {
  try {
    const response = await fetch(`${AUTH_API_URL}/auth/validate`, {
      method: "POST",
      headers: {
        "Accept": "application/json;v=1",
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
      }
    });
    if (!response.ok) {
      logger.error("Token validation failed", new Error(`Token validation failed: ${response.status}`));
      return null;
    }
    const data = await response.json();
    return {
      id: data.userId,
      email: data.email,
      isAdmin: data.role === "admin" || data.isAdmin
    };
  } catch (err) {
    logger.error("Auth API call failed", err);
    return null;
  }
}
const PROTECTED_ROUTES = [
  "/orders",
  "/marketing/*",
  "/services/scheduling"
].map((route) => new RegExp(`/api(/v\\d+)?${route.replace(/\*/g, "(/.*)?")}(/|$)`));
const ADMIN_ROUTES = [
  "/admin/manage/content/*"
].map((route) => new RegExp(`/api(/v\\d+)?${route.replace(/\*/g, "(/.*)?")}(/|$)`));
const handle = async ({ event, resolve }) => {
  try {
    const start = Date.now();
    const path = event.url.pathname;
    const token = event.request.headers.get("Authorization")?.replace("Bearer ", "");
    if (token) {
      event.locals.user = await validateApiToken(token);
    }
    if (isValidRoute(PROTECTED_ROUTES, path) && !event.locals.user) {
      throw error(401, "Unauthorized - Valid token required");
    }
    if (isValidRoute(ADMIN_ROUTES, path) && !event.locals.user?.isAdmin) {
      throw error(403, "Forbidden - Admin access required");
    }
    const response = await resolve(event);
    logger.info("Request completed", {
      method: event.request.method,
      path,
      status: response.status,
      duration: Date.now() - start,
      userId: event.locals.user?.id
    });
    return response;
  } catch (err) {
    logger.error("Request failed", err);
    throw err;
  }
};
const isValidRoute = (routes, path) => routes.some((pattern) => pattern.test(path));

export { handle };
//# sourceMappingURL=hooks.server-DvsLJlfy.js.map
