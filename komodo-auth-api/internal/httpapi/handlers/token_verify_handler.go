package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"komodo-auth-api/internal/logger"
	"komodo-auth-api/internal/thirdparty/aws"
)

func TokenVerifyHandler(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(wtr, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := ""
	if auth := req.Header.Get("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			token = strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		}
	}

	if token == "" {
		var body map[string]string
		if err := json.NewDecoder(req.Body).Decode(&body); err == nil {
			token = body["token"]
		}
	}

	if token == "" {
		http.Error(wtr, "missing token", http.StatusBadRequest)
		return
	}

	val, err := aws.GetCacheItem(token)
	if err != nil {
		logger.Error("failed to read token from cache: "+err.Error(), req)
		http.Error(wtr, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if val == "" {
		http.Error(wtr, "unauthorized", http.StatusUnauthorized)
		return
	}

	wtr.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(wtr).Encode(map[string]interface{}{
		"valid":      true,
		"expires_in": aws.DEFAULT_SESH_TTL,
	})
}