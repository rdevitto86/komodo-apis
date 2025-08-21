package main

import (
	context "context"
	"errors"
	"komodo-address-api/internal/httpapi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	server := http.NewServeMux()

	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpapi.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	server.HandleFunc("/validate", httpapi.Method("POST", httpapi.HandleValidate))
	server.HandleFunc("/normalize", httpapi.Method("POST", httpapi.HandleNormalize))
	server.HandleFunc("/geocode", httpapi.Method("POST", httpapi.HandleGeocode))

	addr := ":7001" // default port

	if p := os.Getenv("PORT"); strings.TrimSpace(p) != "" {
		addr = ":" + p
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           httpapi.LoggingMiddleware(server),
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
