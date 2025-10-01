package model

import (
	evalTypes "komodo-internal-lib-apis-go/config/rules/types"
	httpTypes "komodo-internal-lib-apis-go/http/types"
	"net/http"
)

type RequestRules map[string]map[string]evalTypes.RequestEvalRule
var ValidationRules RequestRules

// Health Check Request
var HealthRequest = httpTypes.Request{
	Method: http.MethodGet,
	Type: []httpTypes.RequestType{
		httpTypes.REQ_TYPE_API_INT,
		httpTypes.REQ_TYPE_API_EXT,
		httpTypes.REQ_TYPE_UI_USER,
		httpTypes.REQ_TYPE_UI_GUEST,
		httpTypes.REQ_TYPE_ADMIN,
	},
	Headers: httpTypes.Headers{
		httpTypes.HEADER_X_REQUESTED_BY: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_USER_AGENT: {
			Type: "string",
			Required: false,
		},
		httpTypes.HEADER_REFERER: {
			Type: "string",
			Required: false,
		},
	},
}

// Login Request
var LoginRequest = httpTypes.Request{
	Method: http.MethodPost,
	Type: []httpTypes.RequestType{
		httpTypes.REQ_TYPE_UI_USER,
		httpTypes.REQ_TYPE_UI_GUEST,
	},
	Headers: httpTypes.Headers{
		httpTypes.HEADER_CONTENT_TYPE: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_X_REQUESTED_BY: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_USER_AGENT: {
			Type: "string",
			Required: false,
		},
		httpTypes.HEADER_REFERER: {
			Type: "string",
			Required: false,
		},
	},
	Body: httpTypes.Body{},
}

// Logout Request
var LogoutRequest = httpTypes.Request{
	Method: http.MethodPost,
	Type: []httpTypes.RequestType{
		httpTypes.REQ_TYPE_API_INT,
		httpTypes.REQ_TYPE_API_EXT,
		httpTypes.REQ_TYPE_UI_GUEST,
	},
	Headers: httpTypes.Headers{
		httpTypes.HEADER_X_REQUESTED_BY: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_USER_AGENT: {
			Type: "string",
			Required: false,
		},
		httpTypes.HEADER_REFERER: {
			Type: "string",
			Required: false,
		},
	},
	Body: httpTypes.Body{},
}

// Token Create Request
var TokenCreateRequest = httpTypes.Request{
	Method: http.MethodPost,
	Type: []httpTypes.RequestType{
		httpTypes.REQ_TYPE_API_INT,
		httpTypes.REQ_TYPE_API_EXT,
		httpTypes.REQ_TYPE_UI_GUEST,
	},
	Headers: httpTypes.Headers{
		httpTypes.HEADER_X_REQUESTED_BY: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_USER_AGENT: {
			Type: "string",
			Required: false,
		},
		httpTypes.HEADER_REFERER: {
			Type: "string",
			Required: false,
		},
	},
	Body: httpTypes.Body{
		"client_id": {
			Type: "string",
			Required: true,
		},
		"client_secret": {
			Type: "string",
			Required: true,
		},
	},
}

// Token Delete Request
var TokenDeleteRequest = httpTypes.Request{
	Method: http.MethodDelete,
	Type: []httpTypes.RequestType{
		httpTypes.REQ_TYPE_API_INT,
		httpTypes.REQ_TYPE_API_EXT,
	},
	Headers: httpTypes.Headers{
		httpTypes.HEADER_AUTH: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_CONTENT_TYPE: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_X_REQUESTED_BY: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_USER_AGENT: {
			Type: "string",
			Required: false,
		},
		httpTypes.HEADER_REFERER: {
			Type: "string",
			Required: false,
		},
	},
	Body: httpTypes.Body{
		"token": {
			Type: "string",
			Required: true,
		},
	},
}

// Token Refresh Request
var TokenRefreshRequest = httpTypes.Request{
	Method: http.MethodPost,
	Type: []httpTypes.RequestType{
		httpTypes.REQ_TYPE_API_INT,
		httpTypes.REQ_TYPE_API_EXT,
		httpTypes.REQ_TYPE_UI_GUEST,
	},
	Headers: httpTypes.Headers{
		httpTypes.HEADER_X_REQUESTED_BY: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_USER_AGENT: {
			Type: "string",
			Required: false,
		},
		httpTypes.HEADER_REFERER: {
			Type: "string",
			Required: false,
		},
	},
	Body: httpTypes.Body{
		"token": {
			Type: "string",
			Required: true,
		},
	},
}

// Well Known JWKS Request
var WellKnownJWKSRequest = httpTypes.Request{
	Method: http.MethodGet,
	Type: []httpTypes.RequestType{
		httpTypes.REQ_TYPE_API_INT,
		httpTypes.REQ_TYPE_API_EXT,
	},
	Headers: httpTypes.Headers{
		httpTypes.HEADER_X_REQUESTED_BY: {
			Type: "string",
			Required: true,
		},
		httpTypes.HEADER_USER_AGENT: {
			Type: "string",
			Required: false,
		},
		httpTypes.HEADER_REFERER: {
			Type: "string",
			Required: false,
		},
	},
}
