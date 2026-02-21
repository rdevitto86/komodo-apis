package utils

import (
	ctxKeys "komodo-forge-sdk-go/http/context"
	"net/http"
)

// ResponseWriter wraps http.ResponseWriter to capture status code and bytes written.
type ResponseWriter struct {
	http.ResponseWriter
	Status       int
	BytesWritten int
	WroteHeader  bool
}

func (rw *ResponseWriter) WriteHeader(code int) {
	if !rw.WroteHeader {
		rw.Status = code
		rw.WroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if !rw.WroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.BytesWritten += n
	return n, err
}

func (rw *ResponseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func GetRequestID(req *http.Request) string {
	if rid, ok := req.Context().Value(ctxKeys.REQUEST_ID_KEY).(string); ok && rid != "" {
		return rid
	}
	return "unknown"
}
