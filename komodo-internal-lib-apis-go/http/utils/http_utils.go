package http

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// ParseURI extracts version, route, path parameters, and query parameters from the request URL.
func GetAPIVersion(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}

	trimmed := strings.TrimPrefix(req.URL.Path, "/")
	segments := strings.Split(trimmed, "/")

	if len(segments) > 0 && len(segments[0]) > 0 && segments[0][0] == 'v' {
		return "/" + segments[0]
	}
	return ""
}

// GetAPIRoute extracts the API route from the request URL, excluding version prefix if present.
func GetAPIRoute(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}

	var base string = req.URL.Path
	if idx := strings.Index(req.URL.Path, "?"); idx != -1 {
		base = req.URL.Path[:idx]
	}

	// Split path and detect version segment if present
	trimmed := strings.TrimPrefix(base, "/")
	segments := strings.Split(trimmed, "/")

	var pathSegments = []string{}

	if len(segments) > 0 && len(segments[0]) > 0 && segments[0][0] == 'v' {
		pathSegments = segments[1:]
	} else {
		pathSegments = segments // No explicit version prefix
	}

	// Route is the path without version
	route := "/" + strings.Join(pathSegments, "/")
	if route == "//" {
		route = "/"
	}
	return route
}

// GetQueryParams extracts the first value of each query parameter from the request URL.
func GetQueryParams(req *http.Request) map[string]string {
	if req == nil || req.URL == nil {
		return map[string]string{}
	}

	out := make(map[string]string)
	values, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil { return out }

	for k, v := range values {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}

// GetClientKey extracts a client identifier from the request, preferring
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
