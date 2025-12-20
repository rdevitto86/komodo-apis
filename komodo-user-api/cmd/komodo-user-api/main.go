package main

import (
	"komodo-forge-apis-go/aws/elasticache"
	"komodo-forge-apis-go/config"
	"komodo-forge-apis-go/http/common/bootstrap"
	mw "komodo-forge-apis-go/http/middleware/gin"
	logger "komodo-forge-apis-go/logging/runtime"
	"komodo-user-api/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	init := bootstrap.Initialize(bootstrap.Options{
		AppName: "komodo-user-api",
		Secrets: []string{
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

	router.GET("/health", handlers.HealthHandler)

	user := router.Group("/users")
	user.Use(mw.AuthMiddleware())

	// User CRUD
	user.POST("/", handlers.CreateUser)
	user.GET("/:user_id", handlers.GetUserByID)
	user.PUT("/:user_id", handlers.UpdateUserByID)
	user.DELETE("/:user_id", handlers.DeleteUserByID)

	me := user.Group("/me")
	me.Use(mw.IdempotencyMiddleware())

	// Profile management
	me.POST("/profile", handlers.GetMyProfile)
	me.PUT("/profile", handlers.UpdateMyProfile)
	me.DELETE("/", handlers.DeleteMyAccount)

	// Addresses management
	me.POST("/addresses/query", handlers.GetMyAddresses)
	me.POST("/addresses", handlers.AddMyAddress)
	me.PUT("/addresses/:addr_id", handlers.UpdateMyAddress)
	me.DELETE("/addresses/:addr_id", handlers.DeleteMyAddress)

	// Preferences management
	me.GET("/preferences", handlers.GetMyPreferences)
	me.PUT("/preferences", handlers.UpdateMyPreferences)

	logger.Info("server starting on port " + port)

	// start server
	if err := router.Run(":" + port); err != nil {
		logger.Fatal("server failed to start", err)
		os.Exit(1)
	}
}
