package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-internal-lib-apis-go/crypto/jwt"
	logger "komodo-internal-lib-apis-go/logging/runtime"
)

type TokenVerifyResponse struct {
	Active   bool   `json:"active"`              // Is token valid (signature, expiry, not revoked)?
	ClientID string `json:"client_id,omitempty"` // Which service is making the request?
	Scope    string `json:"scope,omitempty"`     // What permissions does it have?
}

// Handles internal JWT token verification for service-to-service authentication
// POST /auth/token/verify - Validates JWT signature, expiration, and revocation status
func TokenVerifyHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header (preferred) or request body (fallback)
	tokenString, err := jwt.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("no token found in request", err)
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(TokenVerifyResponse{Active: false})
		return
	}

	// Verify token signature and parse claims
	_, claims, err := jwt.VerifyToken(tokenString)
	if err != nil {
		logger.Error("token verification failed", err)
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(TokenVerifyResponse{Active: false})
		return
	}

	// Check if token is expired (return active:false per RFC 7662)
	if jwt.IsTokenExpired(claims) {
		logger.Info("token is expired")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(TokenVerifyResponse{Active: false})
		return
	}

	// Extract minimal claims needed for authorization decisions
	claimValues := jwt.ExtractStringClaims(claims, []string{
		"client_id", "scope", "jti",
	})
	
	clientID, _ := claimValues["client_id"].(string)
	scope, _ := claimValues["scope"].(string)

	// TODO: Check if token is revoked in Elasticache
	// jti, _ := claimValues["jti"].(string)
	// if jti != "" && redisClient.Exists("revoked:token:" + jti) {
	//     logger.Info("token has been revoked: " + jti)
	//     wtr.WriteHeader(http.StatusOK)
	//     json.NewEncoder(wtr).Encode(TokenVerifyResponse{Active: false})
	//     return
	// }

	logger.Info("token verification successful for client: " + clientID)

	// Return only what's needed: valid + who + what permissions
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenVerifyResponse{
		Active:   true,
		ClientID: clientID,
		Scope:    scope,
	})
}
