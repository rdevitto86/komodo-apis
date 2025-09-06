package middleware

import (
	"context"
	"net/http"
)

type ctxKey string

const (
	authTokenKey    ctxKey = "auth_token"
	sessionTokenKey ctxKey = "session_token"
	csrfTokenKey    ctxKey = "csrf_token"
)

func ObscureDataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Authorization") != "" {
			ctx := context.WithValue(req.Context(), authTokenKey, req.Header.Get("Authorization"))
			req = req.WithContext(ctx)
			req.Header.Set("Authorization", "Bearer [REDACTED]")
		}
		if req.Header.Get("X-Session-Token") != "" {
			ctx := context.WithValue(req.Context(), sessionTokenKey, req.Header.Get("X-Session-Token"))
			req = req.WithContext(ctx)
			req.Header.Set("X-Session-Token", "[REDACTED]")
		}
		if req.Header.Get("X-CSRF") != "" {
			ctx := context.WithValue(req.Context(), csrfTokenKey, req.Header.Get("X-CSRF"))
			req = req.WithContext(ctx)
			req.Header.Set("X-CSRF", "[REDACTED]")
		}

		next.ServeHTTP(wtr, req)
	})
}
