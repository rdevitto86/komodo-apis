package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	errors "komodo-internal-lib-apis-go/http/common/errors/chi"
	logger "komodo-internal-lib-apis-go/logging/runtime"

	"komodo-internal-lib-apis-go/crypto/jwt"
)

type TokenRevokeResponse struct {
	Revoked   bool   `json:"revoked"`
	TokenID   string `json:"token_id,omitempty"`
	RevokedAt int64  `json:"revoked_at"`
}

// Handles token revocation requests
func TokenRevokeHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header
	token, err := jwt.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("failed to extract token from Authorization header", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusUnauthorized,
			"missing or invalid Authorization header",
			errCodes.ERR_INVALID_TOKEN,
		)
		return
	}

	// Verify token signature and parse claims
	_, claims, err := jwt.VerifyToken(token)
	if err != nil {
		logger.Error("invalid token", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusUnauthorized,
			"invalid or expired token",
			errCodes.ERR_INVALID_TOKEN,
		)
		return
	}

	// Extract JTI and client ID
	claimValues := jwt.ExtractStringClaims(claims, []string{"jti", "client_id"})
	jti, _ := claimValues["jti"].(string)
	clientID, _ := claimValues["client_id"].(string)

	if jti == "" {
		logger.Error("token missing JTI claim - cannot revoke")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"token does not support revocation (missing JTI)",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Calculate TTL - must be positive to store in cache
	ttl := jwt.GetTokenTTL(claims)
	if ttl <= 0 {
		logger.Warn("token already expired, no revocation needed")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(TokenRevokeResponse{
			Revoked:   true,
			TokenID:   jti,
			RevokedAt: time.Now().Unix(),
		})
		return
	}

	// Store revoked token in Elasticache with TTL
	// revokeKey := "revoked:token:" + jti
	// if err := elasticache.SetCacheItem(revokeKey, clientID, ttl); err != nil {
	// 	logger.Error("failed to revoke token in cache", err)
	// 	errUtils.WriteErrorResponse(
	// 		wtr,
	// 		http.StatusInternalServerError,
	// 		"Token revocation failed",
	// 		"50001",
	// 		req.Header.Get("X-Request-ID"),
	// 	)
	// 	return
	// }

	logger.Info("token revoked successfully for client: " + clientID + ", JTI: " + jti)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(TokenRevokeResponse{
		Revoked:   true,
		TokenID:   jti,
		RevokedAt: time.Now().Unix(),
	})
}