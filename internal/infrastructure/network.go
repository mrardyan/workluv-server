package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NetworkConfig holds network-related configuration
type NetworkConfig struct {
	HTTPTimeout     time.Duration
	MaxIdleConns    int
	IdleConnTimeout time.Duration
}

// DefaultNetworkConfig returns sensible defaults
func DefaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		HTTPTimeout:     30 * time.Second,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}
}

// NewHTTPClient creates a configured HTTP client
func NewHTTPClient(config *NetworkConfig) *http.Client {
	if config == nil {
		config = DefaultNetworkConfig()
	}

	transport := &http.Transport{
		MaxIdleConns:       config.MaxIdleConns,
		IdleConnTimeout:    config.IdleConnTimeout,
		DisableCompression: true,
	}

	return &http.Client{
		Timeout:   config.HTTPTimeout,
		Transport: transport,
	}
}

// NewDefaultHTTPClient creates an HTTP client with default settings
func NewDefaultHTTPClient() *http.Client {
	return NewHTTPClient(DefaultNetworkConfig())
}

// HTTPRequest represents a generic HTTP request
type HTTPRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
}

// HTTPResponse represents a generic HTTP response
type HTTPResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// DoRequest performs an HTTP request with the given client
func DoRequest(ctx context.Context, client *http.Client, req *HTTPRequest) (*HTTPResponse, error) {
	var bodyReader io.Reader

	if req.Body != nil {
		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	if req.Body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	// Set custom headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}

// Get performs a GET request
func Get(ctx context.Context, client *http.Client, url string, headers map[string]string) (*HTTPResponse, error) {
	req := &HTTPRequest{
		Method:  http.MethodGet,
		URL:     url,
		Headers: headers,
	}
	return DoRequest(ctx, client, req)
}

// Post performs a POST request
func Post(ctx context.Context, client *http.Client, url string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	req := &HTTPRequest{
		Method:  http.MethodPost,
		URL:     url,
		Body:    body,
		Headers: headers,
	}
	return DoRequest(ctx, client, req)
}

// Delete performs a DELETE request
func Delete(ctx context.Context, client *http.Client, url string, headers map[string]string) (*HTTPResponse, error) {
	req := &HTTPRequest{
		Method:  http.MethodDelete,
		URL:     url,
		Headers: headers,
	}
	return DoRequest(ctx, client, req)
}
