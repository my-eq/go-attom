# AreaAPI Implementation Guide

## Overview

The AreaAPI provides geographic hierarchy and boundary data with **6 endpoints** for looking up counties, states, and various geographic boundary types.

## Endpoints

Base path: `/areaapi/v2.0.0/area/`

| Endpoint | Description | Key Parameters |
|----------|-------------|----------------|
| `/county/lookup` | County lookup by state | StateId |
| `/area/hierarchy/lookup` | Geographic hierarchy lookup | wktstring, geotype, latitude, longitude |
| `/area/state/lookup` | State area lookup | geoIdV4, AreaId |
| `/area/boundary/detail` | Area boundary in GeoJSON/WKT | geoIdV4, AreaId, format |
| `/area/geoId/legacyLookup` | Legacy geocode lookup | geoId |
| `/area/geoId/Lookup` | Geocode lookup | geoId, geotype |

## Geographic Types (geotype)

| Code | Description |
|------|-------------|
| `PZ` | Parcel Zone |
| `SB` | School Attendance Boundary |
| `DB` | School District Boundary |
| `ZI` | ZIP Code Tabulation Area |
| `N1` | Macro Neighborhood |
| `N2` | Neighborhood |
| `N3` | Sub-Neighborhood |
| `N4` | Residential Subdivision |
| `ST` | State |
| `CO` | County |
| `PL` | Census Place |
| `CB` | Census Block |

## Key Parameters

- **geoIdV4**: ATTOM's new geographic identifier (SHA-256 hash format)
- **AreaId**: Legacy area identifier (numeric)
- **StateId**: Two-letter state code (e.g., "CA", "TX", "NY")
- **wktstring**: Well-Known Text geometry string (e.g., "POLYGON(...)")
- **geotype**: Geographic boundary type code (see table above)
- **latitude/longitude**: Coordinate pair for point-in-polygon lookup
- **format**: Output format for boundary data ("geojson" or "wkt")

## Implementation Checklist

### Phase 1: Basic Lookups
- [ ] Implement `CountyLookup()` with state parameter
- [ ] Create state code constants (all 50 US states + DC, PR, etc.)
- [ ] Add state code validation

### Phase 2: Hierarchy Lookups
- [ ] Implement `AreaHierarchyLookup()` with WKT support
- [ ] Add WKT string validation helpers
- [ ] Implement coordinate-based hierarchy lookup
- [ ] Create geotype constants (PZ, SB, DB, ZI, N1-N4, ST, CO, PL, CB)

### Phase 3: State & GeoID Lookups
- [ ] Implement `StateLookup()` with geoIdV4
- [ ] Implement legacy `LegacyGeoIDLookup()`
- [ ] Implement new `GeoIDLookup()`
- [ ] Add geoIdV4 format validation (SHA-256 hex string)

### Phase 4: Boundary Details
- [ ] Implement `BoundaryDetail()` with format parameter
- [ ] Add GeoJSON parsing support
- [ ] Add WKT parsing support
- [ ] Create boundary geometry models

## Model Design

### Core Models

```go
package models

// CountyLookup response
type CountyLookup struct {
    Status   *ResponseStatus `json:"status,omitempty"`
    Counties []County        `json:"counties,omitempty"`
}

type County struct {
    FIPS      *string  `json:"fips,omitempty"`
    Name      *string  `json:"name,omitempty"`
    State     *string  `json:"state,omitempty"`
    GeoIDV4   *string  `json:"geoIdV4,omitempty"`
}

// AreaHierarchy response
type AreaHierarchy struct {
    Status     *ResponseStatus  `json:"status,omitempty"`
    Hierarchies []HierarchyLevel `json:"hierarchies,omitempty"`
}

type HierarchyLevel struct {
    GeoType   *string  `json:"geoType,omitempty"`
    GeoID     *string  `json:"geoId,omitempty"`
    GeoIDV4   *string  `json:"geoIdV4,omitempty"`
    Name      *string  `json:"name,omitempty"`
    AreaID    *string  `json:"areaId,omitempty"`
}

// BoundaryDetail response (GeoJSON format)
type BoundaryDetail struct {
    Status   *ResponseStatus `json:"status,omitempty"`
    Type     *string         `json:"type,omitempty"`      // "FeatureCollection"
    Features []Feature       `json:"features,omitempty"`
}

type Feature struct {
    Type       *string                `json:"type,omitempty"`  // "Feature"
    Geometry   *Geometry              `json:"geometry,omitempty"`
    Properties map[string]interface{} `json:"properties,omitempty"`
}

type Geometry struct {
    Type        *string         `json:"type,omitempty"`  // "Polygon", "MultiPolygon"
    Coordinates json.RawMessage `json:"coordinates,omitempty"`
}
```

### Parameter Structs

```go
// HierarchyParams for area hierarchy lookup
type HierarchyParams struct {
    WKTString *string  `json:"-"`
    GeoType   *string  `json:"-"`
    Latitude  *float64 `json:"-"`
    Longitude *float64 `json:"-"`
}

// BoundaryParams for boundary detail lookup
type BoundaryParams struct {
    GeoIDV4 *string `json:"-"`
    AreaID  *string `json:"-"`
    Format  *string `json:"-"`  // "geojson" or "wkt"
}

// StateParams for state lookup
type StateParams struct {
    GeoIDV4 *string `json:"-"`
    AreaID  *string `json:"-"`
}

// GeoIDParams for geocode lookup
type GeoIDParams struct {
    GeoID   *string `json:"-"`
    GeoType *string `json:"-"`
}
```

