package handlers

import (
	"encoding/json"
	"net/http"

	errUtils "komodo-internal-lib-apis-go/http/utils/error"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"

	jwtUtils "komodo-auth-api/internal/httpapi/utils/jwt"
)

type TokenVerifyRequest struct {
	Token string `json:"token,omitempty"`
}

type TokenVerifyResponse struct {
	Valid     bool   `json:"valid"`
	ClientID  string `json:"client_id,omitempty"`
	Scope     string `json:"scope,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
	IssuedAt  int64  `json:"issued_at,omitempty"`
}

// Handles token verification requests
func TokenVerifyHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header or request body
	token, err := jwtUtils.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("Failed to extract token", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Missing token in Authorization header or request body",
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
	clientID, _ := jwtUtils.ExtractStringClaim(claims, "client_id")
	scope, _ := jwtUtils.ExtractStringClaim(claims, "scope")
	tokenUse, _ := jwtUtils.ExtractStringClaim(claims, "token_use")
	expiresAt, _ := jwtUtils.ExtractInt64Claim(claims, "exp")
	issuedAt, _ := jwtUtils.ExtractInt64Claim(claims, "iat")

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
	})
}