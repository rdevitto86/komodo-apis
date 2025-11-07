package main

import (
	"komodo-auth-internal-api/internal/httpapi/handlers"
	jwtUtils "komodo-auth-internal-api/internal/httpapi/utils/jwt"
	elasticache "komodo-internal-lib-apis-go/aws/elasticache"
	secretsManager "komodo-internal-lib-apis-go/aws/secrets-manager"
	"komodo-internal-lib-apis-go/config"
	mw "komodo-internal-lib-apis-go/http/middleware/chi"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	moxtox "komodo-internal-lib-apis-go/test/moxtox"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	logger.InitLogger()

	env := config.GetConfigValue("ENV")
	logger.Info("Starting komodo-auth-internal-api in " + env + " environment")

	switch strings.ToLower(env) {
		case "local":
			logger.Info("Running in local environment - skipping AWS Secrets Manager integration")
		case "dev", "staging", "prod":
			if !secretsManager.IsUsingAWS() {
				logger.Fatal("AWS Secrets Manager integration disabled")
				os.Exit(1)
			}
			logger.Info("AWS Secrets Manager integration enabled")

			secrets := []string{
				"JWT_PUBLIC_KEY",
				"JWT_PRIVATE_KEY",
				"JWT_ENC_KEY",
				"JWT_HMAC_SECRET",
				"IP_WHITELIST",
				"IP_BLACKLIST",
			}

			// load AWS Secrets
			if err := secretsManager.LoadSecrets(secrets); err != nil {
				logger.Fatal("Failed to get secrets", err)
				os.Exit(1)
			}
		default:
			logger.Fatal("Environment variable ENV invalid or not set")
			os.Exit(1)
	}

	// initialize Elasticache client
	elasticache.InitElasticacheClient()

	// initialize JWT keys
	if err := jwtUtils.InitializeKeys(); err != nil {
		logger.Fatal("Failed to initialize JWT keys", err)
		os.Exit(1)
	}
	logger.Info("JWT keys initialized successfully")

	// initialize router
	rtr := chi.NewRouter()

	// initialize global middleware
	rtr.Use(mw.ContextMiddleware)
	rtr.Use(mw.TelemetryMiddleware)
	rtr.Use(mw.SecurityHeadersMiddleware)
	rtr.Use(mw.IPAccessMiddleware)
	rtr.Use(mw.NormalizationMiddleware)
	rtr.Use(mw.SanitizationMiddleware)
	rtr.Use(mw.RuleValidationMiddleware)

	// initialize moxtox response handler
	if env != "prod" && os.Getenv("USE_MOCKS") == "true" {
		logger.Info("Using mocks in non-production environment")
		rtr.Use(moxtox.InitMoxtoxMiddleware(env))
	}

	// unprotected public routes
	rtr.Get("/health", handlers.HealthHandler)

	rtr.Route("/token", func(token chi.Router) {
		// public endpoint(s) that require no auth
		token.With(mw.RateLimiterMiddleware).Post("/", handlers.TokenCreateHandler)
		
		// protected endpoint(s) that require valid JWT
		token.Group(func(protected chi.Router) {
			protected.Use(mw.ClientTypeMiddleware)
			protected.Use(mw.AuthnJWTMiddleware)

			protected.Post("/introspect", handlers.TokenIntrospectHandler)
			protected.Delete("/revoke", handlers.TokenRevokeHandler)
		})
	})

	port := config.GetConfigValue("PORT")
	if port == "" { port =	 "7001" }
	logger.Info("Server starting on port " + port)

	server := &http.Server{
		Addr:         			":" + port,
		Handler:      			rtr,
		ReadTimeout:  			5 * time.Second, // 5 seconds
		WriteTimeout: 			10 * time.Second, // 10 seconds
		IdleTimeout:  		 	60 * time.Second, // 1 minute
		ReadHeaderTimeout: 	2 * time.Second,	// prevents Slowloris attacks
    MaxHeaderBytes:    	1 << 20, // 1MB max headers
	}

	// start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", err)
		os.Exit(1)
  }
}