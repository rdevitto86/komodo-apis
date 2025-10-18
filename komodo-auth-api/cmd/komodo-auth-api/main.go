package main

import (
	"komodo-auth-api/internal/httpapi/handlers"
	secretsManager "komodo-internal-lib-apis-go/aws/secrets-manager"
	"komodo-internal-lib-apis-go/config"
	mw "komodo-internal-lib-apis-go/http/middleware/chi"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	moxtox "komodo-internal-lib-apis-go/test/moxtox"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// initialize logger
	logger.InitLogger()

	env := config.GetConfigValue("ENV")
	logger.Info("starting komodo-auth-api in " + env + " environment")

	// load secrets from AWS Secrets Manager in prod/staging
	switch env {
		case "dev", "staging", "prod":
			// load secrets from AWS Secrets Manager
			if secretsManager.IsUsingAWS() {
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
					logger.Fatal("failed to get secrets", err)
					os.Exit(1)
				}
			} else {
				logger.Warn("AWS Secrets Manager integration disabled")
			}
		default:
			logger.Fatal("environment variable ENV invalid or not set")
			os.Exit(1)
	}

	// initialize Elasticache client
	// elasticache.InitElasticacheClient()

	// initialize router
	rtr := chi.NewRouter()

	// initialize global middleware
	rtr.Use(mw.ContextMiddleware)
	rtr.Use(mw.TelemetryMiddleware)
	rtr.Use(mw.NormalizationMiddleware)
	rtr.Use(mw.SanitizationMiddleware)
	rtr.Use(mw.SecurityHeadersMiddleware)
	rtr.Use(mw.IPAccessMiddleware)
	rtr.Use(mw.RuleValidationMiddleware)

	// initialize moxtox response handler
	if env != "prod" && os.Getenv("USE_MOCKS") == "true" {
		logger.Info("using mocks in non-production environment")
		rtr.Use(moxtox.InitMoxtoxMiddleware(env))
	}

	// unprotected public routes
	rtr.Get("/health", handlers.HealthHandler)

	rtr.Route(("/v" + os.Getenv("VERSION")), func(ver chi.Router) {
		ver.Route("/auth", func(auth chi.Router) {
			// rate limited public endpoints
			auth.With(mw.RateLimiterMiddleware).Post("/login", handlers.LoginHandler)
      auth.With(mw.RateLimiterMiddleware).Post("/token", handlers.TokenCreateHandler)
      auth.With(mw.RateLimiterMiddleware).Post("/token/refresh", handlers.TokenRefreshHandler)

			// protected endpoints
			auth.Group(func(protected chi.Router) {
				protected.Use(mw.AuthnJWTMiddleware)
				protected.Use(mw.CSRFMiddleware)
				protected.Use(mw.IdempotencyMiddleware)
				
				protected.Post("/logout", handlers.LogoutHandler)
				protected.Post("/mfa/disable", handlers.MFADisableHandler)
				protected.Post("/mfa/enable", handlers.MFAEnableHandler)
				protected.Post("/mfa/setup", handlers.MFASetupHandler)
				protected.Post("/mfa/verify", handlers.MFAVerifyHandler)
				protected.Post("/passkey/start", handlers.PasskeyStartHandler)
				protected.Post("/passkey/verify", handlers.PasskeyVerifyHandler)
				protected.Delete("/token", handlers.TokenDeleteHandler)
				protected.Post("/token/verify", handlers.TokenVerifyHandler)
			})
		})
	})

	port := config.GetConfigValue("PORT")
	if port == "" { port = "7001" }
	logger.Info("server starting on port " + port)

	server := &http.Server{
		Addr:         (":" + port),
		Handler:      rtr,
		ReadTimeout:  5 * time.Second, // 5 seconds
		WriteTimeout: 10 * time.Second, // 10 seconds
		IdleTimeout:  60 * time.Second, // 1 minute
	}

	// start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
  }
}