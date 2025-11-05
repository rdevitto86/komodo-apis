package authUtils

import "strings"

// Checks if the provided scope string is valid according to predefined rules.
func IsValidScope(scope string) bool {
	if scope == "" { return false }

	// Split by spaces and commas, validate each part
	for _, part := range strings.Fields(strings.ReplaceAll(scope, ",", " ")) {
		switch part {
			case "read", "write", "delete", "admin", "users:read",
			"users:write", "tokens:create", "tokens:delete":
				continue
			default:
				return false
		}
	}
	return true
}

// Returns a slice of invalid scopes found in the provided scope string.
func GetInvalidScopes(scope string) []string {
	invalidScopes := []string{}
	if scope == "" { return invalidScopes }

	// Split by spaces and commas, validate each part
	for _, part := range strings.Fields(strings.ReplaceAll(scope, ",", " ")) {
		switch part {
			case "read", "write", "delete", "admin", "users:read",
			"users:write", "tokens:create", "tokens:delete":
				continue
			default:
				invalidScopes = append(invalidScopes, part)
		}
	}
	return invalidScopes
}

// Checks if the provided grant type string is valid according to predefined rules.
func IsValidGrantType(grantType string) bool {
	switch grantType {
		case "client_credentials", "authorization_code", "refresh_token":
			return true
		default:
			return false
	}
}