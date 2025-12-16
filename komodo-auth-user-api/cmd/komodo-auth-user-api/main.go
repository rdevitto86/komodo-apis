package main

import (
	"komodo-auth-user-api/internal/handlers"
	elasticache "komodo-internal-lib-apis-go/aws/elasticache"
	"komodo-internal-lib-apis-go/config"
	"komodo-internal-lib-apis-go/crypto/jwt"
	bootstrap "komodo-internal-lib-apis-go/http/common/bootstrap"
	mw "komodo-internal-lib-apis-go/http/middleware/chi"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	moxtox "komodo-internal-lib-apis-go/test/moxtox"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	init := bootstrap.Initialize(bootstrap.Options{
		AppName: "komodo-auth-user-api",
		Secrets: []string{
			"JWT_PUBLIC_KEY",
			"SESSION_ENCRYPTION_KEY",
			"CSRF_SECRET",
			"DATABASE_URL",
			"IP_WHITELIST",
			"IP_BLACKLIST",
			"PASSKEY_RP_ID",
			"MFA_ISSUER",
			"EMAIL_API_KEY",
		},
	})
	env, port := init.Env, init.Port

	// initialize Elasticache client
	elasticache.InitElasticacheClient()

	// initialize JWT keys
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize jwt keys", err)
		os.Exit(1)
	}

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
	if env != "prod" && config.GetConfigValue("USE_MOCKS") == "true" {
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

	logger.Info("server starting on port " + port)

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
