package sanitization

import (
	"net/http"
)

func SanitizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Sanitization logic can be added here

		next.ServeHTTP(wtr, req)
	})
}