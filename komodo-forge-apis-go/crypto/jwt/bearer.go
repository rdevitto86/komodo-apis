package jwt

import (
	"fmt"
	"regexp"
	"strings"
)

var jwtFormatRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$`)

// Performs full validation of a Bearer token including
// format validation, signature verification, and expiration checking.
func ValidateBearerToken(authHeader string) (bool, error) {
	// Check for Bearer prefix
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return false, fmt.Errorf("missing or invalid Authorization header")
	}

	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	if tokenString == "" {
		return false, fmt.Errorf("empty token")
	}

	// Validate JWT format (3 base64url parts)
	if !jwtFormatRegex.MatchString(tokenString) {
		return false, fmt.Errorf("invalid JWT format")
	}

	// Verify signature and parse claims (uses existing JWT utils)
	_, claims, err := VerifyToken(tokenString)
	if err != nil { return false, err }

	if IsTokenExpired(claims) {
		return false, fmt.Errorf("token has expired")
	}
	return true, nil
}
