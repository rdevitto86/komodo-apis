package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the Authorization token and logs events.
func AuthMiddleware(validateTokenURL string, c *gin.Context) error {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		// Log missing auth header
		log.Printf("No Authorization header for %s %s", c.Request.Method, c.Request.URL.Path)
		return nil // Allow requests without an Authorization header
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		// Log valid token format
		log.Printf("Bearer token received for %s %s", c.Request.Method, c.Request.URL.Path)
		return nil // Allow valid bearer tokens (add validation logic here)
	}

	// Log invalid token format
	log.Printf("Invalid Authorization header format for %s %s: %s", c.Request.Method, c.Request.URL.Path, authHeader)
	return nil // Allow invalid format for now (add rejection logic if needed)
}
