package main

import (
	"komodo-auth-api/internal/handlers"
	elasticache "komodo-forge-apis-go/aws/elasticache"
	"komodo-forge-apis-go/config"
	"komodo-forge-apis-go/crypto/jwt"
	bootstrap "komodo-forge-apis-go/http/common/bootstrap"
	mw "komodo-forge-apis-go/http/middleware"
	moxtox "komodo-forge-apis-go/test/moxtox"
	"net/http"
	"os"
	"time"

	logger "komodo-forge-apis-go/loggers/runtime"

	"github.com/go-chi/chi/v5"
)

func main() {
	init := bootstrap.Initialize(bootstrap.Options{
		AppName: "komodo-auth-api",
		Secrets: []string{
			"JWT_PUBLIC_KEY",
			"JWT_PRIVATE_KEY",
			"JWT_ENC_KEY",
			"JWT_HMAC_SECRET",
			"OAUTH_CLIENT_ID",
			"OAUTH_CLIENT_SECRET",
			"OAUTH_ENCRYPTION_KEY",
			"IP_WHITELIST",
			"IP_BLACKLIST",
		},
	})
	env, port := init.Env, init.Port

	// initialize Elasticache client
	elasticache.InitElasticacheClient()

	// initialize JWT keys (for M2M service tokens)
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize JWT keys", err)
		os.Exit(1)
	}

	// initialize router
	rtr := chi.NewRouter()

	// global middleware
	rtr.Use(mw.ContextMiddleware)
	rtr.Use(mw.TelemetryMiddleware)
	rtr.Use(mw.SecurityHeadersMiddleware)
	rtr.Use(mw.IPAccessMiddleware)
	rtr.Use(mw.NormalizationMiddleware)
	rtr.Use(mw.SanitizationMiddleware)
	rtr.Use(mw.RuleValidationMiddleware)

	// moxtox for mocking in non-prod
	if env != "prod" && config.GetConfigValue("USE_MOCKS") == "true" {
		logger.Info("using mocks in non-production environment")
		rtr.Use(moxtox.InitMoxtoxMiddleware(env))
	}

	// health check (public)
	rtr.Get("/health", handlers.HealthHandler)
	rtr.Get("/.well-known/jwks.json", handlers.JWKSHandler)

	// OAuth 2.0 endpoints (for client/session-based auth)
	rtr.Route("/oauth", func(oauth chi.Router) {
		// Public endpoints
		oauth.With(mw.RateLimiterMiddleware).Post("/token", handlers.OAuthTokenHandler)
		oauth.With(mw.RateLimiterMiddleware).Get("/authorize", handlers.OAuthAuthorizeHandler)
		
		// Protected endpoints (require valid OAuth token)
		oauth.Group(func(protected chi.Router) {
			protected.Use(mw.ClientTypeMiddleware)
			protected.Use(mw.AuthMiddleware)
			
			protected.Post("/introspect", handlers.OAuthIntrospectHandler)
			protected.Post("/revoke", handlers.OAuthRevokeHandler)
		})
	})

	logger.Info("server starting on port " + port)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           rtr,
		ReadTimeout:       5 * time.Second,  // 5 seconds
		WriteTimeout:      10 * time.Second, // 10 seconds
		IdleTimeout:       60 * time.Second, // 1 minute
		ReadHeaderTimeout: 2 * time.Second,  // prevents Slowloris attacks
		MaxHeaderBytes:    1 << 20,          // 1MB max headers
	}

	// start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
}
