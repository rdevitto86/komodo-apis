package session

import (
	"context"
	"net/http"
	"strings"

	"komodo-internal-lib-apis-go/crypto/jwt"
	ctxKeys "komodo-internal-lib-apis-go/http/common/context"
	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	errors "komodo-internal-lib-apis-go/http/common/errors/chi"
	logger "komodo-internal-lib-apis-go/logging/runtime"
)

// Extracts session information from incoming requests
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		var (
			userID    string
			sessionID string
		)

		// Extract session_id from Cookie header
		cookies := req.Header.Get("Cookie")
		if cookies != "" {
			// Parse session_id from cookie string
			for _, cookie := range strings.Split(cookies, "; ") {
				parts := strings.SplitN(cookie, "=", 2)
				if len(parts) == 2 && parts[0] == "session_id" {
					sessionID = parts[1]
					break
				}
			}

			if sessionID != "" {
				logger.Debug("session_id extracted from cookie: " + sessionID)

				// TODO: Look up session in Redis
				// redisKey := "session:user:" + sessionID
				// userID = redis.Get(redisKey)
				// if userID == "" {
				//     logger.Warn("session not found in Redis: " + sessionID)
				//     errors.WriteErrorResponse(wtr, req, http.StatusUnauthorized, "invalid or expired session", errors.ERR_SESSION_EXPIRED)
				//     return
				// }

				// TODO: Return test user_id for now
				userID = "12345"
				logger.Info("mock user_id extracted from session: " + userID)

				ctx := context.WithValue(req.Context(), ctxKeys.USER_ID_KEY, userID)
				next.ServeHTTP(wtr, req.WithContext(ctx))
				return
			}
		}

		// Fallback - Try to extract from Authorization Bearer token
		authHeader := req.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			
			// Verify and parse the token
			token, claims, err := jwt.VerifyToken(tokenString)
			if err == nil && token.Valid {
				// Extract user_id from token claims (subject)
				if sub, exists := claims["sub"].(string); exists && sub != "" {
					userID = sub
					ctx := context.WithValue(req.Context(), ctxKeys.USER_ID_KEY, userID)
					next.ServeHTTP(wtr, req.WithContext(ctx))
					return
				}
			} else {
				logger.Warn("invalid or expired JWT Bearer token", err)
			}
		}

		// No valid session or token found
		logger.Warn("no valid session or token found in request")
		errors.WriteErrorResponse(wtr, req, http.StatusUnauthorized, "unauthorized - no valid session", errCodes.ERR_SESSION_NOT_FOUND)
	})
}
