package telemetry

import (
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func TelemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		ww := chimw.NewWrapResponseWriter(wtr, req.ProtoMajor)
		start := time.Now()

		defer func() {
			ms := time.Since(start).Milliseconds()

			// Recover from panics and ensure a 500 is sent if nothing written.
			if rec := recover(); rec != nil {
				// Get request ID safely
				reqID := chimw.GetReqID(req.Context())
				if reqID == "" { reqID = "unknown" }
				
				// Safely check status
				status := 0
				if ww != nil {
					status = ww.Status()
				}
				
				if status == 0 {
					http.Error(wtr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				
				logger.Error("telemetry panicked!", map[string]any{
					"request_id": reqID,
					"err":        rec,
				})
				return // Don't continue after panic
			}

			// Get status safely
			status := 0
			bytesWritten := 0
			if ww != nil {
				status = ww.Status()
				bytesWritten = ww.BytesWritten()
			}
			
			if status == 0 {
				status = http.StatusOK
			}

			// Get request ID safely
			reqID := chimw.GetReqID(req.Context())
			if reqID == "" { reqID = "unknown" }

			payload := map[string]any{
				"request_id": reqID,
				"method":     req.Method,
				"path":       req.URL.Path,
				"query":      req.URL.RawQuery,
				"status":     status,
				"bytes":     	bytesWritten,
				"latency_ms": ms,
				"ip":         req.RemoteAddr,
				"user_agent": req.UserAgent(),
				"referer":    req.Referer(),
				"proto":     	req.Proto,
				"host":      	req.Host,
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
