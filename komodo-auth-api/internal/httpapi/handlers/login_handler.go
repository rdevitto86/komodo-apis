package handlers

import (
	"komodo-auth-api/internal/config"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func LoginHandler(wtr http.ResponseWriter, req *http.Request) {
	sessionToken := uuid.NewString()
	// aws.SetSessionToken(sessionToken)

	timeout, err := strconv.Atoi(config.GetConfigValue("SESSION_TIMEOUT_HOURS"))
	if err != nil {
		http.Error(wtr, "Invalid session expiration", http.StatusInternalServerError)
		return
	}

	http.SetCookie(wtr, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Path:    "/",
		Expires: time.Now().Add(time.Duration(timeout) * time.Hour),
		Secure:  true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode, // CSRF protection
	})

	wtr.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse and validate the request body (get username and password).
	// Verify credentials (check username/password against your user store).
	// Generate a session or JWT token if credentials are valid.
	// Store session/token (in cache or DB, if needed).
	// Return the token and user info in the response.
	// Handle errors (invalid credentials, missing fields, etc.).
}