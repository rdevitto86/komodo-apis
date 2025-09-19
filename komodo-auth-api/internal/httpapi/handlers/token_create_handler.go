package handlers

import (
	"encoding/json"
	"komodo-auth-api/internal/crypto"
	"komodo-auth-api/internal/httpapi/utils"
	"komodo-auth-api/internal/logger"
	"komodo-auth-api/internal/thirdparty/aws"
	"net/http"
)

func TokenCreateHandler(wtr http.ResponseWriter, req *http.Request) {
	token := utils.GenerateToken()
	token, err := crypto.EncryptToken(token)
	if err != nil {
		logger.Error("failed to encrypt token: " + err.Error(), req)
		http.Error(wtr, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// TODO - store token in Elasticache with appropriate key and TTL
	// Example key could be user ID or session ID depending on use case
	// Here we use a placeholder "TODO_KEY"
	// The TTL can be adjusted as needed; using DEFAULT_SESH_TTL for example
	err = aws.SetCacheItem("TODO_KEY", token, aws.DEFAULT_SESH_TTL)
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
		"expires_in": aws.DEFAULT_SESH_TTL,
	})
	if err != nil {
		logger.Error("failed to encode response", err)
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
	}
}