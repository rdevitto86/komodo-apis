package address

import (
	"regexp"
	"strings"
)

func ValidateAddress(a Address) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(a.Street1) == "" {
		errs["street1"] = "required"
	}
	if strings.TrimSpace(a.City) == "" {
		errs["city"] = "required"
	}
	if strings.TrimSpace(a.State) == "" {
		errs["state"] = "required"
	}
	if strings.TrimSpace(a.PostalCode) == "" {
		errs["postalCode"] = "required"
	}

	country := strings.ToUpper(strings.TrimSpace(a.Country))
	if country == "" {
		country = "US"
	}

	// Simple postal validation for US & CA; others pass-through
	zipUS := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
	zipCA := regexp.MustCompile(`^[A-Za-z]\d[A-Za-z][ -]?\d[A-Za-z]\d$`)

	switch country {
	case "US", "USA", "UNITED STATES":
		if !zipUS.MatchString(strings.TrimSpace(a.PostalCode)) {
			errs["postalCode"] = "invalid US ZIP (12345 or 12345-6789)"
		}
	case "CA", "CANADA":
		if !zipCA.MatchString(strings.TrimSpace(a.PostalCode)) {
			errs["postalCode"] = "invalid CA postal code (A1A 1A1)"
		}
	}
	return errs
}
