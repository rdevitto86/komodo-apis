package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var ignoredPaths = map[string]bool{}
var activeMocks = map[string]map[string]string{}
var dataKeys = []string{"body", "body_nested", "path", "query", "header"}
var testDataRoot string = ""

type MoxtoxConfig struct {
	IgnoredRoutes   map[string]bool `json:"ignored_routes"`
	RequestMappings map[string]map[string]struct {
		Keys []map[string]string `json:"keys"`
	} `json:"request_mappings"`
}

// TODO - pass in handler functions instead??
// TODO - make it framework agnostic
// TODO - documentation
func InitMoxtoxMiddleware(
	configPath string,
	testDataPath string,
	areMocksEnabled bool,
) func(http.Handler) http.Handler {
	if !areMocksEnabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	testDataRoot = testDataPath

	var config MoxtoxConfig
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			fmt.Println("Error loading moxtox config:", err)
		}
	}
	if len(config.RequestMappings) == 0 && len(config.IgnoredRoutes) == 0 {
		fmt.Println("Moxtox config has no data - skipping init")
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return moxtoxRequestHandler
}

// TODO - documentation
func moxtoxRequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wtr http.ResponseWriter, req *http.Request) {
		// TODO handle nested routes
		for route := range ignoredPaths {
			if strings.HasSuffix(req.URL.Path, route) {
				next.ServeHTTP(wtr, req)
				return
			}
		}

		// mockData, err := os.ReadFile(fmt.Sprintf("%s/%s.json", testDataRoot, route))
		// if err != nil {
		// 	http.Error(wtr, "Failed to read mock data", http.StatusInternalServerError)
		// 	return
		// }

		// filter := "default"

		switch req.Method {
			case http.MethodGet:
				// Handle GET specific logic
			case http.MethodPost:
				// Handle POST specific logic
			case http.MethodPut:
				// Handle PUT specific logic
			case http.MethodDelete:
				// Handle DELETE specific logic
			case http.MethodPatch:
				// Handle PATCH specific logic
		}

		// Serve the mock response
		wtr.WriteHeader(http.StatusOK)
		// wtr.Write(mockResponse)
	})
}
