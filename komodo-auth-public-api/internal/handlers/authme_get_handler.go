package handlers

import (
	"encoding/json"
	"net/http"

	logger "komodo-internal-lib-apis-go/logging/runtime"
)

type UserProfileResponse struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Phone         string `json:"phone,omitempty"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	EmailVerified bool   `json:"email_verified"`
	MFAEnabled    bool   `json:"mfa_enabled"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Retrieves the current authenticated user's profile
func AuthMeGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")

	// Extract user_id from session token
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		// TODO: Validate session token with Redis/Elasticache
		// sessionData, err := elasticache.GetSession(cookie.Value)
		// if err != nil {
		//     logger.Error("failed to retrieve session from cache", err)
		//     w.WriteHeader(http.StatusUnauthorized)
		//     json.NewEncoder(w).Encode(map[string]string{
		//         "error":   "unauthorized",
		//         "message": "session expired or invalid",
		//     })
		//     return
		// }
		// userID = sessionData.UserID
	}

	logger.Info("fetching profile for user: " + userID)

	// TODO: Call User API to fetch user profile
	// GET http://komodo-user-api/users/{userID}
	//
	// userAPIURL := config.GetConfigValue("USER_API_URL")
	// req, err := http.NewRequestWithContext(r.Context(), "GET", 
	//     fmt.Sprintf("%s/users/%s", userAPIURL, userID), nil)
	// if err != nil {
	//     logger.Error("failed to create user API request", err)
	//     w.WriteHeader(http.StatusInternalServerError)
	//     json.NewEncoder(w).Encode(map[string]string{
	//         "error": "internal_server_error",
	//         "message": "failed to fetch user profile",
	//     })
	//     return
	// }
	//
	// // Add internal service authentication (JWT from Auth Internal API)
	// serviceToken := getServiceToken() // Get M2M token for Auth→User communication
	// req.Header.Set("Authorization", "Bearer " + serviceToken)
	// req.Header.Set("Accept", "application/json;v=1")
	//
	// client := &http.Client{Timeout: 5 * time.Second}
	// resp, err := client.Do(req)
	// if err != nil {
	//     logger.Error("failed to call user API", err)
	//     w.WriteHeader(http.StatusServiceUnavailable)
	//     json.NewEncoder(w).Encode(map[string]string{
	//         "error": "service_unavailable",
	//         "message": "user service temporarily unavailable",
	//     })
	//     return
	// }
	// defer resp.Body.Close()
	//
	// if resp.StatusCode != http.StatusOK {
	//     logger.Error(fmt.Sprintf("user API returned status %d", resp.StatusCode))
	//     w.WriteHeader(resp.StatusCode)
	//     io.Copy(w, resp.Body)
	//     return
	// }
	//
	// var profile UserProfileResponse
	// if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
	//     logger.Error("failed to decode user API response", err)
	//     w.WriteHeader(http.StatusInternalServerError)
	//     json.NewEncoder(w).Encode(map[string]string{
	//         "error": "internal_server_error",
	//         "message": "failed to parse user profile",
	//     })
	//     return
	// }

	// Mock response until User API is implemented
	profile := UserProfileResponse{
		UserID:        userID,
		Email:         "user@example.com",
		Name:          "John Doe",
		Phone:         "+1234567890",
		AvatarURL:     "https://example.com/avatars/user.jpg",
		EmailVerified: true,
		MFAEnabled:    false,
		CreatedAt:     "2025-01-01T00:00:00Z",
		UpdatedAt:     "2025-11-09T00:00:00Z",
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		logger.Error("failed to encode response", err)
	}
}