package address

import (
	"fmt"
	"regexp"
	"strings"
)

// Precompiled regex for ZIP code normalization
var zipDigits = regexp.MustCompile(`[^0-9]`)

// Constants for country codes
const (
	DefaultCountry = "US"
)

// NormalizeAddress cleans and standardizes an address.
func NormalizeAddress(a Address) Address {
	// Helper to clean and trim strings
	clean := func(s string) string {
		return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
	}

	// Helper to capitalize words while preserving specific acronyms
	capWords := func(s string) string {
		words := strings.Fields(strings.ToLower(clean(s)))
		for i, word := range words {
			switch strings.ToUpper(word) {
			case "NE", "NW", "SE", "SW", "N", "S", "E", "W", "US":
				words[i] = strings.ToUpper(word)
			default:
				words[i] = strings.Title(word) // Capitalize first letter
			}
		}
		return strings.Join(words, " ")
	}

	// Normalize country
	country := strings.ToUpper(clean(a.Country))
	if country == "" {
		country = DefaultCountry
	}

	// Normalize state and postal code
	state := strings.ToUpper(clean(a.State))
	postal := normalizePostalCode(clean(a.PostalCode), country)

	// Return normalized address
	return Address{
		Street1:    capWords(a.Street1),
		Street2:    capWords(a.Street2),
		City:       capWords(a.City),
		State:      state,
		PostalCode: postal,
		Country:    country,
	}
}

// normalizePostalCode standardizes postal codes based on the country.
func normalizePostalCode(postal, country string) string {
	// Handle US ZIP codes
	if country == "US" || country == "USA" || country == "UNITED STATES" {
		digits := zipDigits.ReplaceAllString(postal, "")
		if len(digits) >= 9 {
			return fmt.Sprintf("%s-%s", digits[:5], digits[5:9])
		} else if len(digits) >= 5 {
			return digits[:5]
		}
	}
	return postal
}