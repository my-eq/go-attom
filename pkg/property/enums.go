package property

import "fmt"

// AcceptHeader represents valid values for the Accept header in API requests.
const (
	AcceptHeaderJSON = "application/json"
	AcceptHeaderXML  = "application/xml"
)

// Format represents valid values for geographic boundary formats.
const (
	FormatGeoJSON = "geojson"
	FormatWKT     = "wkt"
)

// PropertyType represents valid property type classifications.
// These values can be used with the propertytype parameter in various endpoints.
const (
	PropertyTypeAgriculturalNEC      = "AGRICULTURAL (NEC)"
	PropertyTypeApartment            = "APARTMENT"
	PropertyTypeCabin                = "CABIN"
	PropertyTypeClub                 = "CLUB"
	PropertyTypeCommercialNEC        = "COMMERCIAL (NEC)"
	PropertyTypeCommercialBuilding   = "COMMERCIAL BUILDING"
	PropertyTypeCommercialCondo      = "COMMERCIAL CONDOMINIUM"
	PropertyTypeCommonArea           = "COMMON AREA"
	PropertyTypeCondominium          = "CONDOMINIUM"
	PropertyTypeConvertedResidence   = "CONVERTED RESIDENCE"
	PropertyTypeCountyProperty       = "COUNTY PROPERTY"
	PropertyTypeDuplex               = "DUPLEX"
	PropertyTypeFarms                = "FARMS"
	PropertyTypeFastFoodFranchise    = "FAST FOOD FRANCHISE"
	PropertyTypeFederalProperty      = "FEDERAL PROPERTY"
	PropertyTypeForest               = "FOREST"
	PropertyTypeGroupQuarters        = "GROUP QUARTERS"
	PropertyTypeIndustrialNEC        = "INDUSTRIAL (NEC)"
	PropertyTypeIndustrialPlant      = "INDUSTRIAL PLANT"
	PropertyTypeManufacturedHome     = "MANUFACTURED HOME"
	PropertyTypeMarinaFacility       = "MARINA FACILITY"
	PropertyTypeMiscellaneous        = "MISCELLANEOUS"
	PropertyTypeMobileHome           = "MOBILE HOME"
	PropertyTypeMobileHomeLot        = "MOBILE HOME LOT"
	PropertyTypeMobileHomePark       = "MOBILE HOME PARK"
	PropertyTypeMultiFamilyDwelling  = "MULTI FAMILY DWELLING"
	PropertyTypeNurseryHorticulture  = "NURSERY/HORTICULTURE"
	PropertyTypeOfficeBuilding       = "OFFICE BUILDING"
	PropertyTypePublicNEC            = "PUBLIC (NEC)"
	PropertyTypeReligious            = "RELIGIOUS"
	PropertyTypeResidentialNEC       = "RESIDENTIAL (NEC)"
	PropertyTypeResidentialAcreage   = "RESIDENTIAL ACREAGE"
	PropertyTypeResidentialLot       = "RESIDENTIAL LOT"
	PropertyTypeRetailTrade          = "RETAIL TRADE"
	PropertyTypeSFR                  = "SFR"
	PropertyTypeStateProperty        = "STATE PROPERTY"
	PropertyTypeStoresAndOffices     = "STORES & OFFICES"
	PropertyTypeStoresAndResidential = "STORES & RESIDENTIAL"
	PropertyTypeTaxExempt            = "TAX EXEMPT"
	PropertyTypeTownhouseRowhouse    = "TOWNHOUSE/ROWHOUSE"
	PropertyTypeTriplex              = "TRIPLEX"
	PropertyTypeUtilities            = "UTILITIES"
	PropertyTypeVacantLandNEC        = "VACANT LAND (NEC)"
)

