package middleware

import (
	"net/http"
)

func ResponsePreprocessorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		wtr.Header().Set("Content-Type", "application/json") // default type - override in handlers
		wtr.Header().Set("Cache-Control", "no-store")
		wtr.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(wtr, req)
	})
}
