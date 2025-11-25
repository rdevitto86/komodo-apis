package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"komodo-internal-lib-apis-go/crypto/jwt"
	authServ "komodo-internal-lib-apis-go/domains/auth/service"
	userServ "komodo-internal-lib-apis-go/domains/user"
	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	errors "komodo-internal-lib-apis-go/http/common/errors/chi"
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
}

// Handles user login requests
func LoginHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "no-store")

	var loginReq LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		logger.Error("failed to parse login request payload", err)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusInternalServerError,
			"failed to parse request payload",
			errCodes.ERR_INTERNAL_SERVER,
		)
		return
	}
	if loginReq.Email == "" || loginReq.Password == "" {
		logger.Error("invalid login request payload")
		errors.WriteErrorResponse(
			wtr, req, http.StatusBadRequest, "invalid request payload", errCodes.ERR_INVALID_REQUEST,
		)
		return
	}

	// Get service token for User API
	res := authServ.TokenGenerate(req, &authServ.TokenGenerateRequest{
		ClientID:     os.Getenv("AUTH_USER_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_USER_CLIENT_SECRET"),
		Scope:        "user.read",
	})

	if res.IsError() {
		logger.Error("failed to get service token from internal auth service", res.Error)
		errors.WriteErrorResponse(
			wtr, req, res.Status, "failed to authenticate service", errCodes.ERR_INTERNAL_API_CALL_FAILED,
		)
		return
	}

	// Parse response and extract service token
	bearer, ok := res.BodyParsed.(*authServ.TokenGenerateResponse)
	if !ok || bearer == nil {
		logger.Error("failed to type assert service token response")
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusInternalServerError,
			"failed to parse service token",
			errCodes.ERR_INTERNAL_SERVER,
		)
		return
	}

	// Fetch user details from User API with service token
	res = userServ.GetUserProfile(req, &userServ.UserProfileGetRequest{
		UserID:      loginReq.Email,
		Size:        userServ.ProfileSizeMinimal,
		BearerToken: "Bearer " + bearer.Token,
	})

	if (res.IsError()) {
		logger.Error("failed to fetch user profile", res.Error)
		errors.WriteErrorResponse(
			wtr, req, res.Status, "failed to fetch user profile", errCodes.ERR_INTERNAL_API_CALL_FAILED,
		)
		return
	}

	// Parse user profile response
	profile, ok := res.BodyParsed.(*userServ.UserProfileGetResponseMinimal)
	if !ok || profile == nil {
		logger.Error("failed to type assert user profile response", res.Error)
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusInternalServerError,
			"failed to parse user profile",
			errCodes.ERR_INTERNAL_SERVER,
		)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(profile.PasswordHash), []byte(loginReq.Password)); err != nil {
		logger.Error("invalid password attempt for email: " + loginReq.Email)
		errors.WriteErrorResponse(
			wtr, req, http.StatusUnauthorized, "invalid email or password", errCodes.ERR_ACCESS_DENIED,
		)
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
		errors.WriteErrorResponse(
			wtr,
			req,
			http.StatusInternalServerError,
			"failed to generate authentication token",
			errCodes.ERR_INTERNAL_SERVER,
		)
		return
	}

	response := LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
	}

	wtr.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(wtr).Encode(response); err != nil {
		logger.Error("failed to encode response", err)
	}
}
