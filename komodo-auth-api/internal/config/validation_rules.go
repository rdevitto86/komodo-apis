package middleware

import (
	reqs "komodo-auth-api/internal/model"
	evalTypes "komodo-internal-lib-apis-go/config/rules/types"
	"net/http"
)

var ValidationRules = map[string]map[string]evalTypes.RequestEvalRule{
	"/health": {
		http.MethodGet: {
			Path:          "/health",
			Type:             reqs.HealthRequest.Type,
			Level:            evalTypes.LevelIgnore,
			Headers: 					reqs.HealthRequest.Headers,
			RequireVersion: 	evalTypes.RuleOff,
		},
	},
	"/auth/login": {
		http.MethodPost: {
			Path:            "/auth/login",
			Type:             	reqs.LoginRequest.Type,
			Level:            	evalTypes.LevelLenient,
			Headers: 						reqs.LoginRequest.Headers,
			Body: 							reqs.LoginRequest.Body,
			RequireVersion: 		evalTypes.RuleOn,
		},
	},
	"/auth/logout": {
		http.MethodPost: {
			Path:            "/auth/logout",
			Type:             	reqs.LogoutRequest.Type,
			Level:            	evalTypes.LevelLenient,
			Headers: 						reqs.LogoutRequest.Headers,
			Body: 							reqs.LogoutRequest.Body,
			RequireVersion: 		evalTypes.RuleOn,
		},
	},
	"/auth/token": {
		// POST
		http.MethodPost: {
			Path:						"/auth/token",
			Type: 							reqs.TokenCreateRequest.Type,
			Level: 							evalTypes.LevelStrict,
			Headers: 						reqs.TokenCreateRequest.Headers,
			Body: 							reqs.TokenCreateRequest.Body,
			RequireVersion:			evalTypes.RuleOn,
		},
		// DELETE
		http.MethodDelete: {
			Path:						"/auth/token",
			Type:             	reqs.TokenDeleteRequest.Type,
			Level:            	evalTypes.LevelStrict,
			Headers: 						reqs.TokenDeleteRequest.Headers,
			Body: 							reqs.TokenDeleteRequest.Body,
			RequireVersion: 		evalTypes.RuleOn,
		},
	},
	"/auth/token/refresh": {
		http.MethodPost: {
			Path:						"/auth/token/refresh",
			Type:             	reqs.TokenRefreshRequest.Type,
			Level:            	evalTypes.LevelStrict,
			Headers: 						reqs.TokenRefreshRequest.Headers,
			Body: 							reqs.TokenRefreshRequest.Body,
			RequireVersion: 		evalTypes.RuleOn,
		},
	},
	"/well-known/jwks.json": {
		http.MethodGet: {
			Path:          "/well-known/jwks.json",
			Type:            	reqs.WellKnownJWKSRequest.Type,
			Level:           	evalTypes.LevelIgnore,
			Headers: 					reqs.WellKnownJWKSRequest.Headers,
			RequireVersion: 	evalTypes.RuleOff,
		},
	},
}