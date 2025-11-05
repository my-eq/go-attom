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

- **Comprehensive Property API coverage** – 30+ endpoints for property detail, ownership, mortgages, assessments, AVMs, sales history, trends, school lookups, and more.
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

All endpoints use ATTOM API version `v1.0.0` unless noted otherwise. Descriptions come from the official ATTOM swagger definitions included in this repository.

### Property Profiles & Basics

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetPropertyID` | `/propertyapi/v1.0.0/property/id` | Returns properties that match search criteria and "get a list of property IDs within a specific geography that have a specific number of beds".【F:docs/attom/swagger/propertyapi_property.pretty.json†L117-L144】 |
| `GetPropertyDetail` | `/propertyapi/v1.0.0/property/detail` | Returns property details for a supplied ATTOM ID.【F:docs/attom/swagger/propertyapi_property.pretty.json†L115-L148】 |
| `GetPropertyAddress` | `/propertyapi/v1.0.0/property/address` | Returns properties within a ZIP code and supports narrowing with property type and ordering options.【F:docs/attom/swagger/propertyapi_property.pretty.json†L148-L188】 |
| `GetPropertySnapshot` | `/propertyapi/v1.0.0/property/snapshot` | Returns property snapshots that match filters such as city, size range, and property type.【F:docs/attom/swagger/propertyapi_property.pretty.json†L188-L227】 |
| `GetBasicProfile` | `/propertyapi/v1.0.0/property/basicprofile` | Returns basic property information plus the most recent transaction and tax data for an address.【F:docs/attom/swagger/propertyapi_property.pretty.json†L227-L269】 |
| `GetExpandedProfile` | `/propertyapi/v1.0.0/property/expandedprofile` | Returns detailed property information with the latest transaction and taxes for an address.【F:docs/attom/swagger/propertyapi_property.pretty.json†L269-L309】 |
| `GetBuildingPermits` | `/propertyapi/v1.0.0/property/buildingpermits` | Returns basic property information and detailed building permits for an address.【F:docs/attom/swagger/propertyapi_property.pretty.json†L309-L352】 |
| `GetAllEventsDetail` | `/propertyapi/v1.0.0/allevents/detail` | Returns the full timeline of events that occurred on a property, including cross-domain activity.【F:docs/attom/swagger/propertyapi_allevents.pretty.json†L5-L47】 |

### Ownership, Mortgage, and Schools

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetDetailWithSchools` | `/propertyapi/v1.0.0/property/detailwithschools` | Returns property details together with the schools inside the attendance zones for the address.【F:docs/attom/swagger/propertyapi_school.pretty.json†L28-L70】 |
| `GetDetailMortgage` | `/propertyapi/v1.0.0/property/detailmortgage` | Returns property detail enriched with mortgage information for the provided address.【F:pkg/property/service.go†L268-L288】 |
| `GetDetailOwner` | `/propertyapi/v1.0.0/property/detailowner` | Returns property detail enriched with ownership information for the provided address.【F:pkg/property/service.go†L290-L307】 |
| `GetDetailMortgageOwner` | `/propertyapi/v1.0.0/property/detailmortgageowner` | Returns property detail enriched with combined mortgage and ownership information for the address.【F:pkg/property/service.go†L309-L327】 |
| `SearchSchools` | `/propertyapi/v1.0.0/school/search` | Returns school listings around an address or coordinate search context.【F:docs/attom/swagger/propertyapi_school.pretty.json†L70-L123】 |
| `GetSchoolProfile` | `/propertyapi/v1.0.0/school/profile` | Returns enriched profile information for an individual school.【F:docs/attom/swagger/propertyapi_school.pretty.json†L123-L166】 |
| `GetSchoolDistrict` | `/propertyapi/v1.0.0/school/district` | Returns school district boundaries and related contact data.【F:docs/attom/swagger/propertyapi_school.pretty.json†L166-L209】 |

