package requestvalidation

import (
	"encoding/json"
	"fmt"
	"komodo-internal-lib-apis-go/config"
	"komodo-internal-lib-apis-go/crypto/jwt"
	httpTypes "komodo-internal-lib-apis-go/http/types"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ================ Request Validation Utils ================

func IsValidPathParams(params httpTypes.PathParams, req *http.Request) bool {
	if len(params) == 0 { return false }

	for name, spec := range params {
		val := chi.URLParam(req, name)
		if val == "" && !spec.Required { continue }
		if spec.Required && val == "" {
			return false
		}
		// Add type/pattern checks here if needed
	}
	return true
}

func IsValidQueryParams(params httpTypes.QueryParams, req *http.Request) bool {
	if len(params) == 0 { return false }

	for name, spec := range params {
		vals := req.URL.Query()[name]
		if len(vals) == 0 && !spec.Required { continue }
		if spec.Required && len(vals) == 0 {
			return false
		}
		// Add type/pattern checks here
	}
	return true
}

func IsValidBody(body httpTypes.Body, req *http.Request) bool {
	if len(body) == 0 { return false }

	var data map[string]any
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		fmt.Printf("Failed to decode request body: %v\n", err)
		return false
	}
	for key, spec := range body {
		_, exists := data[key]
		if !exists && spec.Required {
			fmt.Printf("Missing required field: %s\n", key)
			return false
		}
		// Add type/length/pattern checks based on spec
	}
	return true
}

func IsValidMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func IsValidAPIVersion(path string) bool {
	return path == ("/v" + config.GetConfigValue("API_VERSION"))
}

// ================ Header Validation Utils ================

func IsValidBearer(bearer string) bool {
	if bearer == "" { return false }

	bearerSplit := strings.Split(bearer, " ")
	if len(bearerSplit) != 2 || bearerSplit[0] != "Bearer" {
		return false
	}

	valid, err := jwt.VerifyJWT(bearerSplit[1])
	if !valid || err != nil {
		return false
	}
	return true
}

func IsValidContentAcceptType(str string) bool {
	return strings.HasPrefix(str, "application/json") ||
		strings.HasPrefix(str, "application/x-www-form-urlencoded") ||
		strings.HasPrefix(str, "multipart/form-data")
}

func IsValidContentLength(str string) bool {
	if str == "" { return false }

	val, err := strconv.Atoi(str)
	if err != nil { return false }

	getMax := func() int {
		val := config.GetConfigValue("MAX_CONTENT_LENGTH")
		num, err := strconv.Atoi(val)
		if val == "" || err != nil { return 4096 }
		return num
	}

	return val > 0 && val <= getMax()
}

func IsValidSession(session string) bool {
	if session == "" { return false }
	// TODO - format check
	return true
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

func IsValidClientID(s string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9_\-]{16,128}$`).MatchString(s)
}

func IsValidClientSecret(s string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9_\-\.~]{32,256}$`).MatchString(s)
}

