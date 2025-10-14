package handlers

import (
	"net/http"
)

func MFAEnableHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// Implementation for enabling MFA

	wtr.WriteHeader(http.StatusOK)
}
