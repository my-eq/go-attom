package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const testContentTypeJSON = "application/json"

// mockHTTPClient implements HTTPClient for testing.
type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return m.resp, m.err
}

func TestNew_DefaultsToStdClient(t *testing.T) {
	c := New("test-key", nil)
	if c.httpClient == nil {
		t.Error("expected default httpClient to be set")
	}
	if c.apiKey != "test-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "test-key")
	}
	if c.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, DefaultBaseURL)
	}
}

// headerCheckHTTPClient allows inspection of request headers in tests.
type headerCheckHTTPClient struct {
	t       *testing.T
	wantKey string
	called  bool
}

func (m *headerCheckHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.called = true
	if req.Header.Get("apikey") != m.wantKey {
		m.t.Errorf("apikey header = %q, want %q", req.Header.Get("apikey"), m.wantKey)
	}
	return &http.Response{StatusCode: 200}, nil
}

func TestDoRequest_APIKeyInjection(t *testing.T) {
	mock := &headerCheckHTTPClient{t: t, wantKey: "my-key"}
	c := New("my-key", mock)
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = c.DoRequest(req)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !mock.called {
		t.Error("expected mock Do to be called")
	}
}

func TestDoRequest_Errors(t *testing.T) {
	c := New("", &mockHTTPClient{})
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = c.DoRequest(req)
	if !errors.Is(err, ErrInvalidAPIKey) {
		t.Errorf("expected ErrInvalidAPIKey, got %v", err)
	}

	c = New("key", &mockHTTPClient{})
	_, err = c.DoRequest(nil)
	if err == nil || !strings.Contains(err.Error(), "request cannot be nil") {
		t.Errorf("expected error for nil request, got %v", err)
	}

	c = New("key", &mockHTTPClient{err: errors.New("fail")})
	req, err = http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = c.DoRequest(req)
	if err == nil || !strings.Contains(err.Error(), "failed to execute request") {
		t.Errorf("expected wrapped error, got %v", err)
	}
}

func TestWithBaseURL_Option(t *testing.T) {
	custom := "https://custom.example.com/api"
	c := New("key", nil, WithBaseURL(custom))
	expected := "https://custom.example.com/api/" // normalized with trailing slash
	if c.baseURL != expected {
		t.Errorf("baseURL = %q, want %q", c.baseURL, expected)
	}

	// Empty string should keep default
	c2 := New("key", nil, WithBaseURL(""))
	if c2.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want default %q", c2.baseURL, DefaultBaseURL)
	}

	// Already has trailing slash: ensure no double slash
	customTrailing := "https://other.example.com/root/"
	c3 := New("key", nil, WithBaseURL(customTrailing))
	expected3 := "https://other.example.com/root/"
	if c3.baseURL != expected3 {
		t.Errorf("baseURL = %q, want %q", c3.baseURL, expected3)
	}

	// Nil option should not alter defaults or panic
	c4 := New("key", nil, nil)
	if c4.baseURL != DefaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c4.baseURL, DefaultBaseURL)
	}
}

