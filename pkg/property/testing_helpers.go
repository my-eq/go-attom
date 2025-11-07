package property

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
)

// MockHTTPClient is used to mock HTTP requests for endpoint tests.
type mockHTTPClient struct {
	t              *testing.T
	expectedMethod string
	expectedPath   string
	expectedQuery  url.Values
	responseBody   string
	statusCode     int
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.expectedMethod != "" && req.Method != m.expectedMethod {
		m.t.Fatalf("expected method %s, got %s", m.expectedMethod, req.Method)
	}
	if m.expectedPath != "" && req.URL.Path != m.expectedPath {
		m.t.Fatalf("expected path %s, got %s", m.expectedPath, req.URL.Path)
	}
	if m.expectedQuery != nil {
		if diff := diffQuery(m.expectedQuery, req.URL.Query()); diff != "" {
			m.t.Fatalf("query mismatch: %s", diff)
		}
	}
	code := m.statusCode
	if code == 0 {
		code = http.StatusOK
	}
	body := io.NopCloser(strings.NewReader(m.responseBody))
	return &http.Response{StatusCode: code, Body: body, Header: make(http.Header)}, nil
}

// diffQuery compares two url.Values and returns a string describing the difference, or "" if equal.
func diffQuery(expected, actual url.Values) string {
	if len(expected) != len(actual) {
		return "length mismatch"
	}
	for k, v := range expected {
		av, ok := actual[k]
		if !ok {
			return "missing key: " + k
		}
		if strings.Join(v, ",") != strings.Join(av, ",") {
			return "value mismatch for key " + k
		}
	}
	return ""
}

// TestCase represents a common test case structure for service endpoint tests
//
//nolint:govet // fieldalignment optimization not critical for test structs
type TestCase struct {
	expectError           bool
	name                  string
	expectedPath          string
	responseBody          string
	expectedErrorContains string
	call                  func(context.Context, *Service) (interface{}, error)
	expectedQuery         url.Values
	statusCode            int
}

// RunServiceTest executes a service test case with proper error handling and mock setup.
func runServiceTest(ctx context.Context, t *testing.T, tt TestCase) {
	t.Run(tt.name, func(t *testing.T) {
		if tt.expectError {
			// For error cases, set up mock client if status code is specified (HTTP errors)
			// or use nil client for validation errors
			var c *client.Client
			if tt.statusCode != 0 {
				mockClient := &mockHTTPClient{
					t:              t,
					expectedMethod: http.MethodGet,
					expectedPath:   tt.expectedPath,
					expectedQuery:  tt.expectedQuery,
					responseBody:   tt.responseBody,
					statusCode:     tt.statusCode,
				}
				c = client.New("test-key", mockClient, client.WithBaseURL("https://example.com/"))
			} else {
				c = client.New("test-key", nil, client.WithBaseURL("https://example.com/"))
			}
			svc := NewService(c)
			_, err := tt.call(ctx, svc)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.expectedErrorContains)
			}
			if !strings.Contains(err.Error(), tt.expectedErrorContains) {
				t.Fatalf("expected error containing %q, got %q", tt.expectedErrorContains, err.Error())
			}
		} else {
			mockClient := &mockHTTPClient{
				t:              t,
				expectedMethod: http.MethodGet,
				expectedPath:   tt.expectedPath,
				expectedQuery:  tt.expectedQuery,
				responseBody:   tt.responseBody,
				statusCode:     tt.statusCode,
			}
			c := client.New("test-key", mockClient, client.WithBaseURL("https://example.com/"))
			svc := NewService(c)
			_, err := tt.call(ctx, svc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
	})
}

// runEndpointTests runs a collection of endpoint tests with common setup and teardown.
func runEndpointTests(t *testing.T, testName string, tests []TestCase) {
	t.Run(testName, func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		for _, tt := range tests {
			runServiceTest(ctx, t, tt)
		}
	})
}
