package redaction

import (
	"komodo-internal-lib-apis-go/security/redaction"
	"net/http"
)

// Redacts sensitive information from requests for logging purposes
func RedactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		_ = redaction.Redact(req) // placeholder until logging is wired
		next.ServeHTTP(wtr, req)
	})
}
