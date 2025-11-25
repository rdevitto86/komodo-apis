package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"komodo-internal-lib-apis-go/crypto/oauth"
	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	errors "komodo-internal-lib-apis-go/http/common/errors/chi"
	logger "komodo-internal-lib-apis-go/logging/runtime"

	jwt "komodo-internal-lib-apis-go/crypto/jwt"
)

type TokenCreateRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

type TokenCreateResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	IssuedAt     int64  `json:"issued_at"`
}

// Handles token creation requests
func TokenCreateHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// Parse request body
	var reqBody TokenCreateRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse request body", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"failed to parse request body",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Validate client credentials
	if reqBody.ClientID == "" || reqBody.ClientSecret == "" {
		logger.Error("missing client credentials")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"missing client credentials",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// TODO: Validate client_id and client_secret against database
	// Placeholder validation for now
	if reqBody.ClientID != "test-client" || reqBody.ClientSecret != "test-secret" {
		logger.Error("invalid client credentials")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusUnauthorized,
			"invalid client credentials",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Validate grant_type
	if !oauth.IsValidGrantType(reqBody.GrantType) {
		logger.Error("unsupported grant_type: " + reqBody.GrantType)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"unsupported grant_type",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Validate scope
	if !oauth.IsValidScope(reqBody.Scope) {
		logger.Error("invalid scope: " + reqBody.Scope)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid scope",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Token TTL
	now := time.Now()
	accessExpiresIn := int64(jwt.DefaultTokenTTL)

	// Determine client type based on grant_type
	clientType := "browser"
	if reqBody.GrantType == "client_credentials" {
		clientType = "api"
	}

	// Create access token with JTI
	accessClaims := jwt.CreateStandardClaims(
		"komodo-auth-service-api",
		reqBody.ClientID,
		"komodo-apis",
		accessExpiresIn,
		map[string]interface{}{
			"scope":       reqBody.Scope,
			"grant_type":  reqBody.GrantType,
			"client_id":   reqBody.ClientID,
			"token_use":   "access",
			"client_type": clientType,
		},
	)

	// Sign access token
	accessToken, err := jwt.SignToken(accessClaims)
	if err != nil {
		logger.Error("failed to sign access token", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusInternalServerError,
			"Token signing failed",
			errCodes.ERR_INTERNAL_SERVER,
		)
		return
	}

	// Don't generate refresh token for client_credentials (machine-to-machine)
	// Services can just request new tokens with their credentials
	var refreshToken string
	if reqBody.GrantType != "client_credentials" {
		// Determine refresh token TTL with bounds checking (default 7 days)
		refreshExpiresIn := 604800 // 7 days
		if reqBody.ExpiresIn > 0 {
			refreshExpiresIn = jwt.ClampTTL(reqBody.ExpiresIn, reqBody.ClientID)
		}

		// Generate refresh token only for user flows
		refreshClaims := jwt.CreateStandardClaims(
			"komodo-auth-service-api",
			reqBody.ClientID,
			"komodo-apis",
			int64(refreshExpiresIn),
			map[string]interface{}{
				"scope":       reqBody.Scope,
				"client_id":   reqBody.ClientID,
				"token_use":   "refresh",
				"client_type": clientType,
				"grant_type":  reqBody.GrantType, // Include original grant_type
			},
		)

		refreshToken, err = jwt.SignToken(refreshClaims)
		if err != nil {
			logger.Error("failed to sign refresh token", err)
			errors.WriteErrorResponse(
				wtr,
				req,
				http.StatusInternalServerError,
				"Token signing failed",
				errCodes.ERR_INTERNAL_SERVER,
			)
			return
		}
	}

	// TODO: Store tokens in Elasticache/Redis with TTL for tracking and revocation
	// accessJTI, _ := jwtUtils.ExtractStringClaim(accessClaims, "jti")
	// redisClient.Set("token:" + accessJTI, reqBody.ClientID, time.Duration(accessExpiresIn)*time.Second)

	response := TokenCreateResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   int(accessExpiresIn),
		Scope:       reqBody.Scope,
		IssuedAt:    now.Unix(),
	}

	// Only include refresh_token if generated (not for client_credentials)
	if refreshToken != "" {
		response.RefreshToken = refreshToken
	}

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(response)
}
