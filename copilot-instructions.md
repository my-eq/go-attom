# Copilot Instructions for go-attom

This document defines coding standards and development practices for the go-attom client library. **This file should evolve as the project grows—update it when you learn new patterns or best practices.**

## Project Overview

go-attom is a **client library** for interacting with the ATTOM Data API. As a library consumed by other applications, it must be clean, testable, and idiomatic Go.

## Core Principles

### 1. No Logging in Library Code

**NEVER** add logging statements to library code.

- ❌ **Don't**: Use `log.Printf()`, `fmt.Println()`, or any logging framework
- ✅ **Do**: Return descriptive errors; let consuming applications decide how to log
- **Rationale**: Libraries should not dictate logging behavior or pollute consumer output

```go
// ❌ BAD
func (c *Client) GetProperty(id string) (*Property, error) {
    log.Printf("Fetching property %s", id)
    // ...
}

// ✅ GOOD
func (c *Client) GetProperty(id string) (*Property, error) {
    // Return errors; consumers can log them
    if id == "" {
        return nil, fmt.Errorf("property ID cannot be empty")
    }
    // ...
}
```

### 2. Mockable and Injectable Design

All clients and interfaces must be mockable for testing.

- ✅ **Do**: Define interfaces for all major components
- ✅ **Do**: Accept interfaces, return structs
- ✅ **Do**: Use dependency injection (constructor parameters)
- ✅ **Do**: Provide mock implementations in `_test.go` files or `mocks/` package

```go
// ✅ GOOD: Interface for HTTP client
type HTTPClient interface {
    Do(req *http.Request) (*http.Response, error)
}

// ✅ GOOD: Client accepts interface
type Client struct {
    httpClient HTTPClient
    apiKey     string
}

func NewClient(apiKey string, httpClient HTTPClient) *Client {
    if httpClient == nil {
        httpClient = &http.Client{Timeout: 30 * time.Second}
    }
    return &Client{
        httpClient: httpClient,
        apiKey:     apiKey,
    }
}
```

### 3. Consistent Error Handling

Error messages must be clear, actionable, and follow Go conventions.

- ✅ **Do**: Use `fmt.Errorf()` with `%w` verb for error wrapping
- ✅ **Do**: Start error messages with lowercase (unless proper noun)
- ✅ **Do**: Provide context in errors (what failed, why)
- ✅ **Do**: Define sentinel errors as package-level variables for common cases
- ❌ **Don't**: Use generic error messages like "error occurred"

```go
// ✅ GOOD: Sentinel errors
var (
    ErrInvalidAPIKey    = errors.New("invalid or missing API key")
    ErrPropertyNotFound = errors.New("property not found")
    ErrRateLimitExceeded = errors.New("API rate limit exceeded")
)

// ✅ GOOD: Wrapped errors with context
func (c *Client) GetProperty(id string) (*Property, error) {
    if id == "" {
        return nil, fmt.Errorf("property ID cannot be empty")
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch property %s: %w", id, err)
    }
    // ...
}
```

### 4. Idiomatic Go Interfaces

Follow Go best practices for interface design.

- ✅ **Do**: Keep interfaces small (1-3 methods when possible)
- ✅ **Do**: Name single-method interfaces with `-er` suffix (`Getter`, `Fetcher`)
- ✅ **Do**: Define interfaces at the point of use (consumer side)
- ✅ **Do**: Use standard library interfaces where applicable (`io.Reader`, `io.Writer`)
- ❌ **Don't**: Create large "god" interfaces

```go
// ✅ GOOD: Small, focused interfaces
type PropertyGetter interface {
    GetProperty(id string) (*Property, error)
}

type AreaSearcher interface {
    SearchAreas(params *SearchParams) (*AreaResults, error)
}

// ❌ BAD: Large interface
type ATTOMClient interface {
    GetProperty(id string) (*Property, error)
    SearchAreas(params *SearchParams) (*AreaResults, error)
    GetPOI(id string) (*POI, error)
    GetCommunity(id string) (*Community, error)
    // ... 20 more methods
}
```

## Testing Requirements

### 100% Test Coverage

Every public function, method, and type must have comprehensive tests.

- ✅ **Do**: Write table-driven tests for multiple scenarios
- ✅ **Do**: Test happy paths AND error cases
- ✅ **Do**: Use `t.Run()` for subtests
- ✅ **Do**: Mock external dependencies (HTTP calls, file I/O)
- ✅ **Do**: Include edge cases (empty strings, nil pointers, boundary values)

