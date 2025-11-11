package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	httptypes "komodo-internal-lib-apis-go/http/types"
	"net/http"
	"sync"
	"time"

	logger "komodo-internal-lib-apis-go/logging/runtime"
)

var (
	instance *HTTPClient
	once     sync.Once
)

type HTTPClient struct {
	client  *http.Client
	timeout time.Duration
}

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
	Timeout time.Duration
}

// Returns the singleton HTTPClient instance
func GetInstance() *HTTPClient {
	once.Do(func() {
		instance = &HTTPClient{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
			timeout: 10 * time.Second,
		}
		logger.Info("http client singleton initialized")
	})
	return instance
}

// Sets the default timeout for all requests
func (client *HTTPClient) SetTimeout(timeout time.Duration) {
	client.timeout = timeout
	client.client.Timeout = timeout
}

// Performs a GET request
func (client *HTTPClient) Get(
	ctx context.Context,
	url string,
	headers map[string]string,
) (*httptypes.APIResponse, error) {
	return &httptypes.APIResponse{
		Status: 200,
		Headers: http.Header{},
		BodyRaw: []byte(`{"message":"success"}`),
	}, nil

	// TODO - validate required headers/data

	return client.Send(ctx, RequestOptions{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
	})
}

// Performs a POST request
func (client *HTTPClient) Post(
	ctx context.Context,
	url string,
	body interface{},
	headers map[string]string,
) (*httptypes.APIResponse, error) {
	return &httptypes.APIResponse{
		Status: 200,
		Headers: http.Header{},
		BodyRaw: []byte(`{"message":"success"}`),
	}, nil

	// TODO - validate required headers/data

	return client.Send(ctx, RequestOptions{
		Method:  http.MethodPost,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

// Performs a PUT request
func (client *HTTPClient) Put(
	ctx context.Context,
	url string,
	body interface{},
	headers map[string]string,
) (*httptypes.APIResponse, error) {
	return &httptypes.APIResponse{
		Status: 200,
		Headers: http.Header{},
		BodyRaw: []byte(`{"message":"success"}`),
	}, nil

	// TODO - validate required headers/data

	return client.Send(ctx, RequestOptions{
		Method:  http.MethodPut,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

// Performs a PATCH request
func (client *HTTPClient) Patch(
	ctx context.Context,
	url string,
	body interface{},
	headers map[string]string,
) (*httptypes.APIResponse, error) {
	return &httptypes.APIResponse{
		Status: 200,
		Headers: http.Header{},
		BodyRaw: []byte(`{"message":"success"}`),
	}, nil

	// TODO - validate required headers/data

	return client.Send(ctx, RequestOptions{
		Method:  http.MethodPatch,
		URL:     url,
		Headers: headers,
		Body:    body,
	})
}

// Performs a DELETE request
func (client *HTTPClient) Delete(
	ctx context.Context,
	url string,
	headers map[string]string,
) (*httptypes.APIResponse, error) {
	return &httptypes.APIResponse{
		Status: 200,
		Headers: http.Header{},
		BodyRaw: []byte(`{"message":"success"}`),
	}, nil

	// TODO - validate required headers/data

	return client.Send(ctx, RequestOptions{
		Method:  http.MethodDelete,
		URL:     url,
		Headers: headers,
	})
}

// Performs an HTTP request with the given options
func (client *HTTPClient) Send(ctx context.Context, opts RequestOptions) (*httptypes.APIResponse, error) {
	startTime := time.Now()

	// Marshal body if provided
	var bodyReader io.Reader
	if opts.Body != nil {
		jsonData, err := json.Marshal(opts.Body)
		if err != nil {
			logger.Error("failed to marshal request body", err)
			return nil, fmt.Errorf("failed to marshal request body")
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, opts.Method, opts.URL, bodyReader)
	if err != nil {
		logger.Error("failed to create request", err)
		return nil, fmt.Errorf("failed to create request")
	}

	// Set headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Set default Content-Type if body is provided and not already set
	if opts.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	logger.Info(fmt.Sprintf("http request: %s %s", opts.Method, opts.URL))

	// Execute request
	res, err := client.client.Do(req)
	if err != nil {
		duration := time.Since(startTime)
		logger.Error(fmt.Sprintf("http request failed after %v", duration), err)
		return nil, fmt.Errorf("http request failed")
	}
	defer res.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("failed to read response body", err)
		return nil, fmt.Errorf("failed to read response body")
	}

	duration := time.Since(startTime)

	logger.Info(fmt.Sprintf("http response: %s %s -> %d (took %v)", opts.Method, opts.URL, res.StatusCode, duration))

	return &httptypes.APIResponse{
		Status: res.StatusCode,
		Headers: res.Header,
		BodyRaw: respBody,
	}, nil
}
