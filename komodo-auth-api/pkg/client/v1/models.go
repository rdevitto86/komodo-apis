package pkg

type IntrospectResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Iat       int64  `json:"iat,omitempty"`
	Sub       string `json:"sub,omitempty"`
	Aud       string `json:"aud,omitempty"`
}

// type TokenGenerateRequest struct {
// 	ClientID     string
// 	ClientSecret string
// 	Scope        string
// }

// type TokenGenerateResponse struct {
// 	Token     string
// 	ExpiresIn int64
// }

// type TokenVerifyResponse struct {
// 	Active   bool   `json:"active"`               // Is token valid (signature, expiry, not revoked)?
// 	ClientID string `json:"client_id,omitempty"`  // Which service is making the request?
// 	Scope    string `json:"scope,omitempty"`      // What permissions does it have?
// }

// type TokenRevokeResponse struct {
// 	Revoked bool
// }