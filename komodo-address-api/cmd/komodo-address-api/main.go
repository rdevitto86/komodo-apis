package main

import (
	"context"
	"errors"
	"komodo-address-api/internal/httpapi/handlers"
	internal_mw "komodo-address-api/internal/httpapi/middleware"
	"komodo-address-api/thirdparty/aws"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	env := os.Getenv("API_ENV")

  switch strings.ToLower(env) {
    case "prod":
			secret, err := aws.GetSecret("prod/db/password")

			if err != nil {
        log.Fatalf("failed to load secret: %v", err)
      }

      os.Setenv("DB_PASSWORD", secret)
			gin.SetMode(gin.ReleaseMode)
		case "staging":
			gin.SetMode(gin.ReleaseMode)
		case "dev":
			gin.SetMode(gin.DebugMode)
		default:
			log.Fatal("API_ENV is not set")
	}

	router := gin.New()

	// Gin middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Custom authentication middleware
	validateTokenURL := os.Getenv("AUTH_SERVICE_VALIDATE_URL")
	if validateTokenURL == "" {
		log.Fatal("AUTH_SERVICE_VALIDATE_URL is not set")
	}

	// Authentication middleware
	router.Use(func(c *gin.Context) {
		if err := internal_mw.AuthMiddleware(validateTokenURL, c); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	})

	// Define routes
	router.GET("/health", func(c *gin.Context) {
		handlers.HandleHealth(c)
	})
	router.POST("/validate", func(c *gin.Context) {
		handlers.HandleValidate(c)
	})
	router.POST("/normalize", func(c *gin.Context) {
		handlers.HandleNormalize(c)
	})
	router.POST("/geocode", func(c *gin.Context) {
		handlers.HandleGeocode(c)
	})

	serverAddress := ":7010"
	if port := os.Getenv("PORT"); strings.TrimSpace(port) != "" {
		serverAddress = ":" + port
	}

	srv := &http.Server{
		Addr:              serverAddress,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("komodo-address-api listening on %s", serverAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
