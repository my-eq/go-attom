package property

import "encoding/json"

// Status describes the standard ATTOM response status block.
type Status struct {
	Version  *string `json:"version,omitempty"`
	Code     *int    `json:"code,omitempty"`
	Msg      *string `json:"msg,omitempty"`
	Total    *int    `json:"total,omitempty"`
	Page     *int    `json:"page,omitempty"`
	PageSize *int    `json:"pagesize,omitempty"`
}

// Identifier contains core identifiers for a property record.
type Identifier struct {
	AttomID  *string `json:"attomId,omitempty"`
	ID       *string `json:"id,omitempty"`
	FIPS     *string `json:"fips,omitempty"`
	APN      *string `json:"apn,omitempty"`
	ObPropID *string `json:"obPropId,omitempty"`
}

// Address represents a postal address and geographic coordinates.
type Address struct {
	Line1      *string  `json:"line1,omitempty"`
	Line2      *string  `json:"line2,omitempty"`
	City       *string  `json:"city,omitempty"`
	State      *string  `json:"state,omitempty"`
	County     *string  `json:"county,omitempty"`
	Country    *string  `json:"country,omitempty"`
	PostalCode *string  `json:"postalCode,omitempty"`
	UnitNumber *string  `json:"unitNumber,omitempty"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
}

// GeoLocation captures latitude and longitude alongside precision metadata.
type GeoLocation struct {
	Latitude  *float64 `json:"lat,omitempty"`
	Longitude *float64 `json:"lon,omitempty"`
	MatchCode *string  `json:"matchCode,omitempty"`
	Quality   *string  `json:"quality,omitempty"`
}

// Lot describes lot-specific attributes for a property.
type Lot struct {
	Acres          *float64 `json:"acres,omitempty"`
	Depth          *float64 `json:"depth,omitempty"`
	Frontage       *float64 `json:"frontage,omitempty"`
	AreaSquareFeet *float64 `json:"areaSqFt,omitempty"`
	LotNumber      *string  `json:"lotNumber,omitempty"`
	Range          *string  `json:"range,omitempty"`
	Section        *string  `json:"section,omitempty"`
	Township       *string  `json:"township,omitempty"`
	Shape          *string  `json:"shape,omitempty"`
	Zoning         *string  `json:"zoning,omitempty"`
	Pool           *string  `json:"pool,omitempty"`
}

// Summary provides high-level information about a property.
type Summary struct {
	PropertyType            *string  `json:"propertyType,omitempty"`
	PropertyTypeDescription *string  `json:"propertyTypeDescription,omitempty"`
	YearBuilt               *int     `json:"yearBuilt,omitempty"`
	EffectiveYearBuilt      *int     `json:"effectiveYearBuilt,omitempty"`
	Stories                 *float64 `json:"stories,omitempty"`
	UnitsCount              *int     `json:"unitsCount,omitempty"`
	LegalDescription        *string  `json:"legalDescription,omitempty"`
	PropertyIndicator       *int     `json:"propertyIndicator,omitempty"`
}

// Building describes structure-level detail.
type Building struct {
	Construction *Construction    `json:"construction,omitempty"`
	Rooms        *Rooms           `json:"rooms,omitempty"`
	Area         *BuildingArea    `json:"area,omitempty"`
	Interior     *Interior        `json:"interior,omitempty"`
	Exterior     *Exterior        `json:"exterior,omitempty"`
	Summary      *BuildingSummary `json:"summary,omitempty"`
}

// Construction captures construction-specific information.
type Construction struct {
	FrameType        *string `json:"frameType,omitempty"`
	Foundation       *string `json:"foundation,omitempty"`
	RoofCover        *string `json:"roofCover,omitempty"`
	RoofType         *string `json:"roofType,omitempty"`
	WallType         *string `json:"wallType,omitempty"`
	FloorType        *string `json:"floorType,omitempty"`
	CoolingType      *string `json:"coolingType,omitempty"`
	HeatingType      *string `json:"heatingType,omitempty"`
	ConstructionType *string `json:"constructionType,omitempty"`
}

// Rooms captures bedroom and bathroom counts.
type Rooms struct {
	TotalRooms        *int     `json:"totalRooms,omitempty"`
	Beds              *int     `json:"beds,omitempty"`
	BathsFull         *int     `json:"bathsFull,omitempty"`
	BathsHalf         *int     `json:"bathsHalf,omitempty"`
	BathsThreeQuarter *int     `json:"bathsThreeQuarter,omitempty"`
	BathsTotal        *float64 `json:"bathsTotal,omitempty"`
}

// BuildingArea stores various square footage measurements.
type BuildingArea struct {
	LivingSquareFeet   *int `json:"livingSqFt,omitempty"`
	TotalSquareFeet    *int `json:"totalSqFt,omitempty"`
	GarageSquareFeet   *int `json:"garageSqFt,omitempty"`
	BasementSquareFeet *int `json:"basementSqFt,omitempty"`
	AtticSquareFeet    *int `json:"atticSqFt,omitempty"`
}

// Interior captures interior attributes such as fireplaces.
type Interior struct {
	FireplaceCount *int    `json:"fireplaceCount,omitempty"`
	FlooringType   *string `json:"flooringType,omitempty"`
	Laundry        *string `json:"laundry,omitempty"`
}

// Exterior holds exterior feature information.
type Exterior struct {
	GarageType    *string `json:"garageType,omitempty"`
	ParkingSpaces *int    `json:"parkingSpaces,omitempty"`
	PorchType     *string `json:"porchType,omitempty"`
	PatioType     *string `json:"patioType,omitempty"`
}

// BuildingSummary collates additional building-level metrics.
type BuildingSummary struct {
	Quality            *string `json:"quality,omitempty"`
	Condition          *string `json:"condition,omitempty"`
	ArchitecturalStyle *string `json:"style,omitempty"`
	PropClass          *string `json:"propClass,omitempty"`
}

// Assessment represents property tax assessment information.
type Assessment struct {
	AssessedTotalValue       *float64 `json:"assdTtlValue,omitempty"`
	AssessedLandValue        *float64 `json:"assdLandValue,omitempty"`
	AssessedImprovementValue *float64 `json:"assdImpValue,omitempty"`
	MarketTotalValue         *float64 `json:"mktTtlValue,omitempty"`
	MarketLandValue          *float64 `json:"mktLandValue,omitempty"`
	MarketImprovementValue   *float64 `json:"mktImpValue,omitempty"`
	TaxAmount                *float64 `json:"taxAmt,omitempty"`
	TaxYear                  *int     `json:"taxYear,omitempty"`
	TaxRate                  *float64 `json:"taxRate,omitempty"`
	AppraisedValue           *float64 `json:"apprsdTotValue,omitempty"`
}

// AssessmentHistoryRecord contains historical assessment entries.
type AssessmentHistoryRecord struct {
	CalendarYear  *int     `json:"calendarYear,omitempty"`
	AssessedValue *float64 `json:"assdTtlValue,omitempty"`
	TaxAmount     *float64 `json:"taxAmt,omitempty"`
}

// Sale represents a single sale transaction for a property.
type Sale struct {
	SaleDate        *string  `json:"saleDate,omitempty"`
	SaleSearchDate  *string  `json:"saleSearchDate,omitempty"`
	RecordingDate   *string  `json:"recordingDate,omitempty"`
	Amount          *float64 `json:"amount,omitempty"`
	DocumentType    *string  `json:"documentType,omitempty"`
	DocumentNumber  *string  `json:"documentNumber,omitempty"`
	TransactionType *string  `json:"transactionType,omitempty"`
	BuyerName       *string  `json:"buyerName,omitempty"`
	SellerName      *string  `json:"sellerName,omitempty"`
}

// SalesHistoryRecord contains historical sales entries.
type SalesHistoryRecord struct {
	SaleDate       *string  `json:"saleDate,omitempty"`
	SaleAmount     *float64 `json:"saleAmount,omitempty"`
	DocumentType   *string  `json:"documentType,omitempty"`
	DocumentNumber *string  `json:"documentNumber,omitempty"`
	RecordingDate  *string  `json:"recordingDate,omitempty"`
}

// AVM contains automated valuation model data.
type AVM struct {
	Value      *float64 `json:"value,omitempty"`
	High       *float64 `json:"high,omitempty"`
	Low        *float64 `json:"low,omitempty"`
	Percentile *float64 `json:"percentile,omitempty"`
	Score      *float64 `json:"score,omitempty"`
	Confidence *string  `json:"confidence,omitempty"`
	Updated    *string  `json:"updated,omitempty"`
}

// AVMHistoryRecord describes valuation history entries.
type AVMHistoryRecord struct {
	Date  *string  `json:"date,omitempty"`
	Value *float64 `json:"value,omitempty"`
	High  *float64 `json:"high,omitempty"`
	Low   *float64 `json:"low,omitempty"`
}

// RentalAVM represents rental valuation output.
type RentalAVM struct {
	Value       *float64 `json:"value,omitempty"`
	Confidence  *string  `json:"confidence,omitempty"`
	UpdatedDate *string  `json:"updatedDate,omitempty"`
}

// Mortgage contains mortgage-related details for a property.
type Mortgage struct {
	LenderName    *string  `json:"lenderName,omitempty"`
	LoanType      *string  `json:"loanType,omitempty"`
	LoanAmount    *float64 `json:"loanAmount,omitempty"`
	LoanDate      *string  `json:"loanDate,omitempty"`
	InterestRate  *float64 `json:"interestRate,omitempty"`
	MaturityDate  *string  `json:"maturityDate,omitempty"`
	DueDate       *string  `json:"dueDate,omitempty"`
	RecordingDate *string  `json:"recordingDate,omitempty"`
	LoanNumber    *string  `json:"loanNumber,omitempty"`
	MortgageType  *string  `json:"mortgageType,omitempty"`
}

// Ownership represents owner information for a property.
type Ownership struct {
	OwnerType       *string  `json:"ownerType,omitempty"`
	Owner1FirstName *string  `json:"owner1FirstName,omitempty"`
	Owner1LastName  *string  `json:"owner1LastName,omitempty"`
	Owner2FirstName *string  `json:"owner2FirstName,omitempty"`
	Owner2LastName  *string  `json:"owner2LastName,omitempty"`
	MailingAddress  *Address `json:"mailingAddress,omitempty"`
	OccupancyStatus *string  `json:"occupancyStatus,omitempty"`
}

// Tax captures current tax data for a property.
type Tax struct {
	PaidAmount *float64 `json:"paidAmount,omitempty"`
	TaxYear    *int     `json:"taxYear,omitempty"`
	Delinquent *bool    `json:"delinquent,omitempty"`
}

// BuildingPermit represents a single permit record associated with a property.
type BuildingPermit struct {
	PermitNumber *string  `json:"permitNumber,omitempty"`
	PermitType   *string  `json:"permitType,omitempty"`
	PermitDate   *string  `json:"permitDate,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Contractor   *string  `json:"contractor,omitempty"`
	Value        *float64 `json:"value,omitempty"`
}

