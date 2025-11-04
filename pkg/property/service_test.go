package property

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/my-eq/go-attom/pkg/client"
	"github.com/my-eq/go-usps/parser"
)

type mockHTTPClient struct {
	t              *testing.T
	expectedQuery  url.Values
	expectedMethod string
	expectedPath   string
	responseBody   string
	statusCode     int
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

// errorReader simulates a failing io.ReadCloser for testing error paths
type errorReader struct{}

func (errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("simulated read error")
}

func (errorReader) Close() error {
	return nil
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
		call          func(context.Context, *Service) (interface{}, error)
		expectedQuery url.Values
		name          string
		expectedPath  string
		responseBody  string
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
		call func() error
		name string
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

func TestErrorTypes(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var e *Error
		got := e.Error()
		if got != "property: nil error" {
			t.Errorf("expected 'property: nil error', got %q", got)
		}
	})

	t.Run("error with message", func(t *testing.T) {
		e := &Error{Message: "test message"}
		got := e.Error()
		if got != "property: test message" {
			t.Errorf("expected 'property: test message', got %q", got)
		}
	})

	t.Run("error with status message", func(t *testing.T) {
		msg := "status message"
		e := &Error{Status: &Status{Msg: &msg}}
		got := e.Error()
		if got != "property: status message" {
			t.Errorf("expected 'property: status message', got %q", got)
		}
	})

	t.Run("error with status code only", func(t *testing.T) {
		code := 400
		e := &Error{Status: &Status{Code: &code}}
		got := e.Error()
		if got != "property: status code 400" {
			t.Errorf("expected 'property: status code 400', got %q", got)
		}
	})

	t.Run("error with http status only", func(t *testing.T) {
		e := &Error{StatusCode: 500}
		got := e.Error()
		if got != "property: http status 500" {
			t.Errorf("expected 'property: http status 500', got %q", got)
		}
	})
}

func TestNewService(t *testing.T) {
	t.Run("nil client", func(t *testing.T) {
		svc := NewService(nil)
		if svc != nil {
			t.Errorf("expected nil service for nil client")
		}
	})

	t.Run("valid client", func(t *testing.T) {
		c := client.New("test-key", nil)
		svc := NewService(c)
		if svc == nil {
			t.Errorf("expected non-nil service")
		}
	})
}

func TestEnsureClient(t *testing.T) {
	t.Run("nil service", func(t *testing.T) {
		var svc *Service
		err := svc.ensureClient()
		if err == nil {
			t.Errorf("expected error for nil service")
		}
	})

	t.Run("service with nil client", func(t *testing.T) {
		svc := &Service{}
		err := svc.ensureClient()
		if err == nil {
			t.Errorf("expected error for nil client")
		}
	})
}

