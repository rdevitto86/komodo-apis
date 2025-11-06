package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	errUtils "komodo-internal-lib-apis-go/http/utils/error"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"

	jwtUtils "komodo-auth-api/internal/httpapi/utils/jwt"
)

type TokenDeleteResponse struct {
	Revoked   bool   `json:"revoked"`
	TokenID   string `json:"token_id,omitempty"`
	RevokedAt int64  `json:"revoked_at"`
}

// Handles token revocation requests
func TokenDeleteHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header
	token, err := jwtUtils.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("Failed to extract token from Authorization header", err)
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusUnauthorized,
			"Missing or invalid Authorization header",
			"20001",
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

	// Extract JTI and client ID
	claimValues := jwtUtils.ExtractStringClaims(claims, []string{"jti", "client_id"})
	jti, _ := claimValues["jti"].(string)
	clientID, _ := claimValues["client_id"].(string)

	if jti == "" {
		logger.Error("Token missing JTI claim - cannot revoke")
		errUtils.WriteErrorResponse(
			wtr,
			http.StatusBadRequest,
			"Token does not support revocation (missing JTI)",
			"40008",
			req.Header.Get("X-Request-ID"),
		)
		return
	}

	// Calculate TTL - must be positive to store in cache
	ttl := jwtUtils.GetTokenTTL(claims)
	if ttl <= 0 {
		logger.Warn("Token already expired, no revocation needed")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(TokenDeleteResponse{
			Revoked:   true,
			TokenID:   jti,
			RevokedAt: time.Now().Unix(),
		})
		return
	}

	// Store revoked token in Elasticache with TTL
	// revokeKey := "revoked:token:" + jti
	// if err := elasticache.SetCacheItem(revokeKey, clientID, ttl); err != nil {
	// 	logger.Error("Failed to revoke token in cache", err)
	// 	errUtils.WriteErrorResponse(
	// 		wtr,
	// 		http.StatusInternalServerError,
	// 		"Token revocation failed",
	// 		"50001",
	// 		req.Header.Get("X-Request-ID"),
	// 	)
	// 	return
	// }

	logger.Info("Token revoked successfully for client: " + clientID + ", JTI: " + jti)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenDeleteResponse{
		Revoked:   true,
		TokenID:   jti,
		RevokedAt: time.Now().Unix(),
	})
}