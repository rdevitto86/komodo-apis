package gin

import "net/http"

func TelemetryMiddleware(nextFunc func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return nextFunc(next)
	}
}
