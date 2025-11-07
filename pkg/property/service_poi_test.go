package property

import (
	"context"
	"net/url"
	"testing"
)

func TestPOIEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []TestCase{
		{
			name:                  "GetPOI",
			expectedPath:          "/v4/neighborhood/poi",
			expectedQuery:         url.Values{"latitude": {"40.7128"}, "longitude": {"-74.006"}},
			responseBody:          `{"status":{},"poi":[{}]}`,
			expectError:           false,
			expectedErrorContains: "",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPOI(ctx, WithLatitudeLongitude(40.7128, -74.0060))
			},
		},
		{
			name:                  "GetPOICategoryLookup",
			expectedPath:          "/v4/neighborhood/poicategorylookup",
			expectedQuery:         url.Values{},
			responseBody:          `{"status":{},"poiCategory":[{}]}`,
			expectError:           false,
			expectedErrorContains: "",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPOICategoryLookup(ctx)
			},
		},
	}

	for _, tt := range tests {
		runServiceTest(ctx, t, tt)
	}
}
