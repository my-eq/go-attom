package property

import (
	"context"
	"net/url"
	"testing"
)

func TestAVMEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []TestCase{
		{
			name:          "GetAVMSnapshot",
			expectedPath:  "/v4/property/snapshot",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"avm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMSnapshot(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetAVMSnapshot_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMSnapshot(ctx)
			},
		},
		{
			name:          "GetAttomAVMDetail",
			expectedPath:  "/v4/property/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"attomAvm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAttomAVMDetail(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetAttomAVMDetail_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAttomAVMDetail(ctx)
			},
		},
		{
			name:          "GetAVMHistory",
			expectedPath:  "/v4/property/detail",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"avmHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMHistory(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetAVMHistory_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMHistory(ctx)
			},
		},
		{
			name:          "GetRentalAVM",
			expectedPath:  "/v4/property/rentalavm",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"rentalAvm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetRentalAVM(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetRentalAVM_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetRentalAVM(ctx)
			},
		},
		// --- NEW ENDPOINT TESTS ---
		{
			name:          "GetAVMSnapshotGeo",
			expectedPath:  "/v4/property/snapshot",
			expectedQuery: url.Values{"geoIdV4": {"geo-2"}, "minavmvalue": {"100000"}, "maxavmvalue": {"500000"}, "propertytype": {"SFR"}},
			responseBody:  `{"status":{},"avm":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMSnapshotGeo(ctx, "geo-2", "100000", "500000", "SFR")
			},
		},
		{
			name:                  "GetAVMSnapshotGeo_Error_NoGeoID",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "geoIdV4 required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMSnapshotGeo(ctx, "", "100000", "500000", "SFR")
			},
		},
		{
			name:          "GetAVMHistoryByAddress",
			expectedPath:  "/v4/property/detail",
			expectedQuery: url.Values{"address1": {"123 Main St"}, "address2": {"Springfield, IL"}},
			responseBody:  `{"status":{},"avmHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMHistoryByAddress(ctx, "123 Main St", "Springfield, IL")
			},
		},
		{
			name:                  "GetAVMHistoryByAddress_Error_NoAddress",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "address1 and address2 required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAVMHistoryByAddress(ctx, "", "Springfield, IL")
			},
		},
		{
			name:          "GetHomeEquity",
			expectedPath:  "/v4/property/homeequity",
			expectedQuery: url.Values{"address1": {"123 Main St"}, "address2": {"Springfield, IL"}},
			responseBody:  `{"status":{},"homeEquity":150000.50}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetHomeEquity(ctx, "123 Main St", "Springfield, IL")
			},
		},
		{
			name:                  "GetHomeEquity_Error_NoAddress",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "address1 and address2 required",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetHomeEquity(ctx, "", "Springfield, IL")
			},
		},
	}

	for _, tt := range tests {
		runServiceTest(t, ctx, tt)
	}
}

// ...existing code...
// AVM endpoint tests will be moved here.