## Constants

Create `pkg/models/geotypes.go`:

```go
package models

// Geographic boundary type constants
const (
    GeoTypeParcelZone              = "PZ"
    GeoTypeSchoolAttendanceBoundary = "SB"
    GeoTypeSchoolDistrictBoundary   = "DB"
    GeoTypeZIPCode                  = "ZI"
    GeoTypeMacroNeighborhood        = "N1"
    GeoTypeNeighborhood             = "N2"
    GeoTypeSubNeighborhood          = "N3"
    GeoTypeResidentialSubdivision   = "N4"
    GeoTypeState                    = "ST"
    GeoTypeCounty                   = "CO"
    GeoTypeCensusPlace              = "PL"
    GeoTypeCensusBlock              = "CB"
)

// Output format constants
const (
    FormatGeoJSON = "geojson"
    FormatWKT     = "wkt"
)
```

Create `pkg/models/states.go`:

```go
package models

// US state code constants
const (
    StateAlabama    = "AL"
    StateAlaska     = "AK"
    StateArizona    = "AZ"
    // ... all 50 states + DC, PR, etc.
)

// ValidStates for validation
var ValidStates = []string{
    StateAlabama, StateAlaska, StateArizona,
    // ... complete list
}
```

## WKT String Examples

### Point
```
POINT(-118.2437 34.0522)
```

### Polygon
```
POLYGON((-118.5 34.0, -118.0 34.0, -118.0 34.5, -118.5 34.5, -118.5 34.0))
```

### MultiPolygon
```
MULTIPOLYGON(((-118.5 34.0, -118.0 34.0, -118.0 34.5, -118.5 34.5, -118.5 34.0)))
```

## Testing Strategy

1. **County Lookup Tests**: All 50 states + territories
2. **Hierarchy Tests**: Test all 12 geotype values
3. **Boundary Tests**: Both GeoJSON and WKT formats
4. **Coordinate Tests**: Point-in-polygon lookups across US
5. **GeoID Migration Tests**: Verify legacy vs. new geoId compatibility

## Common Pitfalls

1. **GeoIDV4 Format**: Must be valid SHA-256 hex string (64 characters)
2. **WKT Validation**: Coordinates must be in longitude, latitude order (not lat, lon!)
3. **Boundary Size**: Large polygons (e.g., states) can be huge - handle large responses
4. **Coordinate System**: ATTOM uses WGS84 (EPSG:4326) for coordinates
5. **Format Parameter**: Case-sensitive - use lowercase "geojson", "wkt"

## Example Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/my-eq/go-attom/pkg/client"
    "github.com/my-eq/go-attom/pkg/area"
    "github.com/my-eq/go-attom/pkg/models"
)

func main() {
    c := client.NewClient("YOUR_API_KEY")
    svc := area.NewService(c)
    ctx := context.Background()
    
    // Lookup counties in California
    counties, err := svc.CountyLookup(ctx, models.StateCalifornia)
    if err != nil {
        panic(err)
    }
    
    for _, county := range counties.Counties {
        if county.Name != nil && county.FIPS != nil {
            fmt.Printf("%s (FIPS: %s)\n", *county.Name, *county.FIPS)
        }
    }
    
    // Get area hierarchy for a point
    lat, lon := 34.0522, -118.2437  // Los Angeles
    params := &models.HierarchyParams{
        Latitude:  &lat,
        Longitude: &lon,
        GeoType:   stringPtr(models.GeoTypeNeighborhood),
    }
    
    hierarchy, err := svc.AreaHierarchyLookup(ctx, params)
    if err != nil {
        panic(err)
    }
    
    for _, level := range hierarchy.Hierarchies {
        if level.GeoType != nil && level.Name != nil {
            fmt.Printf("%s: %s\n", *level.GeoType, *level.Name)
        }
    }
    
    // Get boundary in GeoJSON
    boundaryParams := &models.BoundaryParams{
        GeoIDV4: stringPtr("abc123..."),  // Valid geoIdV4
        Format:  stringPtr(models.FormatGeoJSON),
    }
    
    boundary, err := svc.BoundaryDetail(ctx, boundaryParams)
    if err != nil {
        panic(err)
    }
    
    // boundary.Features contains GeoJSON features
}

func stringPtr(s string) *string {
    return &s
}
```

## Integration with PropertyAPI

AreaAPI data complements PropertyAPI:

```go
// Get property, then get its neighborhood boundary
property, _ := propertySvc.GetPropertyDetail(ctx, address)

if property.Location != nil && property.Location.GeoIDV4 != nil {
    boundary, _ := areaSvc.BoundaryDetail(ctx, &models.BoundaryParams{
        GeoIDV4: property.Location.GeoIDV4,
        Format:  stringPtr(models.FormatGeoJSON),
    })
    
    // Now you have the property AND its geographic boundary
}
```
