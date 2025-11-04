# GitHub Actions CI/CD Pipeline Summary

This document describes the GitHub Actions workflows implemented for the go-attom project.

## Workflows Implemented

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**

#### Lint
- Uses `golangci-lint` to check code quality
- Runs all configured linters from `.golangci.yml`
- Timeout: 5 minutes

#### Format Check
- Verifies all Go code is properly formatted with `gofmt`
- Fails if any files need formatting
- Provides guidance to run `go fmt ./...`

#### Test
- Runs all tests with race detection
- Generates coverage report
- Uploads coverage to Codecov (when token is configured)
- Displays total coverage percentage

#### Build
- Builds all Go packages
- Verifies `go.mod` and `go.sum` are up to date
- Ensures `go mod tidy` doesn't introduce changes

### 2. README Lint Workflow (`.github/workflows/readme-lint.yml`)

**Triggers:**
- Push to `main` or `develop` branches (when markdown files change)
- Pull requests to `main` or `develop` branches (when markdown files change)

**Jobs:**

#### Markdown Lint
- Uses `markdownlint-cli2` to lint all markdown files
- Configuration in `.markdownlint.json`
- Ensures consistent markdown formatting

### 3. Security Workflow (`.github/workflows/security.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Weekly schedule (Mondays at midnight UTC)

**Jobs:**

#### Gosec Security Scanner
- Scans Go code for security vulnerabilities
- Generates SARIF report
- Uploads results to GitHub Security tab

#### CodeQL Analysis
- Performs static analysis for security issues
- Automatically builds the code
- Uploads findings to GitHub Advanced Security

## Configuration Files

### `.golangci.yml`
Configures golangci-lint with the following linters:
- `errcheck` - Check for unchecked errors
- `gosimple` - Suggest code simplifications
- `govet` - Go vet examination
- `ineffassign` - Detect ineffectual assignments
- `staticcheck` - Advanced static analysis
- `unused` - Find unused code
- `gofmt` - Check formatting
- `goimports` - Check import formatting
- `misspell` - Find spelling mistakes
- `unconvert` - Remove unnecessary conversions
- `dupl` - Find duplicated code
- `goconst` - Find repeated strings
- `gocyclo` - Check cyclomatic complexity
- `gosec` - Security audit
- `stylecheck` - Style violations
- `revive` - Fast linter

### `.markdownlint.json`
Configures markdown linting rules. Currently configured to be lenient with existing files while enforcing basic markdown quality.

### `.gitignore`
Updated to exclude:
- Go build artifacts (`.exe`, `.dll`, `.so`, `.dylib`)
- Test artifacts (`.test`, `.out`)
- Coverage reports (`coverage.out`, `coverage.html`)
- Build directories (`/bin/`, `/dist/`)
- IDE files (`.idea/`, `.swp`, `.swo`)

## README Badges

The README now includes badges for:
- **CI Status** - Shows build status
- **Security Status** - Shows security scan status
- **README Lint Status** - Shows markdown lint status
- **Codecov** - Shows code coverage (requires token configuration)
- **Go Report Card** - Shows code quality grade
- **GoDoc** - Links to documentation
- **License** - Shows MIT license

## Initial Go Module

### `go.mod`
- Module name: `github.com/my-eq/go-attom`
- Go version: 1.24.9

### `pkg/client/client.go`
- Basic package structure
- Version constant: "0.1.0"

### `pkg/client/client_test.go`
- Basic test to ensure workflows pass
- Tests the Version constant

## How to Use

### Running Locally

```bash
# Run all tests
go test ./... -v -race

# Run with coverage
go test ./... -race -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linter
golangci-lint run ./...

# Check formatting
gofmt -l .

# Format code
go fmt ./...

# Lint markdown
markdownlint-cli2 "**/*.md"

# Verify dependencies
go mod tidy
```

### Setting Up Codecov (Optional)

1. Sign up at <https://codecov.io>
2. Add the repository
3. Add `CODECOV_TOKEN` to GitHub Secrets
4. Coverage reports will automatically upload on each test run

### Enabling GitHub Advanced Security (Optional)

For private repositories:
1. Enable GitHub Advanced Security in repository settings
2. CodeQL and Gosec SARIF uploads will populate the Security tab

## Next Steps

1. ✅ CI/CD pipeline is complete and ready
2. Add actual API client implementation
3. Add comprehensive tests as features are implemented
4. Monitor coverage and maintain high quality
5. Review security scan results regularly

## Maintenance

- The security workflow runs weekly to catch new vulnerabilities
- Update `golangci-lint` version periodically
- Review and update linter rules as needed
- Keep dependencies updated with `go get -u ./...`
