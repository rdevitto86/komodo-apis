package model

import "net/http"

// TODO - remap headers using map[string]spec

// Health Check Request
var HealthRequest = Request{
	Method: http.MethodGet,
	Type: []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_USER,
		REQ_TYPE_UI_GUEST,
		REQ_TYPE_ADMIN,
	},
	Headers: Headers{
		HEADER_X_REQUESTED_BY: {
			Type:     "string",
			Required: true,
		},
		HEADER_USER_AGENT: {
			Type:     "string",
			Required: false,
		},
		HEADER_REFERER: {
			Type:     "string",
			Required: false,
		},
	},
}

// Login Request
var LoginRequest = Request{
	Method: http.MethodPost,
	Type: []RequestType{
		REQ_TYPE_UI_USER,
		REQ_TYPE_UI_GUEST,
	},
	Headers: Headers{
		HEADER_CONTENT_TYPE: {
			Type:     "string",
			Required: true,
		},
		HEADER_X_REQUESTED_BY: {
			Type:     "string",
			Required: true,
		},
		HEADER_USER_AGENT: {
			Type:     "string",
			Required: false,
		},
		HEADER_REFERER: {
			Type:     "string",
			Required: false,
		},
	},
	Body: Body{},
}

// Logout Request
var LogoutRequest = Request{
	Method: http.MethodPost,
	Type: []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_GUEST,
	},
	Headers: Headers{
		HEADER_X_REQUESTED_BY: {
			Type:     "string",
			Required: true,
		},
		HEADER_USER_AGENT: {
			Type:     "string",
			Required: false,
		},
		HEADER_REFERER: {
			Type:     "string",
			Required: false,
		},
	},
	Body: Body{},
}

// Token Create Request
var TokenCreateRequest = Request{
	Method: http.MethodPost,
	Type: []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_GUEST,
	},
	Headers: Headers{
		HEADER_X_REQUESTED_BY: {
			Type:     "string",
			Required: true,
		},
		HEADER_USER_AGENT: {
			Type:     "string",
			Required: false,
		},
		HEADER_REFERER: {
			Type:     "string",
			Required: false,
		},
	},
	Body: Body{
		"client_id": {
			Type:     "string",
			Required: true,
		},
		"client_secret": {
			Type:     "string",
			Required: true,
		},
	},
}

// Token Delete Request
var TokenDeleteRequest = Request{
	Method: http.MethodDelete,
	Type: []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
	},
	Headers: Headers{
		HEADER_AUTH: {
			Type:     "string",
			Required: true,
		},
		HEADER_CONTENT_TYPE: {
			Type:     "string",
			Required: true,
		},
		HEADER_X_REQUESTED_BY: {
			Type:     "string",
			Required: true,
		},
		HEADER_USER_AGENT: {
			Type:     "string",
			Required: false,
		},
		HEADER_REFERER: {
			Type:     "string",
			Required: false,
		},
	},
	Body: Body{
		"token": {
			Type:     "string",
			Required: true,
		},
	},
}

// Token Refresh Request
var TokenRefreshRequest = Request{
	Method: http.MethodPost,
	Type: []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_GUEST,
	},
	Headers: Headers{
		HEADER_X_REQUESTED_BY: {
			Type:     "string",
			Required: true,
		},
		HEADER_USER_AGENT: {
			Type:     "string",
			Required: false,
		},
		HEADER_REFERER: {
			Type:     "string",
			Required: false,
		},
	},
	Body: Body{
		"token": {
			Type:     "string",
			Required: true,
		},
	},
}

// Well Known JWKS Request
var WellKnownJWKSRequest = Request{
	Method: http.MethodGet,
	Type: []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
	},
	Headers: Headers{
		HEADER_X_REQUESTED_BY: {
			Type:     "string",
			Required: true,
		},
		HEADER_USER_AGENT: {
			Type:     "string",
			Required: false,
		},
		HEADER_REFERER: {
			Type:     "string",
			Required: false,
		},
	},
}
