package serviceauth

import (
	httpclient "komodo-internal-lib-apis-go/http/client"
	httptypes "komodo-internal-lib-apis-go/http/types"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"net/http"
	"os"
)

var (
	client *httpclient.HTTPClient
	authAPIEndpoint string
)

func init() {
	if client = httpclient.GetInstance(); client == nil {
		logger.Error("failed to initialize http client")
	}
	if authAPIEndpoint = os.Getenv("URL_AUTH_SERVICE_API"); authAPIEndpoint == "" {
		logger.Error("auth api endpoint environment variable not set")
	}
}

// Calls the internal auth service to get a service token
func TokenGenerate(req *http.Request, payload *TokenGenerateRequest) *httptypes.APIResponse{
	// TODO implement internal auth service call to get M2M token

	res := &httptypes.APIResponse{
		Status: 		http.StatusOK,
		BodyRaw:    []byte(`{"token":"mock-service-token","expires_in":3600}`),
		BodyParsed: &TokenGenerateResponse{
			Token:		 "mock-service-token",
		},
	}

	var tokenRes TokenGenerateResponse
	res = res.ParseBody(&tokenRes)

	return res
}

// Calls the internal auth service to verify a service token
func TokenVerify(req *http.Request) *httptypes.APIResponse {
	// TODO implement internal auth service call to verify M2M token

	res := &httptypes.APIResponse{
		Status:     http.StatusOK,
		BodyRaw:    []byte(`{"active":true,"client_id":"mock-service","scope":"read write"}`),
		BodyParsed: &TokenVerifyResponse{
			Active:   true,
			ClientID: "mock-service",
			Scope:    "read write",
		},
	}

	return res
}

func TokenRevoke(req *http.Request) *httptypes.APIResponse {
	// TODO implement internal auth service call to revoke token

	res := &httptypes.APIResponse{
		Status: 		http.StatusOK,
		BodyRaw:    []byte(`{"revoked":true}`),
		BodyParsed: &TokenRevokeResponse{
			Revoked:		 true,
		},
	}

	return res
}
