package requestid

import (
	"context"
	utils "komodo-internal-lib-apis-go/http/utils/http"
	ctxKeys "komodo-internal-lib-apis-go/types/context"
	"net/http"

	chimw "github.com/go-chi/chi/v5/middleware"
)

// RequestIDMiddleware ensures each request has a unique X-Request-ID in both header and context
// Priority: Header (external) > Context (middleware) > Generated (new)
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		var reqID string
		if rid := req.Header.Get("X-Request-ID"); rid != "" {
			reqID = rid
		} else if rid := chimw.GetReqID(req.Context()); rid != "" {
			reqID = rid
		} else {
			reqID = utils.GenerateRequestId()
		}

		req.Header.Set("X-Request-ID", reqID)
		ctx := context.WithValue(req.Context(), chimw.RequestIDKey, reqID)
		ctx = context.WithValue(ctx, ctxKeys.RequestIDKey, reqID)
		wtr.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}