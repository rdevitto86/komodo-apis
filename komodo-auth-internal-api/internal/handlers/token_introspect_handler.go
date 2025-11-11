package handlers

import (
	"encoding/json"
	"net/http"

	logger "komodo-internal-lib-apis-go/logging/runtime"

	jwtUtils "komodo-internal-lib-apis-go/auth/jwt"
)

// OAuth 2.0 Token Introspection Request (RFC 7662)
type TokenIntrospectRequest struct {
	Token         string `json:"token"`           // Required: token to introspect
	TokenTypeHint string `json:"token_type_hint"` // Optional: access_token or refresh_token
}

// OAuth 2.0 Token Introspection Response (RFC 7662)
type TokenIntrospectResponse struct {
	Active     bool   `json:"active"`               // Required: whether token is active
	Scope      string `json:"scope,omitempty"`      // Optional: space-separated scopes
	ClientID   string `json:"client_id,omitempty"`  // Optional: client identifier
	TokenType  string `json:"token_type,omitempty"` // Optional: access_token or refresh_token
	Exp        int64  `json:"exp,omitempty"`        // Optional: expiration timestamp
	Iat        int64  `json:"iat,omitempty"`        // Optional: issued at timestamp
	Sub        string `json:"sub,omitempty"`        // Optional: subject (usually client_id)
	Aud        string `json:"aud,omitempty"`        // Optional: audience
	Iss        string `json:"iss,omitempty"`        // Optional: issuer
	Jti        string `json:"jti,omitempty"`        // Optional: JWT ID
	ClientType string `json:"client_type,omitempty"` // Custom: api or browser
}

// Handles token introspection requests per OAuth 2.0 RFC 7662
// POST /auth/token/introspect with token in request body
func TokenIntrospectHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Parse request body per RFC 7662 (token in body, not header)
	var reqBody TokenIntrospectRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse introspection request body", err)
		wtr.WriteHeader(http.StatusOK) // RFC 7662: Return 200 even for bad requests
		json.NewEncoder(wtr).Encode(TokenIntrospectResponse{Active: false})
		return
	}

	// Validate token present
	if reqBody.Token == "" {
		logger.Error("missing token in introspection request")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(TokenIntrospectResponse{Active: false})
		return
	}

	// Verify token signature and parse claims
	_, claims, err := jwtUtils.VerifyToken(reqBody.Token)
	if err != nil {
		logger.Error("invalid token signature", err)
		wtr.WriteHeader(http.StatusOK) // RFC 7662: Return 200 with active:false
		json.NewEncoder(wtr).Encode(TokenIntrospectResponse{Active: false})
		return
	}

	// Check if token is expired (return active:false per RFC 7662)
	if jwtUtils.IsTokenExpired(claims) {
		logger.Info("token expired during introspection")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(TokenIntrospectResponse{Active: false})
		return
	}

	// Extract all relevant claims
	claimValues := jwtUtils.ExtractStringClaims(claims, []string{
		"client_id", "scope", "token_use", "exp", "iat", "sub", "aud", "iss", "jti", "client_type",
	})
	
	clientID, _ := claimValues["client_id"].(string)
	scope, _ := claimValues["scope"].(string)
	tokenUse, _ := claimValues["token_use"].(string)
	exp, _ := claimValues["exp"].(int64)
	iat, _ := claimValues["iat"].(int64)
	sub, _ := claimValues["sub"].(string)
	aud, _ := claimValues["aud"].(string)
	iss, _ := claimValues["iss"].(string)
	jti, _ := claimValues["jti"].(string)
	clientType, _ := claimValues["client_type"].(string)

	// TODO: Check if token is revoked in Elasticache
	// if jti != "" && redisClient.Exists("revoked:token:" + jti) {
	//     logger.Info("token has been revoked: " + jti)
	//     wtr.WriteHeader(http.StatusOK)
	//     json.NewEncoder(wtr).Encode(TokenIntrospectResponse{Active: false})
	//     return
	// }

	// Map token_use to token_type
	tokenType := "access_token"
	if tokenUse == "refresh" {
		tokenType = "refresh_token"
	}

	logger.Info("token introspection successful for client: " + clientID)

	// Return active:true with all claims per RFC 7662
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenIntrospectResponse{
		Active:     true,
		Scope:      scope,
		ClientID:   clientID,
		TokenType:  tokenType,
		Exp:        exp,
		Iat:        iat,
		Sub:        sub,
		Aud:        aud,
		Iss:        iss,
		Jti:        jti,
		ClientType: clientType,
	})
}