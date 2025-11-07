package handlers

import (
	"net/http"
)

func PasskeyStartHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// TODO: Implement passkey start logic

	wtr.WriteHeader(http.StatusCreated)
}
