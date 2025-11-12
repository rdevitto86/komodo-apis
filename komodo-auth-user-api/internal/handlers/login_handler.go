package handlers

import (
	"encoding/json"
	"net/http"

	"komodo-internal-lib-apis-go/common/errors"
	"komodo-internal-lib-apis-go/crypto/jwt"
	authServ "komodo-internal-lib-apis-go/domains/auth/service"
	userServ "komodo-internal-lib-apis-go/domains/user"
	logger "komodo-internal-lib-apis-go/logging/runtime"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
	User userServ.UserProfileGetResponse `json:"user"`
}

// Handles user login requests
func LoginHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	var loginReq LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		logger.Error("failed to parse login request payload")
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, errors.ERR_INTERNAL_SERVER, "Failed to parse request payload")
		return
	}
	if loginReq.Email == "" || loginReq.Password == "" {
		logger.Error("invalid login request payload")
		errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, errors.ERR_INVALID_REQUEST, "Invalid request payload")
		return
	}

	// Get service token for User API
	res := authServ.GetServiceToken(req, &authServ.ServiceTokenRequest{
		ClientID:     "komodo-auth-user-api",
		ClientSecret: "super-secret-key",
		Scope:        "user.read",
	})
	if res.IsError() {
		logger.Error("failed to get service token from internal auth service", res.Error)
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, errors.ERR_INTERNAL_SERVER, "Failed to authenticate service")
		return
	}

	// Parse response and set bearer token
	bearer, ok := res.BodyParsed.(*authServ.ServiceTokenResponse)
	if !ok || bearer == nil {
		logger.Error("failed to type assert service token response")
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, errors.ERR_INTERNAL_SERVER, "Failed to parse service token")
		return
	}
	wtr.Header().Set("Authorization", "Bearer " + bearer.Token)

	// Fetch user details from User API
	res = userServ.GetUserProfile(req, &userServ.UserProfileGetRequest{
		UserID: loginReq.Email,
		Size:   userServ.ProfileSizeMinimal,
	})
	if (res.IsError()) {
		logger.Error("failed to fetch user profile")
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, errors.ERR_INTERNAL_SERVER, "Failed to fetch user profile")
		return
	}

	// Parse user profile response
	profile, ok := res.BodyParsed.(*userServ.UserProfileGetResponse)
	if !ok || profile == nil {
		logger.Error("failed to type assert user profile response", res.Error)
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, errors.ERR_INTERNAL_SERVER, "Failed to parse user profile")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(profile.PasswordHash), []byte(loginReq.Password)); err != nil {
		logger.Error("invalid password attempt for email: " + loginReq.Email)
		errors.WriteErrorResponse(wtr, req, http.StatusUnauthorized, errors.ERR_ACCESS_DENIED, "invalid email or password")
		return
	}

	logger.Info("user authenticated successfully: " + profile.UserID)

	// Generate JWT token (temporary until Redis session is implemented)
	expiresIn := int64(3600) // 1 hour
	claims := jwt.CreateStandardClaims(
		"komodo-auth-user-api",
		profile.UserID,
		"komodo-web-app",
		expiresIn,
		map[string]interface{}{
			"user_id": profile.UserID,
			"first_name": profile.FirstName,
			"last_name": profile.LastName,
		},
	)

	token, signErr := jwt.SignToken(claims)
	if signErr != nil {
		logger.Error("failed to sign JWT token for login", signErr)
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, errors.ERR_INTERNAL_SERVER, "failed to generate authentication token")
		return
	}

	response := LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
		User:      *profile,
	}

	wtr.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(wtr).Encode(response); err != nil {
		logger.Error("failed to encode response", err)
	}
}
