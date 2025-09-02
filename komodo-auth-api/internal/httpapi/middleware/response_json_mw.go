package middleware

import (
	"net/http"
)

func ResponseJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		wtr.Header().Set("Content-Type", "application/json")
		wtr.Header().Set("Cache-Control", "no-store")
		wtr.Header().Set("X-Content-Type-Options", "nosniff")
		wtr.Header().Set("Access-Control-Allow-Origin", "*") // CORS
		next.ServeHTTP(wtr, req)
	})
}
