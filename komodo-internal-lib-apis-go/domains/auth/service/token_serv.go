package serviceauth

import (
	httptypes "komodo-internal-lib-apis-go/http/types"
	"net/http"
)

// Calls the internal auth service to get a service token
func GetServiceToken(req *http.Request, payload *ServiceTokenRequest) *httptypes.APIResponse{
	// TODO implement internal auth service call to get M2M token
	res := &httptypes.APIResponse{
		Status: 		http.StatusOK,
		BodyRaw:    []byte(`{"token":"mock-service-token","expires_in":3600}`),
		BodyParsed: &ServiceTokenResponse{
			Token:		 "mock-service-token",
		},
	}

	return res
}

// Calls the internal auth service to verify a service token
func VerifyServiceToken(req *http.Request, payload *VerifyTokenRequest) *httptypes.APIResponse {
	// TODO implement internal auth service call to verify M2M token
	res := &httptypes.APIResponse{
		Status: 		http.StatusOK,
		BodyRaw:    []byte(`{"active":true}`),
		BodyParsed: &VerifyTokenResponse{
			Active:		 true,
		},
	}

	return res
}
