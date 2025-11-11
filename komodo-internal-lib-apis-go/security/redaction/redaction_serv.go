package redaction

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Returns a shallow copy of the given request suitable for
// logging: headers, URL query params and common sensitive JSON body fields
// are redacted. The original request is not modified.
func Redact(req *http.Request) *http.Request {
	if req == nil { return nil }

	// shallow copy of request
	r2 := new(http.Request)
	*r2 = *req

	// redact headers
	r2.Header = redactHeaders(req.Header)

	// redact URL (query params)
	if req.URL != nil {
		u := *req.URL
		u.RawQuery = redactQuery(req.URL.Query()).Encode()
		r2.URL = &u
	}

	// redact body (non-destructive for the original)
	if req.Body != nil {
		// read body
		b, err := io.ReadAll(req.Body)
		if err == nil {
			// restore original body for downstream (rewind)
			req.Body = io.NopCloser(bytes.NewReader(b))
			// redact copy for log
			ct := req.Header.Get("Content-Type")
			rb := redactBody(b, ct)
			r2.Body = io.NopCloser(bytes.NewReader(rb))
		}
	}

	return r2
}

func redactHeaders(header http.Header) http.Header {
	if header == nil { return nil }

	out := make(http.Header, len(header))
	sensitiveHeaderRE := regexp.MustCompile(`(?i)authorization|cookie|set-cookie|x-api-key|x-amz-signature`)

	for k, vv := range header {
		// if header name matches sensitive pattern, redact values
		if sensitiveHeaderRE.MatchString(k) {
			out[k] = []string{"REDACTED"}
			continue
		}
		// otherwise copy values but scrub if any value looks like a bearer token or long secret
		nv := make([]string, 0, len(vv))
		for _, v := range vv {
			if looksLikeToken(v) {
				nv = append(nv, "REDACTED")
			} else {
				nv = append(nv, v)
			}
		}
		out[k] = nv
	}
	return out
}

var bearerRE = regexp.MustCompile(`(?i)^\s*bearer\s+[A-Za-z0-9\-\._~\+/]+=*$`)
var longTokenRE = regexp.MustCompile(`[A-Za-z0-9\-\._~\+/]{20,}`)

func looksLikeToken(s string) bool {
	if s == "" {
		return false
	}
	if bearerRE.MatchString(s) {
		return true
	}
	// long base64-like strings are suspicious
	if longTokenRE.MatchString(s) && len(s) > 30 {
		return true
	}
	return false
}

func redactQuery(vals url.Values) url.Values {
	if vals == nil { return nil }

	out := url.Values{}
	for k, v := range vals {
		lowK := strings.ToLower(k)
		if containsSensitiveKey(lowK) {
			out[k] = []string{"REDACTED"}
			continue
		}

		nameVal := make([]string, 0, len(v))
		for _, vv := range v {
			if looksLikeToken(vv) {
				nameVal = append(nameVal, "REDACTED")
			} else {
				nameVal = append(nameVal, vv)
			}
		}

		out[k] = nameVal
	}
	return out
}

var sensitiveKeys = []string{
	"password",
	"passwd",
	"secret",
	"credit_card",
	"creditcard",
	"card_number",
	"ssn",
	"token",
	"access_token",
	"refresh_token",
	"client_secret",
}

func containsSensitiveKey(k string) bool {
	for _, s := range sensitiveKeys {
		if k == s || strings.Contains(k, s) {
			return true
		}
	}
	return false
}

func redactBody(b []byte, contentType string) []byte {
	if len(b) == 0 { return b }

	// only attempt JSON redaction; for others, do a simple token mask
	if strings.Contains(strings.ToLower(contentType), "application/json") {
		var v interface{}
		if err := json.Unmarshal(b, &v); err == nil {
			redactInterface(v)
			if out, err := json.Marshal(v); err == nil {
				return out
			}
		}
	}

	// fallback: mask bearer tokens and long tokens
	rb := bearerRE.ReplaceAllStringFunc(string(b), func(_ string) string { return "REDACTED" })
	rb = longTokenRE.ReplaceAllString(rb, "REDACTED")

	return []byte(rb)
}

func redactInterface(v interface{}) {
	switch t := v.(type) {
		case map[string]interface{}:
			for k, vv := range t {
				lk := strings.ToLower(k)
				if containsSensitiveKey(lk) {
					t[k] = "REDACTED"
					continue
				}
				// recurse
				redactInterface(vv)
			}
		case []interface{}:
			for i := range t {
				redactInterface(t[i])
			}
		case string:
			// nothing to do for plain string here (we mutate containers above)
	}
}
