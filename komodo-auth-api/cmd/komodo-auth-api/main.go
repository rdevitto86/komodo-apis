package main

import (
	"komodo-auth-api/internal/httpapi/handlers"
	"komodo-auth-api/internal/httpapi/middleware"
	"komodo-auth-api/internal/thirdparty/aws"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	env := os.Getenv("API_ENV")

	switch env {
		case "dev":
			break
		case "staging", "prod":
			aws.InitSecretsClient()
			aws.InitElasticacheClient()
		default:
			log.Fatalf("Environment variable API_ENV is not set or invalid")
	}
	log.Printf("Starting Komodo Auth API in %s environment", env)

	// Init HTTP client
	middleware.InitHttpClient()

	// Initialize router
	rtr := chi.NewRouter()

	// Initialize middleware
	rtr.Use(middleware.RequestValidationMiddleware)
	rtr.Use(middleware.ObscurePIIMiddleware)
	rtr.Use(middleware.SecureLoggerMiddleware)
	rtr.Use(middleware.ResponseJSONMiddleware)

	if os.Getenv("USE_MOCKS") == "true" {
    rtr.Use(middleware.InitMoxtoxMiddleware(
      "test/mocks/config/ignored_routes.json",
      "test/mocks/config/request_mapping.json",
      true,
    ))
	}

	// Initialize routes
	rtr.Route(("/" + os.Getenv("API_VERSION")), func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler)
		r.Post("/login", handlers.LoginHandler)
		r.Post("/logout", handlers.LogoutHandler)
		// r.Post("/mfa-disable", handlers.MFADisableHandler)
		// r.Post("/mfa-enable", handlers.MFAEnableHandler)
		// r.Post("/mfa-setup", handlers.MFASetupHandler)
		// r.Post("/mfa-verify", handlers.MFAVerifyHandler)
		r.Post("/session-delete", handlers.SessionDeleteHandler)
		r.Post("/session-create", handlers.SessionCreateHandler)
		// r.Post("/token-refresh", handlers.TokenRefreshHandler)
		// r.Post("/token-revoke", handlers.TokenRevokeHandler)
		// r.Post("/token-verify", handlers.TokenVerifyHandler)
	})

	// Start server
	http.ListenAndServe(":7001", rtr)
}