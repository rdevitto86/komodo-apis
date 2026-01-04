package main

import (
	"komodo-forge-apis-go/aws/dynamodb"
	sm "komodo-forge-apis-go/aws/secrets-manager"
	"komodo-forge-apis-go/config"
	mw "komodo-forge-apis-go/http/middleware"
	logger "komodo-forge-apis-go/logging/runtime"
	"komodo-user-api/internal/handlers"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	// initialize runtime logger
	logger.Init("komodo-user-api", config.GetConfigValue("LOG_LEVEL"))

	// initialize AWS Secrets Manager
	err := sm.Init(sm.Config{
		Region: config.GetConfigValue("AWS_REGION"),
		SecretPrefix: config.GetConfigValue("SECRET_PREFIX"),
		BatchSecretName: config.GetConfigValue("BATCH_SECRET_NAME"),
		SecretKeys: []string{
			"USER_API_CLIENT_ID",
			"USER_API_CLIENT_SECRET",
			"IP_WHITELIST",
			"IP_BLACKLIST",
		},
	})
	if err != nil {
		logger.Fatal("failed to load aws secrets", err)
		os.Exit(1)
	}

	// initialize DynamoDB client
	dynamodb.Init(dynamodb.Config{})

	// initialize chi router
	rtr := chi.NewRouter()

	// Health check endpoint
	rtr.Get("/health", handlers.HealthHandler)

	rtr.Use(
		mw.RequestIDMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.IPAccessMiddleware,
		mw.TelemetryMiddleware,
		mw.CORSMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
	)

	// User routes
	rtr.Route("/users", func(users chi.Router) {
		users.Use(
			mw.AuthMiddleware,
			mw.RateLimiterMiddleware,
		)

		// User CRUD
		users.Post("/", handlers.CreateUser)
		users.Route("/{user_id}", func(user chi.Router) {
			user.Use(mw.SanitizationMiddleware)
			user.Get("", handlers.GetUserByID)
			user.Put("", handlers.UpdateUserByID)
			user.Delete("", handlers.DeleteUserByID)
		})

		// Me routes
		users.Route("/me", func(me chi.Router) {
			me.Use(
				mw.IdempotencyMiddleware,
				mw.CSRFMiddleware,
			)

			// Profile management
			me.Post("/profile", handlers.GetProfile)
			me.Put("/profile", handlers.UpdateProfile)
			me.Delete("/profile", handlers.DeleteProfile)

			// Addresses management
			me.Post("/addresses/query", handlers.GetAddresses)
			me.Post("/addresses", handlers.AddAddress)
			me.Route("/addresses/{addr_id}", func(addr chi.Router) {
				addr.Use(mw.SanitizationMiddleware)
				addr.Put("", handlers.UpdateAddress)
				addr.Delete("", handlers.DeleteAddress)
			})

			// Preferences management
			me.Get("/preferences", handlers.GetPreferences)
			me.Put("/preferences", handlers.UpdatePreferences)
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
