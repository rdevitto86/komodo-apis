package context

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	utils "komodo-internal-lib-apis-go/http/utils"
	ctxKeys "komodo-internal-lib-apis-go/types/context"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		ctx = context.WithValue(ctx, chimw.RequestIDKey, generateRequestId())
		ctx = context.WithValue(ctx, ctxKeys.StartTimeKey, time.Now().UTC())
		ctx = context.WithValue(ctx, ctxKeys.VersionKey, utils.GetAPIVersion(req))
		ctx = context.WithValue(ctx, ctxKeys.UriKey, utils.GetAPIRoute(req))
		ctx = context.WithValue(ctx, ctxKeys.PathParamsKey, getPathParams(req))
		ctx = context.WithValue(ctx, ctxKeys.QueryParamsKey, utils.GetQueryParams(req))

		req = req.WithContext(ctx)
		next.ServeHTTP(wtr, req)
	})
}

func generateRequestId() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

func getPathParams(req *http.Request) map[string]string {
	if req == nil {
		return map[string]string{}
	}
	return map[string]string{} // TODO
}
