package middleware

import (
	"context"
	"komodo-auth-api/internal/httpapi/utils"
	"komodo-auth-api/internal/logger"
	"komodo-auth-api/internal/thirdparty/aws"
	"net/http"
	"os"
	"sync"
	"time"
)

var idemStore sync.Map

type idemCtxKey string
const IdempotencyValidCtxKey idemCtxKey = "Idempotency-Key_valid"

func IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		rules := req.Context().Value(ValidationRuleKey).(ValidationRule)
		isValid := false

		// Only guard unsafe, state-changing methods
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				break
			default:
				next.ServeHTTP(wtr, req)
				return
		}

		// Accept common header names
		key := req.Header.Get("Idempotency-Key")

		// Check if idempotency is required for this route
		if rules.Headers != nil && !rules.Headers["Idempotency-Key"].Required {
			next.ServeHTTP(wtr, req)
			return
		}

		// If no key provided, proceed normally
		if key == "" {
			next.ServeHTTP(wtr, req)
			return
		}

		// Fast validation: keep keys compact and simple
		if !utils.IsValidIdempotencyKey(key) {
			logger.Error("invalid idempotency key: " + key, req)
			http.Error(wtr, "Invalid Idempotency-Key", http.StatusBadRequest)
			return
		}

		isValid = true
		now := time.Now().Unix()

		// Load existing entry
		if exp, ok := idemStore.Load(key); ok {
			// If expired, evict and continue; else reject as duplicate
			if until, _ := exp.(int64); until > now {
				wtr.Header().Set("Idempotency-Replayed", "true")
				logger.Error("duplicate request: " + key, req)
				http.Error(wtr, "Duplicate request", http.StatusConflict)
				return
			}
			idemStore.Delete(key)
		}

		// Mark request as valid for downstream handlers
		req = req.WithContext(context.WithValue(req.Context(), IdempotencyValidCtxKey, isValid))

		// Store key with expiration
		aws.SetCacheItem("idem-"+key, "1", getIdemTTL())

		next.ServeHTTP(wtr, req)
	})
}

func getIdemTTL() int64 {
	// Parse env only once per process would be ideal, but keep simple/fast
	if ttl := os.Getenv("IDEMPOTENCY_TTL_SEC"); ttl != "" {
		if dur, err := time.ParseDuration(ttl + "s"); err == nil {
			if dur <= 0 { return 300 }
			return int64(dur.Seconds())
		}
	}
	return 300 // default: 5 minutes
}
