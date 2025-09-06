package handlers

import (
	"net/http"
)

func TokenVerifyHandler(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(wtr, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}