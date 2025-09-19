package middleware

import (
	"net/http"
)

func IPAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO: Implement IP checking logic (filter + client)
		next.ServeHTTP(wtr, req)
	})
}
