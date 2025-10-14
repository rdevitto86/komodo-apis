package idempotency

import (
	"context"
	elasticacheClient "komodo-internal-lib-apis-go/aws/elasticache"
	"komodo-internal-lib-apis-go/config"
	hdrSrv "komodo-internal-lib-apis-go/services/headers/eval"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	hdrTypes "komodo-internal-lib-apis-go/types/headers"
	"net/http"
	"sync"
	"time"
)

var idemStore sync.Map

type idemCtxKey string
const IdempotencyValidCtxKey idemCtxKey = "Idempotency-Key_valid"

func IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// Only guard unsafe, state-changing methods
		switch req.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				next.ServeHTTP(wtr, req)
				return
		}

		key := req.Header.Get("Idempotency-Key")

		if !hdrSrv.ValidateHeaderValue(hdrTypes.HEADER_IDEMPOTENCY, req) {
			logger.Error("invalid idempotency key: " + key, req)
			http.Error(wtr, "Invalid Idempotency-Key", http.StatusBadRequest)
			return
		} 

		// Load existing entry
		if exp, ok := idemStore.Load(key); ok {
			// If expired, evict and continue; else reject as duplicate
			if until, _ := exp.(int64); until > time.Now().Unix() {
				wtr.Header().Set("Idempotency-Replayed", "true")
				logger.Error("duplicate request: " + key, req)
				http.Error(wtr, "Duplicate request", http.StatusConflict)
				return
			}
			idemStore.Delete(key)
		}

		req = req.WithContext(context.WithValue(
			req.Context(), IdempotencyValidCtxKey, true,
		))

		// Store key with expiration
		elasticacheClient.SetCacheItem("idem-" + key, "1", getIdemTTL())

		next.ServeHTTP(wtr, req)
	})
}

func getIdemTTL() int64 {
	// Parse env only once per process would be ideal, but keep simple/fast
	if ttl := config.GetConfigValue("IDEMPOTENCY_TTL_SEC"); ttl != "" {
		if dur, err := time.ParseDuration(ttl + "s"); err == nil {
			if dur <= 0 { return 300 }
			return int64(dur.Seconds())
		}
	}
	return 300 // default: 5 minutes
}
