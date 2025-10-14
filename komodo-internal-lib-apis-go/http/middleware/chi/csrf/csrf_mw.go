package csrf

import (
	"context"
	hdrSrv "komodo-internal-lib-apis-go/services/headers/eval"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	hdrTypes "komodo-internal-lib-apis-go/types/headers"
	"net/http"
)

type csrfCtxKey string

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				if !hdrSrv.ValidateHeaderValue(hdrTypes.HEADER_X_CSRF_TOKEN, req) {
					logger.Error("invalid or missing CSRF token", req)
					http.Error(wtr, "Invalid CSRF token", http.StatusBadRequest)
					return
				}
		}

		ctx := context.WithValue(req.Context(), csrfCtxKey("X-CSRF"), "")
		ctx = context.WithValue(ctx, csrfCtxKey("X-CSRF_valid"), true)
		req = req.WithContext(ctx)

		next.ServeHTTP(wtr, req)
	})
}
