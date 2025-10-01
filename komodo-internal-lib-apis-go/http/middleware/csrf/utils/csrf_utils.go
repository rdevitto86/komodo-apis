package utils

func IsValidCSRF(csrf string, session string) bool {
	return csrf != "" && csrf == session
}