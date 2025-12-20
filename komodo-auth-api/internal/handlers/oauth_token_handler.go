package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-forge-apis-go/crypto/jwt"
	"komodo-forge-apis-go/crypto/oauth"
	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/chi"
	logger "komodo-forge-apis-go/logging/runtime"
)

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"` // For refresh_token grant
	Code         string `json:"code,omitempty"`          // For authorization_code grant
	RedirectURI  string `json:"redirect_uri,omitempty"`  // For authorization_code grant
	Username     string `json:"username,omitempty"`      // For password grant
	Password     string `json:"password,omitempty"`      // For password grant
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// Unified OAuth 2.0 token endpoint (RFC 6749 Section 3.2).
// Handles all grant types: client_credentials, refresh_token, authorization_code, password.
// All tokens issued are JWTs
func OAuthTokenHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	var reqBody TokenRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse request body", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_request",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	if reqBody.GrantType == "" {
		logger.Error("missing grant_type")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_request: missing grant_type",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}
	if !oauth.IsValidGrantType(reqBody.GrantType) {
		logger.Error("unsupported grant_type: " + reqBody.GrantType)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"unsupported_grant_type",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Route to appropriate grant handler
	switch reqBody.GrantType {
		case "client_credentials":
			handleClientCredentials(wtr, req, &reqBody)
		case "refresh_token":
			handleRefreshToken(wtr, req, &reqBody)
		case "authorization_code":
			handleAuthorizationCode(wtr, req, &reqBody)
		default:
			errors.WriteErrorResponse(
				wtr,
				req,
				http.StatusBadRequest,
				"unsupported_grant_type",
				errCodes.ERR_INVALID_REQUEST,
			)
	}
}

// Handles M2M service authentication (RFC 6749 Section 4.4)
func handleClientCredentials(wtr http.ResponseWriter, req *http.Request, reqBody *TokenRequest) {
	// Validate client credentials
	if reqBody.ClientID == "" || reqBody.ClientSecret == "" {
		logger.Error("missing client credentials")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_client",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// TODO: Validate client_id and client_secret against database/secrets store
	if reqBody.ClientID != "test-client" || reqBody.ClientSecret != "test-secret" {
		logger.Error("invalid client credentials")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusUnauthorized,
			"invalid_client",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}
	if reqBody.Scope != "" && !oauth.IsValidScope(reqBody.Scope) {
		logger.Error("invalid scope: " + reqBody.Scope)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_scope",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Issue access token (JWT) - no refresh token for client_credentials
	accessExpiresIn := int64(jwt.DefaultTokenTTL)

	accessClaims := jwt.CreateStandardClaims(
		"komodo-auth-api",
		reqBody.ClientID,
		"komodo-apis:service",
		accessExpiresIn,
		map[string]interface{}{
			"scope":       reqBody.Scope,
			"grant_type":  "client_credentials",
			"client_id":   reqBody.ClientID,
			"token_use":   "access",
			"client_type": "api",
		},
	)

	accessToken, err := jwt.SignToken(accessClaims)
	if err != nil {
		logger.Error("failed to sign access token", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusInternalServerError,
			"server_error",
			errCodes.ERR_INTERNAL_SERVER,
		)
		return
	}

	// TODO: Store token JTI in Elasticache for tracking/revocation

	response := TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessExpiresIn),
		Scope:       reqBody.Scope,
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(response)

	logger.Info("issued client_credentials token for: " + reqBody.ClientID)
}

// Handles token refresh (RFC 6749 Section 6)
func handleRefreshToken(wtr http.ResponseWriter, req *http.Request, reqBody *TokenRequest) {
	if reqBody.RefreshToken == "" {
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_request: missing refresh_token",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Verify refresh token
	_, claims, err := jwt.VerifyToken(reqBody.RefreshToken)
	if err != nil || jwt.IsTokenExpired(claims) {
		logger.Error("invalid or expired refresh token", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusUnauthorized,
			"invalid_grant",
			errCodes.ERR_INVALID_TOKEN,
		)
		return
	}

	// Ensure token is a refresh token
	tokenUse, _ := claims["token_use"].(string)
	if tokenUse != "refresh" {
		logger.Error("token is not a refresh token")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_grant",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// TODO: Check if refresh token is revoked in Elasticache

	// Extract claims
	claimValues := jwt.ExtractStringClaims(claims, []string{
		"client_id", "scope", "client_type",
	})
	clientID, _ := claimValues["client_id"].(string)
	scope, _ := claimValues["scope"].(string)

	// Issue new access token
	accessExpiresIn := int64(jwt.DefaultTokenTTL)

	accessClaims := jwt.CreateStandardClaims(
		"komodo-auth-api",
		clientID,
		"komodo-apis:user",
		accessExpiresIn,
		map[string]interface{}{
			"scope":       scope,
			"grant_type":  "refresh_token",
			"client_id":   clientID,
			"token_use":   "access",
			"client_type": "browser",
		},
	)

	accessToken, err := jwt.SignToken(accessClaims)
	if err != nil {
		logger.Error("failed to sign access token", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusInternalServerError,
			"server_error",
			errCodes.ERR_INTERNAL_SERVER,
		)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessExpiresIn),
		Scope:        scope,
		RefreshToken: reqBody.RefreshToken, // Can optionally rotate
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(response)

	logger.Info("refreshed token for: " + clientID)
}

// Handles authorization code exchange (RFC 6749 Section 4.1)
func handleAuthorizationCode(wtr http.ResponseWriter, req *http.Request, reqBody *TokenRequest) {
	// TODO: Implement authorization code flow
	// 1. Validate code against stored authorization grants
	// 2. Verify redirect_uri matches original request
	// 3. Verify client credentials
	// 4. Issue access + refresh tokens
	// 5. Delete used authorization code

	logger.Info("authorization_code grant not yet implemented")

	errors.WriteErrorResponse(
		wtr,
		req,
		http.StatusNotImplemented,
		"authorization_code grant not yet implemented",
		errCodes.ERR_INTERNAL_SERVER,
	)
}
