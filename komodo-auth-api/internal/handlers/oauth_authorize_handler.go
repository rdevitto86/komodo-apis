package handlers

import (
	"encoding/json"
	"net/http"

	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/chi"
	logger "komodo-forge-apis-go/loggers/runtime"
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
		logger.Error("missing required parameters")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_request: missing required parameters (response_type, client_id, redirect_uri)",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Only support "code" response type for now
	if responseType != "code" {
		logger.Error("unsupported response_type: " + responseType)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"unsupported_response_type",
			errCodes.ERR_INVALID_REQUEST,
		)
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
