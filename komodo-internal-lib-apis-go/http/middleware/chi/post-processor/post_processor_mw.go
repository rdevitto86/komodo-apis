package postprocessor

import "net/http"

func PostProcessorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Post-processing logic can be added here

		next.ServeHTTP(wtr, req)
	})
}