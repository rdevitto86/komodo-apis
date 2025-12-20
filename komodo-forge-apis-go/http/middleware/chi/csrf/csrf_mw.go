package csrf

import (
	"context"
	ctxKeys "komodo-forge-apis-go/http/common/context"
	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/chi"
	hdrSrv "komodo-forge-apis-go/http/headers/eval"
	hdrTypes "komodo-forge-apis-go/http/types"
	httpUtils "komodo-forge-apis-go/http/utils/http"
	logger "komodo-forge-apis-go/logging/runtime"
	"net/http"
)

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				clientType := req.Context().Value(ctxKeys.CLIENT_TYPE_KEY)
				if clientType == nil {
					clientType = httpUtils.GetClientType(req)
				}

				if clientType == "api" {
					ctx := context.WithValue(req.Context(), ctxKeys.CSRF_TOKEN_KEY, "api-client-exempt")
					ctx = context.WithValue(ctx, ctxKeys.CSRF_VALID_KEY, true)
					req = req.WithContext(ctx)
					next.ServeHTTP(wtr, req)
					return
				}
				
				// Browser client - require CSRF token
				if ok, err := hdrSrv.ValidateHeaderValue(hdrTypes.HEADER_X_CSRF_TOKEN, req); !ok || err != nil {
					logger.Error("invalid or missing CSRF token for browser client", err)
					errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, errCodes.ERR_INVALID_REQUEST, "invalid CSRF token")
					return
				}
		}

		ctx := context.WithValue(req.Context(), ctxKeys.CSRF_TOKEN_KEY, "")
		ctx = context.WithValue(ctx, ctxKeys.CSRF_VALID_KEY, true)
		req = req.WithContext(ctx)

		next.ServeHTTP(wtr, req)
	})
}
