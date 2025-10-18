package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"komodo-internal-lib-apis-go/config"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TokenCreateRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type TokenCreateResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	IssuedAt     int64  `json:"issued_at"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

func TokenCreateHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")
	wtr.Header().Set("Pragma", "no-cache")

	// Parse request body
	var reqBody TokenCreateRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		wtr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wtr).Encode(ErrorResponse{
			Error:            "invalid_request",
			ErrorDescription: "Invalid request body",
		})
		return
	}

	// Validate client credentials
	if reqBody.ClientID == "" || reqBody.ClientSecret == "" {
		wtr.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(wtr).Encode(ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Missing client credentials",
		})
		return
	}

	// TODO: Validate client_id and client_secret against database
	// Placeholder validation for now
	if reqBody.ClientID != "test-client" || reqBody.ClientSecret != "test-secret" {
		wtr.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(wtr).Encode(ErrorResponse{
			Error:            "invalid_client",
			ErrorDescription: "Invalid client credentials",
		})
		return
	}

	// Load RSA private key from config
	privateKeyPEM := config.GetConfigValue("JWT_PRIVATE_KEY")
	fmt.Printf("privateKeyPEM: %s\n", privateKeyPEM)
	if privateKeyPEM == "" {
		logger.Error("JWT_PRIVATE_KEY not configured")
		wtr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(wtr).Encode(ErrorResponse{
			Error:            "server_error",
			ErrorDescription: "Token signing not configured",
		})
		return
	}

	// Parse private key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		logger.Error("failed to parse private key", err)
		wtr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(wtr).Encode(ErrorResponse{
			Error:            "server_error",
			ErrorDescription: "Token signing failed",
		})
		return
	}

	// Create JWT claims
	now := time.Now()
	expiresIn := 3600 // 1hr

	claims := jwt.MapClaims{
	"iss":       "komodo-auth-api",
	"sub":       reqBody.ClientID,
	"aud":       "komodo-apis",
	"exp":       now.Add(time.Duration(expiresIn) * time.Second).Unix(),
	"iat":       now.Unix(),
	"nbf":       now.Unix(),
	"scope":     reqBody.Scope,
	"client_id": reqBody.ClientID,
	}

	// Create and sign token with RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	accessToken, err := token.SignedString(privateKey)
	if err != nil {
		logger.Error("failed to sign token", err)
		wtr.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(wtr).Encode(ErrorResponse{
			Error:            "server_error",
			ErrorDescription: "Token generation failed",
		})
		return
	}

	// TODO: Generate refresh token (optional)
	// TODO: Store token in Elasticache/Redis with TTL for revocation support
	// Example:
	// tokenID := uuid.NewString()
	// claims["jti"] = tokenID
	// elasticache.SetCacheItem("token:"+tokenID, accessToken, time.Duration(expiresIn)*time.Second)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenCreateResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Scope:       reqBody.Scope,
		IssuedAt:    now.Unix(),
	})
}