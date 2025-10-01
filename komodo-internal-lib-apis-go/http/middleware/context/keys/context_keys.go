package config

type ctxKey string

const (
	StartTimeKey 					ctxKey = "start_time"
	ApiVersionKey 				ctxKey = "api_version"
	UriKey       					ctxKey = "uri"
	PathParamsKey 				ctxKey = "path_params"
	QueryParamsKey 				ctxKey = "query_params"
	ValidationRuleKey 		ctxKey = "validation_rule"
	RequestIDKey     			ctxKey = "request_id"
	RequestTimeoutKey    	ctxKey = "request_timeout"
)
