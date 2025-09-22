package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"komodo-auth-api/internal/config"
)

// HTTPDoer is satisfied by *http.Client and test doubles
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a small abstraction over HTTP that can optionally return
// deterministic mock responses (from disk) when USE_MOCKS=true.
//
// Usage:
//   // package-level default initialized via InitHttpClient()
//   resp, err := httpclient.DefaultClient.Do(req)
//
// For tests you can create a custom client and call RegisterMock(...).
type Client struct {
	RealClient HTTPDoer
	UseMocks   bool
	MockDir    string

	mu       sync.RWMutex
	mappings map[string]string // key -> filename
}

// DefaultClient is initialized by InitHttpClient and safe for concurrent use.
var DefaultClient *Client

// InitHttpClient initializes the package-level DefaultClient using env/config.
func InitHttpClient(real HTTPDoer) {
	if real == nil {
		real = http.DefaultClient
	}

	DefaultClient = &Client{
		RealClient: real,
		UseMocks:   config.GetConfigValue("USE_MOCKS") == "true",
		MockDir:    config.GetConfigValue("MOCKS_DIR"),
		mappings:   make(map[string]string),
	}

	// Attach an interceptor RoundTripper to the underlying http.Client so
	// we can short-circuit requests with local mock files when enabled.
	transport := &interceptTransport{parent: DefaultClient, next: http.DefaultTransport}
	DefaultClient.RealClient = &http.Client{Transport: transport}
}

// interceptTransport is a RoundTripper that consults the parent Client for
// mock mappings when UseMocks==true. If a mapping is found it returns the
// file contents as a 200 response; otherwise it delegates to next.
type interceptTransport struct {
	parent *Client
	next   http.RoundTripper
}

func (t *interceptTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t == nil {
		return http.DefaultTransport.RoundTrip(req)
	}

	// If parent indicates mocks are enabled, attempt to find a mapping.
	if t.parent != nil && t.parent.UseMocks {
		key := req.Method + " " + req.URL.Scheme + "://" + req.URL.Host + req.URL.Path

		t.parent.mu.RLock()
		fileName, ok := t.parent.mappings[key]
		t.parent.mu.RUnlock()

		if ok {
			mockPath := filepath.Join(t.parent.MockDir, fileName)
			data, err := os.ReadFile(mockPath)
			if err == nil {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer(data)),
					Header:     make(http.Header),
					Request:    req,
				}, nil
			}
			// If reading mock failed, fall through to real transport.
		}
	}

	if t.next == nil {
		t.next = http.DefaultTransport
	}
	return t.next.RoundTrip(req)
}

// RegisterMock registers a file to return for a particular request key.
// The key should be unique per service and can be any string you agree on
// across tests and mapping setup. A convenient key is `METHOD HOST PATH`.
func (client *Client) RegisterMock(key, filename string) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.mappings[key] = filename
}

// UnregisterMock removes a mapping.
func (client *Client) UnregisterMock(key string) {
	client.mu.Lock()
	defer client.mu.Unlock()
	delete(client.mappings, key)
}

// Do implements the HTTPDoer interface. If mocks are enabled and a mapping
// exists for the request, the mapped file is returned as the response body.
func (client *Client) Do(req *http.Request) (*http.Response, error) {
	if client == nil {
		return nil, http.ErrServerClosed
	}

	if client.UseMocks {
		// Build a key: METHOD + space + URL (without query) to make mappings easier
		key := req.Method + " " + req.URL.Scheme + "://" + req.URL.Host + req.URL.Path

		client.mu.RLock()
		fileName, ok := client.mappings[key]
		client.mu.RUnlock()

		if ok {
			mockPath := filepath.Join(client.MockDir, fileName)
			data, err := os.ReadFile(mockPath)
			if err != nil {
				return nil, err
			}
			// Return a 200 with the mock body. Tests can override headers if needed.
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(data)),
				Header:     make(http.Header),
			}, nil
		}
	}

	// Fallback to the real client
	return client.RealClient.Do(req)
}

// Convenience wrapper to call the package DefaultClient
func Do(req *http.Request) (*http.Response, error) {
	if DefaultClient == nil {
		InitHttpClient(nil)
	}
	return DefaultClient.Do(req)
}

// RequestKey builds the same key used by Do() for looking up mocks.
// Use this when registering mocks for a given request.
func RequestKey(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	return req.Method + " " + req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
}

// RegisterMockForURL registers a mock file for the given full URL and method.
// fullURL should include scheme and host (e.g. "https://user.api.example.com/v1/me").
func (client *Client) RegisterMockForURL(method, fullURL, filename string) {
	// build a key aligned with RequestKey by parsing the URL
	// simplest approach: construct key from method + space + fullURL's scheme://host+path
	key := method + " " + fullURL
	client.RegisterMock(key, filename)
}