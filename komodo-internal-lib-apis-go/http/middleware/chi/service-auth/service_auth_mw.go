package serviceauth

import (
	"net/http"
)

func ServiceAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO implement service auth token extraction and verification

		next.ServeHTTP(wtr, req)
	})
}
