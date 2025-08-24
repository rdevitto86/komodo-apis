package main

import (
	context "context"
	"errors"
	"komodo-address-api/internal/httpapi/handlers"
	internal_mw "komodo-address-api/internal/httpapi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	chi_router "github.com/go-chi/chi/v5"
	chi_mw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi_router.NewRouter()

	// Add Chi middleware
	router.Use(chi_mw.RequestID)
	router.Use(chi_mw.Logger)
	router.Use(chi_mw.Recoverer)

	// Add custom authentication middleware
	validateTokenURL := os.Getenv("AUTH_SERVICE_VALIDATE_URL")

	if validateTokenURL == "" {
		log.Fatal("AUTH_SERVICE_VALIDATE_URL is not set")
	}
	router.Use(internal_mw.AuthMiddleware(validateTokenURL))

	// Define routes
	router.Get("/health", handlers.HandleHealth)
	router.Post("/validate", handlers.HandleValidate)
	router.Post("/normalize", handlers.HandleNormalize)
	router.Post("/geocode", handlers.HandleGeocode)

	serverAddress := ":7010" // default port

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
