package context

import (
	"context"
	ctxKeys "komodo-internal-lib-apis-go/common/context"
	utils "komodo-internal-lib-apis-go/http/utils/http"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// Enriches request context with common values
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		var reqID string
		if rid := req.Header.Get("X-Request-ID"); rid != "" {
			reqID = rid
		} else if rid := chimw.GetReqID(ctx); rid != "" {
			reqID = rid
		} else if rid, ok := ctx.Value(ctxKeys.REQUEST_ID_KEY).(string); ok && rid != "" {
			reqID = rid
		} else {
			reqID = utils.GenerateRequestId()
		}
		ctx = context.WithValue(ctx, chimw.RequestIDKey, reqID)
		ctx = context.WithValue(ctx, ctxKeys.REQUEST_ID_KEY, reqID)
		req.Header.Set("X-Request-ID", reqID)
		wtr.Header().Set("X-Request-ID", reqID)

		ctx = context.WithValue(ctx, ctxKeys.START_TIME_KEY, time.Now().UTC())
		ctx = context.WithValue(ctx, ctxKeys.VERSION_KEY, utils.GetAPIVersion(req))
		ctx = context.WithValue(ctx, ctxKeys.URI_KEY, utils.GetAPIRoute(req))
		ctx = context.WithValue(ctx, ctxKeys.METHOD_KEY, req.Method)
		ctx = context.WithValue(ctx, ctxKeys.PATH_PARAMS_KEY, utils.GetPathParams(req))
		ctx = context.WithValue(ctx, ctxKeys.QUERY_PARAMS_KEY, utils.GetQueryParams(req))
		// ctx = context.WithValue(ctx, ctxKeys.CLIENT_IP_KEY, utils.GetClientIP(req))
		ctx = context.WithValue(ctx, ctxKeys.CLIENT_TYPE_KEY, utils.GetClientType(req))
		ctx = context.WithValue(ctx, ctxKeys.USER_AGENT_KEY, req.UserAgent())

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}
