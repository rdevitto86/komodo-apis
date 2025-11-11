package errors

type ErrorStandard struct {
	Status 		int    `json:"status"`
	Code    	string `json:"code"`
	Message 	string `json:"message"`
	RequestId string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type ErrorVerbose struct {
	Status 		int    `json:"status"`
	Code    	string `json:"code"`
	Message 	string `json:"message"`
	RequestId string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	APIName		string `json:"api_name,omitempty"`
	APIError	any 	 `json:"api_error,omitempty"`
}

// Generic error codes
const (
	ERR_INTERNAL_SERVER 		= "10001" // Generic internal server error
	ERR_SERVICE_UNAVAILABLE = "10002" // Service unavailable error
	ERR_TIMEOUT             = "10003" // Timeout error
	ERR_INVALID_REQUEST     = "10004" // Invalid request error
	ERR_METHOD_NOT_ALLOWED  = "10005" // Method not allowed error
	ERR_NOT_FOUND           = "10006" // Not found error
	ERR_CONFLICT            = "10007" // Conflict error
	ERR_PANIC_RECOVERY      = "10008" // Panic recovery error
)

// Authentication/Authorization error codes
const (
	ERR_INVALID_CLIENT_CREDENTIALS = "20001" // Invalid client credentials
	ERR_INVALID_GRANT_TYPE          = "20002" // Invalid grant type
	ERR_INVALID_SCOPE               = "20003" // Invalid scope
	ERR_INVALID_TOKEN               = "20004" // Invalid token
	ERR_EXPIRED_TOKEN               = "20005" // Expired token
	ERR_UNAUTHORIZED_CLIENT         = "20006" // Unauthorized client
	ERR_UNSUPPORTED_GRANT_TYPE      = "20007" // Unsupported grant type
	ERR_UNSUPPORTED_RESPONSE_TYPE   = "20008" // Unsupported response type
	ERR_INVALID_REDIRECT_URI        = "20009" // Invalid redirect URI
	ERR_ACCESS_DENIED               = "20010" // Access denied
	ERR_INSUFFICIENT_SCOPE          = "20011" // Insufficient scope
)

// Validation error codes
const (
	ERR_VALIDATION_FAILED = "30001" // Validation failed error
)

// Resource service error codes
const (
	ERR_RESOURCE_NOT_FOUND       = "40001" // Resource not found error
	ERR_RESOURCE_CREATION_FAILED = "40002" // Resource creation failed error
	ERR_RESOURCE_UPDATE_FAILED   = "40003" // Resource update failed error
	ERR_RESOURCE_DELETION_FAILED = "40004" // Resource deletion failed error
)

// Rate limiting error codes
const (
	ERR_RATE_LIMIT_EXCEEDED = "50001" // Rate limit exceeded error
)

// File handling error codes
const (
	ERR_FILE_NOT_FOUND       = "60001" // File not found error
	ERR_FILE_READ_FAILED     = "60002" // File read failed error
	ERR_FILE_WRITE_FAILED    = "60003" // File write failed error
	ERR_FILE_UPLOAD_FAILED   = "60004" // File upload failed error
	ERR_FILE_DOWNLOAD_FAILED = "60005" // File download failed error
)

// Network error codes
const (
	ERR_NETWORK_UNREACHABLE 	= "70001" // Network unreachable error
	ERR_DNS_RESOLUTION_FAILED = "70002" // DNS resolution failed error
	ERR_CONNECTION_TIMEOUT    = "70003" // Connection timeout error
	ERR_CONNECTION_REFUSED    = "70004" // Connection refused error
)

// External API error codes
const (
	ERR_EXTERNAL_API_CALL_FAILED = "80001" // External API call failed error
	ERR_EXTERNAL_API_TIMEOUT     = "80002" // External API timeout error
	ERR_EXTERNAL_API_INVALID_RESPONSE = "80003" // External API invalid response error
)

// Database error codes
const (
	ERR_DB_CONNECTION_FAILED = "90001" // Database connection failed error
	ERR_DB_QUERY_FAILED      = "90002" // Database query failed error
	ERR_DB_TRANSACTION_FAILED = "90003" // Database transaction failed error
)

// Business logic error codes
const (
	ERR_INSUFFICIENT_FUNDS    = "100001" // Insufficient funds error
)