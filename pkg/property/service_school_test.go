package property

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
)

func TestSchoolEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		call                  func(context.Context, *Service) (interface{}, error)
		expectedQuery         url.Values
		name                  string
		expectedPath          string
		responseBody          string
		expectError           bool
		expectedErrorContains string
	}{
		{
			name:          "SearchSchools",
			expectedPath:  "/v4/school/search",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.SearchSchools(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:                  "SearchSchools_Error_NoGeoContext",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide address or latitude/longitude",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.SearchSchools(ctx)
			},
		},
		{
			name:          "GetSchoolProfile",
			expectedPath:  "/v4/school/profile",
			expectedQuery: url.Values{"schoolId": {"200"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolProfile(ctx, "200")
			},
		},
		{
			name:          "GetSchoolDistrict",
			expectedPath:  "/v4/school/district",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"district":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDistrict(ctx, "123 Main St")
			},
		},
		{
			name:          "GetSchoolDetailWithSchools",
			expectedPath:  "/v4/school/detailwithschools",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDetailWithSchools(ctx, "123 Main St")
			},
		},
		// --- NEW ENDPOINT TESTS ---
		{
			name:          "GetSchoolSnapshot",
			expectedPath:  "/v4/school/snapshot",
			expectedQuery: url.Values{"latitude": {"40.0"}, "longitude": {"-75.0"}, "radius": {"10"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolSnapshot(ctx, "40.0", "-75.0", "10", "", nil)
			},
		},
		{
			name:          "GetSchoolDetail",
			expectedPath:  "/v4/school/detail",
			expectedQuery: url.Values{"id": {"200"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDetail(ctx, "200")
			},
		},
		{
			name:          "GetSchoolDistrictDetail",
			expectedPath:  "/v4/school/districtdetail",
			expectedQuery: url.Values{"id": {"300"}},
			responseBody:  `{"status":{},"district":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDistrictDetail(ctx, "300")
			},
		},
	}

	for _, tt := range tests {
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
			}
		})
	}
}
