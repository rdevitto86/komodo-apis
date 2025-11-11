package session

import "net/http"

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(wtr, req)
	})
}
