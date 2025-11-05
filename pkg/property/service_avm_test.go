package property

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
)

func TestAVMEndpoints(t *testing.T) {
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
			name:          "GetAVMSnapshot",
			expectedPath:  "/propertyapi/v1.0.0/avm/snapshot",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"avm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMSnapshot(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetAttomAVMDetail",
			expectedPath:  "/propertyapi/v1.0.0/attomavm/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"attomAvm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAttomAVMDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetAVMHistory",
			expectedPath:  "/propertyapi/v1.0.0/avmhistory/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"avmHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMHistory(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetRentalAVM",
			expectedPath:  "/propertyapi/v1.0.0/valuation/rentalavm",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"rentalAvm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetRentalAVM(ctx, WithAttomID("100"))
			},
		},
		// --- NEW ENDPOINT TESTS ---
		{
			name:          "GetAVMSnapshotGeo",
			expectedPath:  "/propertyapi/v1.0.0/avm/snapshot",
			expectedQuery: url.Values{"geoIdV4": {"geo-2"}, "minavmvalue": {"100000"}, "maxavmvalue": {"500000"}, "propertytype": {"SFR"}},
			responseBody:  `{"status":{},"avm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMSnapshotGeo(ctx, "geo-2", "100000", "500000", "SFR")
			},
		},
		{
			name:          "GetAVMHistoryByAddress",
			expectedPath:  "/propertyapi/v1.0.0/avmhistory/detail",
			expectedQuery: url.Values{"address1": {"123 Main St"}, "address2": {"Springfield, IL"}},
			responseBody:  `{"status":{},"avmHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMHistoryByAddress(ctx, "123 Main St", "Springfield, IL")
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

// ...existing code...
// AVM endpoint tests will be moved here.
