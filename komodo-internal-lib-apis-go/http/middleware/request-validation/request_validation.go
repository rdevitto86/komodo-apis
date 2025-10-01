package requestvalidation

import (
	evalTypes "komodo-internal-lib-apis-go/config/rules/types"
	keys "komodo-internal-lib-apis-go/http/middleware/context/keys"
	crsfUtils "komodo-internal-lib-apis-go/http/middleware/csrf/utils"
	idemUtils "komodo-internal-lib-apis-go/http/middleware/idempotency/utils"
	evalUtils "komodo-internal-lib-apis-go/http/middleware/request-validation/utils"
	logger "komodo-internal-lib-apis-go/logger/runtime"
	"net/http"
)

func RequestValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// 1. Lookup validation rule
		reqRule, ok := req.Context().Value(keys.ValidationRuleKey).(evalTypes.RequestEvalRule)

		if !ok {
			logger.Error("no validation rule found for: " + req.Method + " " + req.URL.Path, req)
			return
		}

		// 2. Check for ignored routes
		if reqRule.Level == evalTypes.LevelIgnore {
			next.ServeHTTP(wtr, req)
			return
		}

		// 3. Validate request version
		ver := req.Context().Value(keys.ApiVersionKey).(string)
		if reqRule.RequireVersion == evalTypes.RuleOn && !evalUtils.IsValidAPIVersion(ver) {
			logger.Error("invalid API version: " + ver, req)
			http.Error(wtr, "Invalid API version", http.StatusBadRequest)
			return
		}

		// 5. Validate headers
		for header, spec := range reqRule.Headers {
			val := req.Header.Get(header)

			if val == "" && spec.Required {
				logger.Error("missing mandatory header: " + header, req)
				http.Error(wtr, "Missing mandatory header: " + header, http.StatusBadRequest)
				return
			}
			// TODO check types???
			if !validateHeader(header, req, !spec.Required) {
				logger.Error("invalid header: " + header, req)
				http.Error(wtr, "Invalid header: " + header, http.StatusBadRequest)
				return
			}
		}
		
		// 6. Validate path params
		if len(reqRule.PathParams) > 0 && !evalUtils.IsValidPathParams(reqRule.PathParams, req) {
			logger.Error("invalid path params", req)
			http.Error(wtr, "Invalid path params", http.StatusBadRequest)
			return
		}

		// 7. Validate query params
		if len(reqRule.QueryParams) > 0 && !evalUtils.IsValidQueryParams(reqRule.QueryParams, req) {
			logger.Error("invalid query params", req)
			http.Error(wtr, "Invalid query params", http.StatusBadRequest)
			return
		}

		// 8. Validate body
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				if len(reqRule.Body) > 0 && !evalUtils.IsValidBody(reqRule.Body, req) {
					logger.Error("invalid request body", req)
					http.Error(wtr, "Invalid request body", http.StatusBadRequest)
					return
				}
			default:
		}

		next.ServeHTTP(wtr, req)
	})
}

func validateHeader(header string, req *http.Request, isOptional bool) bool {
	isHeaderValidated := false
	headerVal := req.Header.Get(header)

	if (isOptional && headerVal == "") { return true }

	switch header {
		case "Authorization":
			if req.Context().Value("Authorization_valid").(bool) {
				isHeaderValidated = true // validation done in authn middleware already
			} else {
				isHeaderValidated = evalUtils.IsValidBearer(headerVal)
			}
		case "X-Session-Token":
			if req.Context().Value("X-Session-Token_valid").(bool) {
				isHeaderValidated = true // validation done in authn middleware already
			} else {
				isHeaderValidated = evalUtils.IsValidSession(headerVal)
			}
		case "X-CSRF":
			if req.Context().Value("X-CSRF_valid").(bool) {
				isHeaderValidated = true // validation done in CSRF middleware already
			} else {
				isHeaderValidated = crsfUtils.IsValidCSRF(headerVal, req.Header.Get("X-Session-Token"))
			}
		case "X-Requested-By":
			isHeaderValidated = evalUtils.IsValidRequestedBy(headerVal)
		case "Content-Type", "Accept":
			isHeaderValidated = evalUtils.IsValidContentAcceptType(headerVal)
		case "Content-Length":
			isHeaderValidated = evalUtils.IsValidContentLength(headerVal)
		case "Cookie":
			isHeaderValidated = evalUtils.IsValidCookie(headerVal)
		case "User-Agent":
			isHeaderValidated = evalUtils.IsValidUserAgent(headerVal)
		case "Referer":
			isHeaderValidated = evalUtils.IsValidReferer(headerVal)
		case "Cache-Control":
			isHeaderValidated = evalUtils.IsValidCacheControl(headerVal)
		case "X-Client-Id":
			isHeaderValidated = evalUtils.IsValidClientID(headerVal)
		case "X-Client-Secret":
			isHeaderValidated = evalUtils.IsValidClientSecret(headerVal)
		case "Idempotency-Key":
			if req.Context().Value("Idempotency-Key_valid").(bool) {
				isHeaderValidated = true // validation done in idempotency middleware already
			} else {
				isHeaderValidated = idemUtils.IsValidIdempotencyKey(headerVal)
			}
	}
	return isHeaderValidated
}
