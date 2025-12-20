package clienttype

import (
	"context"
	ctxKeys "komodo-forge-apis-go/http/common/context"
	"net/http"
)

const (
	ClientTypeAPI string = "api"
	ClientTypeBrowser string = "browser"
)

// ClientTypeMiddleware detects whether request is from API client or browser
// and stores the result in context for downstream middleware to use
func ClientTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Detect client type based on request characteristics
		authHeader := req.Header.Get("Authorization")
		hasReferer := req.Header.Get("Referer") != ""
		hasCookie := req.Header.Get("Cookie") != ""
		
		// API client: Has Bearer token but no browser-specific headers
		clientType := ClientTypeBrowser
		if authHeader != "" && !hasReferer && !hasCookie {
			clientType = ClientTypeAPI
		}
		
		// Store in context for downstream middleware
		ctx := context.WithValue(req.Context(), ctxKeys.CLIENT_TYPE_KEY, clientType)
		next.ServeHTTP(wtr, req.WithContext(ctx))
	})
}
