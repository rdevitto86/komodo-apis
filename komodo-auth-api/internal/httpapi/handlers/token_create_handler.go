package handlers

import (
	"encoding/json"
	"komodo-internal-lib-apis-go/aws/elasticache"
	"komodo-internal-lib-apis-go/crypto/encryption"
	logger "komodo-internal-lib-apis-go/logger/runtime"
	"net/http"

	"github.com/google/uuid"
)

func TokenCreateHandler(wtr http.ResponseWriter, req *http.Request) {
	token := uuid.NewString()
	token, err := encryption.EncryptToken(token)
	if err != nil {
		logger.Error("failed to encrypt token: " + err.Error(), req)
		http.Error(wtr, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// TODO - store token in Elasticache with appropriate key and TTL
	// Example key could be user ID or session ID depending on use case
	// Here we use a placeholder "TODO_KEY"
	// The TTL can be adjusted as needed; using DEFAULT_SESH_TTL for example
	err = elasticache.SetCacheItem("TODO_KEY", token, elasticache.DEFAULT_SESH_TTL)
	if err != nil {
		logger.Error("failed to store token in cache: " + err.Error(), req)
		http.Error(wtr, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	wtr.Header().Set("Authorization", "Bearer " + token)
	wtr.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(wtr).Encode(map[string]interface{}{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": elasticache.DEFAULT_SESH_TTL,
	})
	if err != nil {
		logger.Error("failed to encode response", err)
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
	}
}