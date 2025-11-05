package ratelimiting

import (
	utils "komodo-internal-lib-apis-go/http/utils/http"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	rl "komodo-internal-lib-apis-go/services/rate_limiter"
	"net/http"
	"strconv"
)

// RateLimiterMiddleware delegates core logic to services/rate_limiter.
func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		key := utils.GetClientKey(req)
		ctx := req.Context()

		allowed, wait, err := rl.Allow(ctx, key)
		if err != nil {
			logger.Error("rate limiter error: "+err.Error(), req)
			if rl.ShouldFailOpen() {
				logger.Error("rate limiter failing open for client: "+key, req)
			} else {
				http.Error(wtr, "Service unavailable", http.StatusServiceUnavailable)
				return
			}
		} else if !allowed {
			if wait > 0 {
				wtr.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds()+0.5)))
			}
			logger.Error("rate limit exceeded for client: "+key, req)
			http.Error(wtr, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(wtr, req)
	})
}
