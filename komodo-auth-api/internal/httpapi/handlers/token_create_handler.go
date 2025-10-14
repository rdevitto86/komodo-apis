package handlers

import (
	"encoding/json"
	"komodo-internal-lib-apis-go/aws/elasticache"
	"komodo-internal-lib-apis-go/crypto/encryption"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	"net/http"

	"github.com/google/uuid"
)

func TokenCreateHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	var (
		token string
		err   error
	)

	if token, err = encryption.EncryptToken(uuid.NewString()); err == nil {
		// TODO - store token in Elasticache with appropriate key and TTL
		// Example key could be user ID or session ID depending on use case
		// Here we use a placeholder "TODO_KEY"
		// The TTL can be adjusted as needed; using DEFAULT_SESH_TTL for example

		if err = elasticache.SetCacheItem("TODO_KEY", token, elasticache.DEFAULT_SESH_TTL); err == nil {
			json.NewEncoder(wtr).Encode(map[string]interface{}{
				"token":      token,
				"token_type": "Bearer",
				"expires_in": elasticache.DEFAULT_SESH_TTL,
			})

			wtr.Header().Set("Authorization", "Bearer " + token)
			wtr.WriteHeader(http.StatusCreated)
			return
		}
	}

	logger.Error("failed to encrypt token: " + err.Error(), req)
	http.Error(wtr, "failed to create token", http.StatusInternalServerError)
}