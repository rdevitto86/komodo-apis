package utils

import (
	"github.com/google/uuid"
)

func GenerateToken(idemKey, userID, deviceID string) string {
	var token string
	if idemKey != "" {
		name := idemKey
		if userID != "" || deviceID != "" {
			name = idemKey + "|" + userID + "|" + deviceID
		}
		token = uuid.NewSHA1(uuid.NameSpaceURL, []byte(name)).String()
	} else {
		token = uuid.NewString()
	}
	return token
}
