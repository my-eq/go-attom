// Package property provides access to the ATTOM Property API.
//
// This package implements a comprehensive client for interacting with ATTOM's
// Property API endpoints, including property details, assessments, sales data,
// AVM (Automated Valuation Model), school information, and more.
package property

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrMissingParameter indicates that a required parameter was not supplied for a request.
var ErrMissingParameter = errors.New("property: missing required parameter")

// Error represents an ATTOM Property API error response.
type Error struct {
	Status     *Status
	Message    string
	Body       json.RawMessage
	StatusCode int
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e == nil {
		return "property: nil error"
	}
	if e.Message != "" {
		return fmt.Sprintf("property: %s", e.Message)
	}
	if e.Status != nil {
		if e.Status.Msg != nil {
			return fmt.Sprintf("property: %s", *e.Status.Msg)
		}
		if e.Status.Code != nil {
			return fmt.Sprintf("property: status code %d", *e.Status.Code)
		}
	}
	return fmt.Sprintf("property: http status %d", e.StatusCode)
}
