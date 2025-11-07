package property

import (
	"context"
	"net/url"
	"testing"
)

func TestSchoolEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []TestCase{
		{
			name:          "SearchSchools",
			expectedPath:  "/v4/school/search",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.SearchSchools(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:                  "SearchSchools_Error_NoGeoContext",
			expectedPath:          "",
			expectedQuery:         url.Values{},
			responseBody:          "",
			expectError:           true,
			expectedErrorContains: "provide address or latitude/longitude",
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.SearchSchools(ctx)
			},
		},
		{
			name:          "GetSchoolProfile",
			expectedPath:  "/v4/school/profile",
			expectedQuery: url.Values{"schoolId": {"200"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolProfile(ctx, "200")
			},
		},
		{
			name:          "GetSchoolDistrict",
			expectedPath:  "/v4/school/district",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"district":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDistrict(ctx, "123 Main St")
			},
		},
		{
			name:          "GetSchoolDetailWithSchools",
			expectedPath:  "/v4/school/detailwithschools",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDetailWithSchools(ctx, "123 Main St")
			},
		},
		// --- NEW ENDPOINT TESTS ---
		{
			name:          "GetSchoolSnapshot",
			expectedPath:  "/v4/school/snapshot",
			expectedQuery: url.Values{"latitude": {"40.0"}, "longitude": {"-75.0"}, "radius": {"10"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolSnapshot(ctx, "40.0", "-75.0", "10", "", nil)
			},
		},
		{
			name:          "GetSchoolDetail",
			expectedPath:  "/v4/school/detail",
			expectedQuery: url.Values{"id": {"200"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDetail(ctx, "200")
			},
		},
		{
			name:          "GetSchoolDistrictDetail",
			expectedPath:  "/v4/school/districtdetail",
			expectedQuery: url.Values{"id": {"300"}},
			responseBody:  `{"status":{},"district":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDistrictDetail(ctx, "300")
			},
		},
	}

	for _, tt := range tests {
		runServiceTest(ctx, t, tt)
	}
}