// School summarizes a school entity used within school endpoints.
type School struct {
	SchoolID        *string        `json:"schoolId,omitempty"`
	Name            *string        `json:"name,omitempty"`
	Type            *string        `json:"type,omitempty"`
	GradeLow        *string        `json:"gradeLow,omitempty"`
	GradeHigh       *string        `json:"gradeHigh,omitempty"`
	Enrollment      *int           `json:"enrollment,omitempty"`
	Phone           *string        `json:"phone,omitempty"`
	DistanceInMiles *float64       `json:"distanceInMiles,omitempty"`
	Address         *Address       `json:"address,omitempty"`
	Ratings         *SchoolRatings `json:"ratings,omitempty"`
}

// SchoolRatings holds rating information for a school.
type SchoolRatings struct {
	Overall *float64 `json:"overall,omitempty"`
	Test    *float64 `json:"test,omitempty"`
	Equity  *float64 `json:"equity,omitempty"`
}

// SchoolDistrict represents school district details.
type SchoolDistrict struct {
	DistrictID *string `json:"districtId,omitempty"`
	Name       *string `json:"name,omitempty"`
	Type       *string `json:"type,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Enrollment *int    `json:"enrollment,omitempty"`
}

// SalesTrendRecord represents a trend datapoint for a given period.
type SalesTrendRecord struct {
	GeoID      *string  `json:"geoId,omitempty"`
	GeoIDV4    *string  `json:"geoIdV4,omitempty"`
	Period     *string  `json:"periodDate,omitempty"`
	Interval   *string  `json:"interval,omitempty"`
	AvgSaleAmt *float64 `json:"avgSaleAmt,omitempty"`
	MedSaleAmt *float64 `json:"medSaleAmt,omitempty"`
	SaleCount  *int     `json:"saleCount,omitempty"`
}

// AllEventsRecord aggregates cross-domain events for a property.
type AllEventsRecord struct {
	EventType *string         `json:"eventType,omitempty"`
	EventDate *string         `json:"eventDate,omitempty"`
	Raw       json.RawMessage `json:"raw,omitempty"`
}

// Property encapsulates the full property data structure.
type Property struct {
	Identifier *Identifier  `json:"identifier,omitempty"`
	Address    *Address     `json:"address,omitempty"`
	Location   *GeoLocation `json:"location,omitempty"`
	Lot        *Lot         `json:"lot,omitempty"`
	Summary    *Summary     `json:"summary,omitempty"`
	Building   *Building    `json:"building,omitempty"`
	Assessment *Assessment  `json:"assessment,omitempty"`
	Sale       *Sale        `json:"sale,omitempty"`
	AVM        *AVM         `json:"avm,omitempty"`
	Mortgage   []Mortgage   `json:"mortgage,omitempty"`
	Ownership  *Ownership   `json:"ownership,omitempty"`
	Tax        *Tax         `json:"tax,omitempty"`
	Schools    []School     `json:"schools,omitempty"`
}

// IDResponse wraps the /property/id endpoint response.
type IDResponse struct {
	Status     *Status       `json:"status,omitempty"`
	Identifier []*Identifier `json:"identifier,omitempty"`
}

// DetailResponse wraps detailed property data.
type DetailResponse struct {
	Status   *Status     `json:"status,omitempty"`
	Property []*Property `json:"property,omitempty"`
}

// AddressResponse wraps address-only responses.
type AddressResponse struct {
	Status   *Status     `json:"status,omitempty"`
	Property []*Property `json:"property,omitempty"`
}

// SnapshotResponse provides lightweight property summaries.
type SnapshotResponse struct {
	Status   *Status     `json:"status,omitempty"`
	Property []*Property `json:"property,omitempty"`
}

// ProfileResponse contains profile data (basic/expanded).
type ProfileResponse struct {
	Status   *Status     `json:"status,omitempty"`
	Property []*Property `json:"property,omitempty"`
}

// WithSchoolsResponse extends property data with school assignments.
type WithSchoolsResponse struct {
	Status   *Status     `json:"status,omitempty"`
	Property []*Property `json:"property,omitempty"`
	Schools  []*School   `json:"school,omitempty"`
}

// MortgageResponse extends property data with mortgage information.
type MortgageResponse struct {
	Status   *Status     `json:"status,omitempty"`
	Property []*Property `json:"property,omitempty"`
	Mortgage []*Mortgage `json:"mortgage,omitempty"`
}

// OwnerResponse extends property data with ownership information.
type OwnerResponse struct {
	Status   *Status      `json:"status,omitempty"`
	Property []*Property  `json:"property,omitempty"`
	Owners   []*Ownership `json:"owner,omitempty"`
}

// MortgageOwnerResponse combines property, mortgage, and owner data.
type MortgageOwnerResponse struct {
	Status   *Status      `json:"status,omitempty"`
	Property []*Property  `json:"property,omitempty"`
	Mortgage []*Mortgage  `json:"mortgage,omitempty"`
	Owners   []*Ownership `json:"owner,omitempty"`
}

// BuildingPermitsResponse wraps permit data.
type BuildingPermitsResponse struct {
	Status  *Status           `json:"status,omitempty"`
	Permits []*BuildingPermit `json:"buildingPermit,omitempty"`
}

// SaleDetailResponse wraps sale detail data.
type SaleDetailResponse struct {
	Status *Status `json:"status,omitempty"`
	Sale   []*Sale `json:"sale,omitempty"`
}

// SaleSnapshotResponse wraps sale snapshot data.
type SaleSnapshotResponse struct {
	Status *Status `json:"status,omitempty"`
	Sale   []*Sale `json:"sale,omitempty"`
}

// AssessmentDetailResponse wraps assessment detail data.
type AssessmentDetailResponse struct {
	Status     *Status       `json:"status,omitempty"`
	Assessment []*Assessment `json:"assessment,omitempty"`
}

// AssessmentSnapshotResponse wraps snapshot-level assessment data.
type AssessmentSnapshotResponse struct {
	Status     *Status       `json:"status,omitempty"`
	Assessment []*Assessment `json:"assessment,omitempty"`
}

// AssessmentHistoryResponse wraps historical assessment data.
type AssessmentHistoryResponse struct {
	Status  *Status                    `json:"status,omitempty"`
	History []*AssessmentHistoryRecord `json:"assessmentHistory,omitempty"`
}

// AVMSnapshotResponse wraps AVM snapshot data.
type AVMSnapshotResponse struct {
	Status *Status `json:"status,omitempty"`
	AVM    []*AVM  `json:"avm,omitempty"`
}

// AttomAVMDetailResponse wraps ATTOM AVM detail data.
type AttomAVMDetailResponse struct {
	Status *Status `json:"status,omitempty"`
	AVM    []*AVM  `json:"attomAvm,omitempty"`
}

// AVMHistoryResponse wraps AVM history data.
type AVMHistoryResponse struct {
	Status  *Status             `json:"status,omitempty"`
	History []*AVMHistoryRecord `json:"avmHistory,omitempty"`
}

// RentalAVMResponse wraps rental AVM data.
type RentalAVMResponse struct {
	Status *Status      `json:"status,omitempty"`
	Rental []*RentalAVM `json:"rentalAvm,omitempty"`
}

// SalesHistoryResponse provides general sales history data.
type SalesHistoryResponse struct {
	Status *Status               `json:"status,omitempty"`
	Sales  []*SalesHistoryRecord `json:"salesHistory,omitempty"`
}

// SalesTrendSnapshotResponse wraps snapshot trend data.
type SalesTrendSnapshotResponse struct {
	Status *Status             `json:"status,omitempty"`
	Trends []*SalesTrendRecord `json:"salesTrend,omitempty"`
}

// TransactionSalesTrendResponse wraps transaction trend data.
type TransactionSalesTrendResponse struct {
	Status *Status             `json:"status,omitempty"`
	Trends []*SalesTrendRecord `json:"transactionTrend,omitempty"`
}

// SchoolSearchResponse wraps school search results.
type SchoolSearchResponse struct {
	Status *Status   `json:"status,omitempty"`
	School []*School `json:"school,omitempty"`
}

// SchoolProfileResponse wraps school profile data.
type SchoolProfileResponse struct {
	Status *Status   `json:"status,omitempty"`
	School []*School `json:"school,omitempty"`
}

// SchoolDistrictResponse wraps district data.
type SchoolDistrictResponse struct {
	Status   *Status           `json:"status,omitempty"`
	District []*SchoolDistrict `json:"district,omitempty"`
}

// SchoolDetailWithSchoolsResponse wraps property with schools detail.
type SchoolDetailWithSchoolsResponse struct {
	Status   *Status     `json:"status,omitempty"`
	Property []*Property `json:"property,omitempty"`
	Schools  []*School   `json:"school,omitempty"`
}

// SchoolSnapshotResponse wraps /school/snapshot endpoint results.
type SchoolSnapshotResponse struct {
	Status *Status   `json:"status,omitempty"`
	School []*School `json:"school,omitempty"`
}

// SchoolDetailResponse wraps /school/detail endpoint results.
type SchoolDetailResponse struct {
	Status *Status   `json:"status,omitempty"`
	School []*School `json:"school,omitempty"`
}

// SchoolDistrictDetailResponse wraps /school/districtdetail endpoint results.
type SchoolDistrictDetailResponse struct {
	Status   *Status           `json:"status,omitempty"`
	District []*SchoolDistrict `json:"district,omitempty"`
}

// HomeEquityResponse wraps /valuation/homeequity endpoint results.
type HomeEquityResponse struct {
	HomeEquity *float64    `json:"homeEquity,omitempty"`
	Status     *Status     `json:"status,omitempty"`
	Property   []*Property `json:"property,omitempty"`
}

// AVMSnapshotGeoResponse wraps /avm/snapshot geoIdV4 endpoint results.
type AVMSnapshotGeoResponse struct {
	Status *Status `json:"status,omitempty"`
	AVM    []*AVM  `json:"avm,omitempty"`
}

// AllEventsDetailResponse wraps all events data for a property.
type AllEventsDetailResponse struct {
	Status *Status            `json:"status,omitempty"`
	Events []*AllEventsRecord `json:"event,omitempty"`
}

// AllEventsSnapshotResponse wraps snapshot of all events data for a property.
type AllEventsSnapshotResponse struct {
	Status   *Status              `json:"status,omitempty"`
	Snapshot []*AllEventsSnapshot `json:"snapshot,omitempty"`
}

// AllEventsSnapshot represents snapshot event data.
type AllEventsSnapshot struct {
	PropertyID *string `json:"propertyId,omitempty"`
	Address    *string `json:"address,omitempty"`
	EventCount *int    `json:"eventCount,omitempty"`
	LastEvent  *string `json:"lastEvent,omitempty"`
}

// EnumerationsDetail represents enumeration detail data.
type EnumerationsDetail struct {
	Field *string `json:"field,omitempty"`
	Value *string `json:"value,omitempty"`
}

// EnumerationsDetailResponse wraps enumerations detail data.
type EnumerationsDetailResponse struct {
	Status       *Status               `json:"status,omitempty"`
	Enumerations []*EnumerationsDetail `json:"enumeration,omitempty"`
}

// BoundaryResponse wraps area boundary detail data.
type BoundaryResponse struct {
	Status   *Status   `json:"status,omitempty"`
	Boundary *Boundary `json:"boundary,omitempty"`
}

// Boundary represents geographic boundary data.
type Boundary struct {
	GeoID    *string   `json:"geoId,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Type     *string   `json:"type,omitempty"`
	Geometry *Geometry `json:"geometry,omitempty"`
}

