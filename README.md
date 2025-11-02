# go-attom

[![CI](https://github.com/my-eq/go-attom/actions/workflows/ci.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/ci.yml)
[![Security](https://github.com/my-eq/go-attom/actions/workflows/security.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/security.yml)
[![README Lint](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml/badge.svg)](https://github.com/my-eq/go-attom/actions/workflows/readme-lint.yml)
[![codecov](https://codecov.io/gh/my-eq/go-attom/branch/main/graph/badge.svg)](https://codecov.io/gh/my-eq/go-attom)
[![Go Report Card](https://goreportcard.com/badge/github.com/my-eq/go-attom)](https://goreportcard.com/report/github.com/my-eq/go-attom)
[![GoDoc](https://godoc.org/github.com/my-eq/go-attom?status.svg)](https://godoc.org/github.com/my-eq/go-attom)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A Go client library for the [ATTOM Data API](https://api.gateway.attomdata.com/). This library provides idiomatic Go access to ATTOM's comprehensive property data, area information, points of interest, community demographics, and parcel tile services.

## Features

- **Comprehensive API Coverage**: Support for all ATTOM API groups (Property, Area, POI, Community, Parcel Tiles)
- **Idiomatic Go**: Clean, testable, and mockable design following Go best practices
- **Type-Safe**: Strongly-typed models with proper handling of optional fields
- **Well-Tested**: Comprehensive test coverage with table-driven tests
- **No Dependencies**: Uses only the Go standard library
- **Context Support**: All methods support context.Context for cancellation and timeouts

## Installation

```bash
go get github.com/my-eq/go-attom
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/my-eq/go-attom/pkg/client"
)

func main() {
    // Create a new client
    c := client.New("your-api-key")

    // Use the client (implementation coming soon)
    fmt.Println("ATTOM Client version:", client.Version)
}
```

## API Coverage

This library provides access to the following ATTOM API groups:

- **PropertyAPI** (36+ endpoints): Property details, sales, assessments, valuations, schools
- **AreaAPI** (6 endpoints): Geographic boundaries, county/state lookups, hierarchies
- **POIAPI** (5 endpoints): Points of interest, business locations, amenities
- **CommunityAPI** (2 endpoints): Demographics, economics, education, housing, climate, transportation
- **ParcelTilesAPI** (1 endpoint): Parcel boundary raster tiles

## Development Status

This library is currently in active development. See the [API Implementation Summary](API_IMPLEMENTATION_SUMMARY.md) for detailed progress.

## Contributing

Contributions are welcome! Please read the [contributing guidelines](copilot-instructions.md) before submitting PRs.

## Testing

```bash
# Run tests
go test ./... -v

# Run tests with coverage
go test ./... -race -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linters
golangci-lint run ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [ATTOM Data Solutions](https://www.attomdata.com/) for providing the API
