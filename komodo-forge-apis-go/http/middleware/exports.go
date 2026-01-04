package middleware

import (
	"komodo-forge-apis-go/http/middleware/auth"
	clienttype "komodo-forge-apis-go/http/middleware/client-type"
	"komodo-forge-apis-go/http/middleware/context"
	"komodo-forge-apis-go/http/middleware/cors"
	"komodo-forge-apis-go/http/middleware/csrf"
	"komodo-forge-apis-go/http/middleware/idempotency"
	ipaccess "komodo-forge-apis-go/http/middleware/ip-access"
	"komodo-forge-apis-go/http/middleware/normalization"
	ratelimiter "komodo-forge-apis-go/http/middleware/rate-limiter"
	"komodo-forge-apis-go/http/middleware/redaction"
	requestid "komodo-forge-apis-go/http/middleware/request-id"
	rulevalidation "komodo-forge-apis-go/http/middleware/rule-validation"
	"komodo-forge-apis-go/http/middleware/sanitization"
	securityheaders "komodo-forge-apis-go/http/middleware/security-headers"
	telemetry "komodo-forge-apis-go/http/middleware/telemetry"
)

var (
	AuthMiddleware = auth.AuthMiddleware
	ClientTypeMiddleware = clienttype.ClientTypeMiddleware
	ContextMiddleware = context.ContextMiddleware
	CORSMiddleware = cors.CORSMiddleware
	CSRFMiddleware = csrf.CSRFMiddleware
	IdempotencyMiddleware = idempotency.IdempotencyMiddleware
	IPAccessMiddleware = ipaccess.IPAccessMiddleware
	NormalizationMiddleware = normalization.NormalizationMiddleware
	RateLimiterMiddleware = ratelimiter.RateLimiterMiddleware
	RedactionMiddleware = redaction.RedactionMiddleware
	RequestIDMiddleware = requestid.RequestIDMiddleware
	RuleValidationMiddleware = rulevalidation.RuleValidationMiddleware
	SanitizationMiddleware = sanitization.SanitizationMiddleware
	SecurityHeadersMiddleware = securityheaders.SecurityHeadersMiddleware
	TelemetryMiddleware = telemetry.TelemetryMiddleware
)
