package errorutils

import (
	"encoding/json"
	"net/http"
	"time"

	errorType "komodo-internal-lib-apis-go/types/error"
)

// Writes a standardized error response
func WriteErrorResponse(wtr http.ResponseWriter, status int, message string, errCode string, requestId ...string) {
	if len(requestId) == 0 {
		requestId = append(requestId, "unknown")
	}

	wtr.WriteHeader(status)
	json.NewEncoder(wtr).Encode(errorType.ErrorStandard{
		Status: 	 status,
		Code:      errCode,
		Message:   message,
		RequestId: requestId[0],
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Writes a standardized error response with more details
func WriteErrorVerboseResponse(wtr http.ResponseWriter, status int, message string, errCode string, requestId ...string) {
	if len(requestId) == 0 {
		requestId = append(requestId, "unknown")
	}

	wtr.WriteHeader(status)
	json.NewEncoder(wtr).Encode(errorType.ErrorStandard{
		Status: 	 status,
		Code:      errCode,
		Message:   message,
		RequestId: requestId[0],
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}