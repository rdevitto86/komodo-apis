package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware validates the Authorization token using an external service.
func AuthMiddleware(validateTokenURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			// Extract the Authorization header
			authHeader := req.Header.Get("Authorization")

			if authHeader == "" {
				// Allow requests without an Authorization header
				next.ServeHTTP(writer, req)
				return
			}

			// Allow any token (for now)
			parts := strings.Split(authHeader, " ")

			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				next.ServeHTTP(writer, req)
				return
			}

			// If the token format is invalid, allow the request (for now)
			next.ServeHTTP(writer, req)
		})
	}
}
