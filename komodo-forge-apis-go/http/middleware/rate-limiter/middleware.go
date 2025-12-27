package ratelimiting

import (
	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/chi"
	rl "komodo-forge-apis-go/http/services/rate_limiter"
	reqUtils "komodo-forge-apis-go/http/utils/request"
	logger "komodo-forge-apis-go/loggers/runtime"
	"net/http"
	"strconv"
)

// RateLimiterMiddleware delegates core logic to services/rate_limiter.
func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		key := reqUtils.GetClientKey(req)
		ctx := req.Context()

		allowed, wait, err := rl.Allow(ctx, key)
		if err != nil {
			logger.Error("rate limiter error: "+err.Error(), req)
			if rl.ShouldFailOpen() {
				logger.Error("rate limiter failing open for client: "+key, req)
			} else {
				errors.WriteErrorResponse(wtr, req, http.StatusInternalServerError, errCodes.ERR_INTERNAL_SERVER, "internal rate limiter error")
				return
			}
		} else if !allowed {
			if wait > 0 {
				wtr.Header().Set("Retry-After", strconv.Itoa(int(wait.Seconds()+0.5)))
			}
			logger.Error("rate limit exceeded for client: "+key, req)
			errors.WriteErrorResponse(wtr, req, http.StatusTooManyRequests, errCodes.ERR_ACCESS_DENIED, "rate limit exceeded")
			return
		}

		next.ServeHTTP(wtr, req)
	})
}
