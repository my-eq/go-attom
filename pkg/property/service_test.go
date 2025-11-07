package property

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/my-eq/go-attom/pkg/client"
)

// Endpoint tests have been migrated to domain-specific test files: service_property_test.go, service_school_test.go, service_avm_test.go, service_sales_test.go

func TestServiceErrorResponse(t *testing.T) {
	ctx := context.Background()
	mock := &mockHTTPClient{
		t:              t,
		expectedMethod: http.MethodGet,
		expectedPath:   "/v4/property/detail",
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

func TestErrorTypes(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		var e *Error
		got := e.Error()
		if got != "property: nil error" {
			t.Errorf("expected 'property: nil error', got %q", got)
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
		WithRadius(5)(vals)
		if vals.Get("radius") != "5" {
			t.Errorf("expected '5', got %q", vals.Get("radius"))
		}
	})
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

// mockHTTPClientWithErrorBody returns responses with bodies that fail to read.
type mockHTTPClientWithErrorBody struct {
	statusCode int
}

// errorReader implements io.ReadCloser and always returns an error.
type errorReader struct{}

func (e errorReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("mock read error")
}
func (e errorReader) Close() error { return nil }

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
		mock.expectedPath = "/v4/property/id"
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
		mock.expectedPath = "/v4/property/snapshot"
		mock.expectedQuery = url.Values{"postalCode": {"12345"}}

		_, err := svc.GetPropertySnapshot(ctx, WithPostalCode("12345"))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with lat/lon and radius", func(t *testing.T) {
		mock.expectedQuery = url.Values{
			"latitude":  {"40.7128"},
			"longitude": {"-74.006"},
			"radius":    {"5"},
		}
		_, err := svc.GetPropertySnapshot(ctx, WithLatitudeLongitude(40.7128, -74.0060), WithRadius(5))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with lat/lon missing radius", func(t *testing.T) {
		mock.expectedQuery = url.Values{
			"latitude":  {"40.7128"},
			"longitude": {"-74.006"},
		}
		_, err := svc.GetPropertySnapshot(ctx, WithLatitudeLongitude(40.7128, -74.0060))
		if err == nil {
			t.Errorf("expected error for missing radius with lat/lon")
		}
		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
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
