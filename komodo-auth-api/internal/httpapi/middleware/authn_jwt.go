package middleware

import (
	"context"
	"komodo-auth-api/internal/crypto"
	"komodo-auth-api/internal/logger"
	"komodo-auth-api/internal/model"
	"net/http"
	"strings"
)

type authCtxKey string

const (
	AuthValidCtxKey    authCtxKey = "Authorization_valid"
	SessionValidCtxKey authCtxKey = "X-Session-Token_valid"
)

func AuthnJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		// Authorization: Bearer <token>
		if auth := req.Header.Get(model.HEADER_AUTH); auth != "" {
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				logger.Error("invalid Authorization header format", req)
				http.Error(wtr, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(parts[1])
			ok, err := crypto.VerifyJWT(token)

			if err != nil || !ok {
				logger.Error("invalid bearer token: "+err.Error(), req)
				http.Error(wtr, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, AuthValidCtxKey, true)
		}

		// X-Session-Token: optional header for session tokens
		if sess := req.Header.Get(model.HEADER_X_SESSION); sess != "" {
			ok, err := crypto.VerifyJWT(sess)
			if err != nil || !ok {
				logger.Error("invalid session token: "+err.Error(), req)
				http.Error(wtr, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, SessionValidCtxKey, true)
		}

		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}