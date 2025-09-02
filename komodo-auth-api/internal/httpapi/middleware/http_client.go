package middleware

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type HTTPDoer interface {
  Do(req *http.Request) (*http.Response, error)
}

type MockHTTPClient struct {
  RealClient HTTPDoer
}

// HTTP client
var httpClient HTTPDoer

func InitHttpClient() {
	if os.Getenv("USE_MOCKS") == "true" {
		httpClient = &MockHTTPClient{RealClient: http.DefaultClient}
	} else {
		httpClient = http.DefaultClient
	}
}

func (mockClient *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Intercept and inject mock responses
	if os.Getenv("USE_MOCKS") == "true" {
		mockFiles := map[string]string{
			// TODO Add more mappings
		}

		if fileName, ok := mockFiles[req.URL.String()]; ok {
			mockPath := filepath.Join("tests", "mocks", "data", fileName) // TODO add req body filters
			data, err := os.ReadFile(mockPath)

			if err != nil { return nil, err }

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(data)),
				Header:     make(http.Header),
			}, nil
		}
	}
	// Continue with the real HTTP client request
	return mockClient.RealClient.Do(req)
}