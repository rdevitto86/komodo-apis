package handlers

import (
	"net/http"

	"github.com/google/uuid"
)

func TokenRefreshHandler(wtr http.ResponseWriter, req *http.Request) {
	sessionToken := uuid.NewString()
	// aws.SetSessionToken(sessionToken)

	wtr.Header().Set("Authorization", "Bearer " + sessionToken)
}