```go
func TestClient_GetProperty(t *testing.T) {
    tests := []struct {
        name        string
        propertyID  string
        mockResp    *http.Response
        mockErr     error
        want        *Property
        wantErr     bool
        errContains string
    }{
        {
            name:       "successful request",
            propertyID: "12345",
            mockResp:   mockSuccessResponse(),
            want:       &Property{ID: "12345"},
            wantErr:    false,
        },
        {
            name:        "empty property ID",
            propertyID:  "",
            wantErr:     true,
            errContains: "property ID cannot be empty",
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := NewClient("test-key", &mockHTTPClient{
                resp: tt.mockResp,
                err:  tt.mockErr,
            })
            
            got, err := client.GetProperty(tt.propertyID)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error containing %q, got nil", tt.errContains)
                } else if !strings.Contains(err.Error(), tt.errContains) {
                    t.Errorf("error = %v, want error containing %q", err, tt.errContains)
                }
                return
            }
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %+v, want %+v", got, tt.want)
            }
        })
    }
}
```

### Test Coverage Verification

```bash
# Run tests with coverage
go test ./... -race -coverprofile=coverage.out

# View coverage report
go tool cover -html=coverage.out

# Ensure 100% coverage (or close to it)
go tool cover -func=coverage.out
```

## Pre-Commit Workflow

**ALWAYS** run the following checks before committing:

### 0. Version Compatibility Check

**NEVER downgrade dependency versions without explicit justification.**

Before modifying `go.mod`, CI configuration, or any dependency versions:

- ✅ **Do**: Use web search to verify the latest stable release of Go or dependencies
- ✅ **Do**: Check the current date and year to ensure version compatibility
- ✅ **Do**: Upgrade to newer stable versions when addressing compatibility issues
- ✅ **Do**: Verify the version exists and is released before using it
- ❌ **Don't**: Assume a version doesn't exist without checking
- ❌ **Don't**: Downgrade versions to match CI without first checking if CI should be updated
- ❌ **Don't**: Make version changes based on assumptions about release schedules

```bash
# Before changing Go version, verify what's actually released
# Use web search: "latest stable Go version [current date]"

# If CI and go.mod differ, upgrade CI to match go.mod (if go.mod version exists)
# NOT the other way around
```

### 1. Run Linters

```bash
# Install golangci-lint if not present
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linters
golangci-lint run ./...

# Or use go vet as minimum
go vet ./...
```

### 2. Format Code

```bash
# Format all Go files
go fmt ./...

# Or use gofumpt for stricter formatting
# gofumpt -l -w .
```

### 3. Run Tests

```bash
# Run all tests with race detection
go test ./... -race -v

# Verify coverage
go test ./... -race -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

### 4. Build Verification

```bash
# Ensure code builds successfully
go build ./...

# Run go mod tidy to clean dependencies
go mod tidy

# Verify go.mod and go.sum are clean
git diff go.mod go.sum
```

## Commit Standards

### Conventional Commits

All commits **MUST** follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

**Format**: `<type>(<scope>): <description>`

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring without behavior change
- `perf`: Performance improvements
- `chore`: Build, dependencies, tooling
- `style`: Code style changes (formatting, no logic change)

**Examples**:
```
feat(property): add GetPropertyDetails endpoint
fix(client): handle nil pointer when API key is empty
test(area): add edge cases for area search
docs(readme): update installation instructions
refactor(models): simplify Property struct tags
```

### One Feature Per Commit Per PR

- ✅ **Do**: Keep commits focused on a single feature or fix
- ✅ **Do**: Squash related commits before merging
- ✅ **Do**: Write descriptive commit messages with context in the body
- ❌ **Don't**: Mix unrelated changes in one commit
- ❌ **Don't**: Make "WIP" or "fix tests" commits in final PR

```bash
# Squash commits before merge
git rebase -i main
# Mark commits as 'squash' or 'fixup'

# Or use GitHub's "Squash and merge" button
```

## Code Quality Standards

### Never Create Incomplete Code

- ✅ **Do**: Ensure all functions are fully implemented before committing
- ✅ **Do**: Include error handling in all code paths
- ✅ **Do**: Add tests for new code in the same commit
- ❌ **Don't**: Use `// TODO: implement this` or `panic("not implemented")`
- ❌ **Don't**: Leave commented-out code
- ❌ **Don't**: Commit code that doesn't compile or pass tests

### Documentation and Comments

All exported types, functions, and methods **MUST** have GoDoc comments.

