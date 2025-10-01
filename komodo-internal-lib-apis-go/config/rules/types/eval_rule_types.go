package types

import "komodo-internal-lib-apis-go/http/types"

type ValidationLevel int
const (
	LevelIgnore   ValidationLevel = 0
	LevelStrict  	ValidationLevel = 1
	LevelLenient  ValidationLevel = 2
)

type RuleToggle int
const (
	RuleOff  	RuleToggle = 0
	RuleOn  	RuleToggle = 1
	RuleOpt  	RuleToggle = 2
)

type RequestEvalRule struct {
	Path             		string
	Method            	string
	Type              	[]types.RequestType
	Level             	ValidationLevel
	Headers           	types.Headers
	PathParams        	types.PathParams
	QueryParams       	types.QueryParams
	Body              	types.Body
	RequireVersion    	RuleToggle
}
