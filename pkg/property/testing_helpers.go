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

// mockHTTPClient is used to mock HTTP requests for endpoint tests
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
type TestCase struct {
	call                  func(context.Context, *Service) (interface{}, error)
	expectedQuery         url.Values
	name                  string
	expectedPath          string
	responseBody          string
	expectError           bool
	expectedErrorContains string
}

// runServiceTest executes a service test case with proper error handling and mock setup
func runServiceTest(t *testing.T, ctx context.Context, tt TestCase) {
	t.Run(tt.name, func(t *testing.T) {
		if tt.expectError {
			// For error cases, we don't set up the mock client since the error occurs before the HTTP call
			c := client.New("test-key", nil, client.WithBaseURL("https://example.com/"))
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
