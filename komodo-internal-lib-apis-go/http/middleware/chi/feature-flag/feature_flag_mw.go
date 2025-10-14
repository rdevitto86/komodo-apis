package featureflags

import "net/http"

func FeatureFlagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Feature flag logic can be added here

		next.ServeHTTP(wtr, req)
	})
}