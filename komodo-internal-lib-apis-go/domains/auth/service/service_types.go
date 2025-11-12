package serviceauth

type ServiceTokenRequest struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

type ServiceTokenResponse struct {
	Token     string
	ExpiresIn int64
}

type VerifyTokenRequest struct {
	Token string
	Scope string
}

type VerifyTokenResponse struct {
	Active bool
}