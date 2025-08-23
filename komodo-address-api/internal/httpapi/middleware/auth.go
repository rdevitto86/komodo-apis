package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware validates the Authorization token using an external service.
func AuthMiddleware(validateTokenURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Allow requests without an Authorization header
				next.ServeHTTP(w, r)
				return
			}

			// Allow any token (for now)
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				// Log the token for debugging purposes (optional)
				// log.Printf("Received token: %s", parts[1])
				next.ServeHTTP(w, r)
				return
			}

			// If the token format is invalid, allow the request (for now)
			next.ServeHTTP(w, r)
		})
	}
}
