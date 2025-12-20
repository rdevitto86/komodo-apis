package userauth

import (
	"errors"
	"regexp"
)

// Performs full session validation (format + cache lookup + expiry)
// TODO: Implement after Elasticache integration
func ValidateSession(token string) (bool, error) {
	if len(token) < 32 || len(token) > 128 || !regexp.MustCompile(`^[A-Za-z0-9_\-]+$`).MatchString(token) {
		return false, errors.New("invalid session token format")
	}

	// sessionData, err := elasticache.GetSession("session:" + token)
	// if err != nil {
	// 	return nil, errors.New("session not found or expired")
	// }
	
	// expiryUnix := sessionData.ExpiryUnix
	// if expiryUnix > 0 && time.Now().Unix() > expiryUnix {
	// 	// Clean up expired session
	// 	elasticache.DeleteSession("session:" + token)
	// 	return nil, errors.New("session expired")
	// }
	
	// if sessionData.Revoked {
	// 	return nil, errors.New("session has been revoked")
	// }
	
	return true, nil
}
