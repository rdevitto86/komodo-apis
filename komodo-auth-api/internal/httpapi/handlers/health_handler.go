package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(wtr http.ResponseWriter, req *http.Request) {
  json.NewEncoder(wtr).Encode(map[string]string{"status": "OK"})
}