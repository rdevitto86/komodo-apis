package serviceauth

type TokenGenerateRequest struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

type TokenGenerateResponse struct {
	Token     string
	ExpiresIn int64
}

type TokenVerifyResponse struct {
	Active   bool   `json:"active"`               // Is token valid (signature, expiry, not revoked)?
	ClientID string `json:"client_id,omitempty"`  // Which service is making the request?
	Scope    string `json:"scope,omitempty"`      // What permissions does it have?
}

type TokenRevokeResponse struct {
	Revoked bool
}
