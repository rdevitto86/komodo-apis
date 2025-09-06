package utils

import (
	"komodo-auth-api/internal/crypto"
	"komodo-auth-api/internal/thirdparty/aws"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
) 

func IsValidAPIVersion(path string) bool {
	return path == ("/v" + os.Getenv("API_VERSION"))
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
	if str == "" { return false }
	return strings.HasPrefix(str, "application/json")
}

func IsValidContentLength(str string) bool {
	if str == "" { return false }

	val, err := strconv.Atoi(str)
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
	return csrf != "" && csrf == session
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

func IsValidUserAgent(str string) bool {
	if str == "" { return false }
	commonAgents := []string{"Mozilla/", "Chrome/", "Safari/", "Opera/", "Edge/", "Firefox/", "PostmanRuntime/", "curl/"}
	for _, agent := range commonAgents {
		if strings.Contains(str, agent) {
			return true
		}
	}
	return false
}

func IsValidReferer(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func IsValidCacheControl(str string) bool {
	return str == "no-cache" || str == "no-store" || str == "must-revalidate"
}

func IsValidRequestedBy(str string) bool {
	return str == "API_INTERNAL" || str == "API_EXTERNAL" || str == "UI_USER" || str == "UI_GUEST" || str == "ADMIN"
}

func ParseURI(req *http.Request) (string, string) {
	var base string = req.URL.Path
	if idx := strings.Index(req.URL.Path, "?"); idx != -1 {
		base = req.URL.Path[:idx]
	}

	// Split path and detect version segment if present
	trimmed := strings.TrimPrefix(base, "/")
	segments := strings.Split(trimmed, "/")

	version := ""
	route := ""

	if len(segments) > 0 && len(segments[0]) > 0 && segments[0][0] == 'v' {
		// Treat the first segment like v1, v2, etc. as API version
		version = "/" + segments[0]
		if len(segments) > 1 {
			route = "/" + strings.Join(segments[1:], "/")
		} else {
			route = "/"
		}
	} else {
		// No explicit version prefix
		route = "/" + strings.Join(segments, "/")
	}

	// Append the raw query back to the route if present (preserves original casing)
	if raw := req.URL.RawQuery; raw != "" {
		route = route + "?" + raw
	}
	return version, route
}

func ParsePathParams(req *http.Request, filters ...string) map[string]string {
	out := make(map[string]string)

	for _, name := range filters {
		if name == "" {
			continue
		}
		val := chi.URLParam(req, name)
		if val != "" {
			out[name] = val
		}
	}
	return out
}

func ParseQueryParams(req *http.Request) map[string][]string {
	out := make(map[string][]string)

	for k, vals := range req.URL.Query() {
		out[k] = append([]string{}, vals...)
	}
	return out
}
