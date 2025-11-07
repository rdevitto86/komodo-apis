package csrf

import (
	"context"
	ctxKeys "komodo-internal-lib-apis-go/common/context"
	hdrTypes "komodo-internal-lib-apis-go/common/http"
	httpUtils "komodo-internal-lib-apis-go/http/utils/http"
	hdrSrv "komodo-internal-lib-apis-go/services/headers/eval"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	"net/http"
)

type csrfCtxKey string

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				clientType := req.Context().Value(ctxKeys.ClientTypeKey)
				if clientType == nil {
					clientType = httpUtils.GetClientType(req)
				}

				if clientType == "api" {
					ctx := context.WithValue(req.Context(), csrfCtxKey("X-CSRF"), "api-client-exempt")
					ctx = context.WithValue(ctx, csrfCtxKey("X-CSRF_valid"), true)
					req = req.WithContext(ctx)
					next.ServeHTTP(wtr, req)
					return
				}
				
				// Browser client: Require CSRF token
				if !hdrSrv.ValidateHeaderValue(hdrTypes.HEADER_X_CSRF_TOKEN, req) {
					logger.Error("invalid or missing CSRF token for browser client")
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
