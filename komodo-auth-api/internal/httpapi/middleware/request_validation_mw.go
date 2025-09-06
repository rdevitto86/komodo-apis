package middleware

import (
	"komodo-auth-api/internal/httpapi/utils"
	"net/http"
)

func RequestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// 1. Pre-process path
		ver, resource := utils.ParseURI(req)

		// 2. Get the validation rule for the current request
		var rule *ValidationRule = nil

		if methodRules, ok := ValidationRules[resource]; ok {
			if r, ok := methodRules[req.Method]; ok {
				rule = &r
			} else {
				http.Error(wtr, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		}
		if rule == nil {
			http.Error(wtr, "No validation rule found", http.StatusInternalServerError)
			return
		}

		// 3. Validate request version
		if rule.RequireVersion && !utils.IsValidAPIVersion(ver) {
			http.Error(wtr, "Invalid API version", http.StatusBadRequest)
			return
		}

		// 4. Check for ignored routes
		if rule.Level == LevelIgnore {
			next.ServeHTTP(wtr, req)
			return
		}

		// params := utils.ParsePathParams(req)
		// query := utils.ParseQueryParams(req)

		// 5. Validate mandatory headers
		for _, header := range rule.MandatoryHeaders {
			if headerVal := req.Header.Get(header); headerVal == "" {
				http.Error(wtr, "Missing mandatory header: " + header, http.StatusBadRequest)
				return
			}
			if !validateHeader(header, req, false) {
				http.Error(wtr, "Invalid header: " + header, http.StatusBadRequest)
				return
			}
		}

		// 6. Validate optional headers
		for _, header := range rule.OptionalHeaders {
			if !validateHeader(req.Header.Get(header), req, true) {
				http.Error(wtr, "Invalid header: " + header, http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(wtr, req)
	})
}

func validateHeader(header string, req *http.Request, isOptional bool) bool {
	isHeaderValidated := false
	headerVal := req.Header.Get(header)

	switch header {
		case "Authorization":
			isHeaderValidated = utils.IsValidBearer(headerVal)
		case "X-Session-Token":
			isHeaderValidated = utils.IsValidSession(headerVal)
		case "X-CSRF":
			isHeaderValidated = utils.IsValidCSRF(headerVal, req.Header.Get("X-Session-Token"))
		case "X-Requested-By":
			isHeaderValidated = utils.IsValidRequestedBy(headerVal)
		case "Content-Type", "Accept":
			isHeaderValidated = utils.IsValidContentAcceptType(headerVal)
		case "Content-Length":
			isHeaderValidated = utils.IsValidContentLength(headerVal)
		case "Cookie":
			isHeaderValidated = utils.IsValidCookie(headerVal)
		case "User-Agent":
			isHeaderValidated = utils.IsValidUserAgent(headerVal)
		case "Referer":
			isHeaderValidated = utils.IsValidReferer(headerVal)
		case "Cache-Control":
			isHeaderValidated = utils.IsValidCacheControl(headerVal)
	}
	return isHeaderValidated || isOptional
}
