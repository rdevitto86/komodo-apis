package context

import (
	"context"
	utils "komodo-internal-lib-apis-go/http/utils/http"
	ctxKeys "komodo-internal-lib-apis-go/types/context"
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
		} else if rid, ok := ctx.Value(ctxKeys.RequestIDKey).(string); ok && rid != "" {
			reqID = rid
		} else {
			reqID = utils.GenerateRequestId()
		}
		ctx = context.WithValue(ctx, chimw.RequestIDKey, reqID)
		ctx = context.WithValue(ctx, ctxKeys.RequestIDKey, reqID)
		req.Header.Set("X-Request-ID", reqID)
		wtr.Header().Set("X-Request-ID", reqID)

		ctx = context.WithValue(ctx, ctxKeys.StartTimeKey, time.Now().UTC())
		ctx = context.WithValue(ctx, ctxKeys.VersionKey, utils.GetAPIVersion(req))
		ctx = context.WithValue(ctx, ctxKeys.UriKey, utils.GetAPIRoute(req))
		ctx = context.WithValue(ctx, ctxKeys.MethodKey, req.Method)
		ctx = context.WithValue(ctx, ctxKeys.PathParamsKey, utils.GetPathParams(req))
		ctx = context.WithValue(ctx, ctxKeys.QueryParamsKey, utils.GetQueryParams(req))
		// ctx = context.WithValue(ctx, ctxKeys.ClientIPKey, utils.GetClientIP(req))
		ctx = context.WithValue(ctx, ctxKeys.ClientTypeKey, utils.GetClientType(req))
		ctx = context.WithValue(ctx, ctxKeys.UserAgentKey, req.UserAgent())

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}