// Geometry represents geometric data for boundaries.
type Geometry struct {
	Type        *string     `json:"type,omitempty"`
	Coordinates interface{} `json:"coordinates,omitempty"` // Can be various geometry types
}

// HierarchyResponse wraps hierarchy lookup data.
type HierarchyResponse struct {
	Status    *Status      `json:"status,omitempty"`
	Hierarchy []*Hierarchy `json:"hierarchy,omitempty"`
}

// Hierarchy represents hierarchical geographic data.
type Hierarchy struct {
	GeoID *string `json:"geoId,omitempty"`
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
	Level *string `json:"level,omitempty"`
}

// CBSAResponse wraps CBSA lookup data.
type CBSAResponse struct {
	Status *Status `json:"status,omitempty"`
	CBSA   []*CBSA `json:"cbsa,omitempty"`
}

// CBSA represents Core Based Statistical Area data.
type CBSA struct {
	GeoID     *string `json:"geoId,omitempty"`
	Name      *string `json:"name,omitempty"`
	Type      *string `json:"type,omitempty"`
	StateCode *string `json:"stateCode,omitempty"`
}

// CountyResponse wraps county lookup data.
type CountyResponse struct {
	Status   *Status   `json:"status,omitempty"`
	Counties []*County `json:"county,omitempty"`
}

