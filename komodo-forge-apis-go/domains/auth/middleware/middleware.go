package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO: implement JWT handling
		next.ServeHTTP(wtr, req)
	})
}
