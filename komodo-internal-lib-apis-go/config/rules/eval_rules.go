package validationrules

import (
	"fmt"
	"komodo-internal-lib-apis-go/config/rules/types"
	"sync"
)

var (
	RequestRules map[string]map[string]types.RequestEvalRule
	once         sync.Once
)

func init() {
	once.Do(func() {
		RequestRules = make(map[string]map[string]types.RequestEvalRule)
	})
}

func SetRequestRule(rule types.RequestEvalRule) {
	RequestRules[rule.Path] = map[string]types.RequestEvalRule{
		rule.Method: rule,
	}
}

func SetRequestRules(rules []types.RequestEvalRule) {
	for _, rule := range rules {
		RequestRules[rule.Path] = map[string]types.RequestEvalRule{
			rule.Method: rule,
		}
	}
}

func GetRequestRule(key string, method string) (*types.RequestEvalRule, error) {
	if methodRules, ok := RequestRules[key]; ok {
		if rule, ok := methodRules[method]; ok {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no rule found for %s %s", method, key)
}

func GetRequestRules() map[string]map[string]types.RequestEvalRule { return RequestRules }
func ClearRequestRules() { RequestRules = make(map[string]map[string]types.RequestEvalRule) }