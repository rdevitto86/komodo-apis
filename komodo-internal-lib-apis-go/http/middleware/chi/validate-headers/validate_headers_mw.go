package validateheaders

import (
	"komodo-internal-lib-apis-go/common/errors"
	hdrSrv "komodo-internal-lib-apis-go/http/headers/eval"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"net/http"
)

func ValidateHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		headers := req.Header
		if headers == nil {
			logger.Error("no headers found", req)
			errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, errors.ERR_INVALID_REQUEST, "no headers found")
			return
		}

		for header := range headers {
			if ok, err := hdrSrv.ValidateHeaderValue(header, req); !ok || err != nil {
				logger.Error("missing or invalid header: " + header, err)
				errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, errors.ERR_INVALID_REQUEST, "missing or invalid header: " + header)
				return
			}
		}
		next.ServeHTTP(wtr, req)
	})
}