// OrderBy represents valid sorting options for API responses.
// These values can be used with the orderby parameter in various endpoints.
const (
	OrderByCalendarDate        = "calendardate"
	OrderByPublishedDate       = "publisheddate"
	OrderByPropertyType        = "propertytype"
	OrderBySaleAmount          = "saleamt"
	OrderByAVMValue            = "avmvalue"
	OrderByAssessedTotalValue  = "assdttlvalue"
	OrderBySalesSearchDate     = "salesearchdate"
	OrderBySaleTransactionDate = "saletransactiondate"
	OrderByBeds                = "beds"
	OrderByBathsTotal          = "bathstotal"
	OrderByUniversalSize       = "universalsize"
	OrderByLotSize1            = "lotsize1"
	OrderByLotSize2            = "lotsize2"
)

// ValidateAcceptHeader checks if the provided accept header value is valid.
func ValidateAcceptHeader(accept string) error {
	switch accept {
	case AcceptHeaderJSON, AcceptHeaderXML:
		return nil
	default:
		return fmt.Errorf("invalid accept header: %q (must be %q or %q)", accept, AcceptHeaderJSON, AcceptHeaderXML)
	}
}

// ValidateFormat checks if the provided format value is valid.
func ValidateFormat(format string) error {
	switch format {
	case FormatGeoJSON, FormatWKT:
		return nil
	default:
		return fmt.Errorf("invalid format: %q (must be %q or %q)", format, FormatGeoJSON, FormatWKT)
	}
}

// ValidatePropertyType checks if the provided property type is valid.
func ValidatePropertyType(propertyType string) error {
	validTypes := []string{
		PropertyTypeAgriculturalNEC,
		PropertyTypeApartment,
		PropertyTypeCabin,
		PropertyTypeClub,
		PropertyTypeCommercialNEC,
		PropertyTypeCommercialBuilding,
		PropertyTypeCommercialCondo,
		PropertyTypeCommonArea,
		PropertyTypeCondominium,
		PropertyTypeConvertedResidence,
		PropertyTypeCountyProperty,
		PropertyTypeDuplex,
		PropertyTypeFarms,
		PropertyTypeFastFoodFranchise,
		PropertyTypeFederalProperty,
		PropertyTypeForest,
		PropertyTypeGroupQuarters,
		PropertyTypeIndustrialNEC,
		PropertyTypeIndustrialPlant,
		PropertyTypeManufacturedHome,
		PropertyTypeMarinaFacility,
		PropertyTypeMiscellaneous,
		PropertyTypeMobileHome,
		PropertyTypeMobileHomeLot,
		PropertyTypeMobileHomePark,
		PropertyTypeMultiFamilyDwelling,
		PropertyTypeNurseryHorticulture,
		PropertyTypeOfficeBuilding,
		PropertyTypePublicNEC,
		PropertyTypeReligious,
		PropertyTypeResidentialNEC,
		PropertyTypeResidentialAcreage,
		PropertyTypeResidentialLot,
		PropertyTypeRetailTrade,
		PropertyTypeSFR,
		PropertyTypeStateProperty,
		PropertyTypeStoresAndOffices,
		PropertyTypeStoresAndResidential,
		PropertyTypeTaxExempt,
		PropertyTypeTownhouseRowhouse,
		PropertyTypeTriplex,
		PropertyTypeUtilities,
		PropertyTypeVacantLandNEC,
	}

	for _, validType := range validTypes {
		if propertyType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid property type: %q", propertyType)
}

// ValidateOrderBy checks if the provided orderby value is valid.
func ValidateOrderBy(orderBy string) error {
	validOrders := []string{
		OrderByCalendarDate,
		OrderByPublishedDate,
		OrderByPropertyType,
		OrderBySaleAmount,
		OrderByAVMValue,
		OrderByAssessedTotalValue,
		OrderBySalesSearchDate,
		OrderBySaleTransactionDate,
		OrderByBeds,
		OrderByBathsTotal,
		OrderByUniversalSize,
		OrderByLotSize1,
		OrderByLotSize2,
	}

	for _, validOrder := range validOrders {
		if orderBy == validOrder {
			return nil
		}
	}
	return fmt.Errorf("invalid orderby: %q", orderBy)
}
