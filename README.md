# go-attom

[![CI](https://github.com/my-eq/go-attom/actions/workflows/ci.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/ci.yml)
[![Security](https://github.com/my-eq/go-attom/actions/workflows/security.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/security.yml)
[![README Lint](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml)
[![codecov](https://codecov.io/gh/my-eq/go-attom/branch/main/graph/badge.svg)](https://codecov.io/gh/my-eq/go-attom)
[![Go Report Card](https://goreportcard.com/badge/github.com/my-eq/go-attom)](https://goreportcard.com/report/github.com/my-eq/go-attom)
[![Go Reference](https://pkg.go.dev/badge/github.com/my-eq/go-attom.svg)](https://pkg.go.dev/github.com/my-eq/go-attom)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A strongly-typed, mock-friendly Go client for the [ATTOM Data API](https://api.gateway.attomdata.com/). The library focuses on clean, idiomatic Go and provides deep coverage of the Property API including property profiles, sales, assessments, valuations, schools, and historical trends (see the [API implementation summary](API_IMPLEMENTATION_SUMMARY.md#L9-L118) for the full endpoint catalog).

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [Look up an ATTOM ID for an address](#look-up-an-attom-id-for-an-address)
  - [Search for schools near a coordinate](#search-for-schools-near-a-coordinate)
- [Property API Coverage](#property-api-coverage)
  - [Property Profiles & Basics](#property-profiles--basics)
  - [Ownership, Mortgage, and Schools](#ownership-mortgage-and-schools)
  - [Sales, Assessments, and Valuations](#sales-assessments-and-valuations)
  - [Sales History & Trends](#sales-history--trends)
- [Building Requests with Options](#building-requests-with-options)
- [Error Handling](#error-handling)
- [Advanced Usage](#advanced-usage)
- [⚠️ ATTOM Naming Nuances](#%EF%B8%8F-attom-naming-nuances)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Overview

go-attom wraps the ATTOM Property API with a composable client that embraces Go best practices: contexts on every call, functional options for query parameters, dependency injection for HTTP clients, and rich error handling. The goal is to make ATTOM's real estate data easy to integrate in production services.

## Key Features

- **Comprehensive Property API coverage** – 55+ endpoints for property detail, ownership, mortgages, assessments, AVMs, sales history, trends, school lookups, geographic areas, points of interest, community data, parcel tiles, hazards, preforeclosure information, and sale comparables.
- **Functional option builders** – Compose ATTOM query parameters with helpers such as `WithAddress`, `WithAttomID`, `WithLatitudeLongitude`, `WithDateRange`, and `WithPropertyType`.
- **Mockable client** – Inject your own `http.Client` implementation for testing or advanced networking requirements.
- **Consistent error handling** – Detailed errors with contextual wrapping and specialized `property.Error` values when the ATTOM API responds with non-2xx codes.
- **Zero logging in library code** – The client returns rich errors and leaves logging decisions to your application.

## Installation

```bash
go get github.com/my-eq/go-attom
```

## Quick Start

```go
package main

import (
        "context"
        "errors"
        "fmt"
        "log"
        "os"

        "github.com/my-eq/go-attom/pkg/client"
        "github.com/my-eq/go-attom/pkg/property"
)

func main() {
        apiKey := os.Getenv("ATTOM_API_KEY")
        if apiKey == "" {
                log.Fatal("set ATTOM_API_KEY before running the example")
        }

        attomClient := client.New(apiKey, nil)
        propertyService := property.NewService(attomClient)

        ctx := context.Background()
        detail, err := propertyService.GetPropertyDetail(
                ctx,
                property.WithAddress("123 Main St Springfield IL 62704"),
        )
        if err != nil {
                var apiErr *property.Error
                if errors.As(err, &apiErr) {
                        log.Fatalf("ATTOM error: %s", apiErr)
                }
                log.Fatal(err)
        }

        if len(detail.Property) > 0 {
                prop := detail.Property[0]
                if prop.Address != nil {
                        fmt.Printf("City: %s\n", safeString(prop.Address.City))
                }
                if prop.Building != nil && prop.Building.Rooms != nil {
                        fmt.Printf("Bedrooms: %s\n", safeInt(prop.Building.Rooms.Beds))
                }
        }
}

func safeString(v *string) string {
        if v == nil {
                return "(unknown)"
        }
        return *v
}

func safeInt(v *int) string {
        if v == nil {
                return "(unknown)"
        }
        return fmt.Sprint(*v)
}
```

### Look up an ATTOM ID for an address

```go
idResp, err := propertyService.GetPropertyID(
        ctx,
        "1600 Pennsylvania Ave NW Washington DC",
        property.WithPageSize(5),
)
if err != nil {
        return err
}
for _, identifier := range idResp.Identifier {
        fmt.Println("ATTOM ID:", safeString(identifier.AttomID))
}
```

### Search for schools near a coordinate

```go
schools, err := propertyService.SearchSchools(
        ctx,
        property.WithLatitudeLongitude(34.0522, -118.2437),
        property.WithRadius(1.5), // miles
)
if err != nil {
        return err
}
fmt.Printf("Found %d schools\n", len(schools.School))
```

> [!NOTE]
> Helpers like `safeString` are optional—responses expose pointers so you can detect missing data.

## Property API Coverage

All endpoints use ATTOM API version `v4` unless noted otherwise. Some endpoints use specialized versions: sale comparables use `v2`, preforeclosure uses `v3`, and transportation noise hazards use `v1.0.0`. Descriptions come from the official ATTOM swagger definitions included in this repository, updated to reflect the latest API versions as of November 2025.

### Property Profiles & Basics

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetPropertyID` | `/v4/property/id` | Returns properties that match search criteria and "get a list of property IDs within a specific geography that have a specific number of beds".[docs/attom/swagger/propertyapi_property.pretty.json:117-144](docs/attom/swagger/propertyapi_property.pretty.json#L117-L144) |
| `GetPropertyDetail` | `/v4/property/detail` | Returns property details for a supplied ATTOM ID.[docs/attom/swagger/propertyapi_property.pretty.json:115-148](docs/attom/swagger/propertyapi_property.pretty.json#L115-L148) |
| `GetPropertyAddress` | `/v4/property/address` | Returns properties within a ZIP code and supports narrowing with property type and ordering options.[docs/attom/swagger/propertyapi_property.pretty.json:148-188](docs/attom/swagger/propertyapi_property.pretty.json#L148-L188) |
| `GetPropertySnapshot` | `/v4/property/snapshot` | Returns property snapshots that match filters such as city, size range, and property type.[docs/attom/swagger/propertyapi_property.pretty.json:188-227](docs/attom/swagger/propertyapi_property.pretty.json#L188-L227) |
| `GetBasicProfile` | `/v4/property/basicprofile` | Returns basic property information plus the most recent transaction and tax data for an address.[docs/attom/swagger/propertyapi_property.pretty.json:227-269](docs/attom/swagger/propertyapi_property.pretty.json#L227-L269) |
| `GetExpandedProfile` | `/v4/property/expandedprofile` | Returns detailed property information with the latest transaction and taxes for an address.[docs/attom/swagger/propertyapi_property.pretty.json:269-309](docs/attom/swagger/propertyapi_property.pretty.json#L269-L309) |
| `GetBuildingPermits` | `/v4/property/buildingpermits` | Returns basic property information and detailed building permits for an address.[docs/attom/swagger/propertyapi_property.pretty.json:309-352](docs/attom/swagger/propertyapi_property.pretty.json#L309-L352) |
| `GetAllEventsDetail` | `/propertyapi/v1.0.0/allevents/detail` | Returns the full timeline of events that occurred on a property, including cross-domain activity.[docs/attom/swagger/allevents_extended_v4.pretty.json:5-47](docs/attom/swagger/allevents_extended_v4.pretty.json#L5-L47) |

### Ownership, Mortgage, and Schools

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetDetailWithSchools` | `/v4/property/detailwithschools` | Returns property details together with the schools inside the attendance zones for the address.[docs/attom/swagger/propertyapi_school.pretty.json:28-70](docs/attom/swagger/propertyapi_school.pretty.json#L28-L70) |
| `GetDetailMortgage` | `/v4/property/detailmortgage` | Returns property detail enriched with mortgage information for the provided address.[pkg/property/service.go:268-288](pkg/property/service.go#L268-L288) |
| `GetDetailOwner` | `/v4/property/detailowner` | Returns property detail enriched with ownership information for the provided address.[pkg/property/service.go:290-307](pkg/property/service.go#L290-L307) |
| `GetDetailMortgageOwner` | `/v4/property/detailmortgageowner` | Returns property detail enriched with combined mortgage and ownership information for the address.[pkg/property/service.go:309-327](pkg/property/service.go#L309-L327) |
| `SearchSchools` | `/v4/school/search` | Returns school listings around an address or coordinate search context.[docs/attom/swagger/propertyapi_school.pretty.json:70-123](docs/attom/swagger/propertyapi_school.pretty.json#L70-L123) |
| `GetSchoolProfile` | `/v4/school/profile` | Returns enriched profile information for an individual school.[docs/attom/swagger/propertyapi_school.pretty.json:123-166](docs/attom/swagger/propertyapi_school.pretty.json#L123-L166) |
| `GetSchoolDistrict` | `/v4/school/district` | Returns school district boundaries and related contact data.[docs/attom/swagger/propertyapi_school.pretty.json:166-209](docs/attom/swagger/propertyapi_school.pretty.json#L166-L209) |

### Sales, Assessments, and Valuations

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetSaleDetail` | `/v4/sale/detail` | Returns detailed sale transaction information for a property identifier.[docs/attom/swagger/propertyapi_sale.pretty.json:5-51](docs/attom/swagger/propertyapi_sale.pretty.json#L5-L51) |
| `GetSaleSnapshot` | `/v4/sale/snapshot` | Returns a sale snapshot summarizing recent transaction metrics for a property.[docs/attom/swagger/propertyapi_sale.pretty.json:51-95](docs/attom/swagger/propertyapi_sale.pretty.json#L51-L95) |
| `GetAssessmentDetail` | `/v4/assessment/detail` | Returns detailed assessment, tax, and market value data.[docs/attom/swagger/propertyapi_assessment.pretty.json:5-52](docs/attom/swagger/propertyapi_assessment.pretty.json#L5-L52) |
| `GetAssessmentSnapshot` | `/v4/assessment/snapshot` | Returns assessment snapshot metrics for a property identifier.[docs/attom/swagger/propertyapi_assessment.pretty.json:52-95](docs/attom/swagger/propertyapi_assessment.pretty.json#L52-L95) |
| `GetAssessmentHistory` | `/v4/assessmenthistory/detail` | Returns historical assessment records for the property.[docs/attom/swagger/propertyapi_assessmenthistory.pretty.json:5-48](docs/attom/swagger/propertyapi_assessmenthistory.pretty.json#L5-L48) |
| `GetAVMSnapshot` | `/v4/avm/snapshot` | Returns automated valuation model (AVM) snapshot values and confidence scoring.[docs/attom/swagger/propertyapi_avm.pretty.json:5-49](docs/attom/swagger/propertyapi_avm.pretty.json#L5-L49) |
| `GetAttomAVMDetail` | `/v4/attomavm/detail` | Returns ATTOM AVM detail including percentile and scoring metrics.[docs/attom/swagger/propertyapi_attomavm.pretty.json:5-47](docs/attom/swagger/propertyapi_attomavm.pretty.json#L5-L47) |
| `GetAVMHistory` | `/v4/avmhistory/detail` | Returns month-by-month AVM history for the property.[docs/attom/swagger/propertyapi_avmhistory.pretty.json:5-49](docs/attom/swagger/propertyapi_avmhistory.pretty.json#L5-L49) |
| `GetRentalAVM` | `/v4/valuation/rentalavm` | Returns rental AVM valuations and rent ranges.[docs/attom/swagger/propertyapi_valuation.pretty.json:5-46](docs/attom/swagger/propertyapi_valuation.pretty.json#L5-L46) |
| `GetSaleComparablesByAddress` | `/property/v2/salescomparables/address` | Returns comparable sales data for a given address using v2 API.[pkg/property/service.go:329-349](pkg/property/service.go#L329-L349) |
| `GetSaleComparablesByAPN` | `/property/v2/salescomparables/apn` | Returns comparable sales data for a given APN using v2 API.[pkg/property/service.go:351-371](pkg/property/service.go#L351-L371) |
| `GetSaleComparablesByPropID` | `/property/v2/salescomparables/propid` | Returns comparable sales data for a given property ID using v2 API.[pkg/property/service.go:908-915](pkg/property/service.go#L908-L915) |

### Sales History & Trends

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetSalesHistoryDetail` | `/v4/saleshistory/detail` | Returns the full sales history for a property.[docs/attom/swagger/propertyapi_saleshistory.pretty.json:5-50](docs/attom/swagger/propertyapi_saleshistory.pretty.json#L5-L50) |
| `GetSalesHistorySnapshot` | `/v4/saleshistory/snapshot` | Returns a snapshot of historical transactions for quick lookups.[docs/attom/swagger/propertyapi_saleshistory.pretty.json:50-92](docs/attom/swagger/propertyapi_saleshistory.pretty.json#L50-L92) |
| `GetSalesHistoryBasic` | `/v4/saleshistory/basichistory` | Returns a lightweight transaction history for rapid searches.[docs/attom/swagger/propertyapi_saleshistory.pretty.json:92-136](docs/attom/swagger/propertyapi_saleshistory.pretty.json#L92-L136) |
| `GetSalesHistoryExpanded` | `/v4/saleshistory/expandedhistory` | Returns expanded transaction history including document metadata.[docs/attom/swagger/propertyapi_saleshistory.pretty.json:136-180](docs/attom/swagger/propertyapi_saleshistory.pretty.json#L136-L180) |
| `GetSalesTrendSnapshot` | `/v4/salestrend/snapshot` | Returns sales trend metrics for a specified geographic ID.[docs/attom/swagger/propertyapi_salestrend.pretty.json:5-47](docs/attom/swagger/propertyapi_salestrend.pretty.json#L5-L47) |
| `GetTransactionSalesTrend` | `/v4/transaction/salestrend` | Returns transaction-oriented sales trend metrics across geographies.[docs/attom/swagger/propertyapi_transaction.pretty.json:5-47](docs/attom/swagger/propertyapi_transaction.pretty.json#L5-L47) |

## Geographic Area API Coverage

All endpoints use ATTOM API version `v2.0.0` unless noted otherwise.

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetCountyLookup` | `/areaapi/v2.0.0/county/lookup` | Returns county information for a given state.[docs/attom/swagger/propertyapi_area.pretty.json:5-47](docs/attom/swagger/propertyapi_area.pretty.json#L5-L47) |
| `GetAreaHierarchyLookup` | `/areaapi/v2.0.0/area/hierarchy/lookup` | Returns geographic hierarchy information for an area.[docs/attom/swagger/propertyapi_area.pretty.json:47-89](docs/attom/swagger/propertyapi_area.pretty.json#L47-L89) |
| `GetStateAreaLookup` | `/areaapi/v2.0.0/area/state/lookup` | Returns state area information and boundaries.[docs/attom/swagger/propertyapi_area.pretty.json:89-131](docs/attom/swagger/propertyapi_area.pretty.json#L89-L131) |
| `GetAreaBoundaryDetail` | `/areaapi/v2.0.0/area/boundary/detail` | Returns detailed boundary information in GeoJSON or WKT format.[docs/attom/swagger/propertyapi_area.pretty.json:131-173](docs/attom/swagger/propertyapi_area.pretty.json#L131-L173) |
| `GetLegacyGeoIDLookup` | `/areaapi/v2.0.0/area/geoId/legacyLookup` | Returns geographic ID information using legacy geocoding.[docs/attom/swagger/propertyapi_area.pretty.json:173-215](docs/attom/swagger/propertyapi_area.pretty.json#L173-L215) |
| `GetGeoIDLookup` | `/areaapi/v2.0.0/area/geoId/Lookup` | Returns geographic ID information using current geocoding.[docs/attom/swagger/propertyapi_area.pretty.json:215-257](docs/attom/swagger/propertyapi_area.pretty.json#L215-L257) |

## Points of Interest API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `SearchPOIByAddress` | `/poisearch/v2.0.0/poi/Street+Address/` | Searches for points of interest by street address.[docs/attom/swagger/propertyapi_poi.pretty.json:5-47](docs/attom/swagger/propertyapi_poi.pretty.json#L5-L47) |
| `SearchPOIByGeography` | `/poisearch/v2.0.0/poi/Geography/` | Searches for points of interest by geographic area (ZIP code).[docs/attom/swagger/propertyapi_poi.pretty.json:47-89](docs/attom/swagger/propertyapi_poi.pretty.json#L47-L89) |
| `SearchPOIByPoint` | `/poisearch/v2.0.0/poi/Point/` | Searches for points of interest by coordinate point.[docs/attom/swagger/propertyapi_poi.pretty.json:89-131](docs/attom/swagger/propertyapi_poi.pretty.json#L89-L131) |
| `GetNeighborhoodPOI` | `/v4/neighborhood/poi` | Returns points of interest within neighborhood boundaries.[docs/attom/swagger/propertyapi_poi.pretty.json:131-173](docs/attom/swagger/propertyapi_poi.pretty.json#L131-L173) |
| `GetPOICategoryLookup` | `/v4/neighborhood/poi/categorylookup` | Returns available POI categories for filtering.[docs/attom/swagger/propertyapi_poi.pretty.json:173-215](docs/attom/swagger/propertyapi_poi.pretty.json#L173-L215) |

## Community API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetCommunityFull` | `/communityapi/v2.0.0/area/full` | Returns comprehensive community data including demographics, economics, education, housing, and climate.[docs/attom/swagger/propertyapi_community.pretty.json:5-47](docs/attom/swagger/propertyapi_community.pretty.json#L5-L47) |
| `GetNeighborhoodCommunity` | `/v4/neighborhood/community` | Returns community profile data for neighborhood areas.[docs/attom/swagger/propertyapi_community.pretty.json:47-89](docs/attom/swagger/propertyapi_community.pretty.json#L47-L89) |

## Parcel Tiles & Hazard API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetParcelTiles` | `/parceltiles/{z}/{x}/{y}.png` | Returns parcel boundary raster tiles in PNG format.[docs/attom/swagger/propertyapi_parceltile.pretty.json:5-47](docs/attom/swagger/propertyapi_parceltile.pretty.json#L5-L47) |
| `GetHazardDetail` | `/v4/property/hazarddetail` | Returns natural hazard risk data for properties.[docs/attom/swagger/propertyapi_hazard.pretty.json:5-47](docs/attom/swagger/propertyapi_hazard.pretty.json#L5-L47) |
| `GetTransportationNoise` | `/propertyapi/v1.0.0/transportationnoise/detail` | Returns transportation noise data for geographic areas using v1.0.0 API.[pkg/property/service.go:500-518](pkg/property/service.go#L500-L518) |

## Preforeclosure API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetPreforeclosureDetail` | `/property/v3/preforeclosure` | Returns preforeclosure information and notices for properties using v3 API.[docs/attom/swagger/propertyapi_preforeclosure.pretty.json:5-47](docs/attom/swagger/propertyapi_preforeclosure.pretty.json#L5-L47) |

## Utility & Reference API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetEnumerationsDetail` | `/v4/enumerations/detail` | Returns controlled vocabulary values for API parameters. Use the `field` parameter to get valid values for specific fields like `propertytype`, `orderby`, etc.[pkg/property/service.go:720-740](pkg/property/service.go#L720-L740) |

## Building Requests with Options

Property requests accept a flexible list of functional options:

- **Property identifiers** – `WithAttomID`, `WithPropertyID`, `WithFIPSAndAPN`, and `WithAddressLines` for ATTOM identifiers and assessor parcel numbers.[pkg/property/options.go:26-65](pkg/property/options.go#L26-L65)
- **Geographic search** – `WithLatitudeLongitude`, `WithRadius`, `WithPostalCode`, `WithGeoID`, and `WithGeoIDV4` to target coordinates and geographic codes.[pkg/property/options.go:67-117](pkg/property/options.go#L67-L117)
- **Filtering** – `WithBedsRange`, `WithBathsRange`, `WithSaleAmountRange`, `WithPropertyType`, `WithPropertyIndicator`, `WithUniversalSizeRange`, `WithYearBuiltRange`, `WithLotSize1Range`, and `WithLotSize2Range` for ATTOM's numeric filters.[pkg/property/options.go:119-223](pkg/property/options.go#L119-L223)
- **Date windows** – `WithDateRange` (MM/DD) and `WithISODateRange` (YYYY-MM-DD) cover endpoints that expect legacy or ISO date formats.[pkg/property/options.go:225-280](pkg/property/options.go#L225-L280)
- **Pagination and sorting** – `WithPage`, `WithPageSize`, and `WithOrderBy` mirror ATTOM's paging and ordering controls.[pkg/property/options.go:282-327](pkg/property/options.go#L282-L327)
- **Custom parameters** – `WithString`, `WithStringSlice`, and `WithAdditionalParam` map to ATTOM's long tail of specialized filters.[pkg/property/options.go:17-47](pkg/property/options.go#L17-L47)[pkg/property/options.go:329-371](pkg/property/options.go#L329-L371)

These helpers ensure requests include the correct parameter names and formatting required by ATTOM.

## Error Handling

- `property.ErrMissingParameter` is returned when required query parameters are omitted by the caller.[pkg/property/errors.go:5-15](pkg/property/errors.go#L5-L15)
- API errors use `*property.Error`, which captures the HTTP status code, parsed ATTOM status block, and raw response body for debugging.[pkg/property/errors.go:17-67](pkg/property/errors.go#L17-L67)
- All public methods wrap lower-level errors with context using Go's `%w` semantics so you can use `errors.Is`/`errors.As`. The shared `doGet` helper centralizes the behavior.[pkg/property/service.go:52-121](pkg/property/service.go#L52-L121)

## Advanced Usage

### Customize the HTTP client

Inject your own `http.Client` (or any type implementing `client.HTTPClient`) when building the ATTOM client. This is ideal for adding retries, circuit breakers, or observability instrumentation.

```go
httpClient := &http.Client{Timeout: 10 * time.Second}
attomClient := client.New(apiKey, httpClient)
```

The constructor falls back to a 30-second timeout client when you pass `nil`, keeping defaults safe for production.[pkg/client/client.go:36-58](pkg/client/client.go#L36-L58)

### Get controlled vocabulary values

Use `GetEnumerationsDetail` to discover valid values for API parameters. This is especially useful for fields like `propertytype` that have many possible values:

```go
enums, err := propertyService.GetEnumerationsDetail(
        ctx,
        property.WithString("field", "propertytype"),
)
if err != nil {
        return err
}

fmt.Printf("Valid property types:\n")
for _, enum := range enums.Enumerations {
        fmt.Printf("  - %s\n", safeString(enum.Value))
}
```

This endpoint helps ensure your requests use valid enum values and can be used for building dynamic UIs or validation logic.[pkg/property/service.go:720-740](pkg/property/service.go#L720-L740)

### Override the base URL for staging and proxies

ATTOM provides dedicated gateways per environment. Use `client.WithBaseURL` to point at a staging cluster or run through an internal proxy:

```go
attomClient := client.New(apiKey, nil, client.WithBaseURL("https://staging.attomdata.com/"))
```

The option normalizes trailing slashes to keep request construction predictable.[pkg/client/client.go:24-44](pkg/client/client.go#L24-L44)

### Inspect detailed API failures

When ATTOM returns a non-2xx response, go-attom unmarshals the status payload into `property.Error`, preserving the HTTP code, ATTOM status block, and raw JSON to help with support tickets or sandbox debugging.[pkg/property/service.go:74-118](pkg/property/service.go#L74-L118)[pkg/property/errors.go:17-67](pkg/property/errors.go#L17-L67)

## ⚠️ ATTOM Naming Nuances

ATTOM mixes lower-case, camelCase, and uppercase tokens in both query parameters and JSON payloads. The client mirrors those quirks so requests land correctly:

- Query helpers intentionally set parameters like `attomid`, `APN`, `geoIdV4`, and `propertyIndicator` using ATTOM's exact casing.[pkg/property/options.go:49-113](pkg/property/options.go#L49-L113)[pkg/property/options.go:147-177](pkg/property/options.go#L147-L177)
- Response structs keep ATTOM's field names in their `json` tags (`attomId`, `postalCode`, `areaSqFt`) while surfacing idiomatic Go field names (`AttomID`, `PostalCode`, `AreaSquareFeet`).[pkg/property/models.go:7-53](pkg/property/models.go#L7-L53)

When debugging raw JSON, match the ATTOM documentation rather than the Go field names to avoid confusion.

## Development

```bash
# Run tests with race detection
go test ./... -race -v

# Compile the library
go build ./...

# Install golangci-lint exactly like CI (Go 1.25 toolchain requirement)
GOLANGCI_LINT_VERSION=v1.63.1
TMPDIR=$(mktemp -d)
git clone --depth 1 --branch "${GOLANGCI_LINT_VERSION}" https://github.com/golangci/golangci-lint.git "$TMPDIR/golangci-lint"
sed -i "s/go 1\.[0-9][0-9]\.[0-9]/go 1.25.1/" "$TMPDIR/golangci-lint/go.mod"
(cd "$TMPDIR/golangci-lint" && go build -o "$(go env GOBIN)/golangci-lint" ./cmd/golangci-lint)
rm -rf "$TMPDIR"

golangci-lint run --timeout=5m

# Check markdown formatting
markdownlint-cli2 "**/*.md"
```

See [`API_IMPLEMENTATION_SUMMARY.md`](API_IMPLEMENTATION_SUMMARY.md#L9-L118) for a full breakdown of every ATTOM API group and [`docs/GITHUB_ACTIONS_SUMMARY.md`](docs/GITHUB_ACTIONS_SUMMARY.md#L1-L158) for CI/CD pipeline details.

## Contributing

Contributions are welcome! Please read the [project guidelines](.github/copilot-instructions.md) before opening an issue or pull request.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
