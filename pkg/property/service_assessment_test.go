package property

import (
	"context"
	"net/url"
	"testing"
)

func TestAssessmentEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []TestCase{
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
		runServiceTest(t, ctx, tt)
	}
}
