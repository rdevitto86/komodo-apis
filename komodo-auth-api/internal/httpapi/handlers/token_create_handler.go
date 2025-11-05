package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	authUtils "komodo-internal-lib-apis-go/http/utils/auth"
	errUtils "komodo-internal-lib-apis-go/http/utils/error"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	errTypes "komodo-internal-lib-apis-go/types/error"

	jwtUtils "komodo-auth-api/internal/httpapi/utils/jwt"
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
		logger.Error("Failed to parse request body", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Failed to parse request body",
			errTypes.ERR_INVALID_REQUEST,
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Validate client credentials
	if reqBody.ClientID == "" || reqBody.ClientSecret == "" {
		logger.Error("Missing client credentials")
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Missing client credentials",
			errTypes.ERR_INVALID_REQUEST,
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// TODO: Validate client_id and client_secret against database
	// Placeholder validation for now
	if reqBody.ClientID != "test-client" || reqBody.ClientSecret != "test-secret" {
		logger.Error("Invalid client credentials")
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusUnauthorized,
			"Invalid client credentials",
			errTypes.ERR_INVALID_REQUEST,
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Validate grant_type
	if !authUtils.IsValidGrantType(reqBody.GrantType) {
		logger.Error("Unsupported grant_type: " + reqBody.GrantType)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Unsupported grant_type",
			errTypes.ERR_INVALID_REQUEST,
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Validate scope
	if !authUtils.IsValidScope(reqBody.Scope) {
		logger.Error("Invalid scope: " + reqBody.Scope)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Invalid scope",
			errTypes.ERR_INVALID_REQUEST,
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Use default TTL for access token (1 hour)
	now := time.Now()
	accessExpiresIn := int64(jwtUtils.DefaultTokenTTL)

	// Create access token with JTI
	accessClaims := jwtUtils.CreateStandardClaims(
		"komodo-auth-api",
		reqBody.ClientID,
		"komodo-apis",
		accessExpiresIn,
		map[string]interface{}{
			"scope":      reqBody.Scope,
			"grant_type": reqBody.GrantType,
			"client_id":  reqBody.ClientID,
			"token_use":  "access",
		},
	)

	// Sign access token
	accessToken, err := jwtUtils.SignToken(accessClaims)
	if err != nil {
		logger.Error("Failed to sign access token", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusInternalServerError,
			"Token signing failed",
			errTypes.ERR_INTERNAL_SERVER,
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Determine refresh token TTL with bounds checking (default 7 days)
	refreshExpiresIn := 604800 // 7 days default
	if reqBody.ExpiresIn > 0 {
		refreshExpiresIn = jwtUtils.ClampTTL(reqBody.ExpiresIn, reqBody.ClientID)
	}

	// Generate refresh token
	refreshClaims := jwtUtils.CreateStandardClaims(
		"komodo-auth-api",
		reqBody.ClientID,
		"komodo-apis",
		int64(refreshExpiresIn),
		map[string]interface{}{
			"scope":     reqBody.Scope,
			"client_id": reqBody.ClientID,
			"token_use": "refresh",
		},
	)

	refreshToken, err := jwtUtils.SignToken(refreshClaims)
	if err != nil {
		logger.Error("Failed to sign refresh token", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusInternalServerError,
			"Token signing failed",
			errTypes.ERR_INTERNAL_SERVER,
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// TODO: Store tokens in Elasticache/Redis with TTL for tracking and revocation
	// accessJTI, _ := jwtUtils.ExtractStringClaim(accessClaims, "jti")
	// redisClient.Set("token:" + accessJTI, reqBody.ClientID, time.Duration(accessExpiresIn)*time.Second)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenCreateResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessExpiresIn),
		Scope:        reqBody.Scope,
		IssuedAt:     now.Unix(),
	})
}
