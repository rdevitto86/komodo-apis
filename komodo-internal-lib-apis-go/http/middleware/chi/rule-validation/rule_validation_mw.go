package rulevalidation

import (
	ruleServ "komodo-internal-lib-apis-go/services/eval-rules"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	"net/http"
)

func RuleValidationMiddleware(next http.Handler) http.Handler {
	// Ensure config is loaded
	ruleServ.LoadConfig()
	if !ruleServ.IsConfigLoaded() {
		logger.Warn("Validation rules failed to load - service will run without request validation", nil)
	}

	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		if rule := ruleServ.GetRule(req.URL.Path, req.Method); rule != nil {
			if !ruleServ.IsRuleValid(req, rule) {
				logger.Error("request does not comply with validation rule", req)
				http.Error(wtr, "Failed to validate request", http.StatusBadRequest)
				return
			}
		} else {
			logger.Error("no validation rule found", req)
			// http.Error(wtr, "Unable to validate request", http.StatusBadRequest)
		}

		next.ServeHTTP(wtr, req)
	})
}
