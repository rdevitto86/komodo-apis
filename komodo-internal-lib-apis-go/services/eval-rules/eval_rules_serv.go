package evalrules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	httpUtils "komodo-internal-lib-apis-go/http/utils"
	headers "komodo-internal-lib-apis-go/services/headers/eval"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	evalrules "komodo-internal-lib-apis-go/types/eval-rules"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Checks if the request complies with all aspects of the provided EvalRule.
func IsRuleValid(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil { return false }
	if rule.Level == evalrules.LevelIgnore { return true }

	return AreValidHeaders(req, rule) &&
		AreValidPathParams(req, rule) &&
		AreValidQueryParams(req, rule) &&
		IsValidBody(req, rule)
}

// Checks if the request headers comply with the provided EvalRule.
func AreValidHeaders(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil {
		logger.Error("request or rule is nil")
		return false
	}

	// iterate rule headers and validate
	for hName, spec := range rule.Headers {
		val := req.Header.Get(hName)

		// Check if required and missing
		if spec.Required && val == "" {
			logger.Error(fmt.Sprintf("header %q is required but missing", hName))
			return false
		}
		if val == "" { continue }

		// Check exact value match if specified
		if spec.Value != "" {
			// Support wildcard matching for "Bearer *" etc
			if spec.Value[len(spec.Value)-1] == '*' {
				prefix := spec.Value[:len(spec.Value)-1]
				if !strings.HasPrefix(val, prefix) {
					logger.Error(fmt.Sprintf("header %q value %q does not match required prefix %q", hName, val, prefix))
					return false
				}
			} else if val != spec.Value {
				logger.Error(fmt.Sprintf("header %q value %q does not match required value %q", hName, val, spec.Value))
				return false
			}
		}

		// if pattern provided, verify
		if spec.Pattern != "" {
			re, err := regexp.Compile(spec.Pattern)
			if err != nil || !re.MatchString(val) {
				logger.Error(fmt.Sprintf("header %q value %q does not match pattern %q", hName, val, spec.Pattern))
				return false
			}
		}

		// enum check
		if len(spec.Enum) > 0 {
			ok := false
			for _, e := range spec.Enum {
				if e == val {
					ok = true
					break
				}
			}
			if !ok {
				logger.Error(fmt.Sprintf("header %q value %q not in enum %v", hName, val, spec.Enum))
				return false
			}
		}

		// length checks
		if spec.MinLen > 0 && len(val) < spec.MinLen {
			logger.Error(fmt.Sprintf("header %q value length %d is less than minLen %d", hName, len(val), spec.MinLen))
			return false
		}
		if spec.MaxLen > 0 && len(val) > spec.MaxLen {
			logger.Error(fmt.Sprintf("header %q value length %d is greater than maxLen %d", hName, len(val), spec.MaxLen))
			return false
		}
		// header-specific validation (optional - comment out if causing issues)
		if !headers.ValidateHeaderValue(hName, req) {
			logger.Error(fmt.Sprintf("header %q failed ValidateHeaderValue check", hName))
			return false
		}
	}
	logger.Info("All headers passed validation")
	return true
}

// Checks if the request path parameters comply with the provided EvalRule.
func AreValidPathParams(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil {
		logger.Error("request or rule is nil")
		return false
	}

	// find matching pattern and extract params
	_, params := matchRouteAndExtractParams(req.URL.Path)
	if params == nil {
		// no dynamic params in path; ensure rule does not require any
		for k, spec := range rule.PathParams {
			if spec.Required {
				// required param missing
				_ = k
				logger.Error(fmt.Sprintf("path param %q is required but missing", k))
				return false
			}
		}
		return true
	}

	// validate each rule-specified param
	for name, spec := range rule.PathParams {
		val, ok := params[name]
		if !ok || val == "" {
			if spec.Required {
				logger.Error(fmt.Sprintf("path param %q is required but missing", name))
				return false
			}
			continue
		}

		// pattern check
		if spec.Pattern != "" {
			re, err := regexp.Compile(spec.Pattern)
			if err != nil || !re.MatchString(val) {
				logger.Error(fmt.Sprintf("path param %q value %q does not match pattern %q", name, val, spec.Pattern))
				return false
			}
		}

		// enum check
		if len(spec.Enum) > 0 {
			okEnum := false
			for _, e := range spec.Enum {
				if e == val { okEnum = true; break }
			}
			if !okEnum {
				logger.Error(fmt.Sprintf("path param %q value %q not in enum %v", name, val, spec.Enum))
				return false
			}
		}

		// length checks
		if spec.MinLen > 0 && len(val) < spec.MinLen {
			logger.Error(fmt.Sprintf("path param %q value length %d is less than minLen %d", name, len(val), spec.MinLen))
			return false
		}
		if spec.MaxLen > 0 && len(val) > spec.MaxLen {
			logger.Error(fmt.Sprintf("path param %q value length %d is greater than maxLen %d", name, len(val), spec.MaxLen))
			return false
		}

		// simple type validation for common scalar types
		switch spec.Type {
			case "", "string":
				// already a string
			case "int":
				if _, err := strconv.Atoi(val); err != nil {
					logger.Error(fmt.Sprintf("path param %q value %q is not a valid int", name, val))
					return false
				}
			case "bool":
				if val != "true" && val != "false" {
					logger.Error(fmt.Sprintf("path param %q value %q is not a valid bool", name, val))
					return false
				}
			default:
				// unknown types are treated as pass-through for now
		}
	}
	return true
}

