package utils

import "regexp"

func IsValidIdempotencyKey(key string) bool {
	if len(key) == 0 || len(key) > 128 { return false }
	return regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`).MatchString(key)
}