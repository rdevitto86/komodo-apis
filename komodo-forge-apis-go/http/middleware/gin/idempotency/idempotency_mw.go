package idempotency

import (
	"komodo-forge-apis-go/config"
	ctxKeys "komodo-forge-apis-go/http/common/context"
	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/gin"
	hdrSrv "komodo-forge-apis-go/http/headers/eval"
	hdrTypes "komodo-forge-apis-go/http/types"
	httpUtils "komodo-forge-apis-go/http/utils/http"
	logger "komodo-forge-apis-go/logging/runtime"

	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const DEFAULT_IDEM_TTL_SEC int64 = 300 // 5 minutes

// In-memory store for idempotency keys (for local development)
// In production, this should be replaced with Redis/ElastiCache
var idemStore sync.Map

// IdempotencyMiddleware prevents duplicate requests by tracking idempotency keys
// Only applies to state-changing methods (POST, PUT, PATCH, DELETE)
// API clients are exempt from idempotency validation
func IdempotencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only guard unsafe, state-changing methods
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		// Check if client type is already set in context
		clientTypeVal, exists := c.Get(string(ctxKeys.CLIENT_TYPE_KEY))
		var clientType string
		if exists {
			clientType, _ = clientTypeVal.(string)
		} else {
			clientType = httpUtils.GetClientType(c.Request)
		}

		// API clients are exempt from idempotency validation
		if clientType == "api" {
			c.Set(string(ctxKeys.IDEMPOTENCY_VALID_KEY), true)
			c.Next()
			return
		}

		// Browser client - require idempotency key
		key := c.GetHeader("Idempotency-Key")

		if ok, err := hdrSrv.ValidateHeaderValue(hdrTypes.HEADER_IDEMPOTENCY, c.Request); !ok || err != nil {
			logger.Error("invalid idempotency key for browser client: " + key)
			errors.WriteErrorResponse(c, http.StatusBadRequest, "invalid idempotency key", errCodes.ERR_INVALID_REQUEST)
			c.Abort()
			return
		}

		// Load existing entry from in-memory store
		if exp, ok := idemStore.Load(key); ok {
			// If expired, evict and continue; else reject as duplicate
			if until, _ := exp.(int64); until > time.Now().Unix() {
				c.Header("Idempotency-Replayed", "true")
				logger.Error("duplicate request: " + key)
				errors.WriteErrorResponse(c, http.StatusConflict, "duplicate request", errCodes.ERR_ACCESS_DENIED)
				c.Abort()
				return
			}
			idemStore.Delete(key)
		}

		// Set validation flag in context
		c.Set(string(ctxKeys.IDEMPOTENCY_VALID_KEY), true)

		// Store key with expiration in in-memory map
		// TODO: In production, use Redis/ElastiCache instead
		// elasticacheClient.SetCacheItem("idem-" + key, "1", getIdemTTL())
		expiresAt := time.Now().Unix() + getIdemTTL()
		idemStore.Store(key, expiresAt)

		c.Next()
	}
}

// getIdemTTL returns the TTL for idempotency keys from config or default
func getIdemTTL() int64 {
	// Parse env only once per process would be ideal, but keep simple/fast
	if ttl := config.GetConfigValue("IDEMPOTENCY_TTL_SEC"); ttl != "" {
		if dur, err := time.ParseDuration(ttl + "s"); err == nil {
			if dur <= 0 {
				return 300
			}
			return int64(dur.Seconds())
		}
	}
	return DEFAULT_IDEM_TTL_SEC
}
