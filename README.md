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
  - [Sales History and Trends](#sales-history-and-trends)
  - [Events and Permits](#events-and-permits)
- [Composing Requests with Options](#composing-requests-with-options)
- [Error Handling](#error-handling)
- [Warning: ATTOM Naming Conventions](#warning-attom-naming-conventions)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)
- [Advanced Usage](#advanced-usage)

## Overview

`go-attom` is a strongly-typed, mock-friendly Go client for the [ATTOM Data API](https://api.gateway.attomdata.com/). The library focuses on clean, idiomatic Go and provides deep coverage of the Property API including property profiles, sales, assessments, valuations, schools, and historical trends (see the [API implementation summary](API_IMPLEMENTATION_SUMMARY.md#L9-L118) for the full endpoint catalog).

## Key Features

- **Deep Property API coverage** – 30+ endpoints spanning property profiles, ownership, mortgage liens, building permits, assessments, AVMs, sales history, and geo trends. The coverage mirrors the official [ATTOM Property API catalog](API_IMPLEMENTATION_SUMMARY.md).
- **Functional option builders** – Compose ATTOM query parameters with helpers such as `WithAddress`, `WithAttomID`, `WithLatitudeLongitude`, `WithDateRange`, and `WithPropertyType` from [`pkg/property/options.go`](pkg/property/options.go).
- **Mockable foundation** – Inject any `http.Client`-compatible type through `client.New`, or supply your own interface implementation for advanced testing scenarios (see [`pkg/client/client.go`](pkg/client/client.go)).
- **Consistent error propagation** – Library methods wrap errors with context and return structured `*property.Error` values that capture status metadata and raw responses (see [`pkg/property/service.go`](pkg/property/service.go)).
- **No hidden logging** – The client surfaces detailed errors while leaving logging decisions entirely to consuming applications.

## Getting Started

### Prerequisites

- An ATTOM API key with access to the Property API portfolio.
- Go 1.25 or newer (matches the module’s go directive).

### Installation

```bash
go get github.com/my-eq/go-attom
```

### Authenticate and Instantiate

```go
attomClient := client.New(apiKey, nil) // defaults to *http.Client with a 30s timeout
propertyService := property.NewService(attomClient)
```

Passing `nil` uses an internal `*http.Client` with a 30 second timeout; you can inject your own implementation to customize timeouts, retries, or tracing (see [`pkg/client/client.go`](pkg/client/client.go)).

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

## Property API Coverage

All endpoints use ATTOM API version `v1.0.0` unless stated otherwise. Descriptions below align with ATTOM’s official endpoint definitions sourced from the ATTOM Data API documentation dump in [`.github/agents/attom-docs-dump.txt`](.github/agents/attom-docs-dump.txt) and summarized in [API_IMPLEMENTATION_SUMMARY.md](API_IMPLEMENTATION_SUMMARY.md).

### Property Profiles and Core Records

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetPropertyID` | `/propertyapi/v1.0.0/property/id` | Get the ATTOM property ID for a parcel using address or identifier inputs. |
| `GetPropertyDetail` | `/propertyapi/v1.0.0/property/detail` | Detailed property information for structure, lot, tax, ownership, and location. |
| `GetPropertyAddress` | `/propertyapi/v1.0.0/property/address` | Normalized property address and geocoding metadata. |
| `GetPropertySnapshot` | `/propertyapi/v1.0.0/property/snapshot` | Summary snapshot with key property fields for display workflows. |
| `GetBasicProfile` | `/propertyapi/v1.0.0/property/basicprofile` | Basic property profile with marketing-friendly characteristics. |
| `GetExpandedProfile` | `/propertyapi/v1.0.0/property/expandedprofile` | Expanded property profile containing the full ATTOM characteristic set. |
| `GetDetailWithSchools` | `/propertyapi/v1.0.0/property/detailwithschools` | Property detail plus school attendance zone assignments. |
| `GetDetailMortgage` | `/propertyapi/v1.0.0/property/detailmortgage` | Property detail with mortgage loan information. |
| `GetDetailOwner` | `/propertyapi/v1.0.0/property/detailowner` | Property detail with owner mailing, vesting, and occupancy data. |
| `GetDetailMortgageOwner` | `/propertyapi/v1.0.0/property/detailmortgageowner` | Property detail with combined mortgage and owner records. |
| `GetBuildingPermits` | `/propertyapi/v1.0.0/property/buildingpermits` | Building permit history and contractor information for the parcel. |

### Ownership, Mortgage, and Schools

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `SearchSchools` | `/propertyapi/v1.0.0/school/search` | Search for schools near a property or coordinate within a supplied radius. |
| `GetSchoolProfile` | `/propertyapi/v1.0.0/school/profile` | School profile details including ratings, programs, and enrollment. |
| `GetSchoolDistrict` | `/propertyapi/v1.0.0/school/district` | School district boundary and contact information. |
| `GetSchoolDetailWithSchools` | `/propertyapi/v1.0.0/school/detailwithschools` | Property detail including school attendance information. |

### Sales, Assessments, and Valuations

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetSaleDetail` | `/propertyapi/v1.0.0/sale/detail` | Sale transaction details including dates, amounts, and transfer types. |
| `GetSaleSnapshot` | `/propertyapi/v1.0.0/sale/snapshot` | Sales snapshot summary for quick valuation views. |
| `GetAssessmentDetail` | `/propertyapi/v1.0.0/assessment/detail` | Assessment and tax details with appraised and market values. |
| `GetAssessmentSnapshot` | `/propertyapi/v1.0.0/assessment/snapshot` | Assessment summary snapshot. |
| `GetAssessmentHistory` | `/propertyapi/v1.0.0/assessmenthistory/detail` | Historical assessment records over time. |
| `GetAVMSnapshot` | `/propertyapi/v1.0.0/avm/snapshot` | Automated valuation snapshot with confidence metrics. |
| `GetAttomAVMDetail` | `/propertyapi/v1.0.0/attomavm/detail` | Detailed ATTOM AVM data with value ranges and comparables. |
| `GetAVMHistory` | `/propertyapi/v1.0.0/avmhistory/detail` | AVM historical values showing value changes over time. |
| `GetRentalAVM` | `/propertyapi/v1.0.0/valuation/rentalavm` | Rental valuation estimates including rent range and ratios. |

### Sales History and Trends

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetSalesHistoryDetail` | `/propertyapi/v1.0.0/saleshistory/detail` | Detailed sales history with full transaction records. |
| `GetSalesHistorySnapshot` | `/propertyapi/v1.0.0/saleshistory/snapshot` | Sales history snapshot summary. |
| `GetSalesHistoryBasic` | `/propertyapi/v1.0.0/saleshistory/basichistory` | Basic sales history with essential sale facts. |
| `GetSalesHistoryExpanded` | `/propertyapi/v1.0.0/saleshistory/expandedhistory` | Expanded sales history with full document metadata. |
| `GetSalesTrendSnapshot` | `/propertyapi/v1.0.0/salestrend/snapshot` | Sales trends by geography at configurable intervals. |
| `GetTransactionSalesTrend` | `/propertyapi/v1.0.0/transaction/salestrend` | Transaction-based sales trends with property type filtering. |

### Events and Permits

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetAllEventsDetail` | `/propertyapi/v1.0.0/allevents/detail` | All property events combined across sales, liens, AVM, and related feeds. |

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
