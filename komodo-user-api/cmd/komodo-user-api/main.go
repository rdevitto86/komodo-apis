package main

import (
	"komodo-forge-sdk-go/aws/dynamodb"
	awsSM "komodo-forge-sdk-go/aws/secrets-manager"
	"komodo-forge-sdk-go/config"
	mw "komodo-forge-sdk-go/http/middleware"
	logger "komodo-forge-sdk-go/logging/runtime"
	"komodo-user-api/internal/handlers"
	"net/http"
	"os"
	"time"

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
			"DYNAMODB_ENDPOINT",
			"DYNAMODB_ACCESS_KEY",
			"DYNAMODB_SECRET_KEY",
			"USER_API_CLIENT_ID",
			"USER_API_CLIENT_SECRET",
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
	} else {
		logger.Info("aws secrets manager initialized successfully")
	}

	ddbCfg := dynamodb.Config{
		Region: config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("DYNAMODB_ENDPOINT"),
		AccessKey: config.GetConfigValue("DYNAMODB_ACCESS_KEY"),
		SecretKey: config.GetConfigValue("DYNAMODB_SECRET_KEY"),
	}
	if err := dynamodb.Init(ddbCfg); err != nil {
		logger.Fatal("failed to initialize dynamodb", err)
		os.Exit(1)
	} else {
		logger.Info("dynamodb initialized successfully")
	}

	rtr := chi.NewRouter()
	rtr.Get("/health", handlers.HealthHandler)

	rtr.Route("/users", func(users chi.Router) {
		users.Use(
			mw.RequestIDMiddleware,
			mw.TelemetryMiddleware,
			mw.RateLimiterMiddleware,
			mw.IPAccessMiddleware,
			mw.CORSMiddleware,
			mw.SecurityHeadersMiddleware,
			mw.AuthMiddleware,
			mw.NormalizationMiddleware,
			mw.RuleValidationMiddleware,
		)

		users.Post("/", handlers.CreateUser)

		users.Route("/{user_id}", func(user chi.Router) {
			user.Use(mw.SanitizationMiddleware)
			user.Get("", handlers.GetUserByID)
			user.Put("", handlers.UpdateUserByID)
			user.Delete("", handlers.DeleteUserByID)
		})

		users.Route("/me", func(me chi.Router) {
			me.Use(mw.IdempotencyMiddleware, mw.CSRFMiddleware)

			me.Post("/profile", handlers.GetProfile)
			me.Put("/profile", handlers.UpdateProfile)
			me.Delete("/profile", handlers.DeleteProfile)
			me.Post("/addresses/query", handlers.GetAddresses)
			me.Post("/addresses", handlers.AddAddress)
			me.Get("/preferences", handlers.GetPreferences)
			me.Put("/preferences", handlers.UpdatePreferences)

			me.Route("/addresses/{addr_id}", func(addr chi.Router) {
				addr.Use(mw.SanitizationMiddleware)
				addr.Put("", handlers.UpdateAddress)
				addr.Delete("", handlers.DeleteAddress)
			})
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
