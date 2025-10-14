package logging

import "net/http"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Pre-processing logic can be added here

		next.ServeHTTP(wtr, req)

		// Post-processing logic can be added here
	})
}