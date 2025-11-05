# go-attom

[![CI](https://github.com/my-eq/go-attom/actions/workflows/ci.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/ci.yml)
[![Security](https://github.com/my-eq/go-attom/actions/workflows/security.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/security.yml)
[![README Lint](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml)
[![codecov](https://codecov.io/gh/my-eq/go-attom/branch/main/graph/badge.svg)](https://codecov.io/gh/my-eq/go-attom)
[![Go Report Card](https://goreportcard.com/badge/github.com/my-eq/go-attom)](https://goreportcard.com/report/github.com/my-eq/go-attom)
[![Go Reference](https://pkg.go.dev/badge/github.com/my-eq/go-attom.svg)](https://pkg.go.dev/github.com/my-eq/go-attom)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A strongly-typed, mock-friendly Go client for the [ATTOM Data API](https://api.gateway.attomdata.com/). The library focuses on clean, idiomatic Go and provides deep coverage of the Property API including property profiles, sales, assessments, valuations, schools, and historical trends (see the [API implementation summary](API_IMPLEMENTATION_SUMMARY.md#L9-L118) for the full endpoint catalog).

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

## Supported Property API Endpoints

All endpoints use ATTOM API version `v1.0.0` unless noted otherwise.

### Property Profiles & Basics

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetPropertyID` | `/propertyapi/v1.0.0/property/id` | Lookup ATTOM identifiers for an address. |
| `GetPropertyDetail` | `/propertyapi/v1.0.0/property/detail` | Detailed property characteristics, building data, and more. |
| `GetPropertyAddress` | `/propertyapi/v1.0.0/property/address` | Address, geocoding, and identifier metadata. |
| `GetPropertySnapshot` | `/propertyapi/v1.0.0/property/snapshot` | Lightweight property snapshot for summaries and lists. |
| `GetBasicProfile` | `/propertyapi/v1.0.0/property/basicprofile` | Basic profile fields for marketing or quick previews. |
| `GetExpandedProfile` | `/propertyapi/v1.0.0/property/expandedprofile` | Complete profile including 400+ data points. |
| `GetBuildingPermits` | `/propertyapi/v1.0.0/property/buildingpermits` | Historical building permits tied to the property. |
| `GetAllEventsDetail` | `/propertyapi/v1.0.0/allevents/detail` | Cross-domain events (sales, liens, MLS, AVM) for a property. |

### Ownership, Mortgage, and Schools

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetDetailWithSchools` | `/propertyapi/v1.0.0/property/detailwithschools` | Property detail with assigned schools. |
| `GetDetailMortgage` | `/propertyapi/v1.0.0/property/detailmortgage` | Detail enriched with mortgage positions. |
| `GetDetailOwner` | `/propertyapi/v1.0.0/property/detailowner` | Ownership information including vesting and mailing address. |
| `GetDetailMortgageOwner` | `/propertyapi/v1.0.0/property/detailmortgageowner` | Combined ownership and mortgage data. |
| `SearchSchools` | `/propertyapi/v1.0.0/school/search` | Locate schools near an address or coordinate. |
| `GetSchoolProfile` | `/propertyapi/v1.0.0/school/profile` | Enriched school profile data. |
| `GetSchoolDistrict` | `/propertyapi/v1.0.0/school/district` | School district boundaries and contacts. |
| `GetSchoolDetailWithSchools` | `/propertyapi/v1.0.0/school/detailwithschools` | Property and assigned schools in one response. |

### Sales, Assessments, and Valuations

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetSaleDetail` | `/propertyapi/v1.0.0/sale/detail` | Latest sale transaction with buyer/seller data. |
| `GetSaleSnapshot` | `/propertyapi/v1.0.0/sale/snapshot` | High-level sale metrics for quick valuation. |
| `GetAssessmentDetail` | `/propertyapi/v1.0.0/assessment/detail` | Current assessment, tax, and market value data. |
| `GetAssessmentSnapshot` | `/propertyapi/v1.0.0/assessment/snapshot` | Assessment summary including tax amounts and rates. |
| `GetAssessmentHistory` | `/propertyapi/v1.0.0/assessmenthistory/detail` | Historical assessments back to 1985 where available. |
| `GetAVMSnapshot` | `/propertyapi/v1.0.0/avm/snapshot` | Automated valuation snapshot with confidence scores. |
| `GetAttomAVMDetail` | `/propertyapi/v1.0.0/attomavm/detail` | ATTOM AVM detail with percentile and scoring. |
| `GetAVMHistory` | `/propertyapi/v1.0.0/avmhistory/detail` | Monthly AVM history. |
| `GetRentalAVM` | `/propertyapi/v1.0.0/valuation/rentalavm` | Rental valuation and rent range. |

### Sales History & Trends

| Go Method | Endpoint | Description |
|-----------|----------|-------------|
| `GetSalesHistoryDetail` | `/propertyapi/v1.0.0/saleshistory/detail` | Full sales history for the property. |
| `GetSalesHistorySnapshot` | `/propertyapi/v1.0.0/saleshistory/snapshot` | Summary view of historical transactions. |
| `GetSalesHistoryBasic` | `/propertyapi/v1.0.0/saleshistory/basichistory` | Lightweight transaction history for quick lookups. |
| `GetSalesHistoryExpanded` | `/propertyapi/v1.0.0/saleshistory/expandedhistory` | Rich transaction history with document metadata. |
| `GetSalesTrendSnapshot` | `/propertyapi/v1.0.0/salestrend/snapshot` | Geographic sales trends for a specified GeoID. |
| `GetTransactionSalesTrend` | `/propertyapi/v1.0.0/transaction/salestrend` | Transaction-oriented sales trend metrics. |

## Building Requests with Options

Property requests accept a flexible list of functional options:

- **Property identifiers** – `WithAttomID`, `WithPropertyID`, `WithFIPSAndAPN`, and `WithAddressLines`.
- **Geographic search** – `WithLatitudeLongitude`, `WithRadius`, `WithPostalCode`, `WithGeoID`, and `WithGeoIDV4`.
- **Filtering** – `WithBedsRange`, `WithBathsRange`, `WithSaleAmountRange`, `WithPropertyType`, `WithPropertyIndicator`, `WithUniversalSizeRange`, `WithYearBuiltRange`, `WithLotSize1Range`, and `WithLotSize2Range`.
- **Date windows** – `WithDateRange` (MM/DD format) and `WithISODateRange` (YYYY-MM-DD).
- **Pagination and sorting** – `WithPage`, `WithPageSize`, and `WithOrderBy`.
- **Custom parameters** – `WithString`, `WithStringSlice`, and `WithAdditionalParam` for rarely used filters.

These helpers ensure requests include the correct parameter names and formatting required by ATTOM.

## Error Handling

- `property.ErrMissingParameter` is returned when required query parameters are omitted.
- API errors use `*property.Error`, which exposes the HTTP status code, parsed ATTOM status block, and raw response body for debugging.
- All public methods wrap lower-level errors with context using Go's `%w` semantics so you can use `errors.Is`/`errors.As`.

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
