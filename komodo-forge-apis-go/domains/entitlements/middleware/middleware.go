package middleware

import "net/http"

func EntitlementsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Check user entitlements here
		next.ServeHTTP(wtr, req)
	})
}