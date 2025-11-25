package gin

import (
	"komodo-internal-lib-apis-go/http/middleware/gin/context"
	serviceauth "komodo-internal-lib-apis-go/http/middleware/gin/service-auth"
	"komodo-internal-lib-apis-go/http/middleware/gin/telemetry"
)

var ContextMiddleware = context.ContextMiddleware
var ServiceAuthMiddleware = serviceauth.ServiceAuthMiddleware
var TelemetryMiddleware = telemetry.TelemetryMiddleware
