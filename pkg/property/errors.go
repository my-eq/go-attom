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
	StatusCode int
	Status     *Status
	Message    string
	Body       json.RawMessage
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
