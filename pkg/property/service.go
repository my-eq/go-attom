package property

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/my-eq/go-attom/pkg/client"
)

// Service provides access to ATTOM Property API resources.
type Service struct {
	client *client.Client
}

// NewService constructs a Property API service using the provided ATTOM client.
func NewService(c *client.Client) *Service {
	if c == nil {
		return nil
	}
	return &Service{client: c}
}

// endpoint constants for Property API resources.
const (
	propertyBasePath         = "v4/property/"
	saleBasePath             = "v4/transaction/"
	assessmentBasePath       = "v4/property/"
	assessmentHistoryPath    = "v4/property/"
	avmBasePath              = "v4/property/"
	avmHistoryBasePath       = "v4/property/"
	attomAVMPath             = "v4/property/"
	valuationBasePath        = "v4/property/"
	salesHistoryBasePath     = "v4/transaction/"
	salesTrendBasePath       = "v4/transaction/"
	transactionTrendBasePath = "v4/transaction/"
	schoolBasePath           = "v4/school/"
	allEventsBasePath        = "v4/property/"
	saleComparablesBasePath  = "property/v2/salescomparables/"
	hazardBasePath           = "transportationnoise"
	enumerationsBasePath     = "v4/enumerations/"
	areaBasePath             = "v4/area/"
	poiBasePath              = "v4/neighborhood/poi"
	communityBasePath        = "v4/neighborhood/neighborhood/community"
	parcelTilesBasePath      = "v4/parceltiles/"
	preforeclosureBasePath   = "v3/preforeclosuredetails"
)

func (s *Service) ensureClient() error {
	if s == nil || s.client == nil {
		return fmt.Errorf("property: service client is not initialized")
	}
	return nil
}

func (s *Service) doGet(ctx context.Context, endpoint string, query url.Values, out interface{}) (err error) {
	if err = s.ensureClient(); err != nil {
		return err
	}
	var req *http.Request
	req, err = s.client.NewRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return fmt.Errorf("property: failed to build request: %w", err)
	}
	var resp *http.Response
	resp, err = s.client.DoRequest(req)
	if err != nil {
		return fmt.Errorf("property: request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("property: failed to close response body: %w", closeErr)
		}
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		rawBody, readErr := io.ReadAll(resp.Body)
		apiErr := &Error{StatusCode: resp.StatusCode, Body: rawBody}
		if readErr == nil && len(rawBody) > 0 {
			var statusWrapper struct {
				Status  *Status `json:"status,omitempty"`
				Message string  `json:"message,omitempty"`
			}
			if unmarshalErr := json.Unmarshal(rawBody, &statusWrapper); unmarshalErr == nil {
				apiErr.Status = statusWrapper.Status
				apiErr.Message = statusWrapper.Message
			}
		}
		if readErr != nil {
			return fmt.Errorf("property: unable to read error response: %w", readErr)
		}
		return apiErr
	}

	if out == nil {
		// Drain and discard the body when no output is needed
		if _, copyErr := io.Copy(io.Discard, resp.Body); copyErr != nil {
			return fmt.Errorf("property: failed to drain response body: %w", copyErr)
		}
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	if decodeErr := decoder.Decode(out); decodeErr != nil {
		return fmt.Errorf("property: failed to decode response: %w", decodeErr)
	}
	return err
}

func (s *Service) get(ctx context.Context, endpoint string, opts []Option, validator func(url.Values) error, out interface{}) error {
	query := applyOptions(opts)
	if validator != nil {
		if err := validator(query); err != nil {
			return err
		}
	}
	return s.doGet(ctx, endpoint, query, out)
}

func requireAny(values url.Values, keys ...string) error {
	for _, key := range keys {
		if v := values.Get(key); v != "" {
			return nil
		}
	}
	return fmt.Errorf("%w: expected one of %v", ErrMissingParameter, keys)
}

func requireAll(values url.Values, keys ...string) error {
	for _, key := range keys {
		if values.Get(key) == "" {
			return fmt.Errorf("%w: missing %s", ErrMissingParameter, key)
		}
	}
	return nil
}

