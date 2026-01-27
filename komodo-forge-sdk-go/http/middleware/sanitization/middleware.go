package sanitization

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	httpErr "komodo-forge-sdk-go/http/errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Sanitizes HTTP requests from malicious content
func SanitizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		sanitizeHeaders(req)
		sanitizePathParams(req)
		sanitizeQueryParams(req)

		if req.Body != nil && req.Header.Get("Content-Type") == "application/json" {
			sanitizeBody(wtr, req)
			if req.Body == nil { return }
		}

		next.ServeHTTP(wtr, req)
	})
}

// Removes malicious content from HTTP headers
func sanitizeHeaders(req *http.Request) {
	for key, values := range req.Header {
		for i, value := range values {
			req.Header[key][i] = sanitizeString(value)
		}
	}
}

// Sanitizes Chi URL parameters
func sanitizePathParams(req *http.Request) {
	ctx := chi.RouteContext(req.Context())
	if ctx == nil { return }

	for i, param := range ctx.URLParams.Values {
		ctx.URLParams.Values[i] = sanitizeString(param)
	}
}

// Sanitizes URL query parameters
func sanitizeQueryParams(req *http.Request) {
	query := req.URL.Query()
	sanitized := url.Values{}

	for key, values := range query {
		sanitizedKey := sanitizeString(key)
		for _, value := range values {
			sanitized.Add(sanitizedKey, sanitizeString(value))
		}
	}

	req.URL.RawQuery = sanitized.Encode()
}

// Sanitizes JSON request body
func sanitizeBody(wtr http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to read request body"))
		return
	}
	req.Body.Close()

	// Parse JSON
	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.BadRequest, httpErr.WithDetail("failed to parse JSON body"))
		return
	}

	// Sanitize the data recursively
	sanitized := sanitizeJSON(data)

	// Re-encode to JSON
	sanitizedBody, err := json.Marshal(sanitized)
	if err != nil {
		httpErr.SendError(wtr, req, httpErr.Global.Internal, httpErr.WithDetail("failed to marshal JSON body"))
		return
	}

	// Replace request body with sanitized version
	req.Body = io.NopCloser(bytes.NewBuffer(sanitizedBody))
	req.ContentLength = int64(len(sanitizedBody))
}

// Recursively sanitizes JSON data structures
func sanitizeJSON(data interface{}) interface{} {
	switch val := data.(type) {
		case string:
			return sanitizeString(val)
		case map[string]interface{}:
			sanitized := make(map[string]interface{})
			for key, value := range val {
				sanitizedKey := sanitizeString(key)
				sanitized[sanitizedKey] = sanitizeJSON(value)
			}
			return sanitized
		case []interface{}:
			sanitized := make([]interface{}, len(val))
			for i, value := range val {
				sanitized[i] = sanitizeJSON(value)
			}
			return sanitized
		default:
			return val
	}
}

// Sanitizes strings from malicious patterns
func sanitizeString(str string) string {
	str = nullBytePattern.ReplaceAllString(str, "")
	str = pathTraversalPattern.ReplaceAllString(str, "")
	str = html.EscapeString(str)
	str = strings.TrimSpace(str)

	if sqlInjectionPattern.MatchString(str) { str = sqlInjectionPattern.ReplaceAllString(str, "") }
	if xssPattern.MatchString(str) { str = xssPattern.ReplaceAllString(str, "") }

	return str
}
