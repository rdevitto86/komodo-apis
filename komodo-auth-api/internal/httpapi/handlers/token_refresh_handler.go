package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"komodo-internal-lib-apis-go/aws/elasticache"
	logger "komodo-internal-lib-apis-go/logger/runtime"
)

func TokenRefreshHandler(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(wtr, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a real flow you might validate an existing refresh token or session.
	// Here we simply emit a new session token and store it in cache.
	sessionToken := uuid.NewString()

	if err := elasticache.SetCacheItem(sessionToken, "1", elasticache.DEFAULT_SESH_TTL); err != nil {
		logger.Error("failed to store refreshed token: "+err.Error(), req)
		http.Error(wtr, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	wtr.Header().Set("Authorization", "Bearer "+sessionToken)
	wtr.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(wtr).Encode(map[string]interface{}{
		"token":      sessionToken,
		"token_type": "Bearer",
		"expires_in": elasticache.DEFAULT_SESH_TTL,
	})
}