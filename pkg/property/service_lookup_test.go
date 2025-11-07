package property

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
)

func TestLookupEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		expectError           bool
		statusCode            int
		name                  string
		expectedPath          string
		expectedQuery         url.Values
		responseBody          string
		expectedErrorContains string
		call                  func(context.Context, *Service) (interface{}, error)
	}{
		{
			name:          "GetStateLookup",
			expectedPath:  "/v4/area/state/lookup",
			expectedQuery: url.Values{},
			responseBody:  `{"status":{},"state":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetStateLookup(ctx)
			},
		},
		{
			name:                  "GetStateLookup_Error_HTTP",
			expectedPath:          "/v4/area/state/lookup",
			expectedQuery:         url.Values{},
			responseBody:          `{"status":{"version":"1.0","transactionId":"test"},"message":"Internal Server Error"}`,
			statusCode:            http.StatusInternalServerError,
			expectError:           true,
			expectedErrorContains: "Internal Server Error",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetStateLookup(ctx)
			},
		},
		{
			name:          "GetCBSALookup",
			expectedPath:  "/v4/area/cbsa/lookup",
			expectedQuery: url.Values{"StateId": {"CA"}},
			responseBody:  `{"status":{},"cbsa":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetCBSALookup(ctx, "CA")
			},
		},
		{
			name:          "GetCountyLookup",
			expectedPath:  "/v4/area/county/lookup",
			expectedQuery: url.Values{"StateId": {"CA"}},
			responseBody:  `{"status":{},"county":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetCountyLookup(ctx, "CA")
			},
		},
		{
			name:          "GetEnumerationsDetail",
			expectedPath:  "/v4/enumerations/detail",
			expectedQuery: url.Values{},
			responseBody:  `{"status":{},"enumerations":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetEnumerationsDetail(ctx)
			},
		},
		{
			name:                  "GetEnumerationsDetail_Error_HTTP",
			expectedPath:          "/v4/enumerations/detail",
			expectedQuery:         url.Values{},
			responseBody:          `{"status":{"version":"1.0","transactionId":"test"},"message":"Bad Request"}`,
			statusCode:            http.StatusBadRequest,
			expectError:           true,
			expectedErrorContains: "Bad Request",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetEnumerationsDetail(ctx)
			},
		},
		{
			name:          "GetBoundaryDetail",
			expectedPath:  "/v4/area/boundary/detail",
			expectedQuery: url.Values{"geoIdV4": {"geo-123"}},
			responseBody:  `{"status":{},"boundary":{}}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetBoundaryDetail(ctx, "geo-123")
			},
		},
		{
			name:          "GetHierarchyLookup",
			expectedPath:  "/v4/area/hierarchy/lookup",
			expectedQuery: url.Values{"WKTString": {"POINT(-122.4194 37.7749)"}},
			responseBody:  `{"status":{},"hierarchy":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetHierarchyLookup(ctx, "POINT(-122.4194 37.7749)")
			},
		},
		{
			name:          "GetGeoIDLookup",
			expectedPath:  "/v4/area/geoid/lookup/",
			expectedQuery: url.Values{"geoIdV4": {"geo-123"}},
			responseBody:  `{"status":{},"geoid":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetGeoIDLookup(ctx, "geo-123")
			},
		},
		{
			name:          "GetGeoIDLegacyLookup",
			expectedPath:  "/v4/area/geoid/legacyLookup/",
			expectedQuery: url.Values{"geoIdV4": {"geo-123"}},
			responseBody:  `{"status":{},"legacyGeoid":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetGeoIDLegacyLookup(ctx, "geo-123")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				t:             t,
				expectedPath:  tt.expectedPath,
				expectedQuery: tt.expectedQuery,
				responseBody:  tt.responseBody,
				statusCode:    tt.statusCode,
			}
			c := client.New("test-key", mockClient, client.WithBaseURL("https://example.com/"))
			svc := NewService(c)

			_, err := tt.call(ctx, svc)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectedErrorContains)
				} else if !strings.Contains(err.Error(), tt.expectedErrorContains) {
					t.Errorf("expected error containing %q, got %q", tt.expectedErrorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