func TestWithStringSlice(t *testing.T) {
	t.Run("valid slice", func(t *testing.T) {
		vals := url.Values{}
		opt := WithStringSlice("test", []string{"a", "b", "c"}, ",")
		opt(vals)
		if vals.Get("test") != "a,b,c" {
			t.Errorf("expected 'a,b,c', got %q", vals.Get("test"))
		}
	})

	t.Run("empty key", func(t *testing.T) {
		vals := url.Values{}
		opt := WithStringSlice("", []string{"a"}, ",")
		opt(vals)
		if len(vals) != 0 {
			t.Errorf("expected no values for empty key")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		vals := url.Values{}
		opt := WithStringSlice("test", []string{}, ",")
		opt(vals)
		if len(vals) != 0 {
			t.Errorf("expected no values for empty slice")
		}
	})

	t.Run("default separator", func(t *testing.T) {
		vals := url.Values{}
		opt := WithStringSlice("test", []string{"a", "b"}, "")
		opt(vals)
		if vals.Get("test") != "a|b" {
			t.Errorf("expected 'a|b' with default separator, got %q", vals.Get("test"))
		}
	})
}

func TestWithPropertyID(t *testing.T) {
	vals := url.Values{}
	WithPropertyID("123")(vals)
	if vals.Get("id") != "123" {
		t.Errorf("expected '123', got %q", vals.Get("id"))
	}
}

func TestWithFIPSAndAPN(t *testing.T) {
	vals := url.Values{}
	WithFIPSAndAPN("001", "456")(vals)
	if vals.Get("fips") != "001" {
		t.Errorf("expected '001', got %q", vals.Get("fips"))
	}
	if vals.Get("APN") != "456" {
		t.Errorf("expected '456', got %q", vals.Get("APN"))
	}
}

func TestWithAddressLines(t *testing.T) {
	vals := url.Values{}
	WithAddressLines("123 Main St", "City, ST")(vals)
	if vals.Get("address1") != "123 Main St" {
		t.Errorf("expected '123 Main St', got %q", vals.Get("address1"))
	}
	if vals.Get("address2") != "City, ST" {
		t.Errorf("expected 'City, ST', got %q", vals.Get("address2"))
	}
}

func TestWithLatitudeLongitude(t *testing.T) {
	vals := url.Values{}
	WithLatitudeLongitude(40.7128, -74.0060)(vals)
	if vals.Get("latitude") != "40.7128" {
		t.Errorf("expected '40.7128', got %q", vals.Get("latitude"))
	}
	if vals.Get("longitude") != "-74.006" {
		t.Errorf("expected '-74.006', got %q", vals.Get("longitude"))
	}
}

func TestWithRadius(t *testing.T) {
	t.Run("valid radius", func(t *testing.T) {
		vals := url.Values{}
		WithRadius(5.5)(vals)
		if vals.Get("radius") != "5.5" {
			t.Errorf("expected '5.5', got %q", vals.Get("radius"))
		}
	})

	t.Run("zero or negative", func(t *testing.T) {
		vals := url.Values{}
		WithRadius(0)(vals)
		if vals.Get("radius") != "" {
			t.Errorf("expected empty for zero radius")
		}
		WithRadius(-1)(vals)
		if vals.Get("radius") != "" {
			t.Errorf("expected empty for negative radius")
		}
	})
}

func TestWithCityName(t *testing.T) {
	vals := url.Values{}
	WithCityName("Springfield")(vals)
	if vals.Get("cityname") != "Springfield" {
		t.Errorf("expected 'Springfield', got %q", vals.Get("cityname"))
	}
}

func TestWithGeoID(t *testing.T) {
	vals := url.Values{}
	WithGeoID("PL0820000")(vals)
	if vals.Get("geoid") != "PL0820000" {
		t.Errorf("expected 'PL0820000', got %q", vals.Get("geoid"))
	}
}

func TestWithPropertyType(t *testing.T) {
	vals := url.Values{}
	WithPropertyType("SFR")(vals)
	if vals.Get("propertytype") != "SFR" {
		t.Errorf("expected 'SFR', got %q", vals.Get("propertytype"))
	}
}

func TestWithPropertyIndicator(t *testing.T) {
	t.Run("valid indicator", func(t *testing.T) {
		vals := url.Values{}
		WithPropertyIndicator(10)(vals)
		if vals.Get("propertyIndicator") != "10" {
			t.Errorf("expected '10', got %q", vals.Get("propertyIndicator"))
		}
	})

	t.Run("zero or negative", func(t *testing.T) {
		vals := url.Values{}
		WithPropertyIndicator(0)(vals)
		if vals.Get("propertyIndicator") != "" {
			t.Errorf("expected empty for zero indicator")
		}
	})
}

func TestWithBedsRange(t *testing.T) {
	vals := url.Values{}
	WithBedsRange(2, 4)(vals)
	if vals.Get("minBeds") != "2" {
		t.Errorf("expected '2', got %q", vals.Get("minBeds"))
	}
	if vals.Get("maxBeds") != "4" {
		t.Errorf("expected '4', got %q", vals.Get("maxBeds"))
	}
}

func TestWithBathsRange(t *testing.T) {
	vals := url.Values{}
	WithBathsRange(1.5, 3.0)(vals)
	if vals.Get("minBathsTotal") != "1.5" {
		t.Errorf("expected '1.5', got %q", vals.Get("minBathsTotal"))
	}
	if vals.Get("maxBathsTotal") != "3" {
		t.Errorf("expected '3', got %q", vals.Get("maxBathsTotal"))
	}
}

func TestWithSaleAmountRange(t *testing.T) {
	vals := url.Values{}
	WithSaleAmountRange(100000, 500000)(vals)
	if vals.Get("minSaleAmt") != "100000" {
		t.Errorf("expected '100000', got %q", vals.Get("minSaleAmt"))
	}
	if vals.Get("maxSaleAmt") != "500000" {
		t.Errorf("expected '500000', got %q", vals.Get("maxSaleAmt"))
	}
}

func TestWithUniversalSizeRange(t *testing.T) {
	vals := url.Values{}
	WithUniversalSizeRange(1000, 3000)(vals)
	if vals.Get("minUniversalSize") != "1000" {
		t.Errorf("expected '1000', got %q", vals.Get("minUniversalSize"))
	}
	if vals.Get("maxUniversalSize") != "3000" {
		t.Errorf("expected '3000', got %q", vals.Get("maxUniversalSize"))
	}
}

func TestWithYearBuiltRange(t *testing.T) {
	vals := url.Values{}
	WithYearBuiltRange(1990, 2020)(vals)
	if vals.Get("minYearBuilt") != "1990" {
		t.Errorf("expected '1990', got %q", vals.Get("minYearBuilt"))
	}
	if vals.Get("maxYearBuilt") != "2020" {
		t.Errorf("expected '2020', got %q", vals.Get("maxYearBuilt"))
	}
}

func TestWithLotSize1Range(t *testing.T) {
	vals := url.Values{}
	WithLotSize1Range(0.5, 2.0)(vals)
	if vals.Get("minLotSize1") != "0.5" {
		t.Errorf("expected '0.5', got %q", vals.Get("minLotSize1"))
	}
	if vals.Get("maxLotSize1") != "2" {
		t.Errorf("expected '2', got %q", vals.Get("maxLotSize1"))
	}
}

func TestWithLotSize2Range(t *testing.T) {
	vals := url.Values{}
	WithLotSize2Range(5000, 10000)(vals)
	if vals.Get("minLotSize2") != "5000" {
		t.Errorf("expected '5000', got %q", vals.Get("minLotSize2"))
	}
	if vals.Get("maxLotSize2") != "10000" {
		t.Errorf("expected '10000', got %q", vals.Get("maxLotSize2"))
	}
}

func TestWithDateRange(t *testing.T) {
	t.Run("valid range", func(t *testing.T) {
		vals := url.Values{}
		start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
		WithDateRange("SaleDate", start, end)(vals)
		if vals.Get("startSaleDate") != "2020/01/01" {
			t.Errorf("expected '2020/01/01', got %q", vals.Get("startSaleDate"))
		}
		if vals.Get("endSaleDate") != "2020/12/31" {
			t.Errorf("expected '2020/12/31', got %q", vals.Get("endSaleDate"))
		}
	})

	t.Run("zero times", func(t *testing.T) {
		vals := url.Values{}
		WithDateRange("Test", time.Time{}, time.Time{})(vals)
		if vals.Get("startTest") != "" || vals.Get("endTest") != "" {
			t.Errorf("expected empty for zero times")
		}
	})
}

func TestWithISODateRange(t *testing.T) {
	vals := url.Values{}
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	WithISODateRange("CalendarDate", start, end)(vals)
	if vals.Get("startCalendarDate") != "2020-01-01" {
		t.Errorf("expected '2020-01-01', got %q", vals.Get("startCalendarDate"))
	}
	if vals.Get("endCalendarDate") != "2020-12-31" {
		t.Errorf("expected '2020-12-31', got %q", vals.Get("endCalendarDate"))
	}
}

func TestWithPage(t *testing.T) {
	t.Run("valid page", func(t *testing.T) {
		vals := url.Values{}
		WithPage(5)(vals)
		if vals.Get("page") != "5" {
			t.Errorf("expected '5', got %q", vals.Get("page"))
		}
	})

	t.Run("zero or negative", func(t *testing.T) {
		vals := url.Values{}
		WithPage(0)(vals)
		if vals.Get("page") != "" {
			t.Errorf("expected empty for zero page")
		}
	})
}

func TestWithPageSize(t *testing.T) {
	vals := url.Values{}
	WithPageSize(100)(vals)
	if vals.Get("pagesize") != "100" {
		t.Errorf("expected '100', got %q", vals.Get("pagesize"))
	}
}

func TestWithOrderBy(t *testing.T) {
	vals := url.Values{}
	WithOrderBy("saleamt")(vals)
	if vals.Get("orderby") != "saleamt" {
		t.Errorf("expected 'saleamt', got %q", vals.Get("orderby"))
	}
}

func TestWithAdditionalParam(t *testing.T) {
	vals := url.Values{}
	WithAdditionalParam("custom", "value")(vals)
	if vals.Get("custom") != "value" {
		t.Errorf("expected 'value', got %q", vals.Get("custom"))
	}
}

func TestValidatorFunctions(t *testing.T) {
	t.Run("requireAny success", func(t *testing.T) {
		vals := url.Values{"key1": {"value1"}}
		err := requireAny(vals, "key1", "key2")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("requireAny failure", func(t *testing.T) {
		vals := url.Values{}
		err := requireAny(vals, "key1", "key2")
		if err == nil {
			t.Errorf("expected error when no keys present")
		}
		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})

	t.Run("requireAll success", func(t *testing.T) {
		vals := url.Values{"key1": {"value1"}, "key2": {"value2"}}
		err := requireAll(vals, "key1", "key2")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("requireAll failure", func(t *testing.T) {
		vals := url.Values{"key1": {"value1"}}
		err := requireAll(vals, "key1", "key2")
		if err == nil {
			t.Errorf("expected error when key missing")
		}
		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})

	t.Run("requirePropertyIdentifier with attomid", func(t *testing.T) {
		vals := url.Values{"attomid": {"123"}}
		err := requirePropertyIdentifier(vals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("requirePropertyIdentifier with id", func(t *testing.T) {
		vals := url.Values{"id": {"123"}}
		err := requirePropertyIdentifier(vals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("requirePropertyIdentifier with address", func(t *testing.T) {
		vals := url.Values{"address": {"123 Main St"}}
		err := requirePropertyIdentifier(vals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("requirePropertyIdentifier with address1", func(t *testing.T) {
		vals := url.Values{"address1": {"123 Main St"}}
		err := requirePropertyIdentifier(vals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("requirePropertyIdentifier with fips and APN", func(t *testing.T) {
		vals := url.Values{"fips": {"001"}, "APN": {"456"}}
		err := requirePropertyIdentifier(vals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("requirePropertyIdentifier failure", func(t *testing.T) {
		vals := url.Values{}
		err := requirePropertyIdentifier(vals)
		if err == nil {
			t.Errorf("expected error when no identifier present")
		}
		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})

	t.Run("ensureGeoContext with address", func(t *testing.T) {
		vals := url.Values{"address": {"123 Main St"}}
		err := ensureGeoContext(vals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ensureGeoContext with lat/lon", func(t *testing.T) {
		vals := url.Values{"latitude": {"40.7"}, "longitude": {"-74.0"}}
		err := ensureGeoContext(vals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("ensureGeoContext failure", func(t *testing.T) {
		vals := url.Values{}
		err := ensureGeoContext(vals)
		if err == nil {
			t.Errorf("expected error when no geo context present")
		}
		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})
}

func TestDoGetErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("decode error", func(t *testing.T) {
		mock := &mockHTTPClient{
			t:            t,
			statusCode:   http.StatusOK,
			responseBody: `{invalid json}`,
		}
		c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
		svc := NewService(c)

		var resp DetailResponse
		err := svc.doGet(ctx, "property/detail", url.Values{}, &resp)
		if err == nil {
			t.Errorf("expected decode error")
		}
		if !strings.Contains(err.Error(), "failed to decode") {
			t.Errorf("expected decode error message, got %v", err)
		}
	})

	t.Run("http error without readable body", func(t *testing.T) {
		mock := &mockHTTPClient{
			t:            t,
			statusCode:   http.StatusBadRequest,
			responseBody: "",
		}
		c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
		svc := NewService(c)

		var resp DetailResponse
		err := svc.doGet(ctx, "property/detail", url.Values{}, &resp)
		if err == nil {
			t.Errorf("expected error for bad request")
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})

	t.Run("nil output parameter", func(t *testing.T) {
		mock := &mockHTTPClient{
			t:            t,
			statusCode:   http.StatusOK,
			responseBody: `{"status":{}}`,
		}
		c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
		svc := NewService(c)

		err := svc.doGet(ctx, "property/detail", url.Values{}, nil)
		if err != nil {
			t.Errorf("unexpected error for nil output: %v", err)
		}
	})

	t.Run("body read error on http error", func(t *testing.T) {
		// Create a mock client that returns an error status with unreadable body
		mockClient := &mockHTTPClientWithErrorBody{statusCode: http.StatusBadRequest}

		c := client.New("test-key", mockClient, client.WithBaseURL("https://example.com/"))
		svc := NewService(c)

		var resp DetailResponse
		err := svc.doGet(ctx, "property/detail", url.Values{}, &resp)
		if err == nil {
			t.Errorf("expected error for bad request with unreadable body")
		}
		if !strings.Contains(err.Error(), "unable to read error response") {
			t.Errorf("expected 'unable to read error response' in error, got %v", err)
		}
	})
}

// mockHTTPClientWithErrorBody returns responses with bodies that fail to read
type mockHTTPClientWithErrorBody struct {
	statusCode int
}

func (m *mockHTTPClientWithErrorBody) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       errorReader{},
		Header:     make(http.Header),
	}, nil
}

func TestGetPropertyIDValidation(t *testing.T) {
	ctx := context.Background()
	mock := &mockHTTPClient{t: t, responseBody: `{"status":{}}`}
	c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
	svc := NewService(c)

	t.Run("with address1 and address2", func(t *testing.T) {
		mock.expectedPath = "/propertyapi/v1.0.0/property/id"
		mock.expectedQuery = url.Values{
			"address1": {"123 Main St"},
			"address2": {"City, ST"},
		}
		mock.responseBody = `{"status":{},"identifier":[]}`

		_, err := svc.GetPropertyID(ctx, "", WithAddressLines("123 Main St", "City, ST"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing address", func(t *testing.T) {
		_, err := svc.GetPropertyID(ctx, "")
		if err == nil {
			t.Errorf("expected error for missing address")
		}
		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})
}

func TestGetPropertySnapshotValidation(t *testing.T) {
	ctx := context.Background()
	mock := &mockHTTPClient{t: t, responseBody: `{"status":{},"property":[]}`}
	c := client.New("test-key", mock, client.WithBaseURL("https://example.com/"))
	svc := NewService(c)

	t.Run("with postal code", func(t *testing.T) {
		mock.expectedPath = "/propertyapi/v1.0.0/property/snapshot"
		mock.expectedQuery = url.Values{"postalCode": {"12345"}}

		_, err := svc.GetPropertySnapshot(ctx, WithPostalCode("12345"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with lat/lon", func(t *testing.T) {
		mock.expectedQuery = url.Values{
			"latitude":  {"40.7128"},
			"longitude": {"-74.006"},
		}

		_, err := svc.GetPropertySnapshot(ctx, WithLatitudeLongitude(40.7128, -74.0060))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing required params", func(t *testing.T) {
		_, err := svc.GetPropertySnapshot(ctx)
		if err == nil {
			t.Errorf("expected error for missing params")
		}
		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})
}
