package utils

import (
	"komodo-auth-api/internal/crypto"
	"komodo-auth-api/internal/thirdparty/aws"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func IsAPIRequest(req *http.Request) bool {
	return strings.HasPrefix(req.URL.Path, "/api/") ||
		strings.Contains(req.Header.Get("Accept"), "application/json") ||
		req.Header.Get("X-Requested-By") == "API"
}

func IsUIRequest(req *http.Request) bool {
	return strings.HasPrefix(req.URL.Path, "/ui/") ||
		strings.Contains(req.Header.Get("Accept"), "text/html") ||
		req.Header.Get("X-Requested-By") == "UI"
}

func IsValidAPIPath(req *http.Request) bool {
	return strings.HasPrefix(req.URL.Path, ("/" + os.Getenv("API_VERSION")))
}

func IsValidBearer(bearer string) bool {
	if bearer == "" { return false }

	bearerSplit := strings.Split(bearer, " ")
	if len(bearerSplit) != 2 || bearerSplit[0] != "Bearer" {
		return false
	}

	valid, err := crypto.VerifyJWT(bearerSplit[1])
	if !valid || err != nil {
		return false
	}
	return true
}

func IsValidContentAcceptType(str string) bool {
	return strings.HasPrefix(str, "application/json")
}

func IsValidContentLength(cl string) bool {
	if cl == "" { return false }

	val, err := strconv.Atoi(cl)
	if err != nil { return false }

	getMax := func() int {
		val := os.Getenv("MAX_CONTENT_LENGTH")
		num, err := strconv.Atoi(val)
		if val == "" || err != nil { return 4096 }
		return num
	}

	return val > 0 && val <= getMax()
}

func IsValidCSRF(csrf string, session string) bool {
	if csrf == "" || session == "" { return false }
	return csrf == session
}

func IsValidSession(session string) bool {
	if session == "" { return false }

	storedSession, err := aws.GetSessionToken("session_token")
	if err != nil { return false }

	return session == storedSession
}

func IsValidCookie(cookie string) bool {
	// TODO: Implement cookie validation logic (e.g., parse, check signature)
	return cookie != ""
}

func IsValidUserAgent(userAgent string) bool {
	// TODO: Implement User-Agent validation logic (e.g., block bots)
	return userAgent != ""
}

func IsValidCacheControl(cacheControl string) bool {
	switch cacheControl {
	case "no-cache", "no-store", "must-revalidate":
		return true
	default:
		return false
	}
}
