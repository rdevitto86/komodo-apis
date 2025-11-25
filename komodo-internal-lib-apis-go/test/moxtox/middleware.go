package moxtox

import (
	errCodes "komodo-internal-lib-apis-go/http/common/errors"
	errors "komodo-internal-lib-apis-go/http/common/errors/chi"
	logger "komodo-internal-lib-apis-go/logging/runtime"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	dir        string
	config     MoxtoxConfig
	once       sync.Once
	allowMocks = true
)

// InitMoxtoxMiddleware initializes the moxtox middleware for intercepting real HTTP requests and replacing with mocks.
func InitMoxtoxMiddleware(env string, configPath ...string) func(http.Handler) http.Handler {
	once.Do(func() {
		// Determine baseDir from configPath or default
		if len(configPath) == 0 || configPath[0] == "" {
			cwd, err := os.Getwd()
			if err != nil {
				logger.Error("[::Moxtox::] unable to determine current working directory", err)
				allowMocks = false
				return
			}
			dir = filepath.Join(cwd, "test", "moxtox")
		} else {
			dir = configPath[0]
		}

		// Load the Moxtox config
		if data, err := os.ReadFile(filepath.Join(dir, "moxtox_config.yml")); err == nil {
			if err := yaml.Unmarshal(data, &config); err != nil {
				logger.Error("[::Moxtox::] error loading moxtox config", err)
				allowMocks = false
				return
			}
			if !config.EnableMoxtox {
				logger.Info("[::Moxtox::] mocks disabled - using default behavior")
				allowMocks = false
				return
			}
			if !contains(config.AllowedEnvironments, env) {
				logger.Info("[::Moxtox::] mocks not allowed in this environment")
				allowMocks = false
				return
			}

			// build mock data store based on mode
			switch config.PerformanceMode {
				case "quick":
					config.buildHashLookupMap()
				case "dynamic":
					totalScenarios := config.countTotalScenarios()
					if totalScenarios > 10 { // threshold for switching to quick mode
						config.buildHashLookupMap()
					} else {
						config.buildSliceLookupMap()
					}
				default: // "default"
					config.buildSliceLookupMap()
			}

			logger.Info("[::Moxtox::] mocks enabled successfully")
		} else {
			logger.Error("[::Moxtox::] error loading moxtox config", err)
			allowMocks = false
		}
	})

	// ignore if mocks are disabled
	if allowMocks {
		return mockResponseHandler()
	}
	return func(next http.Handler) http.Handler { return next }
}

// mockResponseHandler returns a middleware that injects mock responses based on the LookupMap.
func mockResponseHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for ignored routes first
			if contains(config.IgnoredRoutes, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Use LookupMap for lookup
			if scenario, ok := matchesRequest(r); ok {
				if err := injectMock(w, r, scenario); err != nil {
					errors.WriteErrorResponse(w, r, http.StatusInternalServerError, errCodes.ERR_INTERNAL_SERVER, "Mock injection failed")
					return
				}
				return
			}

			// No match: return 418 Teapot error
			errors.WriteErrorResponse(w, r, http.StatusTeapot, errCodes.ERR_NOT_FOUND, "No mocks found")
		})
	}
}
