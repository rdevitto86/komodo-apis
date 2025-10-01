package chi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	evalRules "komodo-internal-lib-apis-go/config/rules"
	ctxKeys "komodo-internal-lib-apis-go/http/middleware/context/keys"
	httpUtils "komodo-internal-lib-apis-go/http/utils"
	logger "komodo-internal-lib-apis-go/logger/runtime"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ver, base, pathParams, queryParams := httpUtils.ParseURI(req, req.URL.Path)

		ctx = context.WithValue(ctx, chimw.RequestIDKey, generateRequestId())
		ctx = context.WithValue(ctx, ctxKeys.StartTimeKey, time.Now().UTC())

		// Set Request timeout
		// TODO

		ctx = context.WithValue(ctx, ctxKeys.ApiVersionKey, ver)
		ctx = context.WithValue(ctx, ctxKeys.UriKey, base)
		ctx = context.WithValue(ctx, ctxKeys.PathParamsKey, pathParams)
		ctx = context.WithValue(ctx, ctxKeys.QueryParamsKey, queryParams)

		// Set validation rule using the injected rules
		if rule, err := evalRules.GetRequestRule(req.URL.Path, req.Method); err != nil {
			ctx = context.WithValue(ctx, ctxKeys.ValidationRuleKey, rule)
		} else {
			logger.Warn("no validation rule found for: " + req.Method + " " + req.URL.Path, req)
		}
	
		// TODO - check if rule is required via env var
		// if rule == nil {
		// 	logger.Error("no validation rule found for: " + req.Method + " " + req.URL.Path, req)
		// 	http.Error(wtr, "Failed to process request", http.StatusInternalServerError)
		// 	return
		// }

		req = req.WithContext(ctx)
		next.ServeHTTP(wtr, req)
	})
}

func generateRequestId() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		// fallback to timestamp if crypto rand fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
