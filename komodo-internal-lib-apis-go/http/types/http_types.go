package types

const HEADER_AUTH = "Authorization"
const HEADER_CONTENT_TYPE = "Content-Type"
const HEADER_ACCEPT = "Accept"
const HEADER_USER_AGENT = "User-Agent"
const HEADER_REFERER = "Referer"
const HEADER_CONTENT_LENGTH = "Content-Length"
const HEADER_COOKIE = "Cookie"
const HEADER_CACHE_CONTROL = "Cache-Control"
const HEADER_X_REQUESTED_BY = "X-Requested-By"
const HEADER_X_REQUESTED_WITH = "X-Requested-With"
const HEADER_X_SESSION = "X-Session-Token"
const HEADER_X_CSRF_TOKEN = "X-CSRF-Token"

type RequestType string
const (
	REQ_TYPE_API_INT 		RequestType = "API_INTERNAL"
	REQ_TYPE_API_EXT 		RequestType = "API_EXTERNAL"
	REQ_TYPE_UI_USER   	RequestType = "UI_USER"
	REQ_TYPE_UI_GUEST		RequestType = "UI_GUEST"
	REQ_TYPE_ADMIN 			RequestType = "ADMIN"
)

type HeaderSpec struct {
	Type     string // "string","int","uuid",...
	Required bool
	Pattern  string // optional regex
}
type Headers map[string]HeaderSpec

type ParamSpec struct {
	Type     string // "string","int","uuid",...
	Required bool
	Pattern  string // optional regex 
}
type PathParams map[string]ParamSpec

type QueryParamSpec struct {
	Type     string
	Required bool
	Multiple bool
	Default  string
	Pattern  string
}
type QueryParams map[string]QueryParamSpec

type BodySpec struct {
	Type     string   // "string","int","bool","object","array"
	Required bool
	MinLen   int
	MaxLen   int
	Pattern  string
	Enum     []string
	// For nested objects you can embed a schema:
	Props map[string]BodySpec
}
type Body map[string]BodySpec

type Request struct {
	Method           	string
	Type              []RequestType
	Headers           Headers
	PathParams        PathParams
	QueryParams       QueryParams
	Body              Body
}
