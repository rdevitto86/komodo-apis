package main

import (
	"komodo-internal-lib-apis-go/aws/elasticache"
	secretsManager "komodo-internal-lib-apis-go/aws/secrets-manager"
	"komodo-internal-lib-apis-go/config"
	"komodo-internal-lib-apis-go/crypto/jwt"
	mw "komodo-internal-lib-apis-go/http/middleware/gin"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"komodo-user-api/internal/handlers"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()

	env := config.GetConfigValue("ENV")
	logger.Info("starting komodo-user-api in " + env + " environment")

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

			// set gin mode
			if env == "prod" {
				gin.SetMode(gin.ReleaseMode)
			}
		default:
			logger.Fatal("environment variable ENV invalid or not set")
			os.Exit(1)
	}

	// initialize Elasticache client
	elasticache.InitElasticacheClient()

	// initialize JWT keys
	if err := jwt.InitializeKeys(); err != nil {
		logger.Fatal("failed to initialize jwt keys", err)
		os.Exit(1)
	}
	logger.Info("jwt keys initialized successfully")

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

	// internal API routes (JWT or OAuth service tokens)
	m2m := router.Group("/m2m")
	m2m.Use(mw.ServiceAuthMiddleware())
	{
		// m2m.POST("/users", handlers.CreateUser)
		m2m.GET("/users/:id", handlers.GetUserHandler)
		// m2m.PATCH("/users/:id", handlers.UpdateUserByID)
		// m2m.DELETE("/users/:id", handlers.DeleteUserByID)
	}

	// user-facing routes (session-based auth)
	// me := router.Group("/users/me")
	// me.Use(mw.Session())
	// me.Use(mw.CSRF())
	// me.Use(mw.Idempotency())
	// {
	// 	profile management
	// 	me.GET("", handlers.GetMyProfile)
	// 	me.PATCH("", handlers.UpdateMyProfile)
	// 	me.DELETE("", handlers.DeleteMyAccount)

	// 	address management
	// 	me.GET("/addresses", handlers.GetMyAddresses)
	// 	me.POST("/addresses", handlers.AddMyAddress)
	// 	me.PATCH("/addresses/:addr_id", handlers.UpdateMyAddress)
	// 	me.DELETE("/addresses/:addr_id", handlers.DeleteMyAddress)

	// 	preferences management
	// 	me.GET("/preferences", handlers.GetMyPreferences)
	// 	me.PATCH("/preferences", handlers.UpdateMyPreferences)
	// }

	port := config.GetConfigValue("PORT")
	logger.Info("Server starting on port " + port)

	// start server
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
}
