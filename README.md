# go-attom

[![CI](https://github.com/my-eq/go-attom/actions/workflows/ci.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/ci.yml)
[![Security](https://github.com/my-eq/go-attom/actions/workflows/security.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/security.yml)
[![README Lint](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml)
[![codecov](https://codecov.io/gh/my-eq/go-attom/branch/main/graph/badge.svg)](https://codecov.io/gh/my-eq/go-attom)
[![Go Report Card](https://goreportcard.com/badge/github.com/my-eq/go-attom)](https://goreportcard.com/report/github.com/my-eq/go-attom)
[![Go Reference](https://pkg.go.dev/badge/github.com/my-eq/go-attom.svg)](https://pkg.go.dev/github.com/my-eq/go-attom)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

`go-attom` is a strongly typed, mock-friendly Go client for the [ATTOM Data API](https://api.gateway.attomdata.com/). The library emphasizes idiomatic Go ergonomics while delivering complete Property API coverage for property profiles, ownership, mortgages, assessments, valuations, schools, and historical trends.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Authenticate and Instantiate](#authenticate-and-instantiate)
- [Quick Start](#quick-start)
  - [Look up an ATTOM ID for an address](#look-up-an-attom-id-for-an-address)
  - [Search for schools near a coordinate](#search-for-schools-near-a-coordinate)
- [Property API Coverage](#property-api-coverage)
  - [Property Profiles and Core Records](#property-profiles-and-core-records)
  - [Ownership, Mortgage, and Schools](#ownership-mortgage-and-schools)
  - [Sales, Assessments, and Valuations](#sales-assessments-and-valuations)
  - [Sales History & Trends](#sales-history--trends)
- [Composing Requests with Options](#composing-requests-with-options)
- [Error Handling](#error-handling)
- [Warning: ATTOM Naming Conventions](#warning-attom-naming-conventions)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)
- [Advanced Usage](#advanced-usage)

## Overview

`go-attom` is a strongly typed, mock-friendly Go client for the [ATTOM Data API](https://api.gateway.attomdata.com/). The library focuses on clean, idiomatic Go and provides deep coverage of the Property API including property profiles, sales, assessments, valuations, schools, and historical trends (see the [API implementation summary](API_IMPLEMENTATION_SUMMARY.md#L9-L118) for the full endpoint catalog). It wraps the ATTOM Property API with a composable client that embraces Go best practices: contexts on every call, functional options for query parameters, dependency injection for HTTP clients, and rich error handling so ATTOM's real estate data is easy to integrate in production services.

## Key Features

- **Deep Property API coverage** – 30+ endpoints spanning property profiles, ownership, mortgage liens, building permits, assessments, AVMs, sales history, and geo trends. The coverage mirrors the official [ATTOM Property API catalog](API_IMPLEMENTATION_SUMMARY.md).
- **Functional option builders** – Compose ATTOM query parameters with helpers such as `WithAddress`, `WithAttomID`, `WithLatitudeLongitude`, `WithDateRange`, and `WithPropertyType` from [`pkg/property/options.go`](pkg/property/options.go).
- **Mockable foundation** – Inject any `http.Client`-compatible type through `client.New`, or supply your own interface implementation for advanced testing scenarios (see [`pkg/client/client.go`](pkg/client/client.go)).
- **Consistent error propagation** – Library methods wrap errors with context and return structured `*property.Error` values that capture status metadata and raw responses (see [`pkg/property/service.go`](pkg/property/service.go)).
- **No hidden logging** – The client surfaces detailed errors while leaving logging decisions entirely to consuming applications.

## Getting Started

### Prerequisites

- An ATTOM API key with access to the Property API portfolio.
- Go 1.25.3 or newer (matches the module’s go directive).

### Installation

```bash
go get github.com/my-eq/go-attom
```

### Authenticate and Instantiate

```go
attomClient := client.New(apiKey, nil) // defaults to *http.Client with a 30s timeout
propertyService := property.NewService(attomClient)
```

Passing `nil` uses an internal `*http.Client` with a 30-second timeout; you can inject your own implementation to customize timeouts, retries, or tracing (see [`pkg/client/client.go`](pkg/client/client.go)).

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

> [!NOTE]
> Helpers like `safeString` are optional—responses expose pointers so you can detect missing data (see [`pkg/property/models.go`](pkg/property/models.go)).

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

All endpoints use ATTOM API version `v4` unless noted otherwise. Descriptions come from the official ATTOM swagger definitions included in this repository (note: swagger files are historical v1.0.0 specifications but endpoint paths have been updated to reflect current API usage).

### Property Profiles and Core Records

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetPropertyID` | `/v4/property/id` | Returns properties that match search criteria and "get a list of property IDs within a specific geography that have a specific number of beds".【F:docs/attom/swagger/propertyapi_property.pretty.json†L117-L144】 |
| `GetPropertyDetail` | `/v4/property/detail` | Returns property details for a supplied ATTOM ID.【F:docs/attom/swagger/propertyapi_property.pretty.json†L115-L148】 |
| `GetPropertyAddress` | `/v4/property/address` | Returns properties within a ZIP code and supports narrowing with property type and ordering options.【F:docs/attom/swagger/propertyapi_property.pretty.json†L148-L188】 |
| `GetPropertySnapshot` | `/v4/property/snapshot` | Returns property snapshots that match filters such as city, size range, and property type.【F:docs/attom/swagger/propertyapi_property.pretty.json†L188-L227】 |
| `GetBasicProfile` | `/v4/property/basicprofile` | Returns basic property information plus the most recent transaction and tax data for an address.【F:docs/attom/swagger/propertyapi_property.pretty.json†L227-L269】 |
| `GetExpandedProfile` | `/v4/property/expandedprofile` | Returns detailed property information with the latest transaction and taxes for an address.【F:docs/attom/swagger/propertyapi_property.pretty.json†L269-L309】 |
| `GetBuildingPermits` | `/v4/property/buildingpermits` | Returns basic property information and detailed building permits for an address.【F:docs/attom/swagger/propertyapi_property.pretty.json†L309-L352】 |
| `GetAllEventsDetail` | `/v4/property/detail` | Returns the full timeline of events that occurred on a property, including cross-domain activity.【F:docs/attom/swagger/propertyapi_allevents.pretty.json†L5-L47】 |

### Ownership, Mortgage, and Schools

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetDetailWithSchools` | `/v4/property/detailwithschools` | Returns property details together with the schools inside the attendance zones for the address.【F:docs/attom/swagger/propertyapi_school.pretty.json†L28-L70】 |
| `GetDetailMortgage` | `/v4/property/detailmortgage` | Returns property detail enriched with mortgage information for the provided address.【F:pkg/property/service.go†L268-L288】 |
| `GetDetailOwner` | `/v4/property/detailowner` | Returns property detail enriched with ownership information for the provided address.【F:pkg/property/service.go†L290-L307】 |
| `GetDetailMortgageOwner` | `/v4/property/detailmortgageowner` | Returns property detail enriched with combined mortgage and ownership information for the address.【F:pkg/property/service.go†L309-L327】 |
| `SearchSchools` | `/v4/school/search` | Returns school listings around an address or coordinate search context.【F:docs/attom/swagger/propertyapi_school.pretty.json†L70-L123】 |
| `GetSchoolProfile` | `/v4/school/profile` | Returns enriched profile information for an individual school.【F:docs/attom/swagger/propertyapi_school.pretty.json†L123-L166】 |
| `GetSchoolDistrict` | `/v4/school/district` | Returns school district boundaries and related contact data.【F:docs/attom/swagger/propertyapi_school.pretty.json†L166-L209】 |

### Sales, Assessments, and Valuations

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetSaleDetail` | `/v4/sale/detail` | Returns detailed sale transaction information for a property identifier.【F:docs/attom/swagger/propertyapi_sale.pretty.json†L5-L51】 |
| `GetSaleSnapshot` | `/v4/sale/snapshot` | Returns a sale snapshot summarizing recent transaction metrics for a property.【F:docs/attom/swagger/propertyapi_sale.pretty.json†L51-L95】 |
| `GetAssessmentDetail` | `/v4/assessment/detail` | Returns detailed assessment, tax, and market value data.【F:docs/attom/swagger/propertyapi_assessment.pretty.json†L5-L52】 |
| `GetAssessmentSnapshot` | `/v4/assessment/snapshot` | Returns assessment snapshot metrics for a property identifier.【F:docs/attom/swagger/propertyapi_assessment.pretty.json†L52-L95】 |
| `GetAssessmentHistory` | `/v4/assessmenthistory/detail` | Returns historical assessment records for the property.【F:docs/attom/swagger/propertyapi_assessmenthistory.pretty.json†L5-L48】 |
| `GetAVMSnapshot` | `/v4/avm/snapshot` | Returns automated valuation model (AVM) snapshot values and confidence scoring.【F:docs/attom/swagger/propertyapi_avm.pretty.json†L5-L49】 |
| `GetAttomAVMDetail` | `/v4/attomavm/detail` | Returns ATTOM AVM detail including percentile and scoring metrics.【F:docs/attom/swagger/propertyapi_attomavm.pretty.json†L5-L47】 |
| `GetAVMHistory` | `/v4/avmhistory/detail` | Returns month-by-month AVM history for the property.【F:docs/attom/swagger/propertyapi_avmhistory.pretty.json†L5-L49】 |
| `GetRentalAVM` | `/v4/valuation/rentalavm` | Returns rental AVM valuations and rent ranges.【F:docs/attom/swagger/propertyapi_valuation.pretty.json†L5-L46】 |

### Sales History & Trends

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetSalesHistoryDetail` | `/v4/saleshistory/detail` | Returns the full sales history for a property.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L5-L50】 |
| `GetSalesHistorySnapshot` | `/v4/saleshistory/snapshot` | Returns a snapshot of historical transactions for quick lookups.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L50-L92】 |
| `GetSalesHistoryBasic` | `/v4/saleshistory/basichistory` | Returns a lightweight transaction history for rapid searches.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L92-L136】 |
| `GetSalesHistoryExpanded` | `/v4/saleshistory/expandedhistory` | Returns expanded transaction history including document metadata.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L136-L180】 |
| `GetSalesTrendSnapshot` | `/v4/salestrend/snapshot` | Returns sales trend metrics for a specified geographic ID.【F:docs/attom/swagger/propertyapi_salestrend.pretty.json†L5-L47】 |
| `GetTransactionSalesTrend` | `/v4/transaction/salestrend` | Returns transaction-oriented sales trend metrics across geographies.【F:docs/attom/swagger/propertyapi_transaction.pretty.json†L5-L47】 |

## Geographic Area API Coverage

All endpoints use ATTOM API version `v2.0.0` unless noted otherwise.

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetCountyLookup` | `/areaapi/v2.0.0/county/lookup` | Returns county information for a given state.【F:docs/attom/swagger/propertyapi_area.pretty.json†L5-L47】 |
| `GetAreaHierarchyLookup` | `/areaapi/v2.0.0/area/hierarchy/lookup` | Returns geographic hierarchy information for an area.【F:docs/attom/swagger/propertyapi_area.pretty.json†L47-L89】 |
| `GetStateAreaLookup` | `/areaapi/v2.0.0/area/state/lookup` | Returns state area information and boundaries.【F:docs/attom/swagger/propertyapi_area.pretty.json†L89-L131】 |
| `GetAreaBoundaryDetail` | `/areaapi/v2.0.0/area/boundary/detail` | Returns detailed boundary information in GeoJSON or WKT format.【F:docs/attom/swagger/propertyapi_area.pretty.json†L131-L173】 |
| `GetLegacyGeoIDLookup` | `/areaapi/v2.0.0/area/geoId/legacyLookup` | Returns geographic ID information using legacy geocoding.【F:docs/attom/swagger/propertyapi_area.pretty.json†L173-L215】 |
| `GetGeoIDLookup` | `/areaapi/v2.0.0/area/geoId/Lookup` | Returns geographic ID information using current geocoding.【F:docs/attom/swagger/propertyapi_area.pretty.json†L215-L257】 |
| `GetTransportationNoise` | `/areaapi/v2.0.0/area/transportationnoise/detail` | Returns transportation noise data for geographic areas.【F:pkg/property/service.go†L500-L518】 |

## Points of Interest API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `SearchPOIByAddress` | `/poisearch/v2.0.0/poi/Street+Address/` | Searches for points of interest by street address.【F:docs/attom/swagger/propertyapi_poi.pretty.json†L5-L47】 |
| `SearchPOIByGeography` | `/poisearch/v2.0.0/poi/Geography/` | Searches for points of interest by geographic area (ZIP code).【F:docs/attom/swagger/propertyapi_poi.pretty.json†L47-L89】 |
| `SearchPOIByPoint` | `/poisearch/v2.0.0/poi/Point/` | Searches for points of interest by coordinate point.【F:docs/attom/swagger/propertyapi_poi.pretty.json†L89-L131】 |
| `GetNeighborhoodPOI` | `/v4/neighborhood/poi` | Returns points of interest within neighborhood boundaries.【F:docs/attom/swagger/propertyapi_poi.pretty.json†L131-L173】 |
| `GetPOICategoryLookup` | `/v4/neighborhood/poi/categorylookup` | Returns available POI categories for filtering.【F:docs/attom/swagger/propertyapi_poi.pretty.json†L173-L215】 |

## Community API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetCommunityFull` | `/communityapi/v2.0.0/area/full` | Returns comprehensive community data including demographics, economics, education, housing, and climate.【F:docs/attom/swagger/propertyapi_community.pretty.json†L5-L47】 |
| `GetNeighborhoodCommunity` | `/v4/neighborhood/community` | Returns community profile data for neighborhood areas.【F:docs/attom/swagger/propertyapi_community.pretty.json†L47-L89】 |

## Parcel Tiles & Hazard API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetParcelTiles` | `/parceltiles/{z}/{x}/{y}.png` | Returns parcel boundary raster tiles in PNG format.【F:docs/attom/swagger/propertyapi_parceltile.pretty.json†L5-L47】 |
| `GetHazardDetail` | `/v4/property/hazarddetail` | Returns natural hazard risk data for properties.【F:docs/attom/swagger/propertyapi_hazard.pretty.json†L5-L47】 |

## Preforeclosure API Coverage

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetPreforeclosureDetail` | `/v4/property/preforeclosure` | Returns preforeclosure information and notices for properties.【F:docs/attom/swagger/propertyapi_preforeclosure.pretty.json†L5-L47】 |

## Composing Requests with Options

Property methods accept a flexible list of functional options defined in [`pkg/property/options.go`](pkg/property/options.go).

- **Property identifiers** – `WithAttomID`, `WithPropertyID`, `WithFIPSAndAPN`, and `WithAddressLines` automatically populate ATTOM’s identifier parameters.
- **Geographic search** – `WithLatitudeLongitude`, `WithRadius`, `WithPostalCode`, `WithGeoID`, and `WithGeoIDV4` generate latitude/longitude or GeoID queries with ATTOM’s expected units and casing.
- **Filtering** – Helpers such as `WithBedsRange`, `WithBathsRange`, `WithSaleAmountRange`, `WithPropertyType`, `WithPropertyIndicator`, `WithUniversalSizeRange`, `WithYearBuiltRange`, `WithLotSize1Range`, and `WithLotSize2Range` enforce allowed parameter combinations.
- **Date windows** – `WithDateRange` (MM/DD) and `WithISODateRange` (YYYY-MM-DD) translate human-friendly inputs into the `start*`/`end*` parameters expected by ATTOM.
- **Pagination and sorting** – `WithPage`, `WithPageSize`, and `WithOrderBy` manage paging controls and ordering fields.
- **Custom parameters** – `WithString`, `WithStringSlice`, and `WithAdditionalParam` let you add rarely used query strings without manual URL encoding.

These helpers guarantee your requests ship with the correct casing and parameter names—many ATTOM filters are case-sensitive and reject mismatched field names (see [`pkg/property/options.go`](pkg/property/options.go) and [`pkg/property/service.go`](pkg/property/service.go)).

## Error Handling

- `property.ErrMissingParameter` is returned when required query parameters are omitted or incompatible (see [`pkg/property/options.go`](pkg/property/options.go)).
- API failures surface as `*property.Error`, which includes the HTTP status, parsed ATTOM status block, and raw response body for debugging (see [`pkg/property/error.go`](pkg/property/error.go)).
- All public methods wrap underlying errors using Go’s `%w` semantics so you can chain `errors.Is`/`errors.As` checks (see [`pkg/property/service.go`](pkg/property/service.go)).

## Warning: ATTOM Naming Conventions

ATTOM mixes camelCase, lowercase, and legacy uppercase parameters (`fips` + `APN`) across the Property API. The models and option builders in `go-attom` preserve those cases deliberately:

- Struct fields use idiomatic Go names while JSON tags mirror ATTOM’s exact field casing (for example, `Identifier.AttomID` maps to `attomId`, and `Identifier.APN` maps to `apn`). See [`pkg/property/models.go`](pkg/property/models.go) for concrete mappings.
- Query helpers emit ATTOM’s parameter names verbatim (`WithFIPSAndAPN` sets `fips` and `APN`, while `WithGeoIDV4` sets `geoIdV4`). Mixing cases or renaming parameters will trigger `property.ErrMissingParameter` before the request is issued (see [`pkg/property/options.go`](pkg/property/options.go) and [`pkg/property/service.go`](pkg/property/service.go)).

When you introduce new fields or options, copy ATTOM’s casing exactly in JSON tags and query keys to avoid subtle API rejections. For more field-level details, consult the official data dictionary excerpted in [`.github/agents/attom-docs-dump.txt`](.github/agents/attom-docs-dump.txt).

## Development

```bash
# Run tests with race detection
go test ./... -race -v

# Run golangci-lint
golangci-lint run ./...

# Check markdown formatting
markdownlint-cli2 "**/*.md"
```

See [`API_IMPLEMENTATION_SUMMARY.md`](API_IMPLEMENTATION_SUMMARY.md#L9-L118) for a full breakdown of every ATTOM API group and [`docs/GITHUB_ACTIONS_SUMMARY.md`](docs/GITHUB_ACTIONS_SUMMARY.md#L1-L158) for CI/CD pipeline details.

## Contributing

Contributions are welcome! Please read the [project guidelines](.github/copilot-instructions.md) before opening an issue or pull request.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Advanced Usage

- **Inject custom HTTP clients** – Supply a transport that adds retries, tracing, or circuit breaking by passing your implementation to `client.New(apiKey, httpClient)`. The client only requires a `Do(*http.Request)` method (see [`pkg/client/client.go`](pkg/client/client.go)).
- **Override base URLs** – Target sandboxes or mock servers with `client.WithBaseURL` while keeping path composition intact (see [`pkg/client/client.go`](pkg/client/client.go)).
- **Context-driven cancellation** – Every service method accepts a `context.Context`, enabling per-call deadlines or integration with upstream request lifecycles (see [`pkg/property/service.go`](pkg/property/service.go)).
- **Thin mocks for testing** – The `client.HTTPClient` interface makes it straightforward to stub ATTOM responses in unit tests without hitting the network (see [`pkg/client/client.go`](pkg/client/client.go)).
