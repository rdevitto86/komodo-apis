package evalrules

import (
	"fmt"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	evalrules "komodo-internal-lib-apis-go/types/eval-rules"
	"os"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	ruleMap       map[string]map[string]evalrules.EvalRule
	patternRoutes []routePattern // patternRoutes is a list of compiled route patterns for templates (/:id or /{id})
	loadOnce      sync.Once
	configLoaded  bool
)

type routePattern struct {
	template  string
	re        *regexp.Regexp
	methods   map[string]evalrules.EvalRule
	paramKeys []string
}

// LoadConfig loads validation rules from a file path or embedded data
func LoadConfig(path ...string) {
	loadOnce.Do(func() { 
		var data []byte
		var err error
		var source string

		// Try explicit path or env var
		if len(path) > 0 && path[0] != "" {
			data, err = os.ReadFile(path[0])
			source = path[0]
		} else if envPath := os.Getenv("EVAL_RULES_PATH"); envPath != "" {
			data, err = os.ReadFile(envPath)
			source = envPath
		}

		if err != nil || data == nil {
			logger.PrintError(fmt.Sprintf("Failed to load validation rules from %s", source), err)
			configLoaded = false
			return
		}

		rt, patterns, parseErr := parseConfigFromData(data)
		if parseErr != nil {
			logger.PrintError("Failed to parse validation rules")
			configLoaded = false
			return
		}

		ruleMap = rt
		patternRoutes = patterns
		configLoaded = true

		logger.PrintInfo(fmt.Sprintf("Successfully loaded validation rules from %s", source))
	})
}

// LoadConfigWithData loads validation rules from byte data (for embedding in client services)
func LoadConfigWithData(data []byte) {
	loadOnce.Do(func() {
		rt, patterns, err := parseConfigFromData(data)
		if err != nil {
			logger.PrintError("Failed to parse validation rules from embedded config", err)
			configLoaded = false
			return
		}

		ruleMap = rt
		patternRoutes = patterns
		configLoaded = true

		logger.PrintInfo("Successfully loaded validation rules from embedded config")
	})
}
	

// IsConfigLoaded returns true if the config has been successfully loaded
func IsConfigLoaded() bool {
	return configLoaded && ruleMap != nil
}

// IsConfigValid returns true if the current ruleMap is non-nil (loaded).
func IsConfigValid() bool {
	if !configLoaded || ruleMap == nil {
		return false
	}
	// TODO: perform deeper validation of patterns, enums, regexes, etc.
	return true
}

// GetRule returns the EvalRule for a given request path and method.
func GetRule(pKey string, method string) *evalrules.EvalRule {
	if pKey == "" || method == "" || ruleMap == nil {
		return nil
	}

	np := normalizePath(pKey)

	// Direct match
	if rules, ok := ruleMap[np]; ok {
		if rule, exists := rules[method]; exists {
			return &rule
		}
	}

	// Pattern matches
	for _, rp := range patternRoutes {
		if rp.re.MatchString(np) {
			if rule, exists := rp.methods[method]; exists {
				return &rule
			}
		}
	}
	return nil
}

// GetRules returns the full rule map as loaded from config.
func GetRules() evalrules.RuleConfig {
	return ruleMap
}

