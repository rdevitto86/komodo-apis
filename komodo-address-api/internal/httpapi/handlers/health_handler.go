package handlers

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
    Status string `json:"status"`
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
}