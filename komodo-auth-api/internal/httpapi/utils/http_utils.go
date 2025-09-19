package utils

import (
	"net/http"
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

func ParseURI(req *http.Request) (string, string, []string, []string) {
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
	return version, route, []string{}, []string{}
}
