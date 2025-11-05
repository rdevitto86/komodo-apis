package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	errUtils "komodo-internal-lib-apis-go/http/utils/error"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"

	jwtUtils "komodo-auth-api/internal/httpapi/utils/jwt"
)

type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
}

type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	IssuedAt     int64  `json:"issued_at"`
}

// Handles token refresh requests
func TokenRefreshHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Parse request body
	var reqBody TokenRefreshRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("Failed to parse request body", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Invalid request body",
			"40006",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Validate refresh token present
	if reqBody.RefreshToken == "" {
		logger.Error("Missing refresh_token")
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Missing required field: refresh_token",
			"40002",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Validate grant_type (must be refresh_token)
	if reqBody.GrantType != "" && reqBody.GrantType != "refresh_token" {
		logger.Error("Invalid grant_type for refresh: " + reqBody.GrantType)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"grant_type must be 'refresh_token'",
			"90001",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Verify refresh token signature and parse claims
	_, claims, err := jwtUtils.VerifyToken(reqBody.RefreshToken)
	if err != nil {
		logger.Error("Invalid refresh token", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusUnauthorized,
			"Invalid or expired refresh token",
			"20004",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Extract client_id and scope from original token
	clientID, _ := jwtUtils.ExtractStringClaim(claims, "client_id")
	scope, _ := jwtUtils.ExtractStringClaim(claims, "scope")
	tokenUse, _ := jwtUtils.ExtractStringClaim(claims, "token_use")

	if clientID == "" {
		logger.Error("No client_id in refresh token")
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusUnauthorized,
			"Invalid token claims",
			"20004",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Validate token type (must be a refresh token)
	if tokenUse != "" && tokenUse != "refresh" {
		logger.Error("Token is not a refresh token, type: " + tokenUse)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Invalid token type - must be refresh token",
			"20004",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// TODO: Check if refresh token is revoked in Elasticache/database
	// oldJTI, _ := jwtUtils.ExtractStringClaim(claims, "jti")
	// if oldJTI != "" && redisClient.Exists("revoked:token:" + oldJTI) {
	//     logger.Error("Refresh token has been revoked: " + oldJTI)
	//     errUtils.WriteErrorResponse(wtr, http.StatusUnauthorized, "Token has been revoked", "20004", req.Header.Get("X-Request-ID"))
	//     return
	// }

	// Determine access token TTL with bounds checking
	now := time.Now()
	expiresIn := jwtUtils.DefaultTokenTTL
	if reqBody.ExpiresIn > 0 {
		expiresIn = jwtUtils.ClampTTL(reqBody.ExpiresIn, clientID)
	}

	// Create new access token with JTI
	newClaims := jwtUtils.CreateStandardClaims(
		"komodo-auth-api",
		clientID,
		"komodo-apis",
		int64(expiresIn),
		map[string]interface{}{
			"scope":      scope,
			"grant_type": "refresh_token",
			"client_id":  clientID,
			"token_use":  "access",
		},
	)

	// Sign new access token
	accessToken, err := jwtUtils.SignToken(newClaims)
	if err != nil {
		logger.Error("Failed to sign new access token", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusInternalServerError,
			"Token generation failed",
			"20006",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Generate new refresh token (token rotation for security)
	refreshExpiresIn := int64(604800) // 7 days
	refreshClaims := jwtUtils.CreateStandardClaims(
		"komodo-auth-api",
		clientID,
		"komodo-apis",
		refreshExpiresIn,
		map[string]interface{}{
			"scope":      scope,
			"client_id":  clientID,
			"token_use":  "refresh",
		},
	)

	newRefreshToken, err := jwtUtils.SignToken(refreshClaims)
	if err != nil {
		logger.Error("Failed to sign new refresh token", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusInternalServerError,
			"Token generation failed",
			"20006",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// TODO: Revoke old refresh token in Elasticache
	// if oldJTI != "" {
	//     ttl := jwtUtils.GetTokenTTL(claims)
	//     redisClient.Set("revoked:token:" + oldJTI, clientID, time.Duration(ttl)*time.Second)
	// }

	// TODO: Store new tokens in Elasticache with TTL for tracking
	// newAccessJTI, _ := jwtUtils.ExtractStringClaim(newClaims, "jti")
	// redisClient.Set("token:" + newAccessJTI, clientID, time.Duration(expiresIn)*time.Second)

	logger.Info("Token refreshed successfully for client: " + clientID)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenRefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken, // Return new refresh token (rotation)
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		Scope:        scope,
		IssuedAt:     now.Unix(),
	})
}