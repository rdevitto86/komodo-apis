package handlers

import (
	"net/http"
)

func MFADisableHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// Implementation for disabling MFA

	wtr.WriteHeader(http.StatusNoContent)
}
