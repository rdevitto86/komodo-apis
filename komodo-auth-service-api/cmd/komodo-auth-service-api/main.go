package main

import (
	"komodo-auth-service-api/internal/handlers"
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
		AppName: "komodo-auth-service-api",
		Secrets: []string{
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_ENC_KEY",
			"JWT_HMAC_SECRET",
			"IP_WHITELIST",
			"IP_BLACKLIST",
		},
	})
	env, port := init.Env, init.Port

	// initialize Elasticache client
	elasticache.InitElasticacheClient()

	// initialize JWT keys
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
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

	rtr.Route("/token", func(token chi.Router) {
		token.With(mw.RateLimiterMiddleware).Post("/", handlers.TokenCreateHandler)
		
		// protected endpoint(s) that require valid JWT
		token.Group(func(protected chi.Router) {
			protected.Use(mw.ClientTypeMiddleware)
			protected.Use(mw.AuthnJWTMiddleware)

			protected.Post("/verify", handlers.TokenVerifyHandler)
			protected.Delete("/revoke", handlers.TokenRevokeHandler)
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