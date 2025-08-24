package handlers

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
    Status string `json:"status"`
}

func HandleHealth(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(healthResponse{Status: "ok"})
}