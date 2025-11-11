package main

import (
	"komodo-auth-public-api/internal/handlers"
	jwtUtils "komodo-internal-lib-apis-go/auth/jwt"
	elasticache "komodo-internal-lib-apis-go/aws/elasticache"
	secretsManager "komodo-internal-lib-apis-go/aws/secrets-manager"
	"komodo-internal-lib-apis-go/config"
	mw "komodo-internal-lib-apis-go/http/middleware/chi"
	logger "komodo-internal-lib-apis-go/logging/runtime"
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
	logger.Info("starting komodo-auth-public-api in " + env + " environment")

	switch strings.ToLower(env) {
		case "local":
			logger.Info("Running in local environment - skipping AWS Secrets Manager integration")
		case "dev", "staging", "prod":
			if !secretsManager.IsUsingAWS() {
				logger.Fatal("aws secrets manager integration disabled")
				os.Exit(1)
			}
			logger.Info("aws secrets manager integration enabled")

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
		default:
			logger.Fatal("environment variable ENV invalid or not set")
			os.Exit(1)
	}

	// initialize Elasticache client
	elasticache.InitElasticacheClient()

	// initialize JWT keys
	if err := jwtUtils.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize jwt keys", err)
		os.Exit(1)
	}
	logger.Info("jwt keys initialized successfully")

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
		logger.Info("using mocks in non-production environment")
		rtr.Use(moxtox.InitMoxtoxMiddleware(env))
	}

	rtr.Get("/health", handlers.HealthHandler)

	rtr.Route("/auth", func(auth chi.Router) {
		auth.Post("/login", handlers.LoginHandler)
		auth.Post("/register", handlers.RegisterHandler)
		auth.Post("/password/reset", handlers.PasswordResetHandler)
		auth.Post("/password/reset/{token}", handlers.PasswordResetCompleteHandler)
		auth.Post("/passkey/login", handlers.PasskeyLoginHandler)
		auth.Post("/passkey/login/verify", handlers.PasskeyLoginVerifyHandler)

		// Protected routes requiring session + CSRF
		auth.Group(func(protected chi.Router) {
			protected.Use(mw.SessionMiddleware)
			protected.Use(mw.CSRFMiddleware)
			protected.Use(mw.IdempotencyMiddleware)

			protected.Post("/logout", handlers.LogoutHandler)
			protected.Get("/me", handlers.AuthMeGetHandler)
			protected.Patch("/me", handlers.AuthMePatchHandler)
			protected.Get("/sessions", handlers.SessionsListHandler)
			protected.Delete("/sessions/{id}", handlers.SessionsRevokeHandler)
			protected.Post("/mfa/setup", handlers.MFASetupHandler)
			protected.Post("/mfa/enable", handlers.MFAEnableHandler)
			protected.Post("/mfa/disable", handlers.MFADisableHandler)
			protected.Post("/mfa/verify", handlers.MFAVerifyHandler)
			protected.Post("/passkey/register", handlers.PasskeyRegisterHandler)
			protected.Post("/passkey/register/verify", handlers.PasskeyRegisterVerifyHandler)
		})
	})

	port := config.GetConfigValue("PORT")
	if port == "" { port =	 "7002" }
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
		logger.Fatal("server failed to start", err)
		os.Exit(1)
  }
}
