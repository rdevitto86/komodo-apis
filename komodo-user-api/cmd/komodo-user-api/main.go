package main

import (
	"komodo-internal-lib-apis-go/aws/elasticache"
	"komodo-internal-lib-apis-go/config"
	"komodo-internal-lib-apis-go/http/common/bootstrap"
	mw "komodo-internal-lib-apis-go/http/middleware/gin"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"komodo-user-api/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	init := bootstrap.Initialize(bootstrap.Options{
		AppName: "komodo-user-api",
		Secrets: []string{
			"JWT_PUBLIC_KEY",
			"SESSION_ENCRYPTION_KEY",
			"CSRF_SECRET",
			"USER_API_CLIENT_ID",
			"USER_API_CLIENT_SECRET",
			"IP_WHITELIST",
			"IP_BLACKLIST",
		},
	})
	env, port := init.Env, init.Port

	// set gin mode
	if env == "prod" { gin.SetMode(gin.ReleaseMode) }

	// initialize Elasticache client
	elasticache.InitElasticacheClient()

	// initialize gin router
	router := gin.New()
	// router.Use(mw.Logger()) // TODO remove

	// initialize moxtox response handler
	if env != "prod" && config.GetConfigValue("USE_MOCKS") == "true" {
		logger.Info("using mocks in non-production environment")
		// TODO: moxtox needs Gin adapter
		// router.Use(moxtox.InitMoxtoxMiddleware(env))
	}

	// health check (public, no auth)
	router.GET("/health", handlers.HealthHandler)

	// ========================================
	// M2M Routes (Service-to-Service)
	// ========================================
	// Auth: JWT or OAuth service tokens only
	// Use: Internal microservice communication
	// Security: No CSRF, higher rate limits
	m2m := router.Group("/m2m")
	m2m.Use(mw.ServiceAuthMiddleware())
	{
		m2m.POST("/users", handlers.CreateUser)
		m2m.POST("/users/query", handlers.GetUserByID)
		m2m.PUT("/users", handlers.UpdateUserByID)
		m2m.DELETE("/users", handlers.DeleteUserByID)
	}

	// ========================================
	// Client Routes (Browser/Mobile Apps)
	// ========================================
	// Auth: Session-based (cookies)
	// Use: End-user interactions from browsers/mobile apps
	// Security: Full browser protections (CSRF, rate limiting, idempotency)
	me := router.Group("/users/me")
	// me.Use(mw.Session())        // TODO: Enable when session middleware is ready
	// me.Use(mw.CSRF())           // TODO: Enable CSRF protection
	// me.Use(mw.Idempotency())    // TODO: Enable idempotency for mutations
	{
		// Profile management
		me.POST("/profile/query", handlers.GetMyProfile)
		me.PUT("/profile", handlers.UpdateMyProfile)
		me.DELETE("/account", handlers.DeleteMyAccount)

		// Address management
		me.POST("/addresses/query", handlers.GetMyAddresses)
		me.POST("/addresses", handlers.AddMyAddress)
		me.PUT("/addresses", handlers.UpdateMyAddress)
		me.DELETE("/addresses", handlers.DeleteMyAddress)

		// Preferences management
		me.POST("/preferences/query", handlers.GetMyPreferences)
		me.PUT("/preferences", handlers.UpdateMyPreferences)
	}

	logger.Info("server starting on port " + port)

	// start server
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
}
