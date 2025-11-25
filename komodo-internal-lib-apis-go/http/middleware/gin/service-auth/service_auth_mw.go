package serviceauth

import (
	"context"
	"net/http"
	"strings"

	"komodo-internal-lib-apis-go/crypto/jwt"
	authServ "komodo-internal-lib-apis-go/domains/auth/service"
	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	logger "komodo-internal-lib-apis-go/logging/runtime"

	"github.com/gin-gonic/gin"
)

// ServiceAuthMiddleware validates service tokens (JWT or OAuth)
// Strategy 1: Try JWT verification (fast path for internal services)
// Strategy 2: Fall back to OAuth introspection (for external partners)
func ServiceAuthMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		authHeader := gctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn("missing or invalid authorization header")
			gctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			gctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Strategy 1: Try JWT verification first (fast path for internal services)
		jwtToken, claims, err := jwt.VerifyToken(token)
		if err == nil && jwtToken.Valid {
			// Check if it's a service token (has "service" or specific audience)
			if aud, ok := claims["aud"].(string); ok && strings.Contains(aud, "service") {
				// Extract service identity and scopes
				clientID, _ := claims["sub"].(string)
				scope, _ := claims["scope"].(string)

				// Set context values for downstream handlers
				gctx.Set("client_id", clientID)
				gctx.Set("scope", scope)
				gctx.Set("auth_method", "jwt")

				gctx.Next()
				return
			}
		}

		// Create a fake http.Request for the service call (TokenVerify expects *http.Request)
		req, _ := http.NewRequestWithContext(context.Background(), "POST", "", nil)
		req.Header.Set("Authorization", "Bearer " + token)

		res := authServ.TokenVerify(req)
		if res.IsError() {
			logger.Warn("OAuth token verification failed", res.Error)
			gctx.JSON(res.Status, gin.H{
				"error": "invalid or expired token",
				"error_code": errCodes.ERR_INVALID_TOKEN,
			})
			gctx.Abort()
			return
		}

		// Parse OAuth verification response
		verified, ok := res.BodyParsed.(*authServ.TokenVerifyResponse)
		if !ok || verified == nil || !verified.Active {
			logger.Warn("OAuth token not active")
			gctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
				"error_code": errCodes.ERR_INVALID_TOKEN,
			})
			gctx.Abort()
			return
		}

		// Set context values for downstream handlers
		gctx.Set("client_id", verified.ClientID)
		gctx.Set("scope", verified.Scope)
		gctx.Set("auth_method", "oauth")

		gctx.Next()
	}
}
