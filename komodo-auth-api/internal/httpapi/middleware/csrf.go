package middleware

import (
	"context"
	"net/http"

	"komodo-auth-api/internal/httpapi/utils"
	"komodo-auth-api/internal/logger"
)

type csrfCtxKey string

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Default to false
		isValid := false

		// For safe methods we consider CSRF not required and mark as valid
		switch req.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				isValid = true
			default:
				isValid = utils.IsValidCSRF(req.Header.Get("X-CSRF"), req.Header.Get("X-Session-Token"))

				if !isValid {
					logger.Error("invalid or missing CSRF token", req)
					http.Error(wtr, "Invalid CSRF token", http.StatusBadRequest)
					return
				}
		}

		ctx := context.WithValue(req.Context(), csrfCtxKey("X-CSRF"), isValid)
		ctx = context.WithValue(ctx, csrfCtxKey("X-CSRF_valid"), isValid)
		req = req.WithContext(ctx)

		next.ServeHTTP(wtr, req)
	})
}
