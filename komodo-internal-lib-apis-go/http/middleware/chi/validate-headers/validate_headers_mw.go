package validateheaders

import (
	hdrSrv "komodo-internal-lib-apis-go/services/headers/eval"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	"net/http"
)

func ValidateHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		headers := req.Header
		if headers == nil {
			logger.Error("no headers found", req)
			http.Error(wtr, "No headers found", http.StatusBadRequest)
			return
		}

		for header := range headers {
			if !hdrSrv.ValidateHeaderValue(header, req) {
				logger.Error("missing or invalid header: " + header, req)
				http.Error(wtr, "Missing or invalid header: "+header, http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(wtr, req)
	})
}
