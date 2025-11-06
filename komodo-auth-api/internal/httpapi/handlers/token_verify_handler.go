package handlers

import (
	"encoding/json"
	"net/http"

	errUtils "komodo-internal-lib-apis-go/http/utils/error"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"

	jwtUtils "komodo-auth-api/internal/httpapi/utils/jwt"
)

type TokenVerifyResponse struct {
	Valid     bool   `json:"valid"`
	ClientID  string `json:"client_id,omitempty"`
	Scope     string `json:"scope,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
	IssuedAt  int64  `json:"issued_at,omitempty"`
	ClientType string `json:"client_type,omitempty"`
}

// Handles token verification requests
func TokenVerifyHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header only
	token, err := jwtUtils.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("Failed to extract token from Authorization header", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Missing or invalid Authorization header",
			"40002",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Verify token signature and parse claims
	_, claims, err := jwtUtils.VerifyToken(token)
	if err != nil {
		logger.Error("Invalid token", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusUnauthorized,
			"Invalid or expired token",
			"20004",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Extract token details
	claimValues := jwtUtils.ExtractStringClaims(claims, []string{"client_id", "scope", "token_use", "exp", "iat", "client_type"})
	clientID, _ := claimValues["client_id"].(string)
	scope, _ := claimValues["scope"].(string)
	tokenUse, _ := claimValues["token_use"].(string)
	expiresAt, _ := claimValues["exp"].(int64)
	issuedAt, _ := claimValues["iat"].(int64)
	clientType, _ := claimValues["client_type"].(string)

	// Verify token is not expired
	if jwtUtils.IsTokenExpired(claims) {
		logger.Error("Token expired")
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusUnauthorized,
			"Token has expired",
			"20003",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Validate token type (should be access token for API calls)
	if tokenUse != "" && tokenUse != "access" {
		logger.Error("Token is not an access token, type: " + tokenUse)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Invalid token type - must be access token",
			"20004",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// TODO: Check if token is revoked in Elasticache/database
	// jti, _ := jwtUtils.ExtractStringClaim(claims, "jti")
	// if jti != "" && redisClient.Exists("revoked:token:" + jti) {
	//     logger.Error("Token has been revoked: " + jti)
	//     errUtils.WriteErrorResponse(wtr, http.StatusUnauthorized, "Token has been revoked", "20004", req.Header.Get("X-Request-ID"))
	//     return
	// }

	logger.Info("Token verified successfully for client: " + clientID)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenVerifyResponse{
		Valid:     true,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
		ClientType: clientType,
	})
}