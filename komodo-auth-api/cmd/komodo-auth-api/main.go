package main

import (
	"komodo-auth-api/internal/config"
	"komodo-auth-api/internal/httpapi/handlers"
	mw "komodo-auth-api/internal/httpapi/middleware"
	"komodo-auth-api/internal/logger"
	"komodo-auth-api/internal/thirdparty/aws"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	env := config.GetConfigValue("API_ENV")

	// Load secrets from AWS Secrets Manager in prod/staging
	switch env {
		case "dev", "staging", "prod":
			if config.GetConfigValue("USE_MOCKS") == "true" {
				if env == "prod" {
					logger.Fatal("mocks cannot be used in production", nil)
					os.Exit(1)
				} else {
					logger.Warn("using mocks in non-production environment", nil)
					break
				}
			}

			secrets, err := aws.GetSecrets([]string{
				"JWT_PUBLIC_KEY",
				"JWT_PRIVATE_KEY",
				"JWT_ENC_KEY",
				"JWT_HMAC_SECRET",
				"IP_WHITELIST",
				"IP_BLACKLIST",
			})
			if err != nil && env != "dev" {
				logger.Fatal("failed to get secrets: "+err.Error(), nil)
				os.Exit(1)
			}

			for k, v := range secrets {
				config.SetConfigValue(k, v)
			}
		default:
			logger.Fatal("environment variable API_ENV is not set", nil)
			os.Exit(1)
	}
	logger.Info("starting komodo-auth-api in " + env + " environment", nil)

	// Initialize Elasticache client
	aws.InitElasticacheClient()

	// Initialize router
	rtr := chi.NewRouter()

	// Initialize middleware
	rtr.Use(mw.SecurityHeadersMiddleware)
	rtr.Use(mw.CanonicalizeMiddleware)
	rtr.Use(mw.ContextMiddleware)
	rtr.Use(mw.IPAccessMiddleware)
	rtr.Use(mw.RateLimiterMiddleware)
	rtr.Use(mw.AuthnJWTMiddleware)
	rtr.Use(mw.RuleValidationMiddleware)
	rtr.Use(mw.CSRFMiddleware)
	rtr.Use(mw.IdempotencyMiddleware)
	rtr.Use(mw.TelemetryMiddleware)
	rtr.Use(mw.ResponsePreprocessorMiddleware)
	if false && os.Getenv("USE_MOCKS") == "true" {
		rtr.Use(mw.InitMoxtoxMiddleware(
			"test/mocks/config/moxtox.json",
			"test/mocks/data",
			true,
		))
	}

	// Initialize routes
	rtr.Get("/health", handlers.HealthHandler)
	rtr.Get("/.well-known/jwks.json", handlers.JWKSHandler)

	rtr.Route(("/v" + os.Getenv("API_VERSION")), func(ver chi.Router) {
		ver.Route("/auth", func(auth chi.Router) {
			auth.Post("/login", handlers.LoginHandler)
			auth.Post("/logout", handlers.LogoutHandler)
			auth.Post("/mfa/disable", handlers.MFADisableHandler)
			auth.Post("/mfa/enable", handlers.MFAEnableHandler)
			auth.Post("/mfa/setup", handlers.MFASetupHandler)
			auth.Post("/mfa/verify", handlers.MFAVerifyHandler)
			auth.Post("/passkey/start", handlers.PasskeyStartHandler)
			auth.Post("/passkey/verify", handlers.PasskeyVerifyHandler)
			auth.Post("/token", handlers.TokenCreateHandler)
			auth.Delete("/token", handlers.TokenDeleteHandler)
			auth.Post("/token/refresh", handlers.TokenRefreshHandler)
		})
	})

	// Start server
	(&http.Server{
		Addr:         ":7001",
		Handler:      rtr,
		ReadTimeout:  5 * time.Second, // 5 seconds
		WriteTimeout: 10 * time.Second, // 10 seconds
		IdleTimeout:  60 * time.Second, // 1 minute
	}).ListenAndServe()
}