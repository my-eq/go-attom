package property

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/my-eq/go-attom/pkg/client"
	"github.com/my-eq/go-usps/parser"
)

type mockHTTPClient struct {
	t              *testing.T
	expectedMethod string
	expectedPath   string
	expectedQuery  url.Values
	statusCode     int
	responseBody   string
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.expectedMethod != "" && req.Method != m.expectedMethod {
		m.t.Fatalf("expected method %s, got %s", m.expectedMethod, req.Method)
	}
	if m.expectedPath != "" && req.URL.Path != m.expectedPath {
		m.t.Fatalf("expected path %s, got %s", m.expectedPath, req.URL.Path)
	}
	if m.expectedQuery != nil {
		if diff := diffQuery(m.expectedQuery, req.URL.Query()); diff != "" {
			m.t.Fatalf("query mismatch: %s", diff)
		}
	}
	code := m.statusCode
	if code == 0 {
		code = http.StatusOK
	}
	body := io.NopCloser(strings.NewReader(m.responseBody))
	return &http.Response{StatusCode: code, Body: body, Header: make(http.Header)}, nil
}

func diffQuery(expected, actual url.Values) string {
	if len(expected) != len(actual) {
		return "values length mismatch"
	}
	for key, expectedVals := range expected {
		actualVals, ok := actual[key]
		if !ok {
			return "missing key " + key
		}
		if len(expectedVals) != len(actualVals) {
			return "value length mismatch for key " + key
		}
		for i, val := range expectedVals {
			if actualVals[i] != val {
				return "value mismatch for key " + key
			}
		}
	}
	return ""
}

func TestServiceEndpoints(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name          string
		expectedPath  string
		expectedQuery url.Values
		responseBody  string
		call          func(context.Context, *Service) (interface{}, error)
	}{
		{
			name:          "GetPropertyID",
			expectedPath:  "/propertyapi/v1.0.0/property/id",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"identifier":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyID(ctx, "123 Main St")
			},
		},
		{
			name:          "GetPropertyDetail",
			expectedPath:  "/propertyapi/v1.0.0/property/detail",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyDetail(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:          "GetPropertyAddress",
			expectedPath:  "/propertyapi/v1.0.0/property/address",
			expectedQuery: url.Values{"attomid": {"100"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertyAddress(ctx, WithAttomID("100"))
			},
		},
		{
			name:          "GetPropertySnapshot",
			expectedPath:  "/propertyapi/v1.0.0/property/snapshot",
			expectedQuery: url.Values{"postalCode": {"62701"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetPropertySnapshot(ctx, WithPostalCode("62701"))
			},
		},
		{
			name:          "GetBasicProfile",
			expectedPath:  "/propertyapi/v1.0.0/property/basicprofile",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetBasicProfile(ctx, "123 Main St")
			},
		},
		{
			name:          "GetExpandedProfile",
			expectedPath:  "/propertyapi/v1.0.0/property/expandedprofile",
			expectedQuery: url.Values{"geoIdV4": {"geo-1"}},
			responseBody:  `{"status":{},"property":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetExpandedProfile(ctx, WithGeoIDV4("geo-1"))
			},
		},
		{
			name:          "GetDetailWithSchools",
			expectedPath:  "/propertyapi/v1.0.0/property/detailwithschools",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailWithSchools(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailMortgage",
			expectedPath:  "/propertyapi/v1.0.0/property/detailmortgage",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"mortgage":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailMortgage(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailOwner",
			expectedPath:  "/propertyapi/v1.0.0/property/detailowner",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"owner":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailOwner(ctx, "123 Main St")
			},
		},
		{
			name:          "GetDetailMortgageOwner",
			expectedPath:  "/propertyapi/v1.0.0/property/detailmortgageowner",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"mortgage":[{}],"owner":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetDetailMortgageOwner(ctx, "123 Main St")
			},
		},
		{
			name:          "GetBuildingPermits",
			expectedPath:  "/propertyapi/v1.0.0/property/buildingpermits",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"buildingPermit":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetBuildingPermits(ctx, "123 Main St")
			},
		},
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
			name:          "GetAssessmentDetail",
			expectedPath:  "/propertyapi/v1.0.0/assessment/detail",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"assessment":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentDetail(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:          "GetAssessmentSnapshot",
			expectedPath:  "/propertyapi/v1.0.0/assessment/snapshot",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"assessment":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentSnapshot(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:          "GetAssessmentHistory",
			expectedPath:  "/propertyapi/v1.0.0/assessmenthistory/detail",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"assessmentHistory":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetAssessmentHistory(ctx, WithAddress("123 Main St"))
			},
		},
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
		{
			name:          "SearchSchools",
			expectedPath:  "/propertyapi/v1.0.0/school/search",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.SearchSchools(ctx, WithAddress("123 Main St"))
			},
		},
		{
			name:          "GetSchoolProfile",
			expectedPath:  "/propertyapi/v1.0.0/school/profile",
			expectedQuery: url.Values{"schoolId": {"200"}},
			responseBody:  `{"status":{},"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolProfile(ctx, "200")
			},
		},
		{
			name:          "GetSchoolDistrict",
			expectedPath:  "/propertyapi/v1.0.0/school/district",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"district":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDistrict(ctx, "123 Main St")
			},
		},
		{
			name:          "GetSchoolDetailWithSchools",
			expectedPath:  "/propertyapi/v1.0.0/school/detailwithschools",
			expectedQuery: url.Values{"address": {"123 Main St"}},
			responseBody:  `{"status":{},"property":[{}],"school":[{}]}`,
			call: func(ctx context.Context, svc *Service) (interface{}, error) {
				return svc.GetSchoolDetailWithSchools(ctx, "123 Main St")
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mock := &mockHTTPClient{
				t:              t,
				expectedMethod: http.MethodGet,
				expectedPath:   tt.expectedPath,
				expectedQuery:  tt.expectedQuery,
				responseBody:   tt.responseBody,
			}
			c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
			svc := NewService(c)
			result, err := tt.call(ctx, svc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result == nil {
				t.Fatalf("expected result, got nil")
			}
		})
	}
}

