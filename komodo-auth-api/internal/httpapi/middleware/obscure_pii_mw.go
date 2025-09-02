package middleware

import "net/http"

func ObscurePIIMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		req.Header.Set("Authorization", "[REDACTED]")
		next.ServeHTTP(wtr, req)
	})
}