func TestNewRequest(t *testing.T) {
	c := New("key", nil)
	ctx := context.Background()
	params := url.Values{}
	params.Set("foo", "bar")

	req, err := c.NewRequest(ctx, http.MethodGet, "propertyapi/v1.0.0/property/detail", params, nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	expectedURL := "https://api.gateway.attomdata.com/propertyapi/v1.0.0/property/detail?foo=bar"
	if req.URL.String() != expectedURL {
		t.Errorf("URL = %q, want %q", req.URL.String(), expectedURL)
	}

	if accept := req.Header.Get("Accept"); accept != testContentTypeJSON {
		t.Errorf("Accept header = %q, want %s", accept, testContentTypeJSON)
	}
}

func TestNewRequestErrors(t *testing.T) {
	c := New("key", nil)
	ctx := context.Background()

	if _, err := c.NewRequest(ctx, "", "endpoint", nil, nil); err == nil || !strings.Contains(err.Error(), "method cannot be empty") {
		t.Errorf("expected method error, got %v", err)
	}

	if _, err := c.NewRequest(ctx, http.MethodGet, "", nil, nil); err == nil || !strings.Contains(err.Error(), "endpoint cannot be empty") {
		t.Errorf("expected endpoint error, got %v", err)
	}
}

func TestNewRequest_InvalidBaseURL(t *testing.T) {
	c := New("key", nil)
	c.baseURL = "://invalid-url"
	ctx := context.Background()

	_, err := c.NewRequest(ctx, http.MethodGet, "endpoint", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "invalid base URL") {
		t.Errorf("expected invalid base URL error, got %v", err)
	}
}

func TestNewRequest_WithBody(t *testing.T) {
	c := New("key", nil)
	ctx := context.Background()
	body := strings.NewReader(`{"test":"data"}`)

	req, err := c.NewRequest(ctx, http.MethodPost, "endpoint", nil, body)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	if ct := req.Header.Get("Content-Type"); ct != testContentTypeJSON {
		t.Errorf("Content-Type header = %q, want %s", ct, testContentTypeJSON)
	}
}

func TestNewRequest_PreservesExistingHeaders(t *testing.T) {
	c := New("key", nil)
	ctx := context.Background()
	body := strings.NewReader(`{"test":"data"}`)

	// Create request that will have headers set
	req, err := c.NewRequest(ctx, http.MethodPost, "endpoint", nil, body)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	// Verify Accept header is set when not present
	if accept := req.Header.Get("Accept"); accept != testContentTypeJSON {
		t.Errorf("Accept header = %q, want %s", accept, testContentTypeJSON)
	}

	// Verify Content-Type is set for body
	if ct := req.Header.Get("Content-Type"); ct != testContentTypeJSON {
		t.Errorf("Content-Type header = %q, want %s", ct, testContentTypeJSON)
	}
}

func TestNewRequest_EndpointTrimming(t *testing.T) {
	c := New("key", nil)
	ctx := context.Background()

	tests := []struct {
		name     string
		endpoint string
		wantPath string
	}{
		{
			name:     "no leading slash",
			endpoint: "property/detail",
			wantPath: "/property/detail",
		},
		{
			name:     "with leading slash",
			endpoint: "/property/detail",
			wantPath: "/property/detail",
		},
		{
			name:     "with whitespace",
			endpoint: "  /property/detail  ",
			wantPath: "/property/detail",
		},
		{
			name:     "multiple leading slashes",
			endpoint: "///property/detail",
			wantPath: "/property/detail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := c.NewRequest(ctx, http.MethodGet, tt.endpoint, nil, nil)
			if err != nil {
				t.Fatalf("NewRequest returned error: %v", err)
			}
			if req.URL.Path != tt.wantPath {
				t.Errorf("URL path = %q, want %q", req.URL.Path, tt.wantPath)
			}
		})
	}
}

func TestNewRequest_QueryParameters(t *testing.T) {
	c := New("key", nil)
	ctx := context.Background()

	params := url.Values{}
	params.Set("foo", "bar")
	params.Add("multi", "value1")
	params.Add("multi", "value2")

	req, err := c.NewRequest(ctx, http.MethodGet, "endpoint", params, nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	// Verify query encoding
	if !strings.Contains(req.URL.RawQuery, "foo=bar") {
		t.Errorf("query string missing foo=bar: %q", req.URL.RawQuery)
	}
	if !strings.Contains(req.URL.RawQuery, "multi=value1") {
		t.Errorf("query string missing multi=value1: %q", req.URL.RawQuery)
	}
	if !strings.Contains(req.URL.RawQuery, "multi=value2") {
		t.Errorf("query string missing multi=value2: %q", req.URL.RawQuery)
	}
}

func TestNewRequest_NilQuery(t *testing.T) {
	c := New("key", nil)
	ctx := context.Background()

	req, err := c.NewRequest(ctx, http.MethodGet, "endpoint", nil, nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	if req.URL.RawQuery != "" {
		t.Errorf("expected empty query string, got %q", req.URL.RawQuery)
	}
}
