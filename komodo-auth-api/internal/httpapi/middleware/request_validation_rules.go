package middleware

import (
	"komodo-auth-api/internal/model"
	"net/http"
)

type ValidationLevel string

const (
	LevelStrict  	ValidationLevel = "STRICT"
	LevelLenient  ValidationLevel = "LENIENT"
	LevelIgnore   ValidationLevel = "IGNORE"
)

type ValidationRule struct {
	Type              []model.RequestType
	Level             ValidationLevel
	MandatoryHeaders  model.MandatoryHeaders
	OptionalHeaders   model.OptionalHeaders
	PathParams        model.PathParams
	QueryParams       model.QueryParams
	Body              model.Body
	RequireVersion    bool
}

var ValidationRules = map[string]map[string]ValidationRule{
	"/health": {
		http.MethodGet: {
			Type:             model.HealthRequest.Type,
			Level:            LevelIgnore,
			MandatoryHeaders: model.HealthRequest.MandatoryHeaders,
			OptionalHeaders:  model.HealthRequest.OptionalHeaders,
			RequireVersion: 	false,
		},
	},
	"/auth/login": {
		http.MethodPost: {
			Type: 						model.LoginRequest.Type,
			Level:            LevelLenient,
			MandatoryHeaders: model.LoginRequest.MandatoryHeaders,
			OptionalHeaders:  model.LoginRequest.OptionalHeaders,
			Body: 						model.LoginRequest.Body,
			RequireVersion: 	true,
		},
	},
	"/auth/logout": {
		http.MethodPost: {
			Type:             model.LogoutRequest.Type,
			Level:            LevelLenient,
			MandatoryHeaders: model.LogoutRequest.MandatoryHeaders,
			OptionalHeaders:  model.LogoutRequest.OptionalHeaders,
			Body: 						model.LogoutRequest.Body,
			RequireVersion: 	true,
		},
	},
	"/auth/token": {
		http.MethodPost: {
			Type: 						model.TokenCreateRequest.Type,
			Level: 						LevelStrict,
			MandatoryHeaders: model.TokenCreateRequest.MandatoryHeaders,
			OptionalHeaders: 	model.TokenCreateRequest.OptionalHeaders,
			Body: 						model.TokenCreateRequest.Body,
			RequireVersion:		true,
		},
		http.MethodDelete: {
			Type:             model.TokenDeleteRequest.Type,
			Level:            LevelStrict,
			MandatoryHeaders: model.APIDefaultHeaders,
			OptionalHeaders:  model.TokenDeleteRequest.OptionalHeaders,
			Body: 						model.TokenDeleteRequest.Body,
			RequireVersion: 	true,
		},
	},
	"/auth/token/refresh": {
		http.MethodPost: {
			Type:             model.TokenRefreshRequest.Type,
			Level:            LevelStrict,
			MandatoryHeaders: model.APIDefaultHeaders,
			OptionalHeaders:  model.TokenRefreshRequest.OptionalHeaders,
			Body: 						model.TokenRefreshRequest.Body,
			RequireVersion: 	true,
		},
	},
	"/well-known/jwks.json": {
		http.MethodGet: {
			Type:             model.WellKnownJWKSRequest.Type,
			Level:            LevelIgnore,
			MandatoryHeaders: model.WellKnownJWKSRequest.MandatoryHeaders,
			OptionalHeaders:  model.WellKnownJWKSRequest.OptionalHeaders,
			RequireVersion: 	false,
		},
	},
}