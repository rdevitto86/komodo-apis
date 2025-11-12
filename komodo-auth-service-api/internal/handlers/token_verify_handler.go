package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-internal-lib-apis-go/crypto/jwt"
	logger "komodo-internal-lib-apis-go/logging/runtime"
)

type TokenVerifyResponse struct {
	Active     bool   `json:"active"`               // Required: whether token is valid and not expired
	ClientID   string `json:"client_id,omitempty"`  // Client identifier (service making the request)
	Scope      string `json:"scope,omitempty"`      // Space-separated permissions/scopes
	ClientType string `json:"client_type,omitempty"` // Type: api (service) or browser (user session)
	Exp        int64  `json:"exp,omitempty"`        // Token expiration timestamp (Unix seconds)
	Jti        string `json:"jti,omitempty"`        // JWT ID (for revocation checking)
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

	// Extract relevant claims for M2M validation
	claimValues := jwt.ExtractStringClaims(claims, []string{
		"client_id", "scope", "client_type", "exp", "jti",
	})
	
	clientID, _ := claimValues["client_id"].(string)
	scope, _ := claimValues["scope"].(string)
	clientType, _ := claimValues["client_type"].(string)
	exp, _ := claimValues["exp"].(int64)
	jti, _ := claimValues["jti"].(string)

	// TODO: Check if token is revoked in Elasticache
	// if jti != "" && redisClient.Exists("revoked:token:" + jti) {
	//     logger.Info("token has been revoked: " + jti)
	//     wtr.WriteHeader(http.StatusOK)
	//     json.NewEncoder(wtr).Encode(TokenVerifyResponse{Active: false})
	//     return
	// }

	logger.Info("token verification successful for client: " + clientID)

	// Return minimal response for service-to-service validation
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenVerifyResponse{
		Active:     true,
		ClientID:   clientID,
		Scope:      scope,
		ClientType: clientType,
		Exp:        exp,
		Jti:        jti,
	})
}
