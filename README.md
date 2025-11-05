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
- [Advanced Usage](#advanced-usage)
- [Naming Convention Warning](#naming-convention-warning)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Overview

`go-attom` packages the ATTOM Property API into a cohesive toolkit that is easy to integrate, test, and extend. Services wrap ATTOM endpoints with first-class Go models, and helpers ensure queries are constructed with the exact parameters the API expects.

## Key Features

- **Deep Property API coverage** – 30+ endpoints spanning property profiles, ownership, mortgage liens, building permits, assessments, AVMs, sales history, and geo trends. The coverage mirrors the official ATTOM Property API catalog.【F:API_IMPLEMENTATION_SUMMARY.md†L9-L63】
- **Functional option builders** – Compose ATTOM query parameters with helpers such as `WithAddress`, `WithAttomID`, `WithLatitudeLongitude`, `WithDateRange`, and `WithPropertyType` from `pkg/property/options.go`.
- **Mockable foundation** – Inject any `http.Client`-compatible type through `client.New`, or supply your own interface implementation for advanced testing scenarios.【F:pkg/client/client.go†L14-L71】
- **Consistent error propagation** – Library methods wrap errors with context and return structured `*property.Error` values that capture status metadata and raw responses.【F:pkg/property/service.go†L55-L122】
- **No hidden logging** – The client surfaces detailed errors while leaving logging decisions entirely to consuming applications.

## Getting Started

### Prerequisites

- An ATTOM API key with access to the Property API portfolio.
- Go 1.21 or newer (matches the module’s go directive).

### Installation

```bash
go get github.com/my-eq/go-attom
```

### Authenticate and Instantiate

```go
attomClient := client.New(apiKey, nil) // defaults to *http.Client with a 30s timeout
propertyService := property.NewService(attomClient)
```

Passing `nil` uses an internal `*http.Client` with a 30 second timeout; you can inject your own implementation to customize timeouts, retries, or tracing.【F:pkg/client/client.go†L32-L71】

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
> Helpers like `safeString` are optional—responses expose pointers so you can detect missing data.【F:pkg/property/models.go†L6-L129】

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

All endpoints use ATTOM API version `v1.0.0` unless stated otherwise. Descriptions below are aligned with ATTOM’s official endpoint definitions.【F:API_IMPLEMENTATION_SUMMARY.md†L9-L63】

### Property Profiles and Core Records

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetPropertyID` | `/propertyapi/v1.0.0/property/id` | Resolve ATTOM, legacy, or county identifiers for a supplied postal address. |
| `GetPropertyDetail` | `/propertyapi/v1.0.0/property/detail` | Retrieve the full property detail dossier covering structure, lot, owner, tax, and location attributes. |
| `GetPropertyAddress` | `/propertyapi/v1.0.0/property/address` | Access normalized address, geocoding metadata, and linked identifiers. |
| `GetPropertySnapshot` | `/propertyapi/v1.0.0/property/snapshot` | Produce a condensed property profile for list views and portfolio summaries. |
| `GetBasicProfile` | `/propertyapi/v1.0.0/property/basicprofile` | Return marketing-friendly high-level property characteristics. |
| `GetExpandedProfile` | `/propertyapi/v1.0.0/property/expandedprofile` | Deliver ATTOM’s expanded profile with 400+ data points across building, lot, tax, and location dimensions. |
| `GetDetailWithSchools` | `/propertyapi/v1.0.0/property/detailwithschools` | Combine property detail with assigned public and private schools. |
| `GetDetailMortgage` | `/propertyapi/v1.0.0/property/detailmortgage` | Enrich property detail with active mortgage positions, lien holders, and loan characteristics. |
| `GetDetailOwner` | `/propertyapi/v1.0.0/property/detailowner` | Surface vesting, owner mailing addresses, and ownership history alongside property data. |
| `GetDetailMortgageOwner` | `/propertyapi/v1.0.0/property/detailmortgageowner` | Provide a unified record including mortgage liens and ownership insight for skip tracing. |
| `GetBuildingPermits` | `/propertyapi/v1.0.0/property/buildingpermits` | List building permits, contractors, valuations, and work types tied to the property. |

### Ownership, Mortgage, and Schools

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `SearchSchools` | `/propertyapi/v1.0.0/school/search` | Search nearby K–12 schools using address, lat/lon, or geo filters with distance scoring. |
| `GetSchoolProfile` | `/propertyapi/v1.0.0/school/profile` | Retrieve school performance, enrollment, and program details. |
| `GetSchoolDistrict` | `/propertyapi/v1.0.0/school/district` | Return district boundaries, contact information, and governance metadata. |
| `GetSchoolDetailWithSchools` | `/propertyapi/v1.0.0/school/detailwithschools` | Blend property characteristics with the set of assigned schools. |

### Sales, Assessments, and Valuations

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetSaleDetail` | `/propertyapi/v1.0.0/sale/detail` | Fetch recorded sale transactions including parties, price, deed type, and document metadata. |
| `GetSaleSnapshot` | `/propertyapi/v1.0.0/sale/snapshot` | Summarize most recent sale metrics and confidence scores for valuation workflows. |
| `GetAssessmentDetail` | `/propertyapi/v1.0.0/assessment/detail` | Provide assessor valuations, tax burdens, exemptions, and rate information. |
| `GetAssessmentSnapshot` | `/propertyapi/v1.0.0/assessment/snapshot` | Deliver condensed assessment metrics suitable for dashboards or list views. |
| `GetAssessmentHistory` | `/propertyapi/v1.0.0/assessmenthistory/detail` | Supply a historical timeline of annual assessments dating back to 1985 when available. |
| `GetAVMSnapshot` | `/propertyapi/v1.0.0/avm/snapshot` | Return ATTOM automated valuation estimates with ranges and confidence scoring. |
| `GetAttomAVMDetail` | `/propertyapi/v1.0.0/attomavm/detail` | Provide the full ATTOM AVM detail record with percentile, trend, and supporting comps. |
| `GetAVMHistory` | `/propertyapi/v1.0.0/avmhistory/detail` | Retrieve month-over-month valuation history for trending analyses. |
| `GetRentalAVM` | `/propertyapi/v1.0.0/valuation/rentalavm` | Calculate rent range, estimated rent, and rent-to-price metrics for investment workflows. |

### Sales History and Trends

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetSalesHistoryDetail` | `/propertyapi/v1.0.0/saleshistory/detail` | Produce the complete chain of historical sales with deed and financing attributes. |
| `GetSalesHistorySnapshot` | `/propertyapi/v1.0.0/saleshistory/snapshot` | Return a summarized view of transaction history for quick evaluation. |
| `GetSalesHistoryBasic` | `/propertyapi/v1.0.0/saleshistory/basichistory` | Deliver essential sale facts (price, date, deed) for high-volume lookups. |
| `GetSalesHistoryExpanded` | `/propertyapi/v1.0.0/saleshistory/expandedhistory` | Provide exhaustive sale history with document images, mortgage data, and recording details. |
| `GetSalesTrendSnapshot` | `/propertyapi/v1.0.0/salestrend/snapshot` | Offer geo-based trend indicators such as median price, volume, and DOM for a given GeoID. |
| `GetTransactionSalesTrend` | `/propertyapi/v1.0.0/transaction/salestrend` | Produce transaction-centric sales trend analytics for underwriting and forecasting. |

### Events and Permits

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetAllEventsDetail` | `/propertyapi/v1.0.0/allevents/detail` | Aggregate cross-domain events (sales, liens, foreclosure, MLS, AVM) into a single property timeline. |

## Composing Requests with Options

Property methods accept a flexible list of functional options defined in `pkg/property/options.go`.

- **Property identifiers** – `WithAttomID`, `WithPropertyID`, `WithFIPSAndAPN`, and `WithAddressLines` automatically populate ATTOM’s identifier parameters.
- **Geographic search** – `WithLatitudeLongitude`, `WithRadius`, `WithPostalCode`, `WithGeoID`, and `WithGeoIDV4` generate latitude/longitude or GeoID queries with ATTOM’s expected units and casing.
- **Filtering** – Helpers such as `WithBedsRange`, `WithBathsRange`, `WithSaleAmountRange`, `WithPropertyType`, `WithPropertyIndicator`, `WithUniversalSizeRange`, `WithYearBuiltRange`, `WithLotSize1Range`, and `WithLotSize2Range` enforce allowed parameter combinations.【F:pkg/property/options.go†L139-L356】
- **Date windows** – `WithDateRange` (MM/DD) and `WithISODateRange` (YYYY-MM-DD) translate human-friendly inputs into the `start*`/`end*` parameters expected by ATTOM.
- **Pagination and sorting** – `WithPage`, `WithPageSize`, and `WithOrderBy` manage paging controls and ordering fields.
- **Custom parameters** – `WithString`, `WithStringSlice`, and `WithAdditionalParam` let you add rarely used query strings without manual URL encoding.

These helpers guarantee your requests ship with the correct casing and parameter names—many ATTOM filters are case-sensitive and reject mismatched field names.【F:pkg/property/options.go†L19-L121】【F:pkg/property/service.go†L123-L179】

## Error Handling

- `property.ErrMissingParameter` is returned when required query parameters are omitted or incompatible.【F:pkg/property/options.go†L21-L33】
- API failures surface as `*property.Error`, which includes the HTTP status, parsed ATTOM status block, and raw response body for debugging.【F:pkg/property/error.go†L14-L77】
- All public methods wrap underlying errors using Go’s `%w` semantics so you can chain `errors.Is`/`errors.As` checks.【F:pkg/property/service.go†L63-L122】

## Advanced Usage

- **Inject custom HTTP clients** – Supply a transport that adds retries, tracing, or circuit breaking by passing your implementation to `client.New(apiKey, httpClient)`. The client only requires a `Do(*http.Request)` method.【F:pkg/client/client.go†L14-L71】
- **Override base URLs** – Target sandboxes or mock servers with `client.WithBaseURL` while keeping path composition intact.【F:pkg/client/client.go†L24-L47】
- **Context-driven cancellation** – Every service method accepts a `context.Context`, enabling per-call deadlines or integration with upstream request lifecycles.【F:pkg/property/service.go†L40-L122】
- **Thin mocks for testing** – The `client.HTTPClient` interface makes it straightforward to stub ATTOM responses in unit tests without hitting the network.【F:pkg/client/client.go†L14-L71】

## Naming Convention Warning

ATTOM mixes camelCase, lowercase, and legacy uppercase parameters (`fips` + `APN`) across the Property API. The models and option builders in `go-attom` preserve those cases deliberately:

- Struct fields use idiomatic Go names while JSON tags mirror ATTOM’s exact field casing (for example, `Identifier.AttomID` maps to `attomId`, and `Identifier.APN` maps to `apn`).【F:pkg/property/models.go†L14-L47】
- Query helpers emit ATTOM’s parameter names verbatim (`WithFIPSAndAPN` sets `fips` and `APN`, while `WithGeoIDV4` sets `geoIdV4`). Mixing cases or renaming parameters will trigger `property.ErrMissingParameter` before the request is issued.【F:pkg/property/options.go†L139-L305】【F:pkg/property/service.go†L123-L179】

When you introduce new fields or options, copy ATTOM’s casing exactly in JSON tags and query keys to avoid subtle API rejections.

## Development

```bash
# Run tests with race detection
go test ./... -race -v

# Run golangci-lint
golangci-lint run ./...

# Check markdown formatting
markdownlint-cli2 "**/*.md"
```

See [`API_IMPLEMENTATION_SUMMARY.md`](API_IMPLEMENTATION_SUMMARY.md) for a full breakdown of every ATTOM API group and [`docs/GITHUB_ACTIONS_SUMMARY.md`](docs/GITHUB_ACTIONS_SUMMARY.md) for CI/CD pipeline details.

## Contributing

Contributions are welcome! Please read the [project guidelines](.github/copilot-instructions.md) before opening an issue or pull request.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