// County represents county data.
type County struct {
	GeoID     *string `json:"geoId,omitempty"`
	Name      *string `json:"name,omitempty"`
	StateCode *string `json:"stateCode,omitempty"`
	FIPS      *string `json:"fips,omitempty"`
}

// StateResponse wraps state lookup data.
type StateResponse struct {
	Status *Status  `json:"status,omitempty"`
	States []*State `json:"state,omitempty"`
}

// State represents state data.
type State struct {
	GeoID *string `json:"geoId,omitempty"`
	Name  *string `json:"name,omitempty"`
	Code  *string `json:"code,omitempty"`
}

// GeoidResponse wraps geoid lookup data.
type GeoidResponse struct {
	Status *Status  `json:"status,omitempty"`
	Geoids []*Geoid `json:"geoid,omitempty"`
}

// Geoid represents geoid data.
type Geoid struct {
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
	Level *string `json:"level,omitempty"`
}

// LegacyGeoidResponse wraps legacy geoid lookup data.
type LegacyGeoidResponse struct {
	Status       *Status        `json:"status,omitempty"`
	LegacyGeoids []*LegacyGeoid `json:"legacyGeoid,omitempty"`
}

// LegacyGeoid represents legacy geoid data.
type LegacyGeoid struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

// POIResponse wraps point of interest data.
type POIResponse struct {
	Status *Status `json:"status,omitempty"`
	POIs   []*POI  `json:"poi,omitempty"`
}

