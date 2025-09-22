package middleware

import (
	"komodo-auth-api/internal/config"
	"net/http"
	"strings"
)

const (
	maxHeaderSize = 20 << 10 // 20KB total headers and per-value cap
	maxQuerySize  = 8 << 10  // 8KB total query string and path cap
	maxParamSize  = 2 << 10  // 2KB per query param name/value
	maxBodySize   = 100 << 10 // 100KB body cap
)

func CanonicalizeMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		req = sanitize(req)
		req = normalize(req)

		// Size limits: path, query, headers
		if req.URL != nil {
			if len(req.URL.Path) > maxQuerySize {
				http.Error(wtr, "Request URI too long", http.StatusRequestURITooLong)
				return
			}
			if r := req.URL.RawQuery; r != "" && len(r) > maxQuerySize {
				http.Error(wtr, "Query string too long", http.StatusRequestURITooLong)
				return
			}

			// Per-query param guard
			if query := req.URL.Query(); len(query) > 0 {
				total := 0
				for key, vals := range query {
					if len(key) > maxParamSize {
						http.Error(wtr, "Query param too large", http.StatusRequestURITooLong)
						return
					}

					total += len(key)
					for _, val := range vals {
						if len(val) > maxParamSize {
							http.Error(wtr, "Query param too large", http.StatusRequestURITooLong)
							return
						}
						total += len(val)
					}
				}
				if total > maxQuerySize {
					http.Error(wtr, "Query params too large", http.StatusRequestURITooLong)
					return
				}
			}
		}

		// Limit total header size (keys + values)
		if hdr := req.Header; len(hdr) > 0 {
			total := 0
			for k, vals := range hdr {
				total += len(k)

				for _, val := range vals {
					if len(val) > maxHeaderSize {
						http.Error(wtr, "Header value too large", http.StatusRequestHeaderFieldsTooLarge)
						return
					}
					total += len(val)
				}
			}
			if total > maxHeaderSize {
				http.Error(wtr, "Headers too large", http.StatusRequestHeaderFieldsTooLarge)
				return
			}
		}

		// Body size limits: enforce max size and guard reader
		switch req.Method {
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				if req.ContentLength > maxBodySize && req.ContentLength != -1 {
					http.Error(wtr, "Request body too large", http.StatusRequestEntityTooLarge)
					return
				}
				if req.Body != nil {
					req.Body = http.MaxBytesReader(wtr, req.Body, maxBodySize)
				}
		}

		// Security hardening: duplicates, TE/CL, host/proxy
		hdr := req.Header

		// Reject duplicate sensitive headers
		if hdr != nil {
			sensitive := []string{
				"Authorization", "Content-Length", "Transfer-Encoding",
				"X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Proto",
			}
			for _, k := range sensitive {
				if vals, ok := hdr[k]; ok && len(vals) > 1 {
					http.Error(wtr, "Duplicate header: " + k, http.StatusBadRequest)
					return
				}
			}
		}

		// Block TE/CL conflicts and only allow single 'chunked'
		if len(hdr["Transfer-Encoding"]) > 0 && len(hdr["Content-Length"]) > 0 {
			http.Error(wtr, "Conflicting Transfer-Encoding and Content-Length", http.StatusBadRequest)
			return
		}
		if vals := hdr["Transfer-Encoding"]; len(vals) == 1 {
			v := strings.TrimSpace(strings.ToLower(vals[0]))
			if v != "chunked" {
				http.Error(wtr, "Unsupported Transfer-Encoding", http.StatusBadRequest)
				return
			}
		} else if len(vals) > 1 {
			http.Error(wtr, "Multiple Transfer-Encoding headers", http.StatusBadRequest)
			return
		}

		// Validate Host against allowlist if provided
		if allowed := strings.TrimSpace(config.GetConfigValue("ALLOWED_HOSTS")); allowed != "" {
			ok := false
			hostLC := strings.ToLower(strings.TrimSpace(req.Host))

			for _, h := range strings.Split(allowed, ",") {
				h = strings.ToLower(strings.TrimSpace(h))
				if h != "" && hostLC == h {
					ok = true
					break
				}
			}
			if !ok {
				http.Error(wtr, "Invalid Host", http.StatusBadRequest)
				return
			}
		}

		// Proxy handling: honor or strip X-Forwarded-* based on TRUST_PROXY
		if strings.EqualFold(config.GetConfigValue("TRUST_PROXY"), "true") {
			if v := hdr.Get("X-Forwarded-Host"); v != "" {
				req.Host = strings.ToLower(strings.TrimSpace(v))
			}
			if v := hdr.Get("X-Forwarded-Proto"); v != "" {
				req.URL.Scheme = strings.ToLower(strings.TrimSpace(v))
			}
		} else {
			delete(hdr, "X-Forwarded-For")
			delete(hdr, "X-Forwarded-Host")
			delete(hdr, "X-Forwarded-Proto")
		}
	
		next.ServeHTTP(wtr, req)
	})
}

func sanitize(req *http.Request) *http.Request {
	if req == nil || req.URL == nil { return req }

	// Local sanitizer: trim spaces and drop ASCII control chars
	san := func(str string) string {
		if str == "" { return str }
		str = strings.TrimSpace(str)
		if str == "" { return str }
		str = strings.Map(func(r rune) rune {
			if r < 0x20 || r == 0x7F { return -1 }
			return r
		}, str)
		return str
	}

	// Method
	if method := req.Method; method != "" {
		req.Method = strings.ToUpper(strings.TrimSpace(method))
	}

	// Path and query
	req.URL.Path = san(req.URL.Path)
	req.URL.RawQuery = san(req.URL.RawQuery)

	// Headers
	for k, vals := range req.Header {
		dst := vals[:0]

		for _, val := range vals {
			if val := san(val); val != "" {
				dst = append(dst, val)
			}
		}
		if len(dst) == 0 {
			delete(req.Header, k)
		} else {
			req.Header[k] = dst
		}
	}
	return req
}

func normalize(req *http.Request) *http.Request {
	if req == nil || req.URL == nil { return req }

	// Lowercase host fields
	if host := req.Host; host != "" {
		lower := strings.ToLower(host)
		if lower != host {
			req.Host = lower
		}
	}

	// Normalize path: ensure leading '/', collapse duplicate slashes
	if path := req.URL.Path; path != "" {
		// Ensure leading slash
		if path[0] != '/' {
			path = "/" + path
		}

		// Collapse // to /
		if strings.Contains(path, "//") {
			bytes := make([]byte, 0, len(path))
			prevSlash := false

			for i := 0; i < len(path); i++ {
				char := path[i]
				if char == '/' {
					if prevSlash { continue }
					prevSlash = true
				} else {
					prevSlash = false
				}
				bytes = append(bytes, char)
			}
			path = string(bytes)
		}
		if path == "" {
			path = "/"
		}
		if path != req.URL.Path {
			req.URL.Path = path
		}
	} else {
		req.URL.Path = "/"
	}

	// Canonicalize header keys (values already sanitized)
	if hdr := req.Header; len(hdr) > 0 {
		for key, vals := range hdr {
			headerKey := http.CanonicalHeaderKey(key)
			if headerKey != key {
				hdr[headerKey] = append(hdr[headerKey], vals...)
				delete(hdr, key)
			}
		}
	}
	return req
}
