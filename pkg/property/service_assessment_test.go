package property

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
)

func TestAssessmentEndpoints(t *testing.T) {
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
			name:          "GetAssessmentDetail",
			expectedPath:  "/v4/property/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetAssessmentDetail_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentDetail(ctx)
			},
		},
		{
			name:          "GetAssessmentSnapshot",
			expectedPath:  "/v4/property/snapshot",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentSnapshot(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetAssessmentSnapshot_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentSnapshot(ctx)
			},
		},
		{
			name:          "GetAssessmentHistory",
			expectedPath:  "/v4/property/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentHistory(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetAssessmentHistory_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentHistory(ctx)
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
}
