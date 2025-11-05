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

// mockHTTPClient and diffQuery are defined here for test isolation. If already present via import, remove duplicate.

func TestPropertyEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		call          func(context.Context, *Service) (interface{}, error)
		expectedQuery url.Values
		name          string
		expectedPath  string
		responseBody  string
	}{
		{
			name:          "GetPropertyDetail",
			expectedPath:  "/propertyapi/v1.0.0/property/detail",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyDetail(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:          "GetPropertyAddress",
			expectedPath:  "/propertyapi/v1.0.0/property/address",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyAddress(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetPropertySnapshot",
			expectedPath:  "/propertyapi/v1.0.0/property/snapshot",
			expectedQuery: url.Values{"postalCode": {"62701"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertySnapshot(ctx, WithPostalCode("62701"))
			},
		},
		{
			name:          "GetBasicProfile",
			expectedPath:  "/propertyapi/v1.0.0/property/basicprofile",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetBasicProfile(ctx, "123 Main St")
			},
		},
		{
			name:          "GetExpandedProfile",
			expectedPath:  "/propertyapi/v1.0.0/property/expandedprofile",
			expectedQuery: url.Values{"geoIdV4": {"geo-1"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetExpandedProfile(ctx, WithGeoIDV4("geo-1"))
			},
		},
		{
			name:          "GetDetailWithSchools",
			expectedPath:  "/propertyapi/v1.0.0/property/detailwithschools",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailWithSchools(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailMortgage",
			expectedPath:  "/propertyapi/v1.0.0/property/detailmortgage",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"mortgage":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailMortgage(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailOwner",
			expectedPath:  "/propertyapi/v1.0.0/property/detailowner",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"owner":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailOwner(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailMortgageOwner",
			expectedPath:  "/propertyapi/v1.0.0/property/detailmortgageowner",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"mortgage":[{}],"owner":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailMortgageOwner(ctx, "123 Main St")
			},
		},
		{
			name:          "GetBuildingPermits",
			expectedPath:  "/propertyapi/v1.0.0/property/buildingpermits",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"buildingPermit":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetBuildingPermits(ctx, "123 Main St")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockHTTPClient{
				t:              t,
				expectedMethod: http.MethodGet,
				expectedPath:   tt.expectedPath,
				expectedQuery:  tt.expectedQuery,
				responseBody:   tt.responseBody,
			}
			c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
			svc := NewService(c)
			_, err := tt.call(ctx, svc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
