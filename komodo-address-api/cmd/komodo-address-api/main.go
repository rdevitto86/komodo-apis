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

	router "github.com/go-chi/chi/v5"
	chi_mw "github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := router.NewRouter()

    // Add Chi middleware
    r.Use(chi_mw.RequestID)
    r.Use(chi_mw.Logger)
    r.Use(chi_mw.Recoverer)

    // Add custom authentication middleware
    validateTokenURL := os.Getenv("AUTH_SERVICE_VALIDATE_URL")
    if validateTokenURL == "" {
        log.Fatal("AUTH_SERVICE_VALIDATE_URL is not set")
    }
    r.Use(internal_mw.AuthMiddleware(validateTokenURL))

    // Define routes
    r.Get("/health", handlers.HandleHealth)
    r.Post("/validate", handlers.HandleValidate)
    r.Post("/normalize", handlers.HandleNormalize)
    r.Post("/geocode", handlers.HandleGeocode)

    addr := ":7001" // default port

    if p := os.Getenv("PORT"); strings.TrimSpace(p) != "" {
        addr = ":" + p
    }

    srv := &http.Server{
        Addr:              addr,
        Handler:           r, // Use Chi router directly
        ReadHeaderTimeout: 5 * time.Second,
        ReadTimeout:       10 * time.Second,
        WriteTimeout:      15 * time.Second,
        IdleTimeout:       60 * time.Second,
    }

    // Graceful shutdown
    go func() {
        log.Printf("komodo-address-api listening on %s", addr)

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
