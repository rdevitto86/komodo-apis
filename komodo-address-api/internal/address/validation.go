package address

import (
	"regexp"
	"strings"
)

// Constants for error messages
const (
	ErrStreet1Required     = "street1 is required"
	ErrCityRequired        = "city is required"
	ErrStateRequired       = "state is required"
	ErrPostalCodeRequired  = "postalCode is required"
	ErrInvalidUSPostalCode = "invalid US ZIP (12345 or 12345-6789)"
	ErrInvalidCAPostalCode = "invalid CA postal code (A1A 1A1)"
)

// Precompiled regex for postal code validation
var (
	zipUS = regexp.MustCompile(`^\d{5}(-\d{4})?$`)
	zipCA = regexp.MustCompile(`^[A-Za-z]\d[A-Za-z][ -]?\d[A-Za-z]\d$`)
)

// ValidateAddress validates the fields of an Address and returns a map of errors.
func ValidateAddress(a Address) map[string]string {
	errs := map[string]string{}

	// Clean and normalize fields
	street1 := strings.TrimSpace(a.Street1)
	city := strings.TrimSpace(a.City)
	state := strings.TrimSpace(a.State)
	postalCode := strings.TrimSpace(a.PostalCode)
	country := strings.ToUpper(strings.TrimSpace(a.Country))

	if country == "" {
		country = DefaultCountry // Use the constant from the other file
	}

	// Validate required fields
	if street1 == "" {
		errs["street1"] = ErrStreet1Required
	}
	if city == "" {
		errs["city"] = ErrCityRequired
	}
	if state == "" {
		errs["state"] = ErrStateRequired
	}
	if postalCode == "" {
		errs["postalCode"] = ErrPostalCodeRequired
	}

	// Validate postal code based on country
	if postalCode != "" {
		validatePostalCode(postalCode, country, errs)
	}

	return errs
}

// validatePostalCode validates the postal code based on the country and adds errors to the map.
func validatePostalCode(postalCode, country string, errs map[string]string) {
	switch country {
	case "US", "USA", "UNITED STATES":
		if !zipUS.MatchString(postalCode) {
			errs["postalCode"] = ErrInvalidUSPostalCode
		}
	case "CA", "CANADA":
		if !zipCA.MatchString(postalCode) {
			errs["postalCode"] = ErrInvalidCAPostalCode
		}
	}
}
