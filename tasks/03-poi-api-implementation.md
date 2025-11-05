# POIAPI Implementation Guide

## Overview

The POIAPI (Points of Interest API) provides location-based business and facility searches with **5 endpoints** covering both legacy (v2.0.0) and modern (v4 neighborhood) approaches.

## Endpoints

### Legacy v2.0.0 Endpoints (3 total)

| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/poisearch/v2.0.0/poi/Street+Address/` | POI search by street address | streetaddress, recordlimit, searchdistance |
| `/poisearch/v2.0.0/poi/Geography/` | POI search by ZIP centroid | zipcode, recordlimit, searchdistance |
| `/poisearch/v2.0.0/poi/Point/` | POI search by coordinates | point (WKT), radius, recordlimit |

### v4 Neighborhood Endpoints (2 total)

| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/neighborhood/poi` | Search POI by location | point, radius, address, categoryName, zipCode |
| `/neighborhood/poi/categorylookup` | Lookup POI categories | industry, category, lineofbusiness |

## Business Categories (14 total)

| Category | Description |
|----------|-------------|
| `ATTRACTIONS-RECREATION` | Tourist attractions, entertainment venues, recreation facilities |
| `AUTOMOTIVE SERVICES` | Auto repair, gas stations, car washes, dealerships |
| `BANKS-FINANCIAL` | Banks, credit unions, ATMs, financial services |
| `EATING-DRINKING` | Restaurants, bars, cafes, food services |
| `EDUCATION` | Schools, colleges, libraries, tutoring |
| `FARM-RANCH` | Agricultural services, farms, ranches |
| `GOVERNMENT-PUBLIC` | Government offices, post offices, public services |
| `HEALTH CARE SERVICES` | Hospitals, clinics, doctors, dentists |
| `HOSPITALITY` | Hotels, motels, lodging |
| `ORGANIZATIONS-ASSOCIATIONS` | Non-profits, clubs, associations |
| `PERSONAL SERVICES` | Salons, dry cleaners, personal care |
| `PET SERVICES` | Veterinarians, pet stores, grooming |
| `SHOPPING` | Retail stores, malls, shopping centers |
| `TRAVEL` | Travel agencies, airlines, transportation |

## Key Parameters

### Search Parameters
- **streetaddress**: Full street address string (e.g., "123 Main St, Springfield, IL 62701")
- **zipcode**: ZIP code for geographic centroid search
- **point**: Geospatial point in WKT format: `POINT(longitude latitude)`
- **radius**: Search radius in miles (default varies by endpoint)
- **searchdistance**: Maximum distance for search in miles (legacy parameter)
- **recordlimit**: Maximum number of results to return

### Filter Parameters
- **categoryName**: Business category name (see table above)
- **LineOfBusinessName**: Specific business line within category
- **IndustryName**: Industry classification
- **categoryId**: Numeric category identifier

### Category Lookup Parameters
- **industry**: Industry name for category filtering
- **category**: Category name for filtering
- **lineofbusiness**: Specific line of business

## Implementation Checklist

### Phase 1: Legacy v2.0.0 Endpoints
- [ ] Implement `SearchByAddress()` with street address
- [ ] Implement `SearchByGeography()` with ZIP code
- [ ] Implement `SearchByPoint()` with coordinates
- [ ] Add WKT point string generation helper
- [ ] Create POI result models

### Phase 2: v4 Neighborhood Endpoints
- [ ] Implement `SearchNeighborhoodPOI()` with flexible parameters
- [ ] Implement `CategoryLookup()` for category discovery
- [ ] Add category validation helpers
- [ ] Create neighborhood POI models

### Phase 3: Business Category Constants
- [ ] Create constants for all 14 business categories
- [ ] Add category validation function
- [ ] Create category description lookup

### Phase 4: Distance Calculations
- [ ] Add distance calculation helpers (haversine formula)
- [ ] Add distance sorting helpers
- [ ] Create bounding box calculation helpers

## Model Design

### Core Models