func TestServiceValidationErrors(t *testing.T) {
	ctx := context.Background()
	mock := &mockHTTPClient{t: t, responseBody: `{"status":{}}`}
	c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
	svc := NewService(c)

	cases := []struct {
		name string
		call func() error
	}{
		{name: "PropertyID", call: func() error { _, err := svc.GetPropertyID(ctx, ""); return err }},
		{name: "PropertyDetail", call: func() error { _, err := svc.GetPropertyDetail(ctx); return err }},
		{name: "PropertyAddress", call: func() error { _, err := svc.GetPropertyAddress(ctx); return err }},
		{name: "PropertySnapshot", call: func() error { _, err := svc.GetPropertySnapshot(ctx); return err }},
		{name: "BasicProfile", call: func() error { _, err := svc.GetBasicProfile(ctx, ""); return err }},
		{name: "ExpandedProfile", call: func() error { _, err := svc.GetExpandedProfile(ctx); return err }},
		{name: "DetailWithSchools", call: func() error { _, err := svc.GetDetailWithSchools(ctx, ""); return err }},
		{name: "DetailMortgage", call: func() error { _, err := svc.GetDetailMortgage(ctx, ""); return err }},
		{name: "DetailOwner", call: func() error { _, err := svc.GetDetailOwner(ctx, ""); return err }},
		{name: "DetailMortgageOwner", call: func() error { _, err := svc.GetDetailMortgageOwner(ctx, ""); return err }},
		{name: "BuildingPermits", call: func() error { _, err := svc.GetBuildingPermits(ctx, ""); return err }},
		{name: "SaleDetail", call: func() error { _, err := svc.GetSaleDetail(ctx); return err }},
		{name: "SaleSnapshot", call: func() error { _, err := svc.GetSaleSnapshot(ctx); return err }},
		{name: "AssessmentDetail", call: func() error { _, err := svc.GetAssessmentDetail(ctx); return err }},
		{name: "AssessmentSnapshot", call: func() error { _, err := svc.GetAssessmentSnapshot(ctx); return err }},
		{name: "AssessmentHistory", call: func() error { _, err := svc.GetAssessmentHistory(ctx); return err }},
		{name: "AVMSnapshot", call: func() error { _, err := svc.GetAVMSnapshot(ctx); return err }},
		{name: "AttomAVMDetail", call: func() error { _, err := svc.GetAttomAVMDetail(ctx); return err }},
		{name: "AVMHistory", call: func() error { _, err := svc.GetAVMHistory(ctx); return err }},
		{name: "RentalAVM", call: func() error { _, err := svc.GetRentalAVM(ctx); return err }},
		{name: "SalesHistoryDetail", call: func() error { _, err := svc.GetSalesHistoryDetail(ctx); return err }},
		{name: "SalesHistorySnapshot", call: func() error { _, err := svc.GetSalesHistorySnapshot(ctx); return err }},
		{name: "SalesHistoryBasic", call: func() error { _, err := svc.GetSalesHistoryBasic(ctx); return err }},
		{name: "SalesHistoryExpanded", call: func() error { _, err := svc.GetSalesHistoryExpanded(ctx); return err }},
		{name: "SalesTrendSnapshot", call: func() error { _, err := svc.GetSalesTrendSnapshot(ctx); return err }},
		{name: "TransactionSalesTrend", call: func() error { _, err := svc.GetTransactionSalesTrend(ctx); return err }},
		{name: "SearchSchools", call: func() error { _, err := svc.SearchSchools(ctx); return err }},
		{name: "SchoolProfile", call: func() error { _, err := svc.GetSchoolProfile(ctx, ""); return err }},
		{name: "SchoolDistrict", call: func() error { _, err := svc.GetSchoolDistrict(ctx, ""); return err }},
		{name: "SchoolDetailWithSchools", call: func() error { _, err := svc.GetSchoolDetailWithSchools(ctx, ""); return err }},
		{name: "AllEventsDetail", call: func() error { _, err := svc.GetAllEventsDetail(ctx); return err }},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Fatalf("expected error")
			}
			if !errors.Is(err, ErrMissingParameter) {
				t.Fatalf("expected ErrMissingParameter, got %v", err)
			}
		})
	}
}

