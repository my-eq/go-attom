package property

import (
	"context"
	"net/url"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
)

func TestCommunityEndpoints(t *testing.T) {
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
			name:          "GetCommunity",
			expectedPath:  "/v4/neighborhood/neighborhood/community",
			expectedQuery: url.Values{"latitude": {"40.7128"}, "longitude": {"-74.006"}},
			responseBody:  `{"status":{},"community":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetCommunity(ctx, WithLatitudeLongitude(40.7128, -74.0060))
			},
		},
		{
			name:          "GetLocationLookup",
			expectedPath:  "/v4/neighborhood/neighborhood/communitylocation/lookup",
			expectedQuery: url.Values{},
			responseBody:  `{"status":{},"locationLookup":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetLocationLookup(ctx)
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
			}
			c := client.New("test-key", mockClient, client.WithBaseURL("https://example.com/"))
			svc := NewService(c)

			_, err := tt.call(ctx, svc)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