// POI represents point of interest data.
type POI struct {
	ID          *string      `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Category    *string      `json:"category,omitempty"`
	Address     *Address     `json:"address,omitempty"`
	GeoLocation *GeoLocation `json:"geoLocation,omitempty"`
	Distance    *float64     `json:"distance,omitempty"`
}

// POICategoryResponse wraps POI category lookup data.
type POICategoryResponse struct {
	Status     *Status        `json:"status,omitempty"`
	Categories []*POICategory `json:"category,omitempty"`
}

// POICategory represents POI category data.
type POICategory struct {
	ID          *string `json:"id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// CommunityResponse wraps neighborhood community data.
type CommunityResponse struct {
	Status      *Status      `json:"status,omitempty"`
	Communities []*Community `json:"community,omitempty"`
}

// Community represents neighborhood community data.
type Community struct {
	ID          *string      `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Type        *string      `json:"type,omitempty"`
	Description *string      `json:"description,omitempty"`
	GeoLocation *GeoLocation `json:"geoLocation,omitempty"`
	Boundary    *Boundary    `json:"boundary,omitempty"`
}

// LocationLookupResponse wraps location lookup data.
type LocationLookupResponse struct {
	Status    *Status     `json:"status,omitempty"`
	Locations []*Location `json:"location,omitempty"`
}

// Location represents location lookup data.
type Location struct {
	ID          *string      `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Type        *string      `json:"type,omitempty"`
	GeoLocation *GeoLocation `json:"geoLocation,omitempty"`
}

// SaleComparablesResponse wraps sale comparables data.
type SaleComparablesResponse struct {
	Status          *Status           `json:"status,omitempty"`
	SaleComparables []*SaleComparable `json:"saleComparable,omitempty"`
}

// SaleComparable represents sale comparable data.
type SaleComparable struct {
	PropertyID *string  `json:"propertyId,omitempty"`
	Address    *Address `json:"address,omitempty"`
	SaleAmount *float64 `json:"saleAmount,omitempty"`
	SaleDate   *string  `json:"saleDate,omitempty"`
	Distance   *float64 `json:"distance,omitempty"`
	MatchCode  *string  `json:"matchCode,omitempty"`
	Quality    *string  `json:"quality,omitempty"`
}

// TransportationNoiseResponse wraps transportation noise data.
type TransportationNoiseResponse struct {
	Status              *Status                `json:"status,omitempty"`
	TransportationNoise []*TransportationNoise `json:"transportationNoise,omitempty"`
}

// TransportationNoise represents transportation noise data.
type TransportationNoise struct {
	PropertyID *string  `json:"propertyId,omitempty"`
	NoiseLevel *string  `json:"noiseLevel,omitempty"`
	Source     *string  `json:"source,omitempty"`
	Distance   *float64 `json:"distance,omitempty"`
}

// ParcelTilesResponse wraps parcel tiles data.
type ParcelTilesResponse struct {
	Status      *Status       `json:"status,omitempty"`
	ParcelTiles []*ParcelTile `json:"parcelTile,omitempty"`
}

// ParcelTile represents parcel tile data.
type ParcelTile struct {
	TileID *string `json:"tileId,omitempty"`
	Format *string `json:"format,omitempty"`
	Data   []byte  `json:"data,omitempty"`
}

// PreforeclosureResponse wraps pre-foreclosure details data.
type PreforeclosureResponse struct {
	Status         *Status           `json:"status,omitempty"`
	Preforeclosure []*Preforeclosure `json:"preforeclosure,omitempty"`
}

// Preforeclosure represents pre-foreclosure data.
type Preforeclosure struct {
	PropertyID      *string  `json:"propertyId,omitempty"`
	ForeclosureType *string  `json:"foreclosureType,omitempty"`
	Status          *string  `json:"status,omitempty"`
	Amount          *float64 `json:"amount,omitempty"`
	DateFiled       *string  `json:"dateFiled,omitempty"`
}

// PreforeclosureDetailsResponse wraps pre-foreclosure details data.
type PreforeclosureDetailsResponse struct {
	Status                *Status                 `json:"status,omitempty"`
	PreforeclosureDetails []*PreforeclosureDetail `json:"preforeclosureDetail,omitempty"`
}

// PreforeclosureDetail represents pre-foreclosure detail data.
type PreforeclosureDetail struct {
	PropertyID    *string `json:"propertyId,omitempty"`
	ForeclosureID *string `json:"foreclosureId,omitempty"`
	Status        *string `json:"status,omitempty"`
	FilingDate    *string `json:"filingDate,omitempty"`
	Amount        *string `json:"amount,omitempty"`
}
