package model

import "net/http"

// Health Check Request
var HealthRequest = Request{
	Method:            http.MethodGet,
	Type:              []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_USER,
		REQ_TYPE_UI_GUEST,
		REQ_TYPE_ADMIN,
	},
	MandatoryHeaders:  MandatoryHeaders{
		HEADER_X_REQUESTED_BY,
	},
	OptionalHeaders:   OptionalHeaders{
		HEADER_USER_AGENT,
		HEADER_REFERER,
	},
}

// Login Request
var LoginRequest = Request{
	Method:            http.MethodPost,
	Type:              []RequestType{
		REQ_TYPE_UI_USER,
		REQ_TYPE_UI_GUEST,
	},
	MandatoryHeaders:  MandatoryHeaders{
		HEADER_CONTENT_TYPE,
		HEADER_X_REQUESTED_BY,
	},
	OptionalHeaders:   OptionalHeaders{
		HEADER_USER_AGENT,
		HEADER_REFERER,
	},
	PathParams:        PathParams{},
	QueryParams:       QueryParams{},
	Body:              Body{},
}

// Logout Request
var LogoutRequest = Request{
	Method:            http.MethodPost,
	Type:              []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_GUEST,
	},
	MandatoryHeaders:  MandatoryHeaders{
		HEADER_X_REQUESTED_BY,
	},
	OptionalHeaders:   OptionalHeaders{
		HEADER_USER_AGENT,
		HEADER_REFERER,
	},
	PathParams:        PathParams{},
	QueryParams:       QueryParams{},
	Body:              Body{},
}

// Token Create Request
var TokenCreateRequest = Request{
	Method:            http.MethodPost,
	Type:              []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_GUEST,
	},
	MandatoryHeaders:  MandatoryHeaders{
		HEADER_X_REQUESTED_BY,
	},
	OptionalHeaders:   OptionalHeaders{
		HEADER_USER_AGENT,
		HEADER_REFERER,
	},
	PathParams:        PathParams{},
	QueryParams:       QueryParams{},
	Body:              Body{},
}

// Token Delete Request
var TokenDeleteRequest = Request{
	Method:            http.MethodDelete,
	Type:              []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
	},
	MandatoryHeaders:  MandatoryHeaders{
		HEADER_AUTH,
		HEADER_CONTENT_TYPE,
		HEADER_X_REQUESTED_BY,
	},
	OptionalHeaders:   OptionalHeaders{
		HEADER_USER_AGENT,
		HEADER_REFERER,
	},
	PathParams:        PathParams{},
	QueryParams:       QueryParams{},
	Body:              Body{},
}

// Token Refresh Request
var TokenRefreshRequest = Request{
	Method:            http.MethodPost,
	Type:              []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
		REQ_TYPE_UI_GUEST,
	},
	MandatoryHeaders:  MandatoryHeaders{
		HEADER_X_REQUESTED_BY,
	},
	OptionalHeaders:   OptionalHeaders{
		HEADER_USER_AGENT,
		HEADER_REFERER,
	},
	PathParams:        PathParams{},
	QueryParams:       QueryParams{},
	Body: Body{
		"refresh_token": {
			Type:     "string",
			Required: true,
		},
	},
}

// Well Known JWKS Request
var WellKnownJWKSRequest = Request{
	Method:            http.MethodGet,
	Type:              []RequestType{
		REQ_TYPE_API_INT,
		REQ_TYPE_API_EXT,
	},
	MandatoryHeaders:  MandatoryHeaders{
		HEADER_X_REQUESTED_BY,
	},
	OptionalHeaders:   OptionalHeaders{
		HEADER_USER_AGENT,
		HEADER_REFERER,
	},
	PathParams:        PathParams{},
	QueryParams:       QueryParams{},
	Body:              Body{},
}