func requirePropertyIdentifier(values url.Values) error {
	if values.Get("attomid") != "" || values.Get("id") != "" || values.Get("address") != "" || values.Get("address1") != "" {
		return nil
	}
	if values.Get("fips") != "" && values.Get("APN") != "" {
		return nil
	}
	return fmt.Errorf("%w: provide attomid, id, address, address1, or fips+APN", ErrMissingParameter)
}

func ensureGeoContext(values url.Values) error {
	if values.Get("address") != "" || values.Get("address1") != "" || (values.Get("latitude") != "" && values.Get("longitude") != "") {
		return nil
	}
	return fmt.Errorf("%w: provide address or latitude/longitude", ErrMissingParameter)
}

// GetPropertyID retrieves ATTOM property identifiers for a supplied address.
func (s *Service) GetPropertyID(ctx context.Context, address string, opts ...Option) (*IDResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp IDResponse
	err := s.get(ctx, propertyBasePath+"id", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		if values.Get("address1") != "" && values.Get("address2") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPropertyDetail retrieves detailed property information.
func (s *Service) GetPropertyDetail(ctx context.Context, opts ...Option) (*DetailResponse, error) {
	var resp DetailResponse
	err := s.get(ctx, propertyBasePath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPropertyAddress retrieves property address details by identifier.
func (s *Service) GetPropertyAddress(ctx context.Context, opts ...Option) (*AddressResponse, error) {
	var resp AddressResponse
	err := s.get(ctx, propertyBasePath+"address", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPropertySnapshot retrieves a lightweight property snapshot summary.
func (s *Service) GetPropertySnapshot(ctx context.Context, opts ...Option) (*SnapshotResponse, error) {
	validator := func(values url.Values) error {
		// attomId or attomid or id
		if values.Get("attomId") != "" || values.Get("attomid") != "" || values.Get("id") != "" {
			return nil
		}
		// FIPS + APN or apn
		if values.Get("fips") != "" && (values.Get("apn") != "" || values.Get("APN") != "") {
			return nil
		}
		// address (single line)
		if values.Get("address") != "" {
			return nil
		}
		// address1 + address2 (two lines)
		if values.Get("address1") != "" && values.Get("address2") != "" {
			return nil
		}
		// postalCode
		if values.Get("postalCode") != "" {
			return nil
		}
		// latitude + longitude (+ radius required)
		lat := values.Get("latitude")
		lon := values.Get("longitude")
		if lat != "" && lon != "" {
			if values.Get("radius") != "" {
				return nil
			}
			return fmt.Errorf("%w: radius required with latitude/longitude", ErrMissingParameter)
		}
		return fmt.Errorf("%w: valid property identifier required (attomId/attomid, id, FIPS+(APN/apn), address, address1/address2, postalCode, or latitude/longitude+radius)", ErrMissingParameter)
	}
	var resp SnapshotResponse
	err := s.get(ctx, propertyBasePath+"snapshot", opts, validator, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBasicProfile retrieves the basic property profile.
func (s *Service) GetBasicProfile(ctx context.Context, address string, opts ...Option) (*ProfileResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp ProfileResponse
	err := s.get(ctx, propertyBasePath+"basicprofile", allOpts, func(values url.Values) error {
		if values.Get("address") != "" || values.Get("address1") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetExpandedProfile retrieves the expanded property profile.
func (s *Service) GetExpandedProfile(ctx context.Context, opts ...Option) (*ProfileResponse, error) {
	var resp ProfileResponse
	err := s.get(ctx, propertyBasePath+"expandedprofile", opts, func(values url.Values) error {
		if requirePropertyIdentifier(values) == nil || values.Get("geoIdV4") != "" {
			return nil
		}
		return fmt.Errorf("%w: property identifier or geoIdV4 required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDetailWithSchools retrieves property detail including school information.
func (s *Service) GetDetailWithSchools(ctx context.Context, address string, opts ...Option) (*WithSchoolsResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp WithSchoolsResponse
	err := s.get(ctx, propertyBasePath+"detailwithschools", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDetailMortgage retrieves property detail with mortgage information.
func (s *Service) GetDetailMortgage(ctx context.Context, address string, opts ...Option) (*MortgageResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp MortgageResponse
	err := s.get(ctx, propertyBasePath+"detailmortgage", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDetailOwner retrieves property detail with owner information.
func (s *Service) GetDetailOwner(ctx context.Context, address string, opts ...Option) (*OwnerResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp OwnerResponse
	err := s.get(ctx, propertyBasePath+"detailowner", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDetailMortgageOwner retrieves property detail with mortgage and ownership information.
func (s *Service) GetDetailMortgageOwner(ctx context.Context, address string, opts ...Option) (*MortgageOwnerResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp MortgageOwnerResponse
	err := s.get(ctx, propertyBasePath+"detailmortgageowner", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBuildingPermits retrieves building permit records for a property.
func (s *Service) GetBuildingPermits(ctx context.Context, address string, opts ...Option) (*BuildingPermitsResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp BuildingPermitsResponse
	err := s.get(ctx, propertyBasePath+"buildingpermits", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSaleDetail retrieves sale detail information.
func (s *Service) GetSaleDetail(ctx context.Context, opts ...Option) (*SaleDetailResponse, error) {
	var resp SaleDetailResponse
	err := s.get(ctx, saleBasePath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSaleSnapshot retrieves sale snapshot information.
func (s *Service) GetSaleSnapshot(ctx context.Context, opts ...Option) (*SaleSnapshotResponse, error) {
	var resp SaleSnapshotResponse
	err := s.get(ctx, saleBasePath+"snapshot", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAssessmentDetail retrieves assessment detail information.
func (s *Service) GetAssessmentDetail(ctx context.Context, opts ...Option) (*AssessmentDetailResponse, error) {
	var resp AssessmentDetailResponse
	err := s.get(ctx, assessmentBasePath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAssessmentSnapshot retrieves assessment snapshot information.
func (s *Service) GetAssessmentSnapshot(ctx context.Context, opts ...Option) (*AssessmentSnapshotResponse, error) {
	var resp AssessmentSnapshotResponse
	err := s.get(ctx, assessmentBasePath+"snapshot", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAssessmentHistory retrieves historical assessment records.
func (s *Service) GetAssessmentHistory(ctx context.Context, opts ...Option) (*AssessmentHistoryResponse, error) {
	var resp AssessmentHistoryResponse
	err := s.get(ctx, assessmentHistoryPath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAVMSnapshot retrieves AVM snapshot values for a property.
func (s *Service) GetAVMSnapshot(ctx context.Context, opts ...Option) (*AVMSnapshotResponse, error) {
	var resp AVMSnapshotResponse
	err := s.get(ctx, avmBasePath+"snapshot", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAttomAVMDetail retrieves detailed ATTOM AVM information.
func (s *Service) GetAttomAVMDetail(ctx context.Context, opts ...Option) (*AttomAVMDetailResponse, error) {
	var resp AttomAVMDetailResponse
	err := s.get(ctx, attomAVMPath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAVMHistory retrieves historical AVM values.
func (s *Service) GetAVMHistory(ctx context.Context, opts ...Option) (*AVMHistoryResponse, error) {
	var resp AVMHistoryResponse
	err := s.get(ctx, avmHistoryBasePath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRentalAVM retrieves rental AVM valuation data.
func (s *Service) GetRentalAVM(ctx context.Context, opts ...Option) (*RentalAVMResponse, error) {
	var resp RentalAVMResponse
	err := s.get(ctx, valuationBasePath+"rentalavm", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSalesHistoryDetail retrieves detailed sales history data.
func (s *Service) GetSalesHistoryDetail(ctx context.Context, opts ...Option) (*SalesHistoryResponse, error) {
	var resp SalesHistoryResponse
	err := s.get(ctx, salesHistoryBasePath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSalesHistorySnapshot retrieves sales history snapshot data.
func (s *Service) GetSalesHistorySnapshot(ctx context.Context, opts ...Option) (*SalesHistoryResponse, error) {
	var resp SalesHistoryResponse
	err := s.get(ctx, salesHistoryBasePath+"snapshot", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSalesHistoryBasic retrieves the basic sales history data set.
func (s *Service) GetSalesHistoryBasic(ctx context.Context, opts ...Option) (*SalesHistoryResponse, error) {
	var resp SalesHistoryResponse
	err := s.get(ctx, salesHistoryBasePath+"basichistory", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSalesHistoryExpanded retrieves the expanded sales history data set.
func (s *Service) GetSalesHistoryExpanded(ctx context.Context, opts ...Option) (*SalesHistoryResponse, error) {
	var resp SalesHistoryResponse
	err := s.get(ctx, salesHistoryBasePath+"expandedhistory", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSalesTrendSnapshot retrieves geographic sales trend data.
func (s *Service) GetSalesTrendSnapshot(ctx context.Context, opts ...Option) (*SalesTrendSnapshotResponse, error) {
	var resp SalesTrendSnapshotResponse
	err := s.get(ctx, salesTrendBasePath+"snapshot", opts, func(values url.Values) error {
		if values.Get("geoIdV4") == "" {
			return fmt.Errorf("%w: geoIdV4 required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTransactionSalesTrend retrieves transaction-based sales trend data.
func (s *Service) GetTransactionSalesTrend(ctx context.Context, opts ...Option) (*TransactionSalesTrendResponse, error) {
	var resp TransactionSalesTrendResponse
	err := s.get(ctx, transactionTrendBasePath+"salestrend", opts, func(values url.Values) error {
		if values.Get("geoIdV4") == "" {
			return fmt.Errorf("%w: geoIdV4 required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchSchools locates schools near a given context.
func (s *Service) SearchSchools(ctx context.Context, opts ...Option) (*SchoolSearchResponse, error) {
	var resp SchoolSearchResponse
	err := s.get(ctx, schoolBasePath+"search", opts, ensureGeoContext, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSchoolProfile retrieves detailed school profile information.
func (s *Service) GetSchoolProfile(ctx context.Context, schoolID string, opts ...Option) (*SchoolProfileResponse, error) {
	allOpts := append([]Option{WithString("schoolId", schoolID)}, opts...)
	var resp SchoolProfileResponse
	err := s.get(ctx, schoolBasePath+"profile", allOpts, func(values url.Values) error {
		if values.Get("schoolId") == "" {
			return fmt.Errorf("%w: schoolId required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSchoolDistrict retrieves school district information.
func (s *Service) GetSchoolDistrict(ctx context.Context, address string, opts ...Option) (*SchoolDistrictResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp SchoolDistrictResponse
	err := s.get(ctx, schoolBasePath+"district", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSchoolDetailWithSchools retrieves property and associated school information.
func (s *Service) GetSchoolDetailWithSchools(ctx context.Context, address string, opts ...Option) (*SchoolDetailWithSchoolsResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp SchoolDetailWithSchoolsResponse
	err := s.get(ctx, schoolBasePath+"detailwithschools", allOpts, func(values url.Values) error {
		if values.Get("address") != "" {
			return nil
		}
		return fmt.Errorf("%w: address required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSchoolSnapshot retrieves schools within a defined radius from a point (deprecated endpoint).
func (s *Service) GetSchoolSnapshot(ctx context.Context, latitude, longitude, radius string, fileTypeText string, opts ...Option) (*SchoolSnapshotResponse, error) {
	allOpts := append([]Option{
		WithString("latitude", latitude),
		WithString("longitude", longitude),
		WithString("radius", radius),
	}, opts...)
	if fileTypeText != "" {
		allOpts = append(allOpts, WithString("filetypetext", fileTypeText))
	}
	var resp SchoolSnapshotResponse
	err := s.get(ctx, schoolBasePath+"snapshot", allOpts, func(values url.Values) error {
		if values.Get("latitude") != "" && values.Get("longitude") != "" && values.Get("radius") != "" {
			return nil
		}
		return fmt.Errorf("%w: latitude, longitude, and radius required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSchoolDetail retrieves details about a particular school (deprecated endpoint).
//
//nolint:dupl // similar code patterns shared across school endpoints
func (s *Service) GetSchoolDetail(ctx context.Context, schoolID string, opts ...Option) (*SchoolDetailResponse, error) {
	allOpts := append([]Option{WithString("id", schoolID)}, opts...)
	var resp SchoolDetailResponse
	err := s.get(ctx, schoolBasePath+"detail", allOpts, func(values url.Values) error {
		if values.Get("id") != "" {
			return nil
		}
		return fmt.Errorf("%w: school id required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSchoolDistrictDetail retrieves details about a particular school district (deprecated endpoint).
//
//nolint:dupl // similar code patterns shared across school endpoints
func (s *Service) GetSchoolDistrictDetail(ctx context.Context, districtID string, opts ...Option) (*SchoolDistrictDetailResponse, error) {
	allOpts := append([]Option{WithString("id", districtID)}, opts...)
	var resp SchoolDistrictDetailResponse
	err := s.get(ctx, schoolBasePath+"districtdetail", allOpts, func(values url.Values) error {
		if values.Get("id") != "" {
			return nil
		}
		return fmt.Errorf("%w: district id required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetHomeEquity retrieves estimated home equity for a property.
//
//nolint:dupl // pattern duplicated with other address-based endpoints
func (s *Service) GetHomeEquity(ctx context.Context, address1, address2 string, opts ...Option) (*HomeEquityResponse, error) {
	allOpts := append([]Option{
		WithString("address1", address1),
		WithString("address2", address2),
	}, opts...)
	var resp HomeEquityResponse
	err := s.get(ctx, valuationBasePath+"homeequity", allOpts, func(values url.Values) error {
		if values.Get("address1") != "" && values.Get("address2") != "" {
			return nil
		}
		return fmt.Errorf("%w: address1 and address2 required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAVMSnapshotGeo retrieves AVM snapshot values for all properties within a specific geography.
func (s *Service) GetAVMSnapshotGeo(ctx context.Context, geoIDV4, minAVMValue, maxAVMValue, propertyType string, opts ...Option) (*AVMSnapshotGeoResponse, error) {
	allOpts := append([]Option{WithString("geoIdV4", geoIDV4)}, opts...)
	if minAVMValue != "" {
		allOpts = append(allOpts, WithString("minavmvalue", minAVMValue))
	}
	if maxAVMValue != "" {
		allOpts = append(allOpts, WithString("maxavmvalue", maxAVMValue))
	}
	if propertyType != "" {
		allOpts = append(allOpts, WithString("propertytype", propertyType))
	}
	var resp AVMSnapshotGeoResponse
	err := s.get(ctx, avmBasePath+"snapshot", allOpts, func(values url.Values) error {
		if values.Get("geoIdV4") != "" {
			return nil
		}
		return fmt.Errorf("%w: geoIdV4 required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAVMHistoryByAddress retrieves AVM history for a property by address.
//
//nolint:dupl // pattern duplicated with other address-based endpoints
func (s *Service) GetAVMHistoryByAddress(ctx context.Context, address1, address2 string, opts ...Option) (*AVMHistoryResponse, error) {
	allOpts := append([]Option{
		WithString("address1", address1),
		WithString("address2", address2),
	}, opts...)
	var resp AVMHistoryResponse
	err := s.get(ctx, avmHistoryBasePath+"detail", allOpts, func(values url.Values) error {
		if values.Get("address1") != "" && values.Get("address2") != "" {
			return nil
		}
		return fmt.Errorf("%w: address1 and address2 required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAllEventsDetail retrieves all events information for a property.
func (s *Service) GetAllEventsDetail(ctx context.Context, opts ...Option) (*AllEventsDetailResponse, error) {
	var resp AllEventsDetailResponse
	err := s.get(ctx, allEventsBasePath+"detail", opts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAllEventsSnapshot retrieves a snapshot of events for a property.
func (s *Service) GetAllEventsSnapshot(ctx context.Context, address string, opts ...Option) (*AllEventsSnapshotResponse, error) {
	allOpts := append([]Option{WithAddress(address)}, opts...)
	var resp AllEventsSnapshotResponse
	err := s.get(ctx, allEventsBasePath+"snapshot", allOpts, func(values url.Values) error {
		if values.Get("address") == "" {
			return fmt.Errorf("%w: address required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetEnumerationsDetail retrieves enumerations detail information.
func (s *Service) GetEnumerationsDetail(ctx context.Context, opts ...Option) (*EnumerationsDetailResponse, error) {
	var resp EnumerationsDetailResponse
	err := s.get(ctx, enumerationsBasePath+"detail", opts, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBoundaryDetail retrieves boundary details for a geography.
func (s *Service) GetBoundaryDetail(ctx context.Context, geoID string, opts ...Option) (*BoundaryResponse, error) {
	allOpts := append([]Option{WithGeoIDV4(geoID)}, opts...)
	var resp BoundaryResponse
	err := s.get(ctx, areaBasePath+"boundary/detail", allOpts, func(values url.Values) error {
		if values.Get("geoIdV4") == "" {
			return fmt.Errorf("%w: geoIdV4 required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetHierarchyLookup retrieves all boundaries a point falls within.
func (s *Service) GetHierarchyLookup(ctx context.Context, wktString string, opts ...Option) (*HierarchyResponse, error) {
	allOpts := append([]Option{WithWKTString(wktString)}, opts...)
	var resp HierarchyResponse
	err := s.get(ctx, areaBasePath+"hierarchy/lookup", allOpts, func(values url.Values) error {
		if values.Get("WKTString") == "" {
			return fmt.Errorf("%w: WKTString required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCBSALookup retrieves all CBSAs within a state.
func (s *Service) GetCBSALookup(ctx context.Context, stateID string, opts ...Option) (*CBSAResponse, error) {
	allOpts := append([]Option{WithStateID(stateID)}, opts...)
	var resp CBSAResponse
	err := s.get(ctx, areaBasePath+"cbsa/lookup", allOpts, func(values url.Values) error {
		if values.Get("StateId") == "" {
			return fmt.Errorf("%w: StateId required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCountyLookup retrieves all counties within a state.
func (s *Service) GetCountyLookup(ctx context.Context, stateID string, opts ...Option) (*CountyResponse, error) {
	allOpts := append([]Option{WithStateID(stateID)}, opts...)
	var resp CountyResponse
	err := s.get(ctx, areaBasePath+"county/lookup", allOpts, func(values url.Values) error {
		if values.Get("StateId") == "" {
			return fmt.Errorf("%w: StateId required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStateLookup retrieves all states and their IDs.
func (s *Service) GetStateLookup(ctx context.Context, opts ...Option) (*StateResponse, error) {
	var resp StateResponse
	err := s.get(ctx, areaBasePath+"state/lookup", opts, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGeoIDLookup retrieves specific Geo IDs that exist within a specified Geo ID.
func (s *Service) GetGeoIDLookup(ctx context.Context, geoID string, opts ...Option) (*GeoidResponse, error) {
	allOpts := append([]Option{WithGeoIDV4(geoID)}, opts...)
	var resp GeoidResponse
	err := s.get(ctx, areaBasePath+"geoid/lookup/", allOpts, func(values url.Values) error {
		if values.Get("geoIdV4") == "" {
			return fmt.Errorf("%w: geoIdV4 required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetGeoIDLegacyLookup retrieves a translation between legacy codes and new geography identifiers.
func (s *Service) GetGeoIDLegacyLookup(ctx context.Context, geoID string, opts ...Option) (*LegacyGeoidResponse, error) {
	allOpts := append([]Option{WithGeoIDV4(geoID)}, opts...)
	var resp LegacyGeoidResponse
	err := s.get(ctx, areaBasePath+"geoid/legacyLookup/", allOpts, func(values url.Values) error {
		if values.Get("geoIdV4") == "" {
			return fmt.Errorf("%w: geoIdV4 required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPOI retrieves points of interest near a location.
func (s *Service) GetPOI(ctx context.Context, opts ...Option) (*POIResponse, error) {
	var resp POIResponse
	err := s.get(ctx, poiBasePath, opts, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPOICategoryLookup retrieves values used for category, lob, industry.
func (s *Service) GetPOICategoryLookup(ctx context.Context, opts ...Option) (*POICategoryResponse, error) {
	var resp POICategoryResponse
	err := s.get(ctx, poiBasePath+"categorylookup", opts, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCommunity retrieves neighborhood community information.
func (s *Service) GetCommunity(ctx context.Context, opts ...Option) (*CommunityResponse, error) {
	var resp CommunityResponse
	err := s.get(ctx, communityBasePath, opts, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLocationLookup retrieves location lookup information.
func (s *Service) GetLocationLookup(ctx context.Context, opts ...Option) (*LocationLookupResponse, error) {
	var resp LocationLookupResponse
	err := s.get(ctx, communityBasePath+"location/lookup", opts, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSaleComparablesByAddress retrieves sale comparables by address.
func (s *Service) GetSaleComparablesByAddress(ctx context.Context, street, city, county, state, zip string, opts ...Option) (*SaleComparablesResponse, error) {
	allOpts := append([]Option{WithAddress(fmt.Sprintf("%s, %s, %s, %s %s", street, city, county, state, zip))}, opts...)
	var resp SaleComparablesResponse
	err := s.get(ctx, fmt.Sprintf("%saddress/%s/%s/%s/%s/%s", saleComparablesBasePath, url.PathEscape(street), url.PathEscape(city), url.PathEscape(county), url.PathEscape(state), url.PathEscape(zip)), allOpts, func(values url.Values) error {
		if street != "" && city != "" && county != "" && state != "" && zip != "" {
			return nil
		}
		return fmt.Errorf("%w: street, city, county, state, and zip required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSaleComparablesByAPN retrieves sale comparables by APN.
func (s *Service) GetSaleComparablesByAPN(ctx context.Context, apn, county, state string, opts ...Option) (*SaleComparablesResponse, error) {
	allOpts := append([]Option{WithAPN(apn)}, opts...)
	var resp SaleComparablesResponse
	err := s.get(ctx, fmt.Sprintf("%sapn/%s/%s/%s", saleComparablesBasePath, url.PathEscape(apn), url.PathEscape(county), url.PathEscape(state)), allOpts, func(values url.Values) error {
		if values.Get("APN") != "" && county != "" && state != "" {
			return nil
		}
		return fmt.Errorf("%w: APN, county, and state required", ErrMissingParameter)
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSaleComparablesByPropID retrieves sale comparables by property ID.
func (s *Service) GetSaleComparablesByPropID(ctx context.Context, propID string, opts ...Option) (*SaleComparablesResponse, error) {
	allOpts := append([]Option{WithAttomID(propID)}, opts...)
	var resp SaleComparablesResponse
	err := s.get(ctx, saleComparablesBasePath+"propid/"+propID, allOpts, requirePropertyIdentifier, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTransportationNoise retrieves transportation noise information.
func (s *Service) GetTransportationNoise(ctx context.Context, attomID string, opts ...Option) (*TransportationNoiseResponse, error) {
	allOpts := append([]Option{WithAttomID(attomID)}, opts...)
	var resp TransportationNoiseResponse
	err := s.get(ctx, hazardBasePath, allOpts, func(values url.Values) error {
		if values.Get("attomid") == "" {
			return fmt.Errorf("%w: attomid required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetParcelTiles retrieves parcel tiles data.
func (s *Service) GetParcelTiles(ctx context.Context, z, x, y int, format string, opts ...Option) (*ParcelTilesResponse, error) {
	var resp ParcelTilesResponse
	endpoint := fmt.Sprintf("%s%d/%d/%d.%s", parcelTilesBasePath, z, x, y, format)
	err := s.get(ctx, endpoint, opts, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPreforeclosureDetails retrieves pre-foreclosure details for a property.
func (s *Service) GetPreforeclosureDetails(ctx context.Context, attomID string, opts ...Option) (*PreforeclosureResponse, error) {
	allOpts := append([]Option{WithAttomID(attomID)}, opts...)
	var resp PreforeclosureResponse
	err := s.get(ctx, preforeclosureBasePath, allOpts, func(values url.Values) error {
		if values.Get("attomid") == "" {
			return fmt.Errorf("%w: attomid required", ErrMissingParameter)
		}
		return nil
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
