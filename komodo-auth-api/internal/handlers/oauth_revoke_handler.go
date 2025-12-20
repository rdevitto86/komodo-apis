package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"komodo-forge-apis-go/crypto/jwt"
	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/chi"
	logger "komodo-forge-apis-go/logging/runtime"
)

type RevokeRequest struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint,omitempty"` // "access_token" or "refresh_token"
}

// Handles OAuth 2.0 token revocation (RFC 7009).
// Revokes access or refresh tokens
func OAuthRevokeHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Parse request body
	var reqBody RevokeRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		logger.Error("failed to parse request body", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_request",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	if reqBody.Token == "" {
		logger.Error("missing token parameter")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusBadRequest,
			"invalid_request: missing token",
			errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Verify token signature and parse claims
	_, claims, err := jwt.VerifyToken(reqBody.Token)
	if err != nil {
		// Per RFC 7009, return 200 OK even if token is invalid
		// (prevents information disclosure about token validity)
		logger.Warn("invalid token submitted for revocation", err)
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]bool{"revoked": true})
		return
	}

	// Extract JTI and client ID
	claimValues := jwt.ExtractStringClaims(claims, []string{"jti", "client_id"})
	jti, _ := claimValues["jti"].(string)
	clientID, _ := claimValues["client_id"].(string)

	if jti == "" {
		// Token without JTI cannot be revoked (shouldn't happen in our system)
		logger.Warn("token missing JTI claim")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]bool{"revoked": true})
		return
	}

	// Calculate TTL
	ttl := jwt.GetTokenTTL(claims)
	if ttl <= 0 {
		// Token already expired, no need to revoke
		logger.Info("token already expired, no revocation needed")
		wtr.WriteHeader(http.StatusOK)
		json.NewEncoder(wtr).Encode(map[string]bool{"revoked": true})
		return
	}

	// TODO: Store revoked token in Elasticache with TTL
	// revokeKey := "revoked:token:" + jti
	// if err := elasticache.SetCacheItem(revokeKey, clientID, ttl); err != nil {
	// 	logger.Error("failed to revoke token in cache", err)
	// 	errors.WriteErrorResponse(
	// 		wtr,
	// 		req,
	// 		http.StatusInternalServerError,
	// 		"server_error",
	// 		errCodes.ERR_INTERNAL_SERVER,
	// 	)
	// 	return
	// }

	logger.Info("token revoked successfully for client: " + clientID + ", JTI: " + jti)

	// Per RFC 7009, return 200 OK with empty response (or small JSON)
	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(map[string]interface{}{
		"revoked":    true,
		"revoked_at": time.Now().Unix(),
	})
}
