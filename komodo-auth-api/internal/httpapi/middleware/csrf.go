package middleware

import (
	"net/http"
)

type csrfCtxKey string
const CSRFValidCtxKey csrfCtxKey = "X-CSRF_valid"

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO: Implement CSRF protection
		next.ServeHTTP(wtr, req)
	})
}
