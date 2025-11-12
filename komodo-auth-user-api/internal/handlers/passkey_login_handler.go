package handlers

import (
	"net/http"
)

func PasskeyLoginHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// TODO: Implement passkey login logic

	wtr.WriteHeader(http.StatusCreated)
}
