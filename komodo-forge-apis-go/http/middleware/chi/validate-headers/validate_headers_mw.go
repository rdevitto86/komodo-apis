package validateheaders

import (
	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/chi"
	hdrSrv "komodo-forge-apis-go/http/headers/eval"
	logger "komodo-forge-apis-go/logging/runtime"
	"net/http"
)

func ValidateHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		headers := req.Header
		if headers == nil {
			logger.Error("no headers found", req)
			errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, errCodes.ERR_INVALID_REQUEST, "no headers found")
			return
		}

		for header := range headers {
			if ok, err := hdrSrv.ValidateHeaderValue(header, req); !ok || err != nil {
				logger.Error("missing or invalid header: " + header, err)
				errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, errCodes.ERR_INVALID_REQUEST, "missing or invalid header: " + header)
				return
			}
		}
		next.ServeHTTP(wtr, req)
	})
}
