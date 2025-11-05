package property

import (
	"fmt"
	"strings"
)

// ValidateFIPSAndAPN ensures both FIPS and APN identifiers are supplied together.
func ValidateFIPSAndAPN(fips, apn string) error {
	if strings.TrimSpace(fips) == "" || strings.TrimSpace(apn) == "" {
		return fmt.Errorf("%w: both fips and APN are required", ErrMissingParameter)
	}
	return nil
}
