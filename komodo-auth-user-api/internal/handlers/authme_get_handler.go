package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	ctxKeys "komodo-internal-lib-apis-go/common/context"
	"komodo-internal-lib-apis-go/common/errors"
	authServ "komodo-internal-lib-apis-go/domains/auth/service"
	userServ "komodo-internal-lib-apis-go/domains/user"
	logger "komodo-internal-lib-apis-go/logging/runtime"
)

// Retrieves the current authenticated user's profile
func AuthMeGetHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")
	wtr.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")

	// Extract user_id from context
	userID, ok := req.Context().Value(ctxKeys.USER_ID_KEY).(string)
	if !ok || userID == "" {
		logger.Error("user_id not found in request context")
		errors.WriteErrorResponse(wtr, req, http.StatusUnauthorized, "unauthorized", errors.ERR_INVALID_TOKEN)
		return
	}

	logger.Info("fetching profile for user: " + userID)

	// Get service token for User API
	res := authServ.TokenGenerate(req, &authServ.TokenGenerateRequest{
		ClientID:     os.Getenv("AUTH_USER_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_USER_CLIENT_SECRET"),
		Scope:        "user.read",
	})
	if res.IsError() {
		logger.Error("failed to get service token from internal auth service", res.Error)
		errors.WriteErrorResponse(wtr, req, res.Status, "failed to authenticate service", errors.ERR_INTERNAL_API_CALL_FAILED)
		return
	}

	// Parse response and extract service token
	bearer, ok := res.BodyParsed.(*authServ.TokenGenerateResponse)
	if !ok || bearer == nil {
		logger.Error("failed to type assert service token response")
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, "failed to parse service token", errors.ERR_INTERNAL_SERVER)
		return
	}

	logger.Debug("obtained service token for User API: " + bearer.Token)

	// Fetch user profile from User API with service token
	res = userServ.GetUserProfile(req, &userServ.UserProfileGetRequest{
		UserID:      userID,
		Size:        userServ.ProfileSizeBasic,
		BearerToken: "Bearer " + bearer.Token,
	})
	if res.IsError() {
		logger.Error("failed to fetch user profile", res.Error)
		errors.WriteErrorResponse(wtr, req, res.Status, "failed to fetch user profile", errors.ERR_INTERNAL_API_CALL_FAILED)
		return
	}

	// Parse user profile response
	profile, ok := res.BodyParsed.(*userServ.UserProfileGetResponseBasic)
	if !ok || profile == nil {
		logger.Error("failed to type assert user profile response", res.Error)
		errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, "failed to parse user profile", errors.ERR_INTERNAL_SERVER)
		return
	}

	logger.Info("successfully fetched profile for user: " + profile.UserID)

	// Map to response format
	response := userServ.UserProfileGetResponseBasic{
		UserID: "user-id-123",
		FirstName: "Test",
		LastName: "User",
		AvatarURL: "https://example.com/avatar.jpg",
	}

	wtr.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(wtr).Encode(response); err != nil {
		logger.Error("failed to encode response", err)
	}
}