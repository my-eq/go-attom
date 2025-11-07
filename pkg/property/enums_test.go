package property

import (
	"strings"
	"testing"
)

func TestValidateAcceptHeader(t *testing.T) {
	tests := []struct {
		name    string
		accept  string
		wantErr bool
	}{
		{
			name:    "valid json",
			accept:  AcceptHeaderJSON,
			wantErr: false,
		},
		{
			name:    "valid xml",
			accept:  AcceptHeaderXML,
			wantErr: false,
		},
		{
			name:    "invalid accept header",
			accept:  "text/plain",
			wantErr: true,
		},
		{
			name:    "empty string",
			accept:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAcceptHeader(tt.accept)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAcceptHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "valid geojson",
			format:  FormatGeoJSON,
			wantErr: false,
		},
		{
			name:    "valid wkt",
			format:  FormatWKT,
			wantErr: false,
		},
		{
			name:    "invalid format",
			format:  "invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			format:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePropertyType(t *testing.T) {
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

	tests := []struct {
		name         string
		propertyType string
		wantErr      bool
	}{
		{
			name:         "invalid property type",
			propertyType: "INVALID_TYPE",
			wantErr:      true,
		},
		{
			name:         "empty string",
			propertyType: "",
			wantErr:      true,
		},
	}

	// Add tests for all valid types
	for _, validType := range validTypes {
		tests = append(tests, struct {
			name         string
			propertyType string
			wantErr      bool
		}{
			name:         "valid " + validType,
			propertyType: validType,
			wantErr:      false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePropertyType(tt.propertyType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePropertyType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOrderBy(t *testing.T) {
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

	tests := []struct {
		name    string
		orderBy string
		wantErr bool
	}{
		{
			name:    "invalid order by",
			orderBy: "invalid_order",
			wantErr: true,
		},
		{
			name:    "empty string",
			orderBy: "",
			wantErr: true,
		},
	}

	// Add tests for all valid orders
	for _, validOrder := range validOrders {
		tests = append(tests, struct {
			name    string
			orderBy string
			wantErr bool
		}{
			name:    "valid " + validOrder,
			orderBy: validOrder,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderBy(tt.orderBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrderBy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFIPSAndAPN(t *testing.T) {
	tests := []struct {
		wantErr bool
		name    string
		fips    string
		apn     string
		errMsg  string
	}{
		{
			name:    "valid FIPS and APN",
			fips:    "06037",
			apn:     "123-456-789",
			wantErr: false,
		},
		{
			name:    "empty FIPS",
			fips:    "",
			apn:     "123-456-789",
			wantErr: true,
			errMsg:  "both fips and APN are required",
		},
		{
			name:    "empty APN",
			fips:    "06037",
			apn:     "",
			wantErr: true,
			errMsg:  "both fips and APN are required",
		},
		{
			name:    "whitespace FIPS",
			fips:    "   ",
			apn:     "123-456-789",
			wantErr: true,
			errMsg:  "both fips and APN are required",
		},
		{
			name:    "whitespace APN",
			fips:    "06037",
			apn:     "   ",
			wantErr: true,
			errMsg:  "both fips and APN are required",
		},
		{
			name:    "both empty",
			fips:    "",
			apn:     "",
			wantErr: true,
			errMsg:  "both fips and APN are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFIPSAndAPN(tt.fips, tt.apn)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
				} else if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error = %v, want error containing %q", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
