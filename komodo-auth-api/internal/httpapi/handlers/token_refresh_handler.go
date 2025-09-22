package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"komodo-auth-api/internal/logger"
	"komodo-auth-api/internal/thirdparty/aws"
)

func TokenRefreshHandler(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(wtr, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a real flow you might validate an existing refresh token or session.
	// Here we simply emit a new session token and store it in cache.
	sessionToken := uuid.NewString()

	if err := aws.SetCacheItem(sessionToken, "1", aws.DEFAULT_SESH_TTL); err != nil {
		logger.Error("failed to store refreshed token: "+err.Error(), req)
		http.Error(wtr, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	wtr.Header().Set("Authorization", "Bearer "+sessionToken)
	wtr.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(wtr).Encode(map[string]interface{}{
		"token":      sessionToken,
		"token_type": "Bearer",
		"expires_in": aws.DEFAULT_SESH_TTL,
	})
}