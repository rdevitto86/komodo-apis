package redaction

import (
	rsvc "komodo-internal-lib-apis-go/services/redaction"
	"net/http"
)

func RedactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		_ = rsvc.RedactForLogging(req) // placeholder until logging is wired

		// wrap response writer to scrub headers on the way out
		next.ServeHTTP(&sanitizeWriter{ ResponseWriter: wtr }, req)
	})
}

// sanitizeWriter strips sensitive response headers before they are sent.
type sanitizeWriter struct {
	http.ResponseWriter
}

func (s *sanitizeWriter) WriteHeader(code int) {
	// remove sensitive headers
	s.Header().Del("Set-Cookie")
	s.Header().Del("Authorization")
	s.ResponseWriter.WriteHeader(code)
}
