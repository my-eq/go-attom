package property

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
)

func TestSalesEndpoints(t *testing.T) {
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
			name:          "GetSaleDetail",
			expectedPath:  "/propertyapi/v1.0.0/sale/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"sale":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetSaleSnapshot",
			expectedPath:  "/propertyapi/v1.0.0/sale/snapshot",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"sale":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleSnapshot(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetSalesHistoryDetail",
			expectedPath:  "/propertyapi/v1.0.0/saleshistory/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetSalesHistorySnapshot",
			expectedPath:  "/propertyapi/v1.0.0/saleshistory/snapshot",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistorySnapshot(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetSalesHistoryBasic",
			expectedPath:  "/propertyapi/v1.0.0/saleshistory/basichistory",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryBasic(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetSalesHistoryExpanded",
			expectedPath:  "/propertyapi/v1.0.0/saleshistory/expandedhistory",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryExpanded(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetSalesTrendSnapshot",
			expectedPath:  "/propertyapi/v1.0.0/salestrend/snapshot",
			expectedQuery: url.Values{"geoIdV4": {"geo-1"}},
			responseBody:  `{"status":{},"salesTrend":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesTrendSnapshot(ctx, WithGeoIDV4("geo-1"))
			},
		},
		{
			name:          "GetTransactionSalesTrend",
			expectedPath:  "/propertyapi/v1.0.0/transaction/salestrend",
			expectedQuery: url.Values{"geoIdV4": {"geo-1"}},
			responseBody:  `{"status":{},"transactionTrend":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetTransactionSalesTrend(ctx, WithGeoIDV4("geo-1"))
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
