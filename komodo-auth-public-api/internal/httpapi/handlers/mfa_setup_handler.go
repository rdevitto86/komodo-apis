package handlers

import (
	"net/http"
)

func MFASetupHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// Implementation for setting up MFA

	wtr.WriteHeader(http.StatusCreated)
}
