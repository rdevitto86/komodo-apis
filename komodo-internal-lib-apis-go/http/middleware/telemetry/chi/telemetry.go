package telemetry

import (
	"fmt"
	logger "komodo-internal-lib-apis-go/logger/runtime"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func TelemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		wtr.Header().Add("Trailer", "Server-Timing")
		wtr.Header().Add("Trailer", "X-Response-Time")

		ww := chimw.NewWrapResponseWriter(wtr, req.ProtoMajor)
		start := time.Now()

		defer func() {
			ms := time.Since(start).Milliseconds()

			// Recover from panics and ensure a 500 is sent if nothing written.
			if rec := recover(); rec != nil {
				if ww.Status() == 0 {
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				logger.Error("telemetry panicked!", map[string]any{
					"request_id": chimw.GetReqID(req.Context()),
					"err":        rec,
				})
			}

			ww.Header().Set("Server-Timing", fmt.Sprintf("app;dur=%d", ms))
			ww.Header().Set("X-Response-Time", fmt.Sprintf("%dms", ms))

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}

			payload := map[string]any{
				"request_id": chimw.GetReqID(req.Context()),
				"method":     req.Method,
				"path":       req.URL.Path,
				"query":      req.URL.RawQuery,
				"status":     status,
				"bytes":     ww.BytesWritten(),
				"latency_ms": ms,
				"ip":         req.RemoteAddr,
				"user_agent": req.UserAgent(),
				"referer":    req.Referer(),
				"proto":     req.Proto,
				"host":      req.Host,
				"start_time": start.UTC().Format(time.RFC3339Nano),
				"finish_time": time.Now().UTC().Format(time.RFC3339Nano),
			}

			if status >= 400 {
				logger.Error("telemetry request failed", payload)
			} else {
				logger.Info("telemetry request completed", payload)
			}
		}()

		next.ServeHTTP(ww, req)
	})
}
