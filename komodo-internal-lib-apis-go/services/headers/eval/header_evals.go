package evalheaders

import (
	"komodo-internal-lib-apis-go/config"
	"komodo-internal-lib-apis-go/crypto/jwt"
	hdrTypes "komodo-internal-lib-apis-go/types/headers"
	"regexp"
	"strconv"
	"strings"
)

func isValidBearer(s string) bool {
	if s != "" { return false }

	bearerSplit := strings.Split(s, " ")
	if len(bearerSplit) != 2 || bearerSplit[0] != "Bearer" {
		return false
	}

	valid, err := jwt.VerifyJWT(bearerSplit[1])
	if !valid || err != nil {
		return false
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
	re := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\-\._ ]*/\d+(\.\d+)*$`)
	return re.MatchString(strings.TrimSpace(s))
}

func isValidReferer(s string) bool {
	re := regexp.MustCompile(`^https?://[A-Za-z0-9\-.%]+(?::\d{1,5})?(?:/.*)?$`)
	return re.MatchString(strings.TrimSpace(s))
}

func isValidCacheControl(s string) bool {
	return s == "no-cache" || s == "no-store" || s == "must-revalidate"
}

func isValidRequestedBy(s string) bool {
	return s == hdrTypes.REQUESTED_BY_API_INT ||
		s == hdrTypes.REQUESTED_BY_API_EXT ||
		s == hdrTypes.REQUESTED_BY_UI_USER_VERIFIED ||
		s == hdrTypes.REQUESTED_BY_UI_USER_UNVERIFIED ||
		s == hdrTypes.REQUESTED_BY_USER_ADMIN
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
