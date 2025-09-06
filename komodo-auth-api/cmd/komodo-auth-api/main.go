package main

import (
	"komodo-auth-api/internal/httpapi/handlers"
	"komodo-auth-api/internal/httpapi/middleware"
	"komodo-auth-api/internal/logger"
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

	// Initialize logger
	logger.InitLogger()

	// Initialize router
	rtr := chi.NewRouter()

	// Initialize middleware
	rtr.Use(middleware.ObscureDataMiddleware)
	rtr.Use(middleware.RequestValidationMiddleware)
	rtr.Use(middleware.ResponseJSONMiddleware)
	rtr.Use(middleware.InitMoxtoxMiddleware(
		"test/mocks/config/moxtox.json",
		"test/mocks/data",
		os.Getenv("USE_MOCKS") == "true",
	))

	// Initialize HTTP client
	middleware.InitHttpClient()

	// Initialize routes
	rtr.Get("/health", handlers.HealthHandler)
	rtr.Get("/.well-known/jwks.json", handlers.JWKSHandler)

	rtr.Route(("/v" + os.Getenv("API_VERSION")), func(ver chi.Router) {
		ver.Route("/auth", func(auth chi.Router) {
			auth.Post("/login", handlers.LoginHandler)
			auth.Post("/logout", handlers.LogoutHandler)
			auth.Post("/mfa/disable", handlers.MFADisableHandler)
			auth.Post("/mfa/enable", handlers.MFAEnableHandler)
			auth.Post("/mfa/setup", handlers.MFASetupHandler)
			auth.Post("/mfa/verify", handlers.MFAVerifyHandler)
			auth.Post("/passkey/start", handlers.PasskeyStartHandler)
			auth.Post("/passkey/verify", handlers.PasskeyVerifyHandler)
			auth.Post("/token", handlers.TokenCreateHandler)
			auth.Delete("/token", handlers.TokenDeleteHandler)
			auth.Post("/token/refresh", handlers.TokenRefreshHandler)
		})
	})

	// Start server
	http.ListenAndServe(":7001", rtr)
}