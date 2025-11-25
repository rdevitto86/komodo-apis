package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-internal-lib-apis-go/crypto/jwt"
	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	errors "komodo-internal-lib-apis-go/http/common/errors/chi"
	logger "komodo-internal-lib-apis-go/logging/runtime"
)

type LogoutResponse struct {
	Message string `json:"message"`
}

// Handles user logout requests
func LogoutHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	// Extract token from Authorization header
	tokenString, err := jwt.ExtractTokenFromRequest(req)
	if err != nil {
		logger.Error("no token found in logout request", err)
		errors.WriteErrorResponse(
			wtr, req, http.StatusUnauthorized, "missing authorization token", errCodes.ERR_INVALID_TOKEN,
		)
		return
	}

	// Verify token and extract claims
	_, claims, err := jwt.VerifyToken(tokenString)
	if err != nil {
		logger.Error("invalid token in logout request", err)
		errors.WriteErrorResponse(
			wtr, req, http.StatusUnauthorized, "invalid token", errCodes.ERR_INVALID_TOKEN,
		)
		return
	}

	// Extract JTI (JWT ID) for revocation
	claimValues := jwt.ExtractStringClaims(claims, []string{"jti", "sub"})
	jti, _ := claimValues["jti"].(string)
	userID, _ := claimValues["sub"].(string)

	if jti == "" {
		logger.Error("token missing jti claim for revocation")
		errors.WriteErrorResponse(
			wtr, req, http.StatusBadRequest, "token cannot be revoked", errCodes.ERR_INVALID_TOKEN,
		)
		return
	}

	// TODO: Add token to revocation list in Redis/Elasticache
	// Key: "revoked:token:{jti}"
	// Value: "1"
	// TTL: token expiration time
	// Example:
	// exp, _ := claimValues["exp"].(int64)
	// ttl := time.Unix(exp, 0).Sub(time.Now())
	// redisClient.Set("revoked:token:" + jti, "1", ttl)

	// TODO: Delete user session from Redis/Elasticache
	// Key: "session:user:{user_id}"
	// Example:
	// redisClient.Del("session:user:" + userID)

	logger.Info("user logged out successfully: " + userID)

	wtr.WriteHeader(http.StatusOK)
	json.NewEncoder(wtr).Encode(LogoutResponse{
		Message: "logged out successfully",
	})
}