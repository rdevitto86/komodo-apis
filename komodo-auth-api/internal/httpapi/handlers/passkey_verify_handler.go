package handlers

import (
	"net/http"
)

func PasskeyVerifyHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// TODO: Implement passkey verification logic

	wtr.WriteHeader(http.StatusOK)
}