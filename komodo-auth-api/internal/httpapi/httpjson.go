package httpapi

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Data  interface{} `json:"data"`
	Error *Error      `json:"error,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}, err *Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{Data: data, Error: err})
}
