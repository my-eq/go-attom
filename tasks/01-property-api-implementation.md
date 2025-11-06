# PropertyAPI Implementation Guide

## Overview
The PropertyAPI is the largest ATTOM API group with **36+ endpoints** across 9 categories. This API provides comprehensive property data including details, sales, assessments, valuations, and school information.

## Endpoints by Category

### 1. Property Resources (11 endpoints)
Base path: `/v4/property/`

| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/property/id` | Retrieve ATTOM property ID | address, address1/address2 |
| `/property/detail` | Detailed property information | attomId, address, fips+APN |
| `/property/address` | Property address data | attomId |
| `/property/snapshot` | Property summary snapshot | address, postalCode, latitude/longitude+radius |
| `/property/basicprofile` | Basic property profile | address |
| `/property/expandedprofile` | Full property characteristics | address, geoIdV4 |
| `/property/detailwithschools` | Property + school zones | address |
| `/property/detailmortgage` | Property + mortgage data | address |
| `/property/detailowner` | Property + owner info | address, absenteeOwner filter |
| `/property/detailmortgageowner` | Property + mortgage + owner | address |
| `/property/buildingpermits` | Building permit records | address |

### 2. Sale Resources (2 endpoints)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/sale/detail` | Transaction details with dates/amounts | address, startSaleSearchDate/endSaleSearchDate |
| `/sale/snapshot` | Sales summary | address |

### 3. Assessment Resources (3 endpoints)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/assessment/detail` | Tax/assessment details | address |
| `/assessment/snapshot` | Assessment summary | address |
| `/assessmenthistory/detail` | Historical assessments | address, startCalendarDate/endCalendarDate |

### 4. AVM Resources (4 endpoints)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/avm/snapshot` | AVM value snapshot | address |
| `/attomavm/detail` | Detailed ATTOM AVM | address, comparables |
| `/avmhistory/detail` | Historical AVM values | address, date range |
| `/valuation/rentalavm` | Rental property valuation | address |

### 5. Sales History Resources (4 endpoints)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/saleshistory/detail` | Full transaction history | address |
| `/saleshistory/snapshot` | History summary | address |
| `/saleshistory/basichistory` | Essential fields only | address |
| `/saleshistory/expandedhistory` | All available data | address |

### 6. Sales Trend Resources (2 endpoints)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/salestrend/snapshot` | Geographic sales trends | geoIdV4, interval (monthly/quarterly/yearly) |
| `/transaction/salestrend` | Transaction-based trends | geoIdV4, propertyType, interval |

### 7. School Resources (4 endpoints)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/school/search` | Find nearby schools | address, radius |
| `/school/profile` | School information | schoolId |
| `/school/district` | District information | address |
| `/school/detailwithschools` | Property + schools | address |

### 8. Assessment History (1 endpoint)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/assessmenthistory/detail` | Assessment records over time | address, date range |

