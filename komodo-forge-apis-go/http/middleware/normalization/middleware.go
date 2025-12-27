package normalization

import "net/http"

func NormalizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Normalization logic can be added here

		next.ServeHTTP(wtr, req)
	})
}