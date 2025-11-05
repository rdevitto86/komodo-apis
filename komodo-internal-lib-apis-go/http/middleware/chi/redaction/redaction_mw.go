package redaction

import (
	rsvc "komodo-internal-lib-apis-go/services/redaction"
	"net/http"
)

func RedactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		_ = rsvc.RedactForLogging(req) // placeholder until logging is wired

		next.ServeHTTP(wtr, req)
	})
}
