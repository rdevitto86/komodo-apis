package user

import (
	"fmt"

	"komodo-internal-lib-apis-go/common/errors"
	httpclient "komodo-internal-lib-apis-go/http/client"
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
	if userAPIEndpoint = os.Getenv("USER_API_URL"); userAPIEndpoint == "" {
		logger.Error("user api endpoint environment variable not set")
	}
}

// Retrieves the user profile for the authenticated user
func GetUserProfile(req *http.Request, payload *UserProfileGetRequest) *httptypes.APIResponse {
	requestID := req.Header.Get("X-Request-ID")
	if requestID == "" { requestID = "unknown" }

	if client == nil {
		return httptypes.ErrorResponse(
			http.StatusInternalServerError,
			errors.ERR_INTERNAL_SERVER,
			"http client not initialized",
			"",
			requestID,
		)
	}
	if userAPIEndpoint == "" {
		return httptypes.ErrorResponse(
			http.StatusInternalServerError,
			errors.ERR_INTERNAL_SERVER,
			"user api endpoint not configured",
			"",
			requestID,
		)
	}

	bearer := req.Header.Get("Authorization")
	if bearer == "" {
		return httptypes.ErrorResponse(
			http.StatusUnauthorized,
			errors.ERR_INVALID_TOKEN,
			"missing authorization token in request",
			"",
			requestID,
		)
	}

	if payload.Size == "" {
		payload.Size = ProfileSizeBasic
	}

	headers := map[string]string{
		"Authorization": 	bearer,
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
			errors.ERR_EXTERNAL_API_CALL_FAILED,
			"failed to fetch user profile",
			err.Error(),
			requestID,
		)
	}

	var profile UserProfileGetResponse
	res = res.ParseBody(&profile)

	if !res.IsSuccess() {
		logger.Error("failed to fetch user profile", res.Error)
		
		code := errors.ERR_EXTERNAL_API_CALL_FAILED
		switch res.Status {
			case http.StatusNotFound:
				code = errors.ERR_RESOURCE_NOT_FOUND
			case http.StatusUnauthorized:
				code = errors.ERR_INVALID_TOKEN
		}
		return httptypes.ErrorResponse(res.Status, code, res.ErrorMessage(), "", requestID)
	}

	// TODO - move into proper mocking framework until User API is available
	res.BodyParsed = UserProfileGetResponse{
		UserID:       "12345",
		FirstName:    "Test",
		LastName:     "User",
		Email:        "testuser@example.com",
		Phone: 				"+1234567890",		
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "password123"
	}

	logger.Info(fmt.Sprintf("successfully fetched profile for user: %s", profile.UserID))
	return res.Forward("successfully fetched profile", requestID)
}

// Updates the user profile for the authenticated user
func UpdateUserProfile(req *http.Request, payload *UserProfileUpdateRequest) *httptypes.APIResponse {
	requestID := req.Header.Get("X-Request-ID")
	if requestID == "" { requestID = "unknown" }

	if client == nil {
		return httptypes.ErrorResponse(
			http.StatusInternalServerError,
			errors.ERR_INTERNAL_SERVER,
			"http client not initialized",
			"",
			requestID,
		)
	}
	if userAPIEndpoint == "" {
		return httptypes.ErrorResponse(
			http.StatusInternalServerError,
			errors.ERR_INTERNAL_SERVER,
			"user api endpoint not configured",
			"",
			requestID,
		)
	}

	bearer := req.Header.Get("Authorization")
	if bearer == "" {
		return httptypes.ErrorResponse(
			http.StatusUnauthorized,
			errors.ERR_INVALID_TOKEN,
			"missing authorization token in request",
			"",
			requestID,
		)
	}

	headers := map[string]string{
		"Authorization": 	bearer,
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
			errors.ERR_EXTERNAL_API_CALL_FAILED,
			"failed to update user profile",
			err.Error(),
			requestID,
		)
	}
	if !res.IsSuccess() {
		logger.Error("failed to update user profile", res.Error)

		code := errors.ERR_RESOURCE_UPDATE_FAILED
		switch res.Status {
			case http.StatusNotFound:
				code = errors.ERR_RESOURCE_NOT_FOUND
			case http.StatusUnauthorized:
				code = errors.ERR_INVALID_TOKEN
			case http.StatusBadRequest:
				code = errors.ERR_VALIDATION_FAILED
		}
		return httptypes.ErrorResponse(res.Status, code, res.ErrorMessage(), "", requestID)
	}

	logger.Info("successfully updated user profile")
	return res.Forward("successfully updated user profile", requestID)
}
