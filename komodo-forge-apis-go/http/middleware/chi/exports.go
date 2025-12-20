package chi

import (
	"komodo-forge-apis-go/http/middleware/chi/auth"
	"komodo-forge-apis-go/http/middleware/chi/canonicalization"
	clientType "komodo-forge-apis-go/http/middleware/chi/client-type"
	"komodo-forge-apis-go/http/middleware/chi/context"
	"komodo-forge-apis-go/http/middleware/chi/cors"
	"komodo-forge-apis-go/http/middleware/chi/csrf"
	"komodo-forge-apis-go/http/middleware/chi/entitlements"
	featureflag "komodo-forge-apis-go/http/middleware/chi/feature-flag"
	"komodo-forge-apis-go/http/middleware/chi/idempotency"
	ipaccess "komodo-forge-apis-go/http/middleware/chi/ip-access"
	logger "komodo-forge-apis-go/http/middleware/chi/logging"
	"komodo-forge-apis-go/http/middleware/chi/normalization"
	ratelimiter "komodo-forge-apis-go/http/middleware/chi/rate-limiter"
	"komodo-forge-apis-go/http/middleware/chi/redaction"
	evalrules "komodo-forge-apis-go/http/middleware/chi/rule-validation"
	"komodo-forge-apis-go/http/middleware/chi/sanitization"
	securityheaders "komodo-forge-apis-go/http/middleware/chi/security-headers"
	"komodo-forge-apis-go/http/middleware/chi/telemetry"
	evalheaders "komodo-forge-apis-go/http/middleware/chi/validate-headers"
)

var AuthMiddleware = auth.AuthMiddleware
var CanonicalizeMiddleware = canonicalization.CanonicalizationMiddleware
var ClientTypeMiddleware = clientType.ClientTypeMiddleware
var ContextMiddleware = context.ContextMiddleware
var CORSMiddleware = cors.CORSMiddleware
var CSRFMiddleware = csrf.CSRFMiddleware
var EntitlementsMiddleware = entitlements.EntitlementsMiddleware
var FeatureFlagMiddleware = featureflag.FeatureFlagMiddleware
var IdempotencyMiddleware = idempotency.IdempotencyMiddleware
var IPAccessMiddleware = ipaccess.IPAccessMiddleware
var LoggingMiddleware = logger.LoggingMiddleware
var NormalizationMiddleware = normalization.NormalizationMiddleware
var RateLimiterMiddleware = ratelimiter.RateLimiterMiddleware
var RedactionMiddleware = redaction.RedactionMiddleware
var RuleValidationMiddleware = evalrules.RuleValidationMiddleware
var SanitizationMiddleware = sanitization.SanitizationMiddleware
var SecurityHeadersMiddleware = securityheaders.SecurityHeadersMiddleware
var TelemetryMiddleware = telemetry.TelemetryMiddleware
var ValidateHeadersMiddleware = evalheaders.ValidateHeadersMiddleware
