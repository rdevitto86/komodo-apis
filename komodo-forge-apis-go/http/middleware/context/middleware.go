package context

import (
	"context"
	ctxKeys "komodo-forge-apis-go/http/common/context"
	reqUtils "komodo-forge-apis-go/http/utils/request"
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
			reqID = reqUtils.GenerateRequestId()
		}
		ctx = context.WithValue(ctx, chimw.RequestIDKey, reqID)
		ctx = context.WithValue(ctx, ctxKeys.REQUEST_ID_KEY, reqID)
		req.Header.Set("X-Request-ID", reqID)
		wtr.Header().Set("X-Request-ID", reqID)

		ctx = context.WithValue(ctx, ctxKeys.START_TIME_KEY, time.Now().UTC())
		ctx = context.WithValue(ctx, ctxKeys.VERSION_KEY, reqUtils.GetAPIVersion(req))
		ctx = context.WithValue(ctx, ctxKeys.URI_KEY, reqUtils.GetAPIRoute(req))
		ctx = context.WithValue(ctx, ctxKeys.METHOD_KEY, req.Method)
		ctx = context.WithValue(ctx, ctxKeys.PATH_PARAMS_KEY, reqUtils.GetPathParams(req))
		ctx = context.WithValue(ctx, ctxKeys.QUERY_PARAMS_KEY, reqUtils.GetQueryParams(req))
		// ctx = context.WithValue(ctx, ctxKeys.CLIENT_IP_KEY, utils.GetClientIP(req))
		ctx = context.WithValue(ctx, ctxKeys.CLIENT_TYPE_KEY, reqUtils.GetClientType(req))
		ctx = context.WithValue(ctx, ctxKeys.USER_AGENT_KEY, req.UserAgent())

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}
