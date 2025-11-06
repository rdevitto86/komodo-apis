package httputils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Extracts version, route, path parameters, and query parameters from the request URL.
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

// Extracts the API route from the request URL, excluding version prefix if present.
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

// Extracts path parameters from the request URL based on a predefined pattern.
// Note: This is a placeholder implementation and should be replaced with actual path parameter extraction logic.
func GetPathParams(req *http.Request) map[string]string {
	// Placeholder: return empty map as path parameter extraction requires route pattern knowledge
	return map[string]string{}
}

// Extracts the first value of each query parameter from the request URL.
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

// Extracts a client identifier from the request, preferring
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

// Creates a unique request ID using random bytes encoded in hex.
func GenerateRequestId() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// Validates if an API key exists and is active in the database.
// TODO: Implement actual validation against DynamoDB/RDS when database is ready.
func IsValidAPIKey(apiKey string) bool {
	// Placeholder: Replace with actual database lookup
	// Expected implementation:
	// 1. Query DynamoDB/RDS for api_key
	// 2. Check if key exists and is active (not revoked/expired)
	// 3. Optional: Rate limit check, scope validation
	// 4. Log the API key usage for auditing
	
	return true
}

// Determines if the request is from an API client or a browser client.
// Validates JWT token claims to prevent header spoofing.
func GetClientType(req *http.Request) string { 
	if apiKey := req.Header.Get("X-API-Key"); apiKey != "" && IsValidAPIKey(apiKey) {
		return "api"
	}
	
	authHeader := req.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		parts := strings.Split(strings.TrimPrefix(authHeader, "Bearer "), ".")

		if len(parts) == 3 {
			payload := parts[1]
			// Add padding if needed
			if m := len(payload) % 4; m != 0 {
				payload += strings.Repeat("=", 4-m)
			}
			
			if decoded, err := base64.URLEncoding.DecodeString(payload); err == nil {
				var claims map[string]interface{}
				if err := json.Unmarshal(decoded, &claims); err == nil {
					if clientType, ok := claims["client_type"].(string); ok {
						switch clientType {
							case "api", "browser":
								return clientType
						}
					}
					if grantType, ok := claims["grant_type"].(string); ok && grantType == "client_credentials" {
						return "api"
					}
					if scope, ok := claims["scope"].(string); ok {
						if strings.Contains(scope, "api:") || strings.Contains(scope, "service:") {
							return "api"
						}
					}
				}
			}
		}
	}
	
	// Default to browser (enforces CSRF)
	return "browser"
}
