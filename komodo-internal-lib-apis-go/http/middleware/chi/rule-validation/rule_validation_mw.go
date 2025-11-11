package rulevalidation

import (
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
				http.Error(wtr, "request contents invalid", http.StatusBadRequest)
				return
			}
		} else {
			logger.Error("no validation rule found", req)
			http.Error(wtr, "failed to validate request", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(wtr, req)
	})
}
