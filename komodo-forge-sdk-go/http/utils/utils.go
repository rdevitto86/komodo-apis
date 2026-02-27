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

func (wtr *ResponseWriter) WriteHeader(code int) {
	if !wtr.WroteHeader {
		wtr.Status = code
		wtr.WroteHeader = true
	}
	wtr.ResponseWriter.WriteHeader(code)
}

func (wtr *ResponseWriter) Write(b []byte) (int, error) {
	if !wtr.WroteHeader { wtr.WriteHeader(http.StatusOK) }
	num, err := wtr.ResponseWriter.Write(b)
	wtr.BytesWritten += num
	return num, err
}

func (wtr *ResponseWriter) Unwrap() http.ResponseWriter {
	return wtr.ResponseWriter
}

func GetRequestID(req *http.Request) string {
	if rid, ok := req.Context().Value(ctxKeys.REQUEST_ID_KEY).(string); ok && rid != "" {
		return rid
	}
	return "unknown"
}