func TestServiceErrorResponse(t *testing.T) {
	ctx := context.Background()
	mock := &mockHTTPClient{
		t:              t,
		expectedMethod: http.MethodGet,
		expectedPath:   "/propertyapi/v1.0.0/property/detail",
		expectedQuery:  url.Values{"attomid": {"100"}},
		statusCode:     http.StatusBadRequest,
		responseBody:   `{"status":{"msg":"bad request"}}`,
	}
	c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
	svc := NewService(c)

	_, err := svc.GetPropertyDetail(ctx, WithAttomID("100"))
	if err == nil {
		t.Fatalf("expected error")
	}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, apiErr.StatusCode)
	}
}

func TestAddressHelpers(t *testing.T) {
	input := "123 Main St, Springfield, IL 62701"
	parsed, diags := parser.Parse(input)
	if parsed == nil {
		t.Fatalf("expected parsed address, got nil")
	}
	for _, d := range diags {
		if d.Severity == parser.SeverityError {
			t.Fatalf("unexpected parse error: %v", d)
		}
	}
	req := parsed.ToAddressRequest()
	if req == nil {
		t.Fatalf("expected AddressRequest, got nil")
	}
	lines := req.Lines()
	if len(lines) < 1 {
		t.Fatalf("expected at least one line")
	}
	address1 := lines[0]
	address2 := ""
	if len(lines) > 1 {
		address2 = lines[1]
	}
	if address1 != "123 MAIN ST" {
		t.Fatalf("unexpected address1: %s", address1)
	}
	if address2 != "SPRINGFIELD, IL 62701" {
		t.Fatalf("unexpected address2: %s", address2)
	}
	joined := req.String()
	if joined != "123 MAIN ST, SPRINGFIELD, IL 62701" {
		t.Fatalf("unexpected joined value: %s", joined)
	}
	single := "Only Second"
	parsed2, _ := parser.Parse(single)
	if parsed2 == nil {
		t.Fatalf("expected parsed address for single line input")
	}
	if got := parsed2.ToAddressRequest().String(); got != "ONLY SECOND" {
		t.Fatalf("unexpected single line result: %s", got)
	}
	if err := ValidateFIPSAndAPN("", "123"); err == nil {
		t.Fatalf("expected error when FIPS missing")
	}
	if err := ValidateFIPSAndAPN("001", "456"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
