package handlers

import (
	"encoding/json"
	"komodo-auth-api/internal/httpapi/utils"
	"komodo-auth-api/internal/thirdparty/aws"
	"net/http"
)

func TokenCreateHandler(wtr http.ResponseWriter, req *http.Request) {
	idemKey := req.Header.Get("Idempotency-Key") // prevent duplicate requests
	// TODO: Move from local/demo utils to Redis/Dynamo for idempotency mapping and token storage.
	token := utils.GenerateToken(idemKey, req.Header.Get("User-ID"), req.Header.Get("Device-ID"))

	// aws.SetSessionToken(token)

	wtr.Header().Set("Authorization", "Bearer "+token)
	if idemKey != "" {
		wtr.Header().Set("Idempotency-Key", idemKey)
	}

	wtr.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(wtr).Encode(map[string]interface{}{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": aws.DEFAULT_SESH_TTL,
	}); err != nil {
		// TODO - logger
		http.Error(wtr, err.Error(), http.StatusInternalServerError)
	}
}