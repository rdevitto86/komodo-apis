package utils

import (
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func GenerateToken() string {
	token := uuid.NewString()
	return token
}

func Encode(data string) (string, error) {
	// TODO: implement proper encoding
	return "", nil
}

func Decode(data string) (string, error) {
	// TODO: implement proper decoding
	return "", nil
}

// Check a regex to match {param} or :param patterns
func HasDynamicPathParam(route string) bool {
	re := regexp.MustCompile(`\{[^}]+\}|:[^/]+`)
	return re.MatchString(route)
}

// ParseURI extracts version, route, path parameters, and query parameters from the request URL.
func ParseURI(req *http.Request, pattern string) (string, string, map[string]string, map[string]string) {
	var base string = req.URL.Path
	if idx := strings.Index(req.URL.Path, "?"); idx != -1 {
		base = req.URL.Path[:idx]
	}

	// Split path and detect version segment if present
	trimmed := strings.TrimPrefix(base, "/")
	segments := strings.Split(trimmed, "/")

	version := ""
	var pathSegments = []string{}

	if len(segments) > 0 && len(segments[0]) > 0 && segments[0][0] == 'v' {
		// Treat the first segment like v1, v2, etc. as API version
		version = "/" + segments[0]
		pathSegments = segments[1:]
	} else {
		// No explicit version prefix
		pathSegments = segments
	}

	// Route is the path without version
	route := "/" + strings.Join(pathSegments, "/")
	if route == "//" {
		route = "/"
	}

	// Extract path params using the pattern
	pathParams := make(map[string]string)
	if pattern != "" {
		// Convert pattern to regex: replace :param or {param} with capture groups
		rePattern := regexp.MustCompile(`\{([^}]+)\}|:([^/]+)`)
		regexStr := "^" + rePattern.ReplaceAllStringFunc(pattern, func(match string) string {
			return `([^/]+)`
		}) + "$"
		re := regexp.MustCompile(regexStr)

		matches := re.FindStringSubmatch(base) // Match against full base path
		if matches != nil {
			// Extract param names from pattern
			paramNames := []string{}
			for _, submatch := range rePattern.FindAllStringSubmatch(pattern, -1) {
				if submatch[1] != "" {
					paramNames = append(paramNames, submatch[1]) // {param}
				} else if submatch[2] != "" {
					paramNames = append(paramNames, submatch[2]) // :param
				}
			}

			// Map param names to values (skip full match at index 0)
			for i, name := range paramNames {
				if i+1 < len(matches) {
					pathParams[name] = matches[i+1]
				}
			}
		}
	}

	// Parse query parameters
	queryParams := make(map[string]string)
	values, err := url.ParseQuery(req.URL.RawQuery)
	if err == nil {
		for k, v := range values {
			if len(v) > 0 {
				queryParams[k] = v[0] // Take the first value if multiple
			}
		}
	}

	return version, route, pathParams, queryParams
}

func GetClientKey(req *http.Request) string {
	// prefer first X-Forwarded-For entry when present
	if xf := req.Header.Get("X-Forwarded-For"); xf != "" {
		parts := strings.Split(xf, ",")
		if len(parts) > 0 {
			if ip := strings.TrimSpace(parts[0]); ip != "" {
				return ip
			}
		}
	}
	// fallback to remote addr host
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return req.RemoteAddr
}
