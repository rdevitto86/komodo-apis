package rulevalidation

import (
	httpErr "komodo-forge-apis-go/http/errors"
	evalRules "komodo-forge-apis-go/http/rules"
	logger "komodo-forge-apis-go/logging/runtime"
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
				httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("request contents invalid"))
				return
			}
		} else {
			logger.Error("no validation rule found", req)
			httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to validate request"))
			return
		}
		next.ServeHTTP(wtr, req)
	})
}
