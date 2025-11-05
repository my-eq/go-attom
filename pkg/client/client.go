// Package client provides the base ATTOM API client implementation.
//
// The Client is designed to be mockable and injectable, accepting an API key and an HTTP client interface.
// It does not perform any logging and returns errors with context.
package client

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// HTTPClient defines the minimal interface for making HTTP requests.
// It is satisfied by *http.Client and can be mocked in tests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client provides methods for interacting with the ATTOM Data API.
// It handles authentication and request execution.
type Client struct {
	httpClient HTTPClient
	apiKey     string
	baseURL    string
}

// Option represents a functional configuration option for Client.
type Option func(*Client)

// DefaultBaseURL is the default root ATTOM API URL used when no override is supplied.
const DefaultBaseURL = "https://api.gateway.attomdata.com/"

// WithBaseURL sets a custom base URL for the API client. Trailing slashes are normalized.
// If an empty string is provided, the option is ignored and DefaultBaseURL remains.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL == "" {
			return
		}
		normalized := strings.TrimRight(baseURL, "/") + "/"
		c.baseURL = normalized
	}
}

// New creates a new ATTOM API client.
//
// If httpClient is nil, a default *http.Client with 30s timeout is used.
// The apiKey must be a valid ATTOM API key.
func New(apiKey string, httpClient HTTPClient, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	c := &Client{
		httpClient: httpClient,
		apiKey:     apiKey,
		baseURL:    DefaultBaseURL,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}
	return c
}

// ErrInvalidAPIKey is returned when the API key is missing or invalid.
var ErrInvalidAPIKey = errors.New("invalid or missing API key")

// DoRequest executes an HTTP request with the API key injected.
//
// The req must be non-nil and will have the API key added as a header.
// Returns an error with context if the request fails.
func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if c.apiKey == "" {
		return nil, ErrInvalidAPIKey
	}
	req.Header.Set("apikey", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	return resp, nil
}