// parseConfigFromData parses YAML data and returns the rule map plus compiled templates.
func parseConfigFromData(data []byte) (map[string]map[string]evalrules.EvalRule, []routePattern, error) {
	var root struct {
		Rules evalrules.RuleConfig `yaml:"rules"`
	}

	if err := yaml.Unmarshal(data, &root); err != nil {
		logger.PrintError("Failed to parse validation rules from embedded config", err)
		return nil, nil, fmt.Errorf("invalid yaml: %w", err)
	}

	cfg := root.Rules
	patterns := make([]routePattern, 0)

	// Build patterns for templates with dynamic segments
	for tpl, methods := range cfg {
		if strings.Contains(tpl, ":") || strings.Contains(tpl, "{") || strings.Contains(tpl, "*") {
			reStr, keys := templateToRegex(tpl)
			re, err := regexp.Compile("^" + reStr + "$")
			if err != nil {
				logger.PrintError(fmt.Sprintf("Invalid route pattern %s", tpl), err)
				return nil, nil, fmt.Errorf("invalid route pattern %s: %w", tpl, err)
			}

			patterns = append(patterns, routePattern{
				template:  tpl,
				re:        re,
				methods:   methods,
				paramKeys: keys,
			})
		}
	}

	// Sort patterns by specificity
	specificityScore := func(tpl string) int {
		parts := strings.Split(strings.TrimPrefix(tpl, "/"), "/")
		literal, wild := 0, 0

		for _, p := range parts {
			if p == "*" || strings.HasPrefix(p, ":") || (strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}")) {
				wild++
			} else if p != "" {
				literal++
			}
		}
		return literal*10 - wild
	}

	for i := 0; i < len(patterns); i++ {
		for j := i + 1; j < len(patterns); j++ {
			if specificityScore(patterns[j].template) > specificityScore(patterns[i].template) {
				patterns[i], patterns[j] = patterns[j], patterns[i]
			}
		}
	}

	return cfg, patterns, nil
}

// templateToRegex converts a route template into a regex string
func templateToRegex(tpl string) (string, []string) {
	// ensure we only operate on path part (strip query if present)
	if idx := strings.Index(tpl, "?"); idx != -1 {
		tpl = tpl[:idx]
	}
	// trim trailing slash handling
	tpl = strings.TrimSuffix(tpl, "/")

	parts := strings.Split(strings.TrimPrefix(tpl, "/"), "/")
	regexParts := make([]string, 0, len(parts))
	keys := make([]string, 0)

	for _, p := range parts {
		if p == "*" {
			regexParts = append(regexParts, ".*")
			continue
		}
		// :param or {param}
		if strings.HasPrefix(p, ":") {
			key := strings.TrimPrefix(p, ":")
			keys = append(keys, key)
			regexParts = append(regexParts, `(?P<`+key+`>[^/]+)`)
			continue
		}
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			key := strings.TrimSuffix(strings.TrimPrefix(p, "{"), "}")
			keys = append(keys, key)
			regexParts = append(regexParts, `(?P<`+key+`>[^/]+)`)
			continue
		}
		// literal segment - escape regexp metacharacters
		regexParts = append(regexParts, regexp.QuoteMeta(p))
	}
	return "/" + strings.Join(regexParts, "/"), keys
}

// normalizePath strips version prefixes and ensures leading slash
func normalizePath(p string) string {
	if p == "" {
		return p
	}
	// drop query
	if idx := strings.Index(p, "?"); idx != -1 {
		p = p[:idx]
	}
	p = strings.TrimSpace(p)
	if p == "" {
		return "/"
	}

	// remove trailing slash (but keep root)
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		p = strings.TrimSuffix(p, "/")
	}

	// strip version prefix like /v1 or /v1.2
	trimmed := strings.TrimPrefix(p, "/")
	segs := strings.Split(trimmed, "/")

	if len(segs) > 0 && len(segs[0]) > 0 && segs[0][0] == 'v' {
		// basic check: next characters are digits or digit+dot
		if regexp.MustCompile(`^v[0-9]`).MatchString(segs[0]) {
			// remove first segment
			segs = segs[1:]
			p = "/" + strings.Join(segs, "/")
		}
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

// matchRouteAndExtractParams finds the first pattern that matches and extracts params
func matchRouteAndExtractParams(path string) (*routePattern, map[string]string) {
	np := normalizePath(path)

	for _, rp := range patternRoutes {
		if !rp.re.MatchString(np) {
			continue
		}
		matches := rp.re.FindStringSubmatch(np)
		names := rp.re.SubexpNames()
		params := make(map[string]string)

		for i, name := range names {
			if i != 0 && name != "" && i < len(matches) {
				params[name] = matches[i]
			}
		}
		return &rp, params
	}
	return nil, nil
}
