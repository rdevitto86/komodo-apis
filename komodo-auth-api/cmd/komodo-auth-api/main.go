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
	env := config.GetConfigValue("ENV")

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

			secrets, err := secretsManager.GetSecrets([]string{
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
			logger.Fatal("environment variable ENV is not set", nil)
			os.Exit(1)
	}
	logger.Info("starting komodo-auth-api in " + env + " environment", nil)

	// Initialize Elasticache client
	// elasticache.InitElasticacheClient()

	// Initialize router
	rtr := chi.NewRouter()

	// Initialize middleware
	rtr.Use(mw.ContextMiddleware)
	rtr.Use(mw.TelemetryMiddleware)
	rtr.Use(mw.NormalizationMiddleware)
	rtr.Use(mw.SanitizationMiddleware)
	rtr.Use(mw.SecurityHeadersMiddleware)
	rtr.Use(mw.IPAccessMiddleware)
	rtr.Use(mw.RateLimiterMiddleware)
	rtr.Use(mw.AuthnJWTMiddleware)
	rtr.Use(mw.CSRFMiddleware)
	rtr.Use(mw.IdempotencyMiddleware)
	rtr.Use(mw.RuleValidationMiddleware)

	// Initialize moxtox response handler
	if env != "prod" && os.Getenv("USE_MOCKS") == "true" {
		rtr.Use(moxtox.InitMoxtoxMiddleware(env))
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
	server := &http.Server{
		Addr:         ":7001",
		Handler:      rtr,
		ReadTimeout:  5 * time.Second, // 5 seconds
		WriteTimeout: 10 * time.Second, // 10 seconds
		IdleTimeout:  60 * time.Second, // 1 minute
	}
	server.ListenAndServe()
}