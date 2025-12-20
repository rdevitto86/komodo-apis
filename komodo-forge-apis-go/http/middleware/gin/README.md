# Gin Middlewares

This directory contains Gin-compatible HTTP middlewares for the Komodo API platform.

## Available Middlewares

### Session Middleware
**Location:** `session/session_mw.go`

Extracts and validates user sessions from incoming requests. Supports two authentication methods:
1. **Cookie-based sessions** - Extracts `session_id` from cookies and validates against Redis/ElastiCache
2. **JWT Bearer tokens** - Validates JWT tokens from the `Authorization` header

**Usage:**
```go
import ginmw "komodo-forge-apis-go/http/middleware/gin"

router.Use(ginmw.SessionMiddleware())
```

**Context Keys Set:**
- `USER_ID_KEY` - The authenticated user's ID

**Notes:**
- Redis/ElastiCache lookups are commented out for local development
- Currently returns mock user ID `"12345"` for testing

---

### CSRF Middleware
**Location:** `csrf/csrf_mw.go`

Validates CSRF tokens for state-changing requests (POST, PUT, PATCH, DELETE) from browser clients.

**Usage:**
```go
import ginmw "komodo-forge-apis-go/http/middleware/gin"

router.Use(ginmw.CSRFMiddleware())
```

**Behavior:**
- **API clients** - Exempt from CSRF validation
- **Browser clients** - Must include valid `X-CSRF-Token` header for state-changing requests
- **Safe methods** (GET, HEAD, OPTIONS) - No validation required

**Context Keys Set:**
- `CSRF_TOKEN_KEY` - The CSRF token value
- `CSRF_VALID_KEY` - Boolean indicating if CSRF validation passed

---

### Idempotency Middleware
**Location:** `idempotency/idempotency_mw.go`

Prevents duplicate requests by tracking idempotency keys for state-changing operations.

**Usage:**
```go
import ginmw "komodo-forge-apis-go/http/middleware/gin"

router.Use(ginmw.IdempotencyMiddleware())
```

**Behavior:**
- **API clients** - Exempt from idempotency validation
- **Browser clients** - Must include valid `Idempotency-Key` header for state-changing requests
- **Safe methods** (GET, HEAD, OPTIONS) - No validation required
- **Duplicate detection** - Returns `409 Conflict` if key is reused within TTL window

**Configuration:**
- Default TTL: 300 seconds (5 minutes)
- Override via env: `IDEMPOTENCY_TTL_SEC`

**Context Keys Set:**
- `IDEMPOTENCY_VALID_KEY` - Boolean indicating if idempotency validation passed

**Response Headers:**
- `Idempotency-Replayed: true` - Set when a duplicate request is detected

**Notes:**
- Uses in-memory storage for local development
- Redis/ElastiCache integration commented out for production use

---

## Common Patterns

### Full Middleware Stack Example
```go
import (
    "github.com/gin-gonic/gin"
    ginmw "komodo-forge-apis-go/http/middleware/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.New()
    
    // Apply middlewares in order
    router.Use(ginmw.ContextMiddleware())
    router.Use(ginmw.TelemetryMiddleware())
    router.Use(ginmw.SessionMiddleware())
    router.Use(ginmw.CSRFMiddleware())
    router.Use(ginmw.IdempotencyMiddleware())
    
    // Define routes...
    
    return router
}
```

### Selective Application
```go
// Apply only to specific route groups
api := router.Group("/api")
api.Use(ginmw.SessionMiddleware())
api.Use(ginmw.CSRFMiddleware())
{
    api.POST("/orders", createOrder)
    api.PUT("/orders/:id", updateOrder)
}
```

---

## Production Deployment

Before deploying to production, uncomment and configure the following:

### Session Middleware
- Uncomment Redis/ElastiCache session lookup
- Remove mock user ID return

### Idempotency Middleware
- Uncomment Redis/ElastiCache storage
- Remove in-memory `sync.Map` usage
- Configure `IDEMPOTENCY_TTL_SEC` environment variable

---

## Testing

Each middleware includes comprehensive test coverage. Run tests with:

```bash
go test ./http/middleware/gin/...
```

---

## Error Codes

All middlewares use standardized error codes from `komodo-forge-apis-go/http/common/errors`:

- `ERR_SESSION_NOT_FOUND` - No valid session or token found
- `ERR_SESSION_EXPIRED` - Session expired or invalid
- `ERR_INVALID_REQUEST` - Invalid CSRF token or idempotency key
- `ERR_ACCESS_DENIED` - Duplicate request detected
