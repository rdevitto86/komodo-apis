package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Writes a standardized error response
func WriteErrorResponse(wtr http.ResponseWriter, req *http.Request, status int, message string, errCode string) {
	key := http.CanonicalHeaderKey("X-Request-ID")
	if len(req.Header[key]) == 0 {
		req.Header[key] = []string{"unknown"}
	}

	wtr.WriteHeader(status)
	json.NewEncoder(wtr).Encode(ErrorStandard{
		Status: 	 status,
		Code:      errCode,
		Message:   message,
		RequestId: req.Header[key][0],
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Writes a standardized error response with more details
func WriteErrorVerboseResponse(wtr http.ResponseWriter, req *http.Request, status int, message string, errCode string, apiError any) {
	key := http.CanonicalHeaderKey("X-Request-ID")
	if len(req.Header[key]) == 0 {
		req.Header[key] = []string{"unknown"}
	}

	wtr.WriteHeader(status)
	json.NewEncoder(wtr).Encode(ErrorVerbose{
		Status: 	 status,
		Code:      errCode,
		Message:   message,
		APIName:   req.URL.Path,
		APIError:  fmt.Sprintf("%v", apiError),
		RequestId: req.Header[key][0],
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
