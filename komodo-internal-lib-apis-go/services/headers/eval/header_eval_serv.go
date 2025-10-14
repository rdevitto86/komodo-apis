package evalheaders

import (
	"net/http"
	"strings"
)

// ValidateHeaderValue runs lightweight validation for known header names.
func ValidateHeaderValue(hdr string, req *http.Request) bool {
	val := req.Header.Get(hdr)
	switch strings.ToLower(hdr) {
		case "access-control-allow-origin":
			return isValidCORS(val)
		case "authorization":
			return isValidBearer(val)
		case "cache-control":
			return isValidCacheControl(val)
		case "cookie":
			return isValidCookie(val)
		case "content-type", "accept":
			return isValidContentAcceptType(val)
		case "content-length":
			return isValidContentLength(val)
		case "idempotency-key":
			return isValidIdempotencyKey(val)
		case "referer", "referrer":
			return isValidReferer(val)
		case "session", "x-session-token":
			return isValidSession(val)
		case "user-agent":
			return isValidUserAgent(val)
		case "x-csrf-token":
			return isValidCSRF(val)
		case "x-requested-by":
			return isValidRequestedBy(val)
		default:
			return val != ""
	}
}
