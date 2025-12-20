package gin

import (
	"komodo-forge-apis-go/http/middleware/gin/auth"
	"komodo-forge-apis-go/http/middleware/gin/context"
	"komodo-forge-apis-go/http/middleware/gin/csrf"
	"komodo-forge-apis-go/http/middleware/gin/idempotency"
	"komodo-forge-apis-go/http/middleware/gin/telemetry"
)

var ContextMiddleware = context.ContextMiddleware
var AuthMiddleware = auth.AuthMiddleware
var TelemetryMiddleware = telemetry.TelemetryMiddleware
var CSRFMiddleware = csrf.CSRFMiddleware
var IdempotencyMiddleware = idempotency.IdempotencyMiddleware
