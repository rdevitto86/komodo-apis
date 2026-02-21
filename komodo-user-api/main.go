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
)

func init() {
	logger.Init(
		config.GetConfigValue("APP_NAME"),
		config.GetConfigValue("LOG_LEVEL"),
		config.GetConfigValue("ENV"),
	)
}

// chain applies middleware in order: first listed = outermost wrapper.
func chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func main() {
	smCfg := awsSM.Config{
		Region:   config.GetConfigValue("AWS_REGION"),
		Endpoint: config.GetConfigValue("AWS_ENDPOINT"),
		Prefix:   config.GetConfigValue("AWS_SECRET_PREFIX"),
		Batch:    config.GetConfigValue("AWS_SECRET_BATCH"),
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
		Region:    config.GetConfigValue("AWS_REGION"),
		Endpoint:  config.GetConfigValue("DYNAMODB_ENDPOINT"),
		AccessKey: config.GetConfigValue("DYNAMODB_ACCESS_KEY"),
		SecretKey: config.GetConfigValue("DYNAMODB_SECRET_KEY"),
	}
	if err := dynamodb.Init(ddbCfg); err != nil {
		logger.Fatal("failed to initialize dynamodb", err)
		os.Exit(1)
	} else {
		logger.Info("dynamodb initialized successfully")
	}

	// Middleware stack for all /me routes
	meMW := []func(http.Handler) http.Handler{
		mw.RequestIDMiddleware,
		mw.TelemetryMiddleware,
		mw.RateLimiterMiddleware,
		mw.IPAccessMiddleware,
		mw.CORSMiddleware,
		mw.SecurityHeadersMiddleware,
		mw.AuthMiddleware,
		mw.CSRFMiddleware,
		mw.NormalizationMiddleware,
		mw.RuleValidationMiddleware,
		mw.SanitizationMiddleware,
		mw.IdempotencyMiddleware,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	mux.Handle("POST /me/profile", chain(http.HandlerFunc(handlers.GetProfile), meMW...))
	mux.Handle("PUT /me/profile", chain(http.HandlerFunc(handlers.UpdateProfile), meMW...))
	mux.Handle("DELETE /me/profile", chain(http.HandlerFunc(handlers.DeleteProfile), meMW...))
	mux.Handle("POST /me/profile/create", chain(http.HandlerFunc(handlers.CreateUser), meMW...))

	mux.Handle("POST /me/addresses/query", chain(http.HandlerFunc(handlers.GetAddresses), meMW...))
	mux.Handle("POST /me/addresses/create", chain(http.HandlerFunc(handlers.AddAddress), meMW...))
	mux.Handle("PUT /me/addresses/update", chain(http.HandlerFunc(handlers.UpdateAddress), meMW...))
	mux.Handle("DELETE /me/addresses/delete", chain(http.HandlerFunc(handlers.DeleteAddress), meMW...))

	mux.Handle("POST /me/orders", chain(http.HandlerFunc(handlers.GetOrders), meMW...))
	mux.Handle("PUT /me/orders", chain(http.HandlerFunc(handlers.UpdateOrder), meMW...))
	mux.Handle("POST /me/orders/create", chain(http.HandlerFunc(handlers.CreateOrder), meMW...))
	mux.Handle("POST /me/orders/cancel", chain(http.HandlerFunc(handlers.CancelOrder), meMW...))
	mux.Handle("POST /me/orders/return", chain(http.HandlerFunc(handlers.ReturnOrder), meMW...))

	mux.Handle("POST /me/payments", chain(http.HandlerFunc(handlers.GetPayments), meMW...))
	mux.Handle("PUT /me/payments", chain(http.HandlerFunc(handlers.UpsertPayment), meMW...))
	mux.Handle("DELETE /me/payments", chain(http.HandlerFunc(handlers.DeletePayment), meMW...))

	mux.Handle("GET /me/preferences", chain(http.HandlerFunc(handlers.GetPreferences), meMW...))
	mux.Handle("PUT /me/preferences", chain(http.HandlerFunc(handlers.UpdatePreferences), meMW...))

	server := &http.Server{
		Addr:              ":" + config.GetConfigValue("PORT"),
		Handler:           mux,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
	logger.Info("server started successfully")
}
