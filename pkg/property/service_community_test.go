package property

import (
	"context"
	"net/url"
	"testing"
)

func TestCommunityEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []TestCase{
		{
			name:                  "GetCommunity",
			expectedPath:          "/v4/neighborhood/neighborhood/community",
			expectedQuery:         url.Values{"latitude": {"40.7128"}, "longitude": {"-74.006"}},
			responseBody:          `{"status":{},"community":[{}]}`,
			expectError:           false,
			expectedErrorContains: "",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetCommunity(ctx, WithLatitudeLongitude(40.7128, -74.0060))
			},
		},
		{
			name:                  "GetLocationLookup",
			expectedPath:          "/v4/neighborhood/neighborhood/communitylocation/lookup",
			expectedQuery:         url.Values{},
			responseBody:          `{"status":{},"locationLookup":[{}]}`,
			expectError:           false,
			expectedErrorContains: "",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetLocationLookup(ctx)
			},
		},
	}

	for _, tt := range tests {
		runServiceTest(ctx, t, tt)
	}
}
