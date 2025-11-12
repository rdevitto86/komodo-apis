package idempotency

import (
	"context"
	elasticacheClient "komodo-internal-lib-apis-go/aws/elasticache"
	ctxKeys "komodo-internal-lib-apis-go/common/context"
	"komodo-internal-lib-apis-go/common/errors"
	"komodo-internal-lib-apis-go/config"
	hdrSrv "komodo-internal-lib-apis-go/http/headers/eval"
	hdrTypes "komodo-internal-lib-apis-go/http/types"
	httpUtils "komodo-internal-lib-apis-go/http/utils/http"
	logger "komodo-internal-lib-apis-go/logging/runtime"
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
		
		clientType := req.Context().Value(ctxKeys.ClientTypeKey)
		if clientType == nil {
			clientType = httpUtils.GetClientType(req)
		}
		
		if clientType == "api" {
			ctx := context.WithValue(req.Context(), IdempotencyValidCtxKey, true)
			req = req.WithContext(ctx)
			next.ServeHTTP(wtr, req)
			return
		}

		key := req.Header.Get("Idempotency-Key")

		if ok, err := hdrSrv.ValidateHeaderValue(hdrTypes.HEADER_IDEMPOTENCY, req); !ok || err != nil {
			logger.Error("invalid idempotency key for browser client: " + key, err)
			errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, errors.ERR_INVALID_REQUEST, "invalid idempotency key")
			return
		} 

		// Load existing entry
		if exp, ok := idemStore.Load(key); ok {
			// If expired, evict and continue; else reject as duplicate
			if until, _ := exp.(int64); until > time.Now().Unix() {
				wtr.Header().Set("Idempotency-Replayed", "true")
				logger.Error("duplicate request: " + key, req)
				errors.WriteErrorResponse(wtr, req, http.StatusConflict, errors.ERR_ACCESS_DENIED, "duplicate request")
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
