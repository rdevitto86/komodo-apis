package evalheaders

import (
	"komodo-internal-lib-apis-go/config"
	"regexp"
	"strconv"
	"strings"
)

func isValidBearer(s string) bool {
	if s == "" { return false }

	bearerSplit := strings.Split(s, " ")
	if len(bearerSplit) != 2 || bearerSplit[0] != "Bearer" {
		return false
	}

	token := bearerSplit[1]
	parts := strings.Split(token, ".")
	if len(parts) != 3 { return false }

	// Ensure each part has content
	for _, part := range parts {
		if len(part) == 0 { return false }
	}
	return true
}

func isValidContentAcceptType(s string) bool {
	return strings.HasPrefix(s, "application/json") ||
		strings.HasPrefix(s, "application/x-www-form-urlencoded") ||
		strings.HasPrefix(s, "multipart/form-data")
}

func isValidContentLength(s string) bool {
	if s == "" { return false }

	val, err := strconv.Atoi(s)
	if err != nil { return false }

	max := (func() int {
		val := config.GetConfigValue("MAX_CONTENT_LENGTH")
		num, err := strconv.Atoi(val)
		if val == "" || err != nil { return 4096 }
		return num
	})()

	return val > 0 && val <= max
}

func isValidSession(s string) bool {
	if s == "" { return false }
	// TODO - format check
	return true
}

func isValidCookie(s string) bool {
	// TODO: Implement cookie validation logic (e.g., parse, check signature)
	return s != ""
}

func isValidUserAgent(s string) bool {
	if s == "" { return false }
	s = strings.TrimSpace(s)
	if len(s) > 256 { return false } // max length
	re := regexp.MustCompile(`^[A-Za-z0-9\-\._ /(),:;]+$`)
	return re.MatchString(s)
}

func isValidReferer(s string) bool {
	re := regexp.MustCompile(`^https?://[A-Za-z0-9\-.%]+(?::\d{1,5})?(?:/.*)?$`)
	return re.MatchString(strings.TrimSpace(s))
}

func isValidCacheControl(s string) bool {
	return s == "no-cache" || s == "no-store" || s == "must-revalidate"
}

func isValidRequestedBy(s string) bool {
  return s != "" && len(s) <= 64 && regexp.MustCompile(`^[A-Za-z0-9_\-/]+$`).MatchString(s)
}

func isValidClientID(s string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9_\-]{16,128}$`).MatchString(s)
}

func isValidClientSecret(s string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9_\-\.~]{32,256}$`).MatchString(s)
}

func isValidIdempotencyKey(s string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9_\-]{8,64}$`).MatchString(s)
}

func isValidCSRF(s string) bool {
	if s == "" { return false }
	// TODO - implement CSRF token validation logic
	return true
}

func isValidCORS(s string) bool {
	if s == "" { return false }
	if s == "*" { return true }
	re := regexp.MustCompile(`^https?://[A-Za-z0-9\-.%]+(?::\d{1,5})?(?:/.*)?$`)
	return re.MatchString(strings.TrimSpace(s))
}
