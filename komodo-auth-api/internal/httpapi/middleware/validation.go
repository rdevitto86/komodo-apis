package middleware

import (
	"context"
	"komodo-auth-api/internal/httpapi/utils"
	"komodo-auth-api/internal/logger"
	"net/http"
)

func RuleValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// 1. Lookup validation rule
		rules := req.Context().Value(ValidationRuleKey).(ValidationRule)

		// 2. Check for ignored routes
		if rules.Level == LevelIgnore {
			next.ServeHTTP(wtr, req)
			return
		}

		// 3. Validate request version
		ver := req.Context().Value(ApiVersionKey).(string)
		if rules.RequireVersion == RuleOn && !utils.IsValidAPIVersion(ver) {
			logger.Error("invalid API version: " + ver, req)
			http.Error(wtr, "Invalid API version", http.StatusBadRequest)
			return
		}

		// 5. Validate headers
		for header, spec := range rules.Headers {
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
		if len(rules.PathParams) > 0 && !utils.IsValidPathParams(rules.PathParams, req) {
			logger.Error("invalid path params", req)
			http.Error(wtr, "Invalid path params", http.StatusBadRequest)
			return
		}

		// 7. Validate query params
		if len(rules.QueryParams) > 0 && !utils.IsValidQueryParams(rules.QueryParams, req) {
			logger.Error("invalid query params", req)
			http.Error(wtr, "Invalid query params", http.StatusBadRequest)
			return
		}

		// 8. Validate body
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				if len(rules.Body) > 0 && !utils.IsValidBody(rules.Body, req) {
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
				isHeaderValidated = utils.IsValidBearer(headerVal)
			}
		case "X-Session-Token":
			if req.Context().Value("X-Session-Token_valid").(bool) {
				isHeaderValidated = true // validation done in authn middleware already
			} else {
				isHeaderValidated = utils.IsValidSession(headerVal)
			}
		case "X-CSRF":
			if req.Context().Value("X-CSRF_valid").(bool) {
				isHeaderValidated = true // validation done in CSRF middleware already
			} else {
				isHeaderValidated = utils.IsValidCSRF(headerVal, req.Header.Get("X-Session-Token"))
			}
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
		case "X-Client-Id":
			isHeaderValidated = utils.IsValidClientID(headerVal)
		case "X-Client-Secret":
			isHeaderValidated = utils.IsValidClientSecret(headerVal)
		case "Idempotency-Key":
			if req.Context().Value("Idempotency-Key_valid").(bool) {
				isHeaderValidated = true // validation done in idempotency middleware already
			} else {
				isHeaderValidated = utils.IsValidIdempotencyKey(headerVal)
			}
	}
	return isHeaderValidated
}

func GetValidationRule(ctx context.Context) (interface{}, bool) {
	val := ctx.Value(ValidationRuleKey)
	if rule, ok := val.(ValidationRule); ok {
		return rule, true
	}
	return nil, false
}
