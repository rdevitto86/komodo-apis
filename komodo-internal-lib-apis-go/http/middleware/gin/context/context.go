package gin

import "net/http"

func ContextMiddleware(nextFunc func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return nextFunc
}