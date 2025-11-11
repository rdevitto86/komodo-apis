package validateheaders

import (
	hdrSrv "komodo-internal-lib-apis-go/http/headers/eval"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"net/http"
)

func ValidateHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		headers := req.Header
		if headers == nil {
			logger.Error("no headers found", req)
			http.Error(wtr, "no headers found", http.StatusBadRequest)
			return
		}

		for header := range headers {
			if ok, err := hdrSrv.ValidateHeaderValue(header, req); !ok || err != nil {
				logger.Error("missing or invalid header: " + header, err)
				http.Error(wtr, "missing or invalid header: " + header, http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(wtr, req)
	})
}