- ✅ **Do**: Start comments with the name of the thing being documented
- ✅ **Do**: Write complete sentences with proper punctuation
- ✅ **Do**: Include usage examples in package documentation
- ✅ **Do**: Document parameters, return values, and error conditions
- ❌ **Don't**: State the obvious (`// GetProperty gets a property`)
- ❌ **Don't**: Leave exported items undocumented

```go
// ✅ GOOD: Detailed, actionable documentation
// Client provides methods for interacting with the ATTOM Data API.
// It handles authentication, request formatting, and response parsing.
//
// Example usage:
//
//	client := attom.NewClient("your-api-key", nil)
//	property, err := client.GetProperty("12345")
//	if err != nil {
//	    log.Fatal(err)
//	}
type Client struct {
    httpClient HTTPClient
    apiKey     string
    baseURL    string
}

// GetProperty retrieves detailed property information by ID.
//
// The property ID must be a valid ATTOM property identifier.
// Returns ErrPropertyNotFound if the property does not exist.
// Returns ErrInvalidAPIKey if authentication fails.
//
// Example:
//
//	property, err := client.GetProperty("12345")
//	if errors.Is(err, attom.ErrPropertyNotFound) {
//	    // handle not found
//	}
func (c *Client) GetProperty(id string) (*Property, error) {
    // ...
}

### Naming Conventions (Go Idioms)

Consistent naming improves readability and maintainability. Follow established Go conventions:

- ✅ Exported identifiers use PascalCase: `Client`, `New`, `ErrInvalidAPIKey`
- ✅ Unexported identifiers use camelCase: `httpClient`, `baseURL`, `doRequest`
- ✅ Acronyms are capitalized consistently: `API`, `HTTP`, `URL` (e.g. `apiKey`, `baseURL`, `HTTPClient`)
- ✅ Avoid stutter: If the package is `client`, exported types should not repeat: prefer `Client` over `ClientClient`
- ✅ Single-purpose option types may use `Option`; if multiple option types emerge, prefix with the domain: `ClientOption`
- ✅ Short receiver names: Use one or two letters (`c`, `r`, `s`) unless ambiguity arises
- ✅ Prefer descriptive names over ambiguous: `baseURL` not `u`, `normalized` not `trimmed`
- ✅ Keep test variable names clear: `got`, `want`, `tt` for table rows
- ❌ Don't use Hungarian notation or type hints in names (e.g. `strName`)
- ❌ Don't use snake_case
- ❌ Don't over-abbreviate (`propID` -> prefer `propertyID` unless widely understood)

```go
// ✅ GOOD
const DefaultBaseURL = "https://api.attomdata.com/v1/"

func WithBaseURL(baseURL string) Option {
    return func(c *Client) {
        if baseURL == "" { return }
        c.baseURL = strings.TrimRight(baseURL, "/") + "/"
    }
}

// ❌ BAD (ambiguous variable names)
func WithBaseURL(u string) Option { /* ... */ }
```

Refactor on sight any ambiguous variable names introduced by agents; update this section if new patterns emerge (e.g., generics, context handling, retries).
```

## Package Structure

Follow Go package best practices:

- ✅ **Do**: Keep packages focused and cohesive
- ✅ **Do**: Use `pkg/models` for shared data structures
- ✅ **Do**: Use `internal/` for non-exported implementation details
- ✅ **Do**: Keep files under 500 lines when possible
- ❌ **Don't**: Create circular dependencies

## Evolving These Instructions

**This file is a living document.** When you encounter new patterns, best practices, or project-specific requirements:

1. **Update this file** with new guidelines
2. **Document WHY** (rationale, not just the rule)
3. **Provide examples** (good and bad)
4. **Commit changes** with: `docs(copilot): update coding standards`

### When to Update

- You find a better way to structure code
- You identify a recurring bug pattern to avoid
- You add new tools or linters to the workflow
- You establish new naming conventions
- You create reusable test utilities or patterns

## Quick Reference Checklist

Before committing, ensure:

- [ ] No version downgrades without web search verification
- [ ] No logging statements in library code
- [ ] All dependencies are injected and mockable
- [ ] Errors are wrapped with context using `%w`
- [ ] All exported items have GoDoc comments
- [ ] Tests cover happy path and error cases
- [ ] Test coverage is 100% (or documented why not)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Linters pass (`golangci-lint run ./...`)
- [ ] Tests pass with race detection (`go test ./... -race`)
- [ ] Code builds successfully (`go build ./...`)
- [ ] Commit message follows Conventional Commits format
- [ ] Commit contains one focused change
- [ ] No incomplete code or TODOs

---

**Remember**: Quality over speed. Take the time to do it right the first time.
