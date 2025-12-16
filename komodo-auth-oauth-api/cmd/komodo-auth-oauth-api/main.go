package main

import (
	elasticache "komodo-internal-lib-apis-go/aws/elasticache"
)

func main() {
	// init := bootstrap.Initialize(bootstrap.Options{
	// 	AppName: "komodo-auth-oauth-api",
	// 	Secrets: []string{
	// 		"JWT_PUBLIC_KEY",
	// 		"OAUTH_CLIENT_ID",
	// 		"OAUTH_CLIENT_SECRET",
	// 		"IP_WHITELIST",
	// 		"IP_BLACKLIST",
	// 	},
	// })
	// env, port := init.Env, init.Port
	
	// initialize Elasticache client
	elasticache.InitElasticacheClient()
	
	// initialize router
	// rtr := chi.NewRouter()

	// TODO global middleware
}
