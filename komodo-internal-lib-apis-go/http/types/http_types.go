package types

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIResponse struct {
	Status 			int     		`json:"status"`
	Headers			http.Header `json:"headers,omitempty"`
	Message   	string      `json:"message,omitempty"`
	BodyRaw     []byte			`json:"body_raw,omitempty"`
	BodyParsed  interface{} `json:"body,omitempty"`
	Error     	*ErrorInfo  `json:"error,omitempty"`
	RequestID 	string      `json:"request_id,omitempty"`
	Timestamp 	string      `json:"timestamp"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Creates a success APIResponse with data
// If data is []byte, it goes to BodyRaw. If it's a struct/object, it goes to BodyParsed and gets marshaled to BodyRaw
func SuccessResponse(status int, data interface{}, message string, requestID string) *APIResponse {
	resp := &APIResponse{
		Status:    status,
		Message:   message,
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if data != nil {
		switch v := data.(type) {
			case []byte:
				resp.BodyRaw = v
			default:
				resp.BodyParsed = data
				if raw, err := json.Marshal(data); err == nil {
					resp.BodyRaw = raw
				}
		}
	}

	return resp
}

// Creates an error APIResponse
func ErrorResponse(status int, code string, message string, details string, requestID string) *APIResponse {
	return &APIResponse{
		Status: status,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func (res *APIResponse) IsSuccess() bool {
	return res.Status >= 200 && res.Status < 300 && res.Error == nil
}

func (res *APIResponse) IsError() bool {
	return res.Status >= 400 || res.Error != nil
}

func (res *APIResponse) ParseBody(parsed interface{}) *APIResponse {
	if len(res.BodyRaw) == 0 {
		return ErrorResponse(http.StatusInternalServerError, "", "Failed to parse body", "", res.RequestID)
	}
	if err := json.Unmarshal(res.BodyRaw, parsed); err != nil {
		return ErrorResponse(http.StatusInternalServerError, "", "Failed to parse body", err.Error(), res.RequestID)
	}
	res.BodyParsed = parsed
	return res
}

func (res *APIResponse) Stringify() string {
	return string(res.BodyRaw)
}

func (res *APIResponse) ErrorMessage() string {
	if res.Error != nil { return res.Error.Message }
	return ""
}

// Creates a new APIResponse by forwarding/copying this response with optional overrides
// Useful when an API call returns an APIResponse that you want to pass through or modify slightly
func (res *APIResponse) Forward(overrideMessage string, overrideRequestID string) *APIResponse {
	resp := &APIResponse{
		Status:     res.Status,
		Headers:    res.Headers,
		Message:    res.Message,
		BodyRaw:    res.BodyRaw,
		BodyParsed: res.BodyParsed,
		Error:      res.Error,
		RequestID:  res.RequestID,
		Timestamp:  res.Timestamp,
	}

	if overrideMessage != "" { resp.Message = overrideMessage }
	if overrideRequestID != "" { resp.RequestID = overrideRequestID }

	return resp
}

// Sets the request ID on this response (useful for chaining)
func (res *APIResponse) WithRequestID(requestID string) *APIResponse {
	res.RequestID = requestID
	return res
}
