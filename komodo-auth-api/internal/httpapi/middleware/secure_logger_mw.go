package middleware

import (
	"fmt"
	"komodo-auth-api/internal/thirdparty/grafana"
	"net/http"
)

func SecureLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		grafana.LogToLoki(fmt.Sprintf("Received request: %s %s", req.Method, req.URL.Path))
		next.ServeHTTP(wtr, req)
	})
}