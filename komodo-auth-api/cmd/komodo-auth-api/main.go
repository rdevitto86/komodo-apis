package main

import (
	"komodo-auth-api/internal/handlers"
	awsEC "komodo-forge-apis-go/aws/elasticache"
	awsSM "komodo-forge-apis-go/aws/secrets-manager"
	"komodo-forge-apis-go/config"
	"komodo-forge-apis-go/crypto/jwt"
	mw "komodo-forge-apis-go/http/middleware"
	"net/http"
	"os"
	"time"

	logger "komodo-forge-apis-go/logging/runtime"

	"github.com/go-chi/chi/v5"
)

func main() {
	// initialize runtime logger
	logger.Init(config.GetConfigValue("APP_NAME"), config.GetConfigValue("LOG_LEVEL"))

	smCfg := awsSM.Config{
		Region: config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Keys: []string{
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_AUDIENCE",
			"JWT_ISSUER",
			"JWT_KID",
			"AWS_ELASTICACHE_PASSWORD",
			"IP_WHITELIST",
			"IP_BLACKLIST",
		},
		Prefix: config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch: config.GetConfigValue("AWS_BATCH_SECRET_NAME"),
	}

	// initialize AWS Secrets Manager
	if err := awsSM.Bootstrap(smCfg); err != nil {
		logger.Fatal("failed to initialize aws secrets manager", err)
		os.Exit(1)
	}

	// load JWT keys into ENV
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	ecCfg := awsEC.Config{
		Endpoint: config.GetConfigValue("AWS_ELASTICACHE_ENDPOINT"),
		Password: config.GetConfigValue("AWS_ELASTICACHE_PASSWORD"),
		DB: config.GetConfigValue("AWS_ELASTICACHE_DB"),
	}

	// initialize elasticache
	if err := awsEC.Init(ecCfg); err != nil {
		logger.Fatal("failed to initialize elasticache", err)
		os.Exit(1)
	}

	rtr := chi.NewRouter()

	// health check (public - no middleware)
	rtr.Get("/health", handlers.HealthHandler)
	rtr.Get("/.well-known/jwks.json", handlers.JWKSHandler)

	// OAuth 2.0 endpoints (for client/session-based auth)
	rtr.Route("/oauth", func(oauth chi.Router) {
		oauth.Use(
			mw.ContextMiddleware,
			mw.SecurityHeadersMiddleware,
			mw.RequestIDMiddleware,
			mw.TelemetryMiddleware,
			mw.IPAccessMiddleware,
			mw.NormalizationMiddleware,
			mw.SanitizationMiddleware,
			mw.RuleValidationMiddleware,
			mw.RateLimiterMiddleware,
		)
		
		// Public endpoints
		oauth.Post("/token", handlers.OAuthTokenHandler)
		oauth.Get("/authorize", handlers.OAuthAuthorizeHandler)
		
		// Protected endpoints (require valid OAuth token)
		oauth.Group(func(protected chi.Router) {
			protected.Use(mw.ClientTypeMiddleware, mw.AuthMiddleware)

			protected.Post("/introspect", handlers.OAuthIntrospectHandler)
			protected.Post("/revoke", handlers.OAuthRevokeHandler)
		})
	})

	server := &http.Server{
		Addr: ":" + config.GetConfigValue("PORT"),
		Handler: rtr,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
	logger.Info("server started successfully")
}
