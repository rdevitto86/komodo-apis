package rulevalidation

import (
	errors "komodo-internal-lib-apis-go/common/errors"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	evalRules "komodo-internal-lib-apis-go/rule-validation"
	"net/http"
)

// Enforces request validation rules based on predefined configurations.
func RuleValidationMiddleware(next http.Handler) http.Handler {
	// Ensure config is loaded
	if !evalRules.LoadConfig() {
		logger.Error("validation rules failed to load - service will run without request validations")
	}

	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		if rule := evalRules.GetRule(req.URL.Path, req.Method); rule != nil {
			if !evalRules.IsRuleValid(req, rule) {
				logger.Error("request does not comply with validation rule", req)
				errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, "request contents invalid", errors.ERR_INVALID_REQUEST)
				return
			}
		} else {
			logger.Error("no validation rule found", req)
			errors.WriteErrorResponse(wtr, req, http.StatusBadRequest, "failed to validate request", errors.ERR_INVALID_REQUEST)
			return
		}
		next.ServeHTTP(wtr, req)
	})
}