```go
package models

// POIResults response (legacy v2.0.0)
type POIResults struct {
    Status *ResponseStatus `json:"status,omitempty"`
    POIs   []POI           `json:"pois,omitempty"`
}

// NeighborhoodPOI response (v4)
type NeighborhoodPOI struct {
    Status *ResponseStatus `json:"status,omitempty"`
    POIs   []POI           `json:"pois,omitempty"`
}

// POI represents a single point of interest
type POI struct {
    Name           *string  `json:"name,omitempty"`
    BusinessName   *string  `json:"businessName,omitempty"`
    Address        *Address `json:"address,omitempty"`
    Category       *string  `json:"category,omitempty"`
    CategoryID     *int     `json:"categoryId,omitempty"`
    Industry       *string  `json:"industry,omitempty"`
    LineOfBusiness *string  `json:"lineOfBusiness,omitempty"`
    Phone          *string  `json:"phone,omitempty"`
    Website        *string  `json:"website,omitempty"`
    Distance       *float64 `json:"distance,omitempty"`  // In miles
    Latitude       *float64 `json:"latitude,omitempty"`
    Longitude      *float64 `json:"longitude,omitempty"`
}

// CategoryLookup response
type CategoryLookup struct {
    Status     *ResponseStatus `json:"status,omitempty"`
    Categories []Category      `json:"categories,omitempty"`
}

type Category struct {
    CategoryID     *int    `json:"categoryId,omitempty"`
    CategoryName   *string `json:"categoryName,omitempty"`
    Industry       *string `json:"industry,omitempty"`
    LineOfBusiness *string `json:"lineOfBusiness,omitempty"`
    Description    *string `json:"description,omitempty"`
}
```

### Parameter Structs

```go
// SearchParams for legacy v2.0.0 searches
type SearchParams struct {
    RecordLimit    *int     `json:"-"`
    SearchDistance *float64 `json:"-"`  // Miles
    CategoryFilter *string  `json:"-"`
}

// NeighborhoodSearchParams for v4 searches
type NeighborhoodSearchParams struct {
    Point              *string  `json:"-"`  // WKT format
    Radius             *float64 `json:"-"`  // Miles
    Address            *string  `json:"-"`
    CategoryName       *string  `json:"-"`
    LineOfBusinessName *string  `json:"-"`
    IndustryName       *string  `json:"-"`
    CategoryID         *int     `json:"-"`
    ZipCode            *string  `json:"-"`
}

// CategoryParams for category lookup
type CategoryParams struct {
    Industry       *string `json:"-"`
    Category       *string `json:"-"`
    LineOfBusiness *string `json:"-"`
}
```

## Constants

Create `pkg/models/poi_categories.go`:

```go
package models

// POI business category constants
const (
    CategoryAttractionsRecreation    = "ATTRACTIONS-RECREATION"
    CategoryAutomotiveServices       = "AUTOMOTIVE SERVICES"
    CategoryBanksFinancial           = "BANKS-FINANCIAL"
    CategoryEatingDrinking           = "EATING-DRINKING"
    CategoryEducation                = "EDUCATION"
    CategoryFarmRanch                = "FARM-RANCH"
    CategoryGovernmentPublic         = "GOVERNMENT-PUBLIC"
    CategoryHealthCareServices       = "HEALTH CARE SERVICES"
    CategoryHospitality              = "HOSPITALITY"
    CategoryOrganizationsAssociations = "ORGANIZATIONS-ASSOCIATIONS"
    CategoryPersonalServices         = "PERSONAL SERVICES"
    CategoryPetServices              = "PET SERVICES"
    CategoryShopping                 = "SHOPPING"
    CategoryTravel                   = "TRAVEL"
)

// AllPOICategories for validation
var AllPOICategories = []string{
    CategoryAttractionsRecreation,
    CategoryAutomotiveServices,
    CategoryBanksFinancial,
    CategoryEatingDrinking,
    CategoryEducation,
    CategoryFarmRanch,
    CategoryGovernmentPublic,
    CategoryHealthCareServices,
    CategoryHospitality,
    CategoryOrganizationsAssociations,
    CategoryPersonalServices,
    CategoryPetServices,
    CategoryShopping,
    CategoryTravel,
}

// ValidateCategory checks if a category is valid
func ValidateCategory(category string) bool {
    for _, c := range AllPOICategories {
        if c == category {
            return true
        }
    }
    return false
}
```

## Helper Functions

### WKT Point Generation

```go
package poi

import "fmt"

// MakePoint creates a WKT POINT string from coordinates
func MakePoint(longitude, latitude float64) string {
    return fmt.Sprintf("POINT(%f %f)", longitude, latitude)
}
```

### Distance Calculation