### 9. All Events (1 endpoint)
| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/allevents/detail` | All property events combined | address |

## Common Parameters

### Property Identification
- `attomId` / `id` - ATTOM's unique property identifier
- `address` - Full address string: "123 Main St, City, ST 12345"
- `address1` / `address2` - Split format: "123 Main St" / "City, ST 12345"
- `fips` + `APN` - County FIPS code + Assessor's Parcel Number
- `geoIdV4` - ATTOM's new geographic identifier (SHA hash)

### Geographic Search
- `latitude` / `longitude` + `radius` - Geospatial search (radius in miles)
- `postalCode` - ZIP code search

### Property Filters
- `propertyType` - Standardized type: "apartment", "single family", "condo", etc.
- `propertyIndicator` - Numeric type code (0-90)
- `minBeds` / `maxBeds` - Bedroom range
- `minBathsTotal` / `maxBathsTotal` - Bathroom range
- `minYearBuilt` / `maxYearBuilt` - Construction year
- `minSaleAmt` / `maxSaleAmt` - Price range

### Date Ranges
- `startCalendarDate` / `endCalendarDate` - Record update dates
- `startAddedDate` / `endAddedDate` - Record addition dates
- `startSaleSearchDate` / `endSaleSearchDate` - Sale transaction dates

### Results Control
- `orderBy` - Sort field
- `page` / `pagesize` - Pagination

## Implementation Checklist

### Phase 1: Core Property Methods
- [ ] Implement `GetPropertyByID()` - Most basic lookup
- [ ] Implement `GetPropertyDetail()` - Main property data method
- [ ] Implement `GetPropertySnapshot()` - Lightweight summary
- [ ] Add address parsing helpers (split/join address1/address2)
- [ ] Add FIPS+APN validation helpers

### Phase 2: Enhanced Property Methods
- [ ] Implement `GetBasicProfile()` / `GetExpandedProfile()`
- [ ] Implement `GetDetailWithSchools()`
- [ ] Implement `GetDetailMortgage()` / `GetDetailOwner()` / `GetDetailMortgageOwner()`
- [ ] Implement `GetBuildingPermits()`
- [ ] Create comprehensive property models with all optional fields

### Phase 3: Sales & Assessment Methods
- [ ] Implement sale endpoints (detail, snapshot)
- [ ] Implement assessment endpoints (detail, snapshot, history)
- [ ] Add date range parameter builders
- [ ] Create sales/assessment models

### Phase 4: AVM & Valuation Methods
- [ ] Implement all 4 AVM endpoints
- [ ] Add confidence score interpretation helpers
- [ ] Create AVM models with value ranges

### Phase 5: Sales History & Trends
- [ ] Implement 4 sales history variants
- [ ] Implement 2 sales trend endpoints
- [ ] Add interval parameter support (monthly/quarterly/yearly)
- [ ] Create trend aggregation models

### Phase 6: School Methods
- [ ] Implement school search with radius
- [ ] Implement school profile lookup
- [ ] Implement district lookup
- [ ] Create school rating models

### Phase 7: All Events
- [ ] Implement `/allevents/detail` endpoint
- [ ] Create combined event model
- [ ] Add event filtering helpers

## Model Design Considerations

### Required Patterns
```go
// ALL fields must be pointers or have omitempty tags
type PropertyDetail struct {
    AttomID       *string          `json:"attomId,omitempty"`
    Address       *Address         `json:"address,omitempty"`
    Lot           *Lot             `json:"lot,omitempty"`
    Building      *Building        `json:"building,omitempty"`
    Assessment    *Assessment      `json:"assessment,omitempty"`
    // ... many more optional fields
}

type Address struct {
    Line1     *string `json:"line1,omitempty"`
    Line2     *string `json:"line2,omitempty"`
    City      *string `json:"city,omitempty"`
    State     *string `json:"state,omitempty"`
    ZipCode   *string `json:"zipCode,omitempty"`
    Latitude  *float64 `json:"latitude,omitempty"`
    Longitude *float64 `json:"longitude,omitempty"`
}
```

### Field Access Patterns
```go
// ALWAYS validate before dereferencing
if detail.Address != nil && detail.Address.Line1 != nil {
    fmt.Println(*detail.Address.Line1)
}

// Consider helper methods
func (a *Address) GetLine1() string {
    if a != nil && a.Line1 != nil {
        return *a.Line1
    }
    return ""
}
```

## Testing Strategy

1. **Unit Tests**: Mock HTTP responses for each endpoint
2. **Integration Tests**: Use real API with test addresses
3. **Field Coverage Tests**: Verify all optional fields parse correctly
4. **Error Handling Tests**: Test missing field scenarios
5. **Parameter Validation Tests**: Test all parameter combinations

## Common Pitfalls

1. **Missing Fields**: NEVER assume a field exists - always check for nil
2. **Casing Variations**: JSON field names may vary - use explicit tags
3. **Date Formats**: ATTOM uses multiple date formats - create parsers
4. **Property Types**: Use constants for propertyType values
5. **Rate Limiting**: Implement exponential backoff for 429 responses

## Example Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/my-eq/go-attom/pkg/client"
    "github.com/my-eq/go-attom/pkg/property"
)

func main() {
    // Create client
    c := client.NewClient("YOUR_API_KEY")
    svc := property.NewService(c)
    
    // Get property detail
    ctx := context.Background()
    detail, err := svc.GetPropertyDetail(ctx, "123 Main St, Springfield, IL 62701")
    if err != nil {
        panic(err)
    }
    
    // Safe field access
    if detail.Address != nil && detail.Address.Line1 != nil {
        fmt.Printf("Address: %s\n", *detail.Address.Line1)
    }
    
    if detail.Building != nil && detail.Building.Rooms != nil && 
       detail.Building.Rooms.Beds != nil {
        fmt.Printf("Bedrooms: %d\n", *detail.Building.Rooms.Beds)
    }
}
```
