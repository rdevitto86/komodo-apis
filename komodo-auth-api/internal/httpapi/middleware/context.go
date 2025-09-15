package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"komodo-auth-api/internal/httpapi/utils"
	"komodo-auth-api/internal/logger"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

type ctxKey string

const (
	StartTimeKey 					ctxKey = "start_time"
	ApiVersionKey 				ctxKey = "api_version"
	UriKey       					ctxKey = "uri"
	PathParamsKey 				ctxKey = "path_params"
	QueryParamsKey 				ctxKey = "query_params"
	ValidationRuleKey 		ctxKey = "validation_rule"
	RequestIDKey     			ctxKey = "request_id"
	RequestTimeoutKey    	ctxKey = "request_timeout"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ver, base, pathParams, queryParams := utils.ParseURI(req)

		ctx = context.WithValue(ctx, chimw.RequestIDKey, generateRequestId())
		ctx = context.WithValue(ctx, StartTimeKey, time.Now().UTC())

		// Set Request timeout
		// TODO

		ctx = context.WithValue(ctx, ApiVersionKey, ver)
		ctx = context.WithValue(ctx, UriKey, base)
		ctx = context.WithValue(ctx, PathParamsKey, pathParams)
		ctx = context.WithValue(ctx, QueryParamsKey, queryParams)

		// Set validation rule
		if methodRules, ok := ValidationRules[base]; ok {
			if r, ok := methodRules[req.Method]; ok {
				ctx = context.WithValue(ctx, ValidationRuleKey, r)
			} else {
				logger.Error("method not allowed: " + req.Method + " " + req.URL.Path, req)
				http.Error(wtr, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
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