package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"komodo-internal-lib-apis-go/aws/elasticache"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
)

func TokenDeleteHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

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

	if err := elasticache.DeleteCacheItem(token); err != nil {
		logger.Error("failed to delete token from cache: " + err.Error(), req)
		http.Error(wtr, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	wtr.WriteHeader(http.StatusNoContent)
}