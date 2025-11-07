package property

import (
	"context"
	"net/url"
	"testing"
)

func TestPropertyEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []TestCase{
		{
			name:          "GetPropertyDetail",
			expectedPath:  "/v4/property/detail",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyDetail(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:          "GetPropertyAddress",
			expectedPath:  "/v4/property/address",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyAddress(ctx, WithAttomID("100"))
			},
		},
		{
			name:                  "GetPropertyAddress_Error_NoIdentifier",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide attomid, id, address, address1, or fips+APN",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyAddress(ctx)
			},
		},
		{
			name:          "GetPropertySnapshot",
			expectedPath:  "/v4/property/snapshot",
			expectedQuery: url.Values{"postalCode": {"62701"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertySnapshot(ctx, WithPostalCode("62701"))
			},
		},
		{
			name:          "GetBasicProfile",
			expectedPath:  "/v4/property/basicprofile",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetBasicProfile(ctx, "123 Main St")
			},
		},
		{
			name:          "GetExpandedProfile",
			expectedPath:  "/v4/property/expandedprofile",
			expectedQuery: url.Values{"geoIdV4": {"geo-1"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetExpandedProfile(ctx, WithGeoIDV4("geo-1"))
			},
		},
		{
			name:          "GetDetailWithSchools",
			expectedPath:  "/v4/property/detailwithschools",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailWithSchools(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailMortgage",
			expectedPath:  "/v4/property/detailmortgage",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"mortgage":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailMortgage(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailOwner",
			expectedPath:  "/v4/property/detailowner",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"owner":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailOwner(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailMortgageOwner",
			expectedPath:  "/v4/property/detailmortgageowner",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"mortgage":[{}],"owner":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailMortgageOwner(ctx, "123 Main St")
			},
		},
		{
			name:          "GetBuildingPermits",
			expectedPath:  "/v4/property/buildingpermits",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"buildingPermit":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetBuildingPermits(ctx, "123 Main St")
			},
		},
		{
			name:          "GetParcelTiles",
			expectedPath:  "/v4/parceltiles/10/512/341.png",
			expectedQuery: url.Values{},
			responseBody:  `{"status":{},"parcelTiles":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetParcelTiles(ctx, 10, 512, 341, "png")
			},
		},
		{
			name:          "GetPreforeclosureDetails",
			expectedPath:  "/property/v3/preforeclosuredetails",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"preforeclosure":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPreforeclosureDetails(ctx, "100")
			},
		},
	}

	for _, tt := range tests {
		runServiceTest(ctx, t, tt)
	}
}
