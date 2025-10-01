package utils

// Contains checks if a slice of strings contains a specific value.
func Contains (slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}