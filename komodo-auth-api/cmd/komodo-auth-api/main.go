package main

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

	"github.com/go-chi/chi/v5"
)

func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)
}

func main() {
	smCfg := awsSM.Config{
		Region: config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Prefix: config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch: config.GetConfigValue("AWS_SECRET_BATCH"),
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
		DB: config.GetConfigValue("AWS_ELASTICACHE_DB"),
	}

	if err := awsEC.Init(ecCfg); err != nil {
		logger.Fatal("failed to initialize elasticache", err)
		os.Exit(1)
	}

	rtr := chi.NewRouter()
	rtr.Get("/health", handlers.HealthHandler)
	rtr.Get("/.well-known/jwks.json", handlers.JWKSHandler)

	rtr.Route("/oauth", func(oauth chi.Router) {
		oauth.Use(
			mw.RequestIDMiddleware,
			mw.TelemetryMiddleware,
			mw.RateLimiterMiddleware,
			mw.IPAccessMiddleware,
			mw.SecurityHeadersMiddleware,
			mw.NormalizationMiddleware,
			mw.SanitizationMiddleware,
			mw.RuleValidationMiddleware,
		)
		
		oauth.Post("/token", handlers.OAuthTokenHandler)
		oauth.Get("/authorize", handlers.OAuthAuthorizeHandler)
		
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

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
	logger.Info("server started successfully")
}
