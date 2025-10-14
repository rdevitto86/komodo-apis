package handlers

import "net/http"

func LogoutHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	// TODO Implement logout logic

	wtr.WriteHeader(http.StatusNoContent)
}