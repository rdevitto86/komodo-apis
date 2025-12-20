package auth

import (
	"context"
	"net/http"
	"strings"

	"komodo-forge-apis-go/crypto/jwt"
	authServ "komodo-forge-apis-go/domains/auth/service"
	errCodes "komodo-forge-apis-go/http/common/errors"
	logger "komodo-forge-apis-go/logging/runtime"

	"github.com/gin-gonic/gin"
)

// Validates service tokens (JWT)
func AuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		authHeader := g.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Warn("missing or invalid authorization header")
			g.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			g.Abort()
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
				g.Set("client_id", clientID)
				g.Set("scope", scope)
				g.Set("auth_method", "jwt")

				g.Next()
				return
			}
		}

		// Create a fake http.Request for the service call (TokenVerify expects *http.Request)
		req, _ := http.NewRequestWithContext(context.Background(), "POST", "", nil)
		req.Header.Set("Authorization", "Bearer " + token)

		res := authServ.TokenVerify(req)
		if res.IsError() {
			logger.Warn("OAuth token verification failed", res.Error)
			g.JSON(res.Status, gin.H{
				"error": "invalid or expired token",
				"error_code": errCodes.ERR_INVALID_TOKEN,
			})
			g.Abort()
			return
		}

		// Parse OAuth verification response
		verified, ok := res.BodyParsed.(*authServ.TokenVerifyResponse)
		if !ok || verified == nil || !verified.Active {
			logger.Warn("OAuth token not active")
			g.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
				"error_code": errCodes.ERR_INVALID_TOKEN,
			})
			g.Abort()
			return
		}

		// Set context values for downstream handlers
		g.Set("client_id", verified.ClientID)
		g.Set("scope", verified.Scope)
		g.Set("auth_method", "oauth")

		g.Next()
	}
}
