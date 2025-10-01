package redaction

import "net/http"

func RedactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO
		next.ServeHTTP(wtr, req)
	})
}