// Checks if the request query parameters comply with the provided EvalRule.
func AreValidQueryParams(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil {
		logger.Error("request or rule is nil")
		return false
	}

	params := httpUtils.GetQueryParams(req)

	for name, spec := range rule.QueryParams {
		val, ok := params[name]
		if !ok || val == "" {
			if spec.Required {
				logger.Error(fmt.Sprintf("query param %q is required but missing", name))
				return false
			}
			continue
		}

		if spec.Pattern != "" {
			re, err := regexp.Compile(spec.Pattern)
			if err != nil || !re.MatchString(val) {
				logger.Error(fmt.Sprintf("query param %q value %q does not match pattern %q", name, val, spec.Pattern))
				return false
			}
		}

		if len(spec.Enum) > 0 {
			okv := false
			for _, e := range spec.Enum {
				if e == val { okv = true; break }
			}
			if !okv {
				logger.Error(fmt.Sprintf("query param %q value %q not in enum %v", name, val, spec.Enum))
				return false
			}
		}

		if spec.MinLen > 0 && len(val) < spec.MinLen {
			logger.Error(fmt.Sprintf("query param %q value length %d is less than minLen %d", name, len(val), spec.MinLen))
			return false
		}
		if spec.MaxLen > 0 && len(val) > spec.MaxLen {
			logger.Error(fmt.Sprintf("query param %q value length %d is greater than maxLen %d", name, len(val), spec.MaxLen))
			return false
		}
	}
	return true
}

// Checks if the request body complies with the provided EvalRule.
func IsValidBody(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil {
		logger.Error("request or rule is nil")
		return false
	}

	switch req.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			return true
	}

	// Read the body (it can only be read once)
	const maxBody = 1 << 20 // 1 MiB
	bodyBytes, err := io.ReadAll(io.LimitReader(req.Body, maxBody))
	if err != nil {
		logger.Error("failed to read request body", err)
		return false
	}

	// Restore the original body for downstream handlers
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// If body is empty, that's valid for some requests
	if len(bodyBytes) == 0 { return true }

	// Parse JSON
	var bodyMap map[string]any
	dec := json.NewDecoder(bytes.NewReader(bodyBytes))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&bodyMap); err != nil {
		logger.Error("failed to decode request body as JSON", err)
		return false
	}

	// Validate body fields against rule
	for name, spec := range rule.Body {
		v, ok := bodyMap[name]
		if !ok {
			if spec.Required {
				logger.Error(fmt.Sprintf("body field %q is required but missing", name))
				return false
			}
			continue
		}

		// basic type checks
		switch spec.Type {
			case "", "string":
				if _, ok := v.(string); !ok {
					logger.Error(fmt.Sprintf("body field %q is not a string", name))
					return false
				}
			case "int":
				// JSON numbers are float64 by default
				if _, ok := v.(float64); !ok {
					logger.Error(fmt.Sprintf("body field %q is not a number", name))
					return false
				}
			case "bool":
				if _, ok := v.(bool); !ok {
					logger.Error(fmt.Sprintf("body field %q is not a bool", name))
					return false
				}
		}
	}
	return true
}
