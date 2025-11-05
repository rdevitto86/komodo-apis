package rulevalidation

import (
	ruleServ "komodo-internal-lib-apis-go/services/eval-rules"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	"net/http"
)

// Enforces request validation rules based on predefined configurations.
func RuleValidationMiddleware(next http.Handler) http.Handler {
	// Ensure config is loaded
	if !ruleServ.LoadConfig() {
		logger.Error("Validation rules failed to load - service will run without request validations")
	}

	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		if rule := ruleServ.GetRule(req.URL.Path, req.Method); rule != nil {
			if !ruleServ.IsRuleValid(req, rule) {
				logger.Error("Request does not comply with validation rule", req)
				http.Error(wtr, "Request contents invalid", http.StatusBadRequest)
				return
			}
		} else {
			logger.Error("No validation rule found", req)
			http.Error(wtr, "Failed to validate request", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(wtr, req)
	})
}
