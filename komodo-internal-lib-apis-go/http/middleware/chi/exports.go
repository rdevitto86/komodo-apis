package chi

import (
	authjwt "komodo-internal-lib-apis-go/http/middleware/chi/auth-jwt"
	"komodo-internal-lib-apis-go/http/middleware/chi/canonicalization"
	"komodo-internal-lib-apis-go/http/middleware/chi/context"
	"komodo-internal-lib-apis-go/http/middleware/chi/cors"
	"komodo-internal-lib-apis-go/http/middleware/chi/csrf"
	"komodo-internal-lib-apis-go/http/middleware/chi/entitlements"
	featureflag "komodo-internal-lib-apis-go/http/middleware/chi/feature-flag"
	"komodo-internal-lib-apis-go/http/middleware/chi/idempotency"
	ipaccess "komodo-internal-lib-apis-go/http/middleware/chi/ip-access"
	logger "komodo-internal-lib-apis-go/http/middleware/chi/logging"
	"komodo-internal-lib-apis-go/http/middleware/chi/normalization"
	postprocessor "komodo-internal-lib-apis-go/http/middleware/chi/post-processor"
	ratelimiter "komodo-internal-lib-apis-go/http/middleware/chi/rate-limiter"
	"komodo-internal-lib-apis-go/http/middleware/chi/redaction"
	evalrules "komodo-internal-lib-apis-go/http/middleware/chi/rule-validation"
	"komodo-internal-lib-apis-go/http/middleware/chi/sanitization"
	securityheaders "komodo-internal-lib-apis-go/http/middleware/chi/security-headers"
	"komodo-internal-lib-apis-go/http/middleware/chi/telemetry"
	evalheaders "komodo-internal-lib-apis-go/http/middleware/chi/validate-headers"
)

var AuthnJWTMiddleware = authjwt.AuthnJWTMiddleware
var CanonicalizeMiddleware = canonicalization.CanonicalizationMiddleware
var ContextMiddleware = context.ContextMiddleware
var CORSMiddleware = cors.CORSMiddleware
var CSRFMiddleware = csrf.CSRFMiddleware
var EntitlementsMiddleware = entitlements.EntitlementsMiddleware
var FeatureFlagMiddleware = featureflag.FeatureFlagMiddleware
var IdempotencyMiddleware = idempotency.IdempotencyMiddleware
var IPAccessMiddleware = ipaccess.IPAccessMiddleware
var LoggingMiddleware = logger.LoggingMiddleware
var NormalizationMiddleware = normalization.NormalizationMiddleware
var PostProcessorMiddleware = postprocessor.PostProcessorMiddleware
var RateLimiterMiddleware = ratelimiter.RateLimiterMiddleware
var RedactionMiddleware = redaction.RedactionMiddleware
var RuleValidationMiddleware = evalrules.RuleValidationMiddleware
var SanitizationMiddleware = sanitization.SanitizationMiddleware
var SecurityHeadersMiddleware = securityheaders.SecurityHeadersMiddleware
var TelemetryMiddleware = telemetry.TelemetryMiddleware
var ValidateHeadersMiddleware = evalheaders.ValidateHeadersMiddleware
