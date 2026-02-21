package komodoauthapi

import (
	"komodo-auth-api/internal/handlers"
	awsEC "komodo-forge-sdk-go/aws/elasticache"
	awsSM "komodo-forge-sdk-go/aws/secrets-manager"
	"komodo-forge-sdk-go/config"
	"komodo-forge-sdk-go/crypto/jwt"
	mw "komodo-forge-sdk-go/http/middleware"
	"net/http"
	"os"
	"time"

	logger "komodo-forge-sdk-go/logging/runtime"
)

func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)
}

// chain applies middleware in order: first listed = outermost wrapper.
func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func main() {
	smCfg := awsSM.Config{
		Region:   config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Prefix:   config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch:    config.GetConfigValue("AWS_SECRET_BATCH"),
		Keys: []string{
			"AWS_ELASTICACHE_ENDPOINT",
			"AWS_ELASTICACHE_PASSWORD",
			"AWS_ELASTICACHE_DB",
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_AUDIENCE",
			"JWT_ISSUER",
			"JWT_KID",
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"MAX_CONTENT_LENGTH",
			"IDEMPOTENCY_TTL_SEC",
			"RATE_LIMIT_RPS",
			"RATE_LIMIT_BURST",
			"BUCKET_TTL_SECOND",
		},
	}

	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize aws secrets manager", err)
		os.Exit(1)
	}

	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	ecCfg := awsEC.Config{
		Endpoint: config.GetConfigValue("AWS_ELASTICACHE_ENDPOINT"),
		Password: config.GetConfigValue("AWS_ELASTICACHE_PASSWORD"),
		DB:       config.GetConfigValue("AWS_ELASTICACHE_DB"),
	}

	if err := awsEC.Init(ecCfg); err != nil {
		logger.Fatal("failed to initialize elasticache", err)
		os.Exit(1)
	}

	// Shared middleware stack for all /oauth routes
	oauthMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.IPAccessMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.NormalizationMiddleware,
		mw.SanitizationMiddleware,
		mw.RuleValidationMiddleware,
	}

	// Extended middleware stack for protected /oauth routes
	protectedMW := append(oauthMW, mw.ClientTypeMiddleware, mw.AuthMiddleware)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("GET /.well-known/jwks.json", handlers.JWKSHandler)

	mux.Handle("POST /oauth/token", chain(http.HandlerFunc(handlers.OAuthTokenHandler), oauthMW...))
	mux.Handle("GET /oauth/authorize", chain(http.HandlerFunc(handlers.OAuthAuthorizeHandler), oauthMW...))

	mux.Handle("POST /oauth/introspect", chain(http.HandlerFunc(handlers.OAuthIntrospectHandler), protectedMW...))
	mux.Handle("POST /oauth/revoke", chain(http.HandlerFunc(handlers.OAuthRevokeHandler), protectedMW...))

	server := &http.Server{
		Addr:              ":" + config.GetConfigValue("PORT"),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
	logger.Info("server started successfully")
}
