package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var enableMocks = false
var ignoredMockPaths = map[string]bool{}
var requestToMockMappings = map[string]map[string]string{}

/*
	Mapping of request methods to their corresponding data keys
*/
var dataKeys = []string{"body", "body_nested", "path", "query", "header", "cookie", "default"}

// TODO - pass in handler functions instead??
// TODO - make it framework agnostic
// TODO - documentation
func InitMoxtoxMiddleware(
	ignoredRoutesPath string,
	requestMappingPath string,
	areMocksEnabled bool,
) func(http.Handler) http.Handler {
	if !areMocksEnabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	enableMocks = true

	if data, err := os.ReadFile(ignoredRoutesPath); err == nil {
		if err := json.Unmarshal(data, &ignoredMockPaths); err != nil {
			fmt.Println("Error loading ignored routes:", err)
		}
	}
	if data, err := os.ReadFile(requestMappingPath); err == nil {
		if err := json.Unmarshal(data, &requestToMockMappings); err != nil {
			fmt.Println("Error loading request mappings:", err)
		}
	}
	if len(ignoredMockPaths) == 0 || len(requestToMockMappings) == 0 {
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
		path := req.URL.Path
		idx := strings.LastIndex(path, "/")
		last := path
		if idx != -1 { last = path[idx+1:] }

		if !enableMocks || ignoredMockPaths[last] {
			next.ServeHTTP(wtr, req)
			return
		}

		data, err := os.ReadFile(fmt.Sprintf("test/mocks/data/%s.json", last))

		if err != nil {	
			http.Error(wtr, "Failed to read mock data", http.StatusInternalServerError)
			return
		}

		var jsonData map[string]interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			http.Error(wtr, "Failed to parse mock JSON", http.StatusInternalServerError)
			return
		}

		key, err := parseJSONKey(req)
		if err != nil {
			http.Error(wtr, "Failed to parse key for mock JSON", http.StatusBadRequest)
			return
		}

		if obj, ok := jsonData[key]; ok {
			res, err := json.Marshal(obj)
			if err != nil {
				http.Error(wtr, "Failed to marshal mock object", http.StatusInternalServerError)
				return
			}
			wtr.WriteHeader(http.StatusOK)
			wtr.Write(res)
			return
		}

		http.Error(wtr, "Mock object not found", http.StatusNotFound)
	})
}

// TODO - documentation
func parseJSONKey(req *http.Request) (string, error) {
	key := "default"

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

	// TODO use requestToMockMappings and grab the specific request prop needed in each JSON
	return key, nil
}
