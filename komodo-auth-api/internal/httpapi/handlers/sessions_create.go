package handlers

import (
	"komodo-auth-api/internal/thirdparty/aws"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SessionRequestBody struct {
	UserID string `json:"user_id,omitempty"`
	CartID string `json:"cart_id,omitempty"`
}

func SessionCreateHandler(wtr http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(wtr, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	timeout, err := strconv.Atoi(os.Getenv("SESSION_EXPIRATION_HR"))
	if err != nil {
		http.Error(wtr, "Invalid session expiration", http.StatusInternalServerError)
		return
	}

	sessionToken := uuid.NewString()

	aws.SetSessionToken("key", sessionToken)

	if req.Context().Value("isUI").(bool) {
		http.SetCookie(wtr, &http.Cookie{
			Name: "session_token",
			Value: sessionToken,
			Path: "/",
			Expires: time.Now().Add(time.Duration(timeout) * time.Hour),
			Secure: true,
			HttpOnly: true,
		})
	} else {
		wtr.Header().Set("X-Session-Token", sessionToken)
	}
}