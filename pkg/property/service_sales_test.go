package property

import (
	"context"
	"net/url"
	"testing"
)

func TestSalesEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []TestCase{
		{
			name:          "GetSaleDetail",
			expectedPath:  "/v4/transaction/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"sale":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetSaleDetail_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleDetail(ctx)
			},
		},
		{
			name:          "GetSaleSnapshot",
			expectedPath:  "/v4/transaction/snapshot",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"sale":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleSnapshot(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetSaleSnapshot_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleSnapshot(ctx)
			},
		},
		{
			name:          "GetSalesHistoryDetail",
			expectedPath:  "/v4/transaction/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetSalesHistoryDetail_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryDetail(ctx)
			},
		},
		{
			name:          "GetSalesHistorySnapshot",
			expectedPath:  "/v4/transaction/snapshot",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistorySnapshot(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetSalesHistorySnapshot_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistorySnapshot(ctx)
			},
		},
		{
			name:          "GetSalesHistoryBasic",
			expectedPath:  "/v4/transaction/basichistory",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryBasic(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetSalesHistoryBasic_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryBasic(ctx)
			},
		},
		{
			name:          "GetSalesHistoryExpanded",
			expectedPath:  "/v4/transaction/expandedhistory",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"salesHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryExpanded(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetSalesHistoryExpanded_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesHistoryExpanded(ctx)
			},
		},
		{
			name:          "GetSalesTrendSnapshot",
			expectedPath:  "/v4/transaction/snapshot",
			expectedQuery: url.Values{"geoIdV4": {"geo-1"}},
			responseBody:  `{"status":{},"salesTrend":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesTrendSnapshot(ctx, WithGeoIDV4("geo-1"))
			},
		},
		{
			name:                  "GetSalesTrendSnapshot_Error_NoGeoID",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "geoIdV4 required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSalesTrendSnapshot(ctx)
			},
		},
		{
			name:          "GetTransactionSalesTrend",
			expectedPath:  "/v4/transaction/salestrend",
			expectedQuery: url.Values{"geoIdV4": {"geo-1"}},
			responseBody:  `{"status":{},"transactionTrend":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetTransactionSalesTrend(ctx, WithGeoIDV4("geo-1"))
			},
		},
		{
			name:                  "GetTransactionSalesTrend_Error_NoGeoID",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "geoIdV4 required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetTransactionSalesTrend(ctx)
			},
		},
		{
			name:          "GetAllEventsDetail",
			expectedPath:  "/propertyapi/v1.0.0/allevents/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"event":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAllEventsDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetAllEventsDetail_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAllEventsDetail(ctx)
			},
		},
		{
			name:          "GetAllEventsSnapshot",
			expectedPath:  "/propertyapi/v1.0.0/allevents/snapshot",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"snapshot":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAllEventsSnapshot(ctx, "123 Main St")
			},
		},
		{
			name:                  "GetAllEventsSnapshot_Error_NoAddress",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "address required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAllEventsSnapshot(ctx, "")
			},
		},
		{
			name:          "GetSaleComparablesByAddress",
			expectedPath:  "/property/v2/salescomparables/address/123%20Main%20St/Springfield/Cook/IL/62701",
			expectedQuery: url.Values{"address": {"123 Main St, Springfield, Cook, IL 62701"}},
			responseBody:  `{"status":{},"saleComparables":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleComparablesByAddress(ctx, "123 Main St", "Springfield", "Cook", "IL", "62701")
			},
		},
		{
			name:                  "GetSaleComparablesByAddress_Error_MissingComponents",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "address components required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleComparablesByAddress(ctx, "", "Springfield", "Cook", "IL", "62701")
			},
		},
		{
			name:          "GetSaleComparablesByAPN",
			expectedPath:  "/property/v2/salescomparables/apn/123456789/Cook/IL",
			expectedQuery: url.Values{"APN": {"123456789"}},
			responseBody:  `{"status":{},"saleComparables":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleComparablesByAPN(ctx, "123456789", "Cook", "IL")
			},
		},
		{
			name:                  "GetSaleComparablesByAPN_Error_MissingComponents",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "APN, county, and state required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleComparablesByAPN(ctx, "", "Cook", "IL")
			},
		},
		{
			name:          "GetSaleComparablesByPropID",
			expectedPath:  "/property/v2/salescomparables/propid/100",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"saleComparables":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleComparablesByPropID(ctx, "100")
			},
		},
		{
			name:                  "GetSaleComparablesByPropID_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSaleComparablesByPropID(ctx, "")
			},
		},
		{
			name:          "GetTransportationNoise",
			expectedPath:  "/propertyapi/v1.0.0/transportationnoise",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"transportationNoise":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetTransportationNoise(ctx, "100")
			},
		},
		{
			name:                  "GetTransportationNoise_Error_NoAttomID",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "attomid required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetTransportationNoise(ctx, "")
			},
		},
	}

	for _, tt := range tests {
		runServiceTest(t, ctx, tt)
	}
}