```go
package poi

import "math"

// Haversine calculates the great-circle distance between two points
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
    const earthRadiusMiles = 3958.8

    lat1Rad := lat1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    deltaLat := (lat2 - lat1) * math.Pi / 180
    deltaLon := (lon2 - lon1) * math.Pi / 180

    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
        math.Cos(lat1Rad)*math.Cos(lat2Rad)*
        math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    return earthRadiusMiles * c
}
```

## Testing Strategy

1. **Category Tests**: Test all 14 business categories
2. **Geographic Tests**: Test searches across different US regions
3. **Distance Tests**: Verify distance calculations and sorting
4. **Legacy vs v4**: Compare results between legacy and v4 endpoints
5. **Edge Cases**: Zero results, very large radius, invalid coordinates

## Common Pitfalls

1. **WKT Format**: MUST be `POINT(lon lat)` NOT `POINT(lat lon)` - longitude first!
2. **Distance Units**: Always in miles, not kilometers
3. **Category Names**: Case-sensitive and must exactly match constants
4. **Missing Fields**: POI data frequently missing phone, website, etc.
5. **Coordinate Precision**: Use float64, not float32, for coordinates

## Example Usage

### Example 1: Find restaurants near an address

```go
package main

import (
    "context"
    "fmt"
    "github.com/my-eq/go-attom/pkg/client"
    "github.com/my-eq/go-attom/pkg/poi"
    "github.com/my-eq/go-attom/pkg/models"
)

func main() {
    c := client.NewClient("YOUR_API_KEY")
    svc := poi.NewService(c)
    ctx := context.Background()
    
    // Search for restaurants within 2 miles
    params := &models.NeighborhoodSearchParams{
        Address:      stringPtr("123 Main St, Springfield, IL 62701"),
        Radius:       float64Ptr(2.0),
        CategoryName: stringPtr(models.CategoryEatingDrinking),
    }
    
    results, err := svc.SearchNeighborhoodPOI(ctx, params)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Found %d restaurants:\n", len(results.POIs))
    for _, p := range results.POIs {
        if p.Name != nil && p.Distance != nil {
            fmt.Printf("- %s (%.2f miles)\n", *p.Name, *p.Distance)
        }
    }
}
```

### Example 2: Find all POIs by coordinates

```go
// Los Angeles coordinates
lat, lon := 34.0522, -118.2437

params := &models.NeighborhoodSearchParams{
    Point:  stringPtr(poi.MakePoint(lon, lat)),  // Note: lon, lat order!
    Radius: float64Ptr(5.0),
}

results, err := svc.SearchNeighborhoodPOI(ctx, params)
```

### Example 3: Lookup available categories

```go
// Get all health care categories
params := &models.CategoryParams{
    Industry: stringPtr("HEALTH CARE"),
}

categories, err := svc.CategoryLookup(ctx, params)
if err != nil {
    panic(err)
}

for _, cat := range categories.Categories {
    if cat.CategoryName != nil && cat.LineOfBusiness != nil {
        fmt.Printf("%s - %s\n", *cat.CategoryName, *cat.LineOfBusiness)
    }
}
```

### Example 4: Legacy search by ZIP code

```go
searchParams := &models.SearchParams{
    RecordLimit:    intPtr(20),
    SearchDistance: float64Ptr(3.0),
    CategoryFilter: stringPtr(models.CategoryShopping),
}

results, err := svc.SearchByGeography(ctx, "90210", searchParams)
```

## Helper Functions

```go
func stringPtr(s string) *string {
    return &s
}

func intPtr(i int) *int {
    return &i
}

func float64Ptr(f float64) *float64 {
    return &f
}
```

## Integration with PropertyAPI

Combine property and POI data:

```go
// Get property first
property, _ := propertySvc.GetPropertyDetail(ctx, address)

// Then find nearby schools
if property.Location != nil && 
   property.Location.Latitude != nil && 
   property.Location.Longitude != nil {
    
    point := poi.MakePoint(*property.Location.Longitude, 
                          *property.Location.Latitude)
    
    params := &models.NeighborhoodSearchParams{
        Point:        &point,
        Radius:       float64Ptr(1.0),
        CategoryName: stringPtr(models.CategoryEducation),
    }
    
    schools, _ := poiSvc.SearchNeighborhoodPOI(ctx, params)
    
    fmt.Printf("Found %d schools within 1 mile\n", len(schools.POIs))
}
```