### Sales, Assessments, and Valuations

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetSaleDetail` | `/propertyapi/v1.0.0/sale/detail` | Returns detailed sale transaction information for a property identifier.【F:docs/attom/swagger/propertyapi_sale.pretty.json†L5-L51】 |
| `GetSaleSnapshot` | `/propertyapi/v1.0.0/sale/snapshot` | Returns a sale snapshot summarizing recent transaction metrics for a property.【F:docs/attom/swagger/propertyapi_sale.pretty.json†L51-L95】 |
| `GetAssessmentDetail` | `/propertyapi/v1.0.0/assessment/detail` | Returns detailed assessment, tax, and market value data.【F:docs/attom/swagger/propertyapi_assessment.pretty.json†L5-L52】 |
| `GetAssessmentSnapshot` | `/propertyapi/v1.0.0/assessment/snapshot` | Returns assessment snapshot metrics for a property identifier.【F:docs/attom/swagger/propertyapi_assessment.pretty.json†L52-L95】 |
| `GetAssessmentHistory` | `/propertyapi/v1.0.0/assessmenthistory/detail` | Returns historical assessment records for the property.【F:docs/attom/swagger/propertyapi_assessmenthistory.pretty.json†L5-L48】 |
| `GetAVMSnapshot` | `/propertyapi/v1.0.0/avm/snapshot` | Returns automated valuation model (AVM) snapshot values and confidence scoring.【F:docs/attom/swagger/propertyapi_avm.pretty.json†L5-L49】 |
| `GetAttomAVMDetail` | `/propertyapi/v1.0.0/attomavm/detail` | Returns ATTOM AVM detail including percentile and scoring metrics.【F:docs/attom/swagger/propertyapi_attomavm.pretty.json†L5-L47】 |
| `GetAVMHistory` | `/propertyapi/v1.0.0/avmhistory/detail` | Returns month-by-month AVM history for the property.【F:docs/attom/swagger/propertyapi_avmhistory.pretty.json†L5-L49】 |
| `GetRentalAVM` | `/propertyapi/v1.0.0/valuation/rentalavm` | Returns rental AVM valuations and rent ranges.【F:docs/attom/swagger/propertyapi_valuation.pretty.json†L5-L46】 |

### Sales History & Trends

| Go Method | Endpoint | ATTOM Description |
|-----------|----------|-------------------|
| `GetSalesHistoryDetail` | `/propertyapi/v1.0.0/saleshistory/detail` | Returns the full sales history for a property.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L5-L50】 |
| `GetSalesHistorySnapshot` | `/propertyapi/v1.0.0/saleshistory/snapshot` | Returns a snapshot of historical transactions for quick lookups.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L50-L92】 |
| `GetSalesHistoryBasic` | `/propertyapi/v1.0.0/saleshistory/basichistory` | Returns a lightweight transaction history for rapid searches.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L92-L136】 |
| `GetSalesHistoryExpanded` | `/propertyapi/v1.0.0/saleshistory/expandedhistory` | Returns expanded transaction history including document metadata.【F:docs/attom/swagger/propertyapi_saleshistory.pretty.json†L136-L180】 |
| `GetSalesTrendSnapshot` | `/propertyapi/v1.0.0/salestrend/snapshot` | Returns sales trend metrics for a specified geographic ID.【F:docs/attom/swagger/propertyapi_salestrend.pretty.json†L5-L47】 |
| `GetTransactionSalesTrend` | `/propertyapi/v1.0.0/transaction/salestrend` | Returns transaction-oriented sales trend metrics across geographies.【F:docs/attom/swagger/propertyapi_transaction.pretty.json†L5-L47】 |

## Building Requests with Options

Property requests accept a flexible list of functional options:

- **Property identifiers** – `WithAttomID`, `WithPropertyID`, `WithFIPSAndAPN`, and `WithAddressLines` for ATTOM identifiers and assessor parcel numbers.【F:pkg/property/options.go†L26-L65】
- **Geographic search** – `WithLatitudeLongitude`, `WithRadius`, `WithPostalCode`, `WithGeoID`, and `WithGeoIDV4` to target coordinates and geographic codes.【F:pkg/property/options.go†L67-L117】
- **Filtering** – `WithBedsRange`, `WithBathsRange`, `WithSaleAmountRange`, `WithPropertyType`, `WithPropertyIndicator`, `WithUniversalSizeRange`, `WithYearBuiltRange`, `WithLotSize1Range`, and `WithLotSize2Range` for ATTOM's numeric filters.【F:pkg/property/options.go†L119-L223】
- **Date windows** – `WithDateRange` (MM/DD) and `WithISODateRange` (YYYY-MM-DD) cover endpoints that expect legacy or ISO date formats.【F:pkg/property/options.go†L225-L280】
- **Pagination and sorting** – `WithPage`, `WithPageSize`, and `WithOrderBy` mirror ATTOM's paging and ordering controls.【F:pkg/property/options.go†L282-L327】
- **Custom parameters** – `WithString`, `WithStringSlice`, and `WithAdditionalParam` map to ATTOM's long tail of specialized filters.【F:pkg/property/options.go†L17-L47】【F:pkg/property/options.go†L329-L371】

These helpers ensure requests include the correct parameter names and formatting required by ATTOM.

## Error Handling

- `property.ErrMissingParameter` is returned when required query parameters are omitted by the caller.【F:pkg/property/errors.go†L5-L15】
- API errors use `*property.Error`, which captures the HTTP status code, parsed ATTOM status block, and raw response body for debugging.【F:pkg/property/errors.go†L17-L67】
- All public methods wrap lower-level errors with context using Go's `%w` semantics so you can use `errors.Is`/`errors.As`. The shared `doGet` helper centralizes the behavior.【F:pkg/property/service.go†L52-L121】

## Advanced Usage

### Customize the HTTP client

Inject your own `http.Client` (or any type implementing `client.HTTPClient`) when building the ATTOM client. This is ideal for adding retries, circuit breakers, or observability instrumentation.

```go
httpClient := &http.Client{Timeout: 10 * time.Second}
attomClient := client.New(apiKey, httpClient)
```

The constructor falls back to a 30-second timeout client when you pass `nil`, keeping defaults safe for production.【F:pkg/client/client.go†L36-L58】

### Override the base URL for staging and proxies

ATTOM provides dedicated gateways per environment. Use `client.WithBaseURL` to point at a staging cluster or run through an internal proxy:

```go
attomClient := client.New(apiKey, nil, client.WithBaseURL("https://staging.attomdata.com/"))
```

The option normalizes trailing slashes to keep request construction predictable.【F:pkg/client/client.go†L24-L44】

### Inspect detailed API failures

When ATTOM returns a non-2xx response, go-attom unmarshals the status payload into `property.Error`, preserving the HTTP code, ATTOM status block, and raw JSON to help with support tickets or sandbox debugging.【F:pkg/property/service.go†L74-L118】【F:pkg/property/errors.go†L17-L67】

## ⚠️ ATTOM Naming Nuances

ATTOM mixes lower-case, camelCase, and uppercase tokens in both query parameters and JSON payloads. The client mirrors those quirks so requests land correctly:

- Query helpers intentionally set parameters like `attomid`, `APN`, `geoIdV4`, and `propertyIndicator` using ATTOM's exact casing.【F:pkg/property/options.go†L49-L113】【F:pkg/property/options.go†L147-L177】
- Response structs keep ATTOM's field names in their `json` tags (`attomId`, `postalCode`, `areaSqFt`) while surfacing idiomatic Go field names (`AttomID`, `PostalCode`, `AreaSquareFeet`).【F:pkg/property/models.go†L7-L53】

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
