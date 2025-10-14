package evalrules

import (
	"encoding/json"
	"io"
	httpUtils "komodo-internal-lib-apis-go/http/utils"
	headers "komodo-internal-lib-apis-go/services/headers/eval"
	evalrules "komodo-internal-lib-apis-go/types/eval-rules"
	"net/http"
	"regexp"
	"strconv"
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
	if req == nil || rule == nil { return false }

	// iterate rule headers and validate
	for hName, spec := range rule.Headers {
		val := req.Header.Get(hName)

		if spec.Required && val == "" { return false }
		if val == "" { continue }

		// if pattern provided, verify
		if spec.Pattern != "" {
			re, err := regexp.Compile(spec.Pattern)
			if err != nil || !re.MatchString(val) {
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
				return false
			}
		}

		// length checks
		if spec.MinLen > 0 && len(val) < spec.MinLen { return false }
		if spec.MaxLen > 0 && len(val) > spec.MaxLen { return false }

		// header-specific validation
		if !headers.ValidateHeaderValue(hName, req) { return false }
	}

	return true
}

// Checks if the request path parameters comply with the provided EvalRule.
func AreValidPathParams(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil { return false }

	// find matching pattern and extract params
	_, params := matchRouteAndExtractParams(req.URL.Path)
	if params == nil {
		// no dynamic params in path; ensure rule does not require any
		for k, spec := range rule.PathParams {
			if spec.Required {
				// required param missing
				_ = k
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
				return false
			}
			continue
		}

		// pattern check
		if spec.Pattern != "" {
			re, err := regexp.Compile(spec.Pattern)
			if err != nil || !re.MatchString(val) {
				return false
			}
		}

		// enum check
		if len(spec.Enum) > 0 {
			okEnum := false
			for _, e := range spec.Enum {
				if e == val { okEnum = true; break }
			}
			if !okEnum { return false }
		}

		// length checks
		if spec.MinLen > 0 && len(val) < spec.MinLen { return false }
		if spec.MaxLen > 0 && len(val) > spec.MaxLen { return false }

		// simple type validation for common scalar types
		switch spec.Type {
			case "", "string":
				// already a string
			case "int":
				if _, err := strconv.Atoi(val); err != nil { return false }
			case "bool":
				if val != "true" && val != "false" { return false }
			default:
				// unknown types are treated as pass-through for now
		}
	}

	return true
}

// Checks if the request query parameters comply with the provided EvalRule.
func AreValidQueryParams(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil { return false }

	params := httpUtils.GetQueryParams(req)

	for name, spec := range rule.QueryParams {
		val, ok := params[name]
		if !ok || val == "" {
			if spec.Required {
				return false
			}
			continue
		}

		if spec.Pattern != "" {
			re, err := regexp.Compile(spec.Pattern)
			if err != nil || !re.MatchString(val) {
				return false
			}
		}

		if len(spec.Enum) > 0 {
			okv := false
			for _, e := range spec.Enum {
				if e == val { okv = true; break }
			}
			if !okv { return false }
		}

		if spec.MinLen > 0 && len(val) < spec.MinLen { return false }
		if spec.MaxLen > 0 && len(val) > spec.MaxLen { return false }
	}

	return true
}

// Checks if the request body complies with the provided EvalRule.
func IsValidBody(req *http.Request, rule *evalrules.EvalRule) bool {
	if req == nil || rule == nil { return false }

	switch req.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			return true
	}

	// limit body size to avoid DoS
	const maxBody = 1 << 20 // 1 MiB
	rdr := io.LimitReader(req.Body, maxBody)
	defer req.Body.Close()

	var bodyMap map[string]any
	dec := json.NewDecoder(rdr)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&bodyMap); err != nil {
		return false
	}

	for name, spec := range rule.Body {
		v, ok := bodyMap[name]
		if !ok {
			if spec.Required { return false }
			continue
		}
		// basic type checks
		switch spec.Type {
		case "", "string":
			if _, ok := v.(string); !ok { return false }
		case "int":
			// JSON numbers are float64 by default
			if _, ok := v.(float64); !ok { return false }
		case "bool":
			if _, ok := v.(bool); !ok { return false }
		}
	}

	return true
}
