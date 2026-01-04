package handlers

import (
	"encoding/json"
	"net/http"

	httpErr "komodo-forge-apis-go/http/errors"
	logger "komodo-forge-apis-go/logging/runtime"
)

// Handles OAuth 2.0 authorization endpoint (RFC 6749 Section 3.1).
// Authorizes client applications and issues authorization codes
func OAuthAuthorizeHandler(wtr http.ResponseWriter, req *http.Request) {
	// Parse query parameters
	query := req.URL.Query()
	responseType := query.Get("response_type")
	clientID := query.Get("client_id")
	redirectURI := query.Get("redirect_uri")
	scope := query.Get("scope")
	state := query.Get("state")

	// Validate required parameters
	if responseType == "" || clientID == "" || redirectURI == "" {
		logger.Error("missing required oauth parameters")
		httpErr.SendError(wtr, req, httpErr.Auth.AccessDenied, httpErr.WithDetail("missing required oauth parameters"))
		return
	}

	// Only support "code" response type for now
	if responseType != "code" {
		logger.Error("unsupported oauth response type: " + responseType)
		httpErr.SendError(wtr, req, httpErr.Auth.AccessDenied, httpErr.WithDetail("unsupported oauth response type"))
		return
	}

	// TODO: Implement authorization code flow
	// 1. Validate client_id against database/client registry
	// 2. Validate redirect_uri is registered for this client
	// 3. Check if user is authenticated (session/cookie)
	//    - If not authenticated: redirect to login page with return URL
	// 4. Show consent screen (if needed) asking user to approve scopes
	// 5. Generate authorization code (short-lived, single-use)
	// 6. Store code with client_id, redirect_uri, scope, user_id in cache
	// 7. Redirect back to redirect_uri with code and state:
	//    redirect_uri?code=<authorization_code>&state=<state>

	logger.Info("authorization endpoint called",
		"client_id", clientID,
		"redirect_uri", redirectURI,
		"scope", scope,
		"state", state,
	)

	wtr.Header().Set("Content-Type", "application/json")
	wtr.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(wtr).Encode(map[string]string{
		"error":             "not_implemented",
		"error_description": "Authorization code flow requires login UI implementation",
		"client_id":         clientID,
		"redirect_uri":      redirectURI,
	})
}
