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

// External API error codes (third-party services)
const (
	ERR_EXTERNAL_API_CALL_FAILED      = "80001" // External API call failed error
	ERR_EXTERNAL_API_TIMEOUT          = "80002" // External API timeout error
	ERR_EXTERNAL_API_INVALID_RESPONSE = "80003" // External API invalid response error
)

// Internal API error codes (service-to-service M2M communication)
const (
	ERR_INTERNAL_API_CALL_FAILED      = "81001" // Internal API call failed error
	ERR_INTERNAL_API_TIMEOUT          = "81002" // Internal API timeout error
	ERR_INTERNAL_API_INVALID_RESPONSE = "81003" // Internal API invalid response error
	ERR_INTERNAL_API_AUTH_FAILED      = "81004" // Internal API authentication failed
	ERR_INTERNAL_API_UNAVAILABLE      = "81005" // Internal API service unavailable
)

// Database error codes
const (
	ERR_DB_CONNECTION_FAILED  = "90001" // Database connection failed error
	ERR_DB_QUERY_FAILED       = "90002" // Database query failed error
	ERR_DB_TRANSACTION_FAILED = "90003" // Database transaction failed error
	ERR_DB_RECORD_NOT_FOUND   = "90004" // Database record not found error
	ERR_DB_DUPLICATE_ENTRY    = "90005" // Database duplicate entry error
)

// Cache/Redis error codes
const (
	ERR_CACHE_CONNECTION_FAILED = "91001" // Cache connection failed error
	ERR_CACHE_READ_FAILED       = "91002" // Cache read failed error
	ERR_CACHE_WRITE_FAILED      = "91003" // Cache write failed error
	ERR_CACHE_DELETE_FAILED     = "91004" // Cache delete failed error
	ERR_CACHE_KEY_NOT_FOUND     = "91005" // Cache key not found error
)

// User/Account error codes
const (
	ERR_USER_NOT_FOUND          = "100001" // User not found error
	ERR_USER_ALREADY_EXISTS     = "100002" // User already exists error
	ERR_ACCOUNT_LOCKED          = "100003" // Account locked error
	ERR_ACCOUNT_SUSPENDED       = "100004" // Account suspended error
	ERR_EMAIL_NOT_VERIFIED      = "100005" // Email not verified error
	ERR_PHONE_NOT_VERIFIED      = "100006" // Phone not verified error
	ERR_INVALID_CREDENTIALS     = "100007" // Invalid username/password error
	ERR_PASSWORD_EXPIRED        = "100008" // Password expired error
	ERR_WEAK_PASSWORD           = "100009" // Weak password error
	ERR_MFA_REQUIRED            = "100010" // Multi-factor authentication required
	ERR_INVALID_MFA_CODE        = "100011" // Invalid MFA code error
)

// Session error codes
const (
	ERR_SESSION_NOT_FOUND   = "110001" // Session not found error
	ERR_SESSION_EXPIRED     = "110002" // Session expired error
	ERR_SESSION_INVALID     = "110003" // Session invalid error
	ERR_SESSION_CREATE_FAILED = "110004" // Session creation failed error
	ERR_SESSION_REVOKED     = "110005" // Session revoked error
)

// Payment/Transaction error codes
const (
	ERR_INSUFFICIENT_FUNDS       = "120001" // Insufficient funds error
	ERR_PAYMENT_DECLINED         = "120002" // Payment declined error
	ERR_PAYMENT_METHOD_INVALID   = "120003" // Payment method invalid error
	ERR_TRANSACTION_FAILED       = "120004" // Transaction failed error
	ERR_REFUND_FAILED            = "120005" // Refund failed error
	ERR_PAYMENT_PROVIDER_ERROR   = "120006" // Payment provider error
)
