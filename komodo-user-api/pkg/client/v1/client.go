package v1

import (
	"context"
	"fmt"

	ctxKeys "komodo-forge-sdk-go/http/context"
	httpReq "komodo-forge-sdk-go/http/request"
	httpRes "komodo-forge-sdk-go/http/response"
	"net/http"
)

const userAPIEndpoint string = "http://localhost:8080"

// Retrieves the user profile for the authenticated user
func GetUserProfile(ctx context.Context, payload UserProfileGetRequest) (UserProfileGetResponse, error) {
	var res UserProfileGetResponse

	if payload.BearerToken == "" {
		return res, fmt.Errorf("missing bearer token for internal API authentication")
	}
	if payload.Size == "" {
		payload.Size = ProfileSizeBasic
	}

	reqId := ctx.Value(ctxKeys.REQUEST_ID_KEY).(string)
	if reqId == "" { reqId = "unknown" }

	url := fmt.Sprintf("%s/profile?size=%s", userAPIEndpoint, payload.Size)
	headers := map[string]string{
		"Authorization": payload.BearerToken,
		"X-Request-ID": reqId,
		"Accept": "application/json",
		"Content-Type": "application/json",
		"X-Requested-By": "komodo-user-api",
	}

	req, err := httpReq.NewRequest("POST", url, payload, headers, ctx)
	if err != nil { return res, err }

	// Send request to User API
	reply, err := http.DefaultClient.Do(req)

	if err != nil { return res, err }
	defer reply.Body.Close()

	if !httpRes.IsSuccess(reply.StatusCode) {
		return res, fmt.Errorf("failed to get user profile")
	}

	switch payload.Size {
		case ProfileSizeBasic:
			res = UserProfileGetResponse{
				UserID:    "12345",
				FirstName: "Test",
				LastName:  "User",
			}
		case ProfileSizeMinimal:
			res = UserProfileGetResponse{
				UserID:       "12345",
				Email:        "testuser@example.com",
				Phone:        "+1234567890",
				FirstName:    "Test",
				LastName:     "User",
				PasswordHash: "$2a$10$zlk59g8mBlztf7E7EPI7r.dRKHckONCsmGo6vUv6SoGiTSbjG842K", // "password123"
			}
		case ProfileSizeFull:
			res = UserProfileGetResponse{
				UserID:       "12345",
				Username:     "testuser",
				Email:        "testuser@example.com",
				Phone:        "+1234567890",
				FirstName:    "Test",
				MiddleInitial: "T",
				LastName:     "User",
				PasswordHash: "$2a$10$zlk59g8mBlztf7E7EPI7r.dRKHckONCsmGo6vUv6SoGiTSbjG842K", // "password123"
			}
	}
	return res, nil
}

// Updates the user profile for the authenticated user
func UpdateUserProfile(ctx context.Context, payload UserProfileUpdateRequest) (UserProfileUpdateResponse, error) {
	var res UserProfileUpdateResponse

	if payload.BearerToken == "" {
		return res, fmt.Errorf("missing bearer token for internal API authentication")
	}
	if payload.UserID == "" {
		return res, fmt.Errorf("missing user id")
	}

	requestID := ctx.Value(ctxKeys.REQUEST_ID_KEY).(string)
	if requestID == "" { requestID = "unknown" }

	url := fmt.Sprintf("%s/profile", userAPIEndpoint)
	headers := map[string]string{
		"Authorization": 	payload.BearerToken,
		"X-Request-ID":  	requestID,
		"Accept":        	"application/json",
		"Content-Type":  	"application/json",
		"X-Requested-By": "komodo-user-api",
	}

	req, err := httpReq.NewRequest("POST", url, payload, headers, ctx)
	if err != nil { return res, err }

	// Send request to User API
	reply, err := http.DefaultClient.Do(req)

	if err != nil { return res, err }
	defer reply.Body.Close()

	if !httpRes.IsSuccess(reply.StatusCode) {
		return res, fmt.Errorf("failed to update user profile")
	}

	res = UserProfileUpdateResponse{
		UserID: payload.UserID,
	}
	return res, nil
}
