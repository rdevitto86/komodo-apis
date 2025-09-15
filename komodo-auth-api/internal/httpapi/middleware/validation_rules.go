package middleware

import (
	"komodo-auth-api/internal/model"
	"net/http"
)

type ValidationLevel int
const (
	LevelIgnore   ValidationLevel = 0
	LevelStrict  	ValidationLevel = 1
	LevelLenient  ValidationLevel = 2
)

type RuleToggle int
const (
	RuleOff  	RuleToggle = 0
	RuleOn  	RuleToggle = 1
	RuleOpt  	RuleToggle = 2
)

type ValidationRule struct {
	Type              	[]model.RequestType
	Level             	ValidationLevel
	Headers           	model.Headers
	PathParams        	model.PathParams
	QueryParams       	model.QueryParams
	Body              	model.Body
	RequireVersion    	RuleToggle
}

var ValidationRules = map[string]map[string]ValidationRule{
	"/health": {
		http.MethodGet: {
			Type:             model.HealthRequest.Type,
			Level:            LevelIgnore,
			Headers: 					model.HealthRequest.Headers,
			RequireVersion: 	RuleOff,
		},
	},
	"/auth/login": {
		http.MethodPost: {
			Type:             	model.LoginRequest.Type,
			Level:            	LevelLenient,
			Headers: 						model.LoginRequest.Headers,
			Body: 							model.LoginRequest.Body,
			RequireVersion: 		RuleOn,
		},
	},
	"/auth/logout": {
		http.MethodPost: {
			Type:             	model.LogoutRequest.Type,
			Level:            	LevelLenient,
			Headers: 						model.LogoutRequest.Headers,
			Body: 							model.LogoutRequest.Body,
			RequireVersion: 		RuleOn,
		},
	},
	"/auth/token": {
		// POST
		http.MethodPost: {
			Type: 							model.TokenCreateRequest.Type,
			Level: 							LevelStrict,
			Headers: 						model.TokenCreateRequest.Headers,
			Body: 							model.TokenCreateRequest.Body,
			RequireVersion:			RuleOn,
		},
		// DELETE
		http.MethodDelete: {
			Type:             	model.TokenDeleteRequest.Type,
			Level:            	LevelStrict,
			Headers: 						model.TokenDeleteRequest.Headers,
			Body: 							model.TokenDeleteRequest.Body,
			RequireVersion: 		RuleOn,
		},
	},
	"/auth/token/refresh": {
		http.MethodPost: {
			Type:             	model.TokenRefreshRequest.Type,
			Level:            	LevelStrict,
			Headers: 						model.TokenRefreshRequest.Headers,
			Body: 							model.TokenRefreshRequest.Body,
			RequireVersion: 		RuleOn,
		},
	},
	"/well-known/jwks.json": {
		http.MethodGet: {
			Type:             model.WellKnownJWKSRequest.Type,
			Level:            LevelIgnore,
			Headers: 					model.WellKnownJWKSRequest.Headers,
			RequireVersion: 	RuleOff,
		},
	},
}