package user

import (
	"fmt"

	httpclient "komodo-internal-lib-apis-go/http/client"
	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	httptypes "komodo-internal-lib-apis-go/http/types"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"net/http"
	"os"
)

var (
	client *httpclient.HTTPClient
	userAPIEndpoint string
)

func init() {
	if client = httpclient.GetInstance(); client == nil {
		logger.Error("failed to initialize http client")
	}
	if userAPIEndpoint = os.Getenv("URL_USER_API"); userAPIEndpoint == "" {
		logger.Error("user api endpoint environment variable not set")
	}
}

// Retrieves the user profile for the authenticated user
func GetUserProfile(req *http.Request, payload *UserProfileGetRequest) *httptypes.APIResponse {
	requestID := req.Header.Get("X-Request-ID")
	if requestID == "" { requestID = "unknown" }

	if client == nil || userAPIEndpoint == "" {
		return httptypes.ErrorResponse(
			http.StatusInternalServerError,
			errCodes.ERR_INTERNAL_SERVER,
			"service not initialized",
			"",
			requestID,
		)
	}
	if payload.BearerToken == "" {
		return httptypes.ErrorResponse(
			http.StatusUnauthorized,
			errCodes.ERR_INVALID_TOKEN,
			"missing bearer token for internal API authentication",
			"",
			requestID,
		)
	}
	if payload.Size == "" {
		payload.Size = ProfileSizeBasic
	}

	headers := map[string]string{
		"Authorization": 	payload.BearerToken,
		"X-Request-ID":  	requestID,
		"Accept":        	req.Header.Get("Accept"),
		"Content-Type":  	req.Header.Get("Content-Type"),
		"X-Requested-By": req.Header.Get("X-Requested-By"),
	}
	url := fmt.Sprintf("%s/profile?size=%s", userAPIEndpoint, payload.Size)

	// Send request to User API
	res, err := client.Post(req.Context(), url, payload, headers)

	if err != nil {
		logger.Error("failed to call User API", err)
		return httptypes.ErrorResponse(
			http.StatusServiceUnavailable,
			errCodes.ERR_EXTERNAL_API_CALL_FAILED,
			"failed to fetch user profile",
			err.Error(),
			requestID,
		)
	}
	if !res.IsSuccess() {
		logger.Error("failed to fetch user profile", res.Error)

		code := errCodes.ERR_EXTERNAL_API_CALL_FAILED
		switch res.Status {
			case http.StatusNotFound:
				code = errCodes.ERR_RESOURCE_NOT_FOUND
			case http.StatusUnauthorized:
				code = errCodes.ERR_INVALID_TOKEN
		}
		return httptypes.ErrorResponse(res.Status, code, res.ErrorMessage(), "", requestID)
	}

	// Parse into the appropriate struct based on size
	var userID string
	switch payload.Size {
		case ProfileSizeBasic:
			// TODO - move to proper mocking framework until User API is available
			profile := &UserProfileGetResponseBasic{
				UserID:    "12345",
				FirstName: "Test",
				LastName:  "User",
			}
			res.BodyParsed = profile
			userID = profile.UserID
		case ProfileSizeMinimal:
			// TODO - move to proper mocking framework until User API is available
			profile := &UserProfileGetResponseMinimal{
				UserID:       "12345",
				Email:        "testuser@example.com",
				Phone:        "+1234567890",
				FirstName:    "Test",
				LastName:     "User",
				PasswordHash: "$2a$10$zlk59g8mBlztf7E7EPI7r.dRKHckONCsmGo6vUv6SoGiTSbjG842K", // "password123"
			}
			res.BodyParsed = profile
			userID = profile.UserID
		case ProfileSizeFull:
			// TODO - move to proper mocking framework until User API is available
			profile := &UserProfileGetResponseFull{
				UserID:       "12345",
				Username:     "testuser",
				Email:        "testuser@example.com",
				Phone:        "+1234567890",
				FirstName:    "Test",
				MiddleInitial: "T",
				LastName:     "User",
				PasswordHash: "$2a$10$zlk59g8mBlztf7E7EPI7r.dRKHckONCsmGo6vUv6SoGiTSbjG842K", // "password123"
			}
			res.BodyParsed = profile
			userID = profile.UserID
	}

	logger.Info(fmt.Sprintf("successfully fetched profile for user: %s", userID))
	return res.Forward("successfully fetched profile", requestID)
}

// Updates the user profile for the authenticated user
func UpdateUserProfile(req *http.Request, payload *UserProfileUpdateRequest) *httptypes.APIResponse {
	requestID := req.Header.Get("X-Request-ID")
	if requestID == "" { requestID = "unknown" }

	if client == nil {
		return httptypes.ErrorResponse(
			http.StatusInternalServerError,
			errCodes.ERR_INTERNAL_SERVER,
			"http client not initialized",
			"",
			requestID,
		)
	}
	if userAPIEndpoint == "" {
		return httptypes.ErrorResponse(
			http.StatusInternalServerError,
			errCodes.ERR_INTERNAL_SERVER,
			"user api endpoint not configured",
			"",
			requestID,
		)
	}
	if payload.BearerToken == "" {
		return httptypes.ErrorResponse(
			http.StatusUnauthorized,
			errCodes.ERR_INVALID_TOKEN,
			"missing bearer token for internal API authentication",
			"",
			requestID,
		)
	}

	headers := map[string]string{
		"Authorization": 	payload.BearerToken,
		"X-Request-ID":  	requestID,
		"Accept":        	req.Header.Get("Accept"),
		"Content-Type":  	req.Header.Get("Content-Type"),
		"X-Requested-By": req.Header.Get("X-Requested-By"),
	}
	url := fmt.Sprintf("%s/profile", userAPIEndpoint)

	// Send request to User API
	res, err := client.Post(req.Context(), url, payload, headers)

	if err != nil {
		logger.Error("failed to call User API", err)
		return httptypes.ErrorResponse(
			http.StatusServiceUnavailable,
			errCodes.ERR_EXTERNAL_API_CALL_FAILED,
			"failed to update user profile",
			err.Error(),
			requestID,
		)
	}
	if !res.IsSuccess() {
		logger.Error("failed to update user profile", res.Error)

		code := errCodes.ERR_RESOURCE_UPDATE_FAILED
		switch res.Status {
			case http.StatusNotFound:
				code = errCodes.ERR_RESOURCE_NOT_FOUND
			case http.StatusUnauthorized:
				code = errCodes.ERR_INVALID_TOKEN
			case http.StatusBadRequest:
				code = errCodes.ERR_VALIDATION_FAILED
		}
		return httptypes.ErrorResponse(res.Status, code, res.ErrorMessage(), "", requestID)
	}

	logger.Info("successfully updated user profile")
	return res.Forward("successfully updated user profile", requestID)
}
