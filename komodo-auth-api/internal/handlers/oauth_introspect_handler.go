package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-forge-apis-go/crypto/jwt"
	logger "komodo-forge-apis-go/loggers/runtime"
)

type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
	Sub       string `json:"sub,omitempty"`
	Aud       string `json:"aud,omitempty"`
}

// Handles OAuth 2.0 token introspection (RFC 7662).
// Returns token metadata if active, or {"active": false} if invalid/expired/revoked
func OAuthIntrospectHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header or request body
	tokenString, err := jwt.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("no token found in request", err)
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
		return
	}

	// Verify token signature and parse claims
	_, claims, err := jwt.VerifyToken(tokenString)
	if err != nil {
		logger.Error("token verification failed", err)
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
		return
	}

	// Check if token is expired
	if jwt.IsTokenExpired(claims) {
		logger.Info("token is expired")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
		return
	}

	// Extract claims
	claimValues := jwt.ExtractStringClaims(claims, []string{
		"client_id", "scope", "jti", "sub", "aud",
	})

	clientID, _ := claimValues["client_id"].(string)
	scope, _ := claimValues["scope"].(string)
	sub, _ := claimValues["sub"].(string)
	aud, _ := claimValues["aud"].(string)

	exp, _ := claims["exp"].(float64)
	iat, _ := claims["iat"].(float64)

	// TODO: Check if token is revoked in Elasticache
	// jti, _ := claimValues["jti"].(string)
	// if jti != "" && elasticache.Exists("revoked:token:" + jti) {
	//     logger.Info("token has been revoked: " + jti)
	//     wtr.WriteHeader(http.StatusOK)
	//     json.NewEncoder(wtr).Encode(IntrospectResponse{Active: false})
	//     return
	// }

	logger.Info("token introspection successful for client: " + clientID)

	// Return token metadata per RFC 7662
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(IntrospectResponse{
		Active:    true,
		Scope:     scope,
		ClientID:  clientID,
		TokenType: "Bearer",
		Exp:       int64(exp),
		Iat:       int64(iat),
		Sub:       sub,
		Aud:       aud,
	})
}
