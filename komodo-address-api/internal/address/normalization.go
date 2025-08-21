package address

import (
	"fmt"
	"regexp"
	"strings"
)

func NormalizeAddress(a Address) Address {
	clean := func(s string) string {
		s = strings.TrimSpace(s)
		s = strings.Join(strings.Fields(s), " ") // collapse whitespace
		return s
	}

	capWords := func(s string) string {
		s = strings.ToLower(clean(s))
		parts := strings.Split(s, " ")
	
		for i, p := range parts {
			if len(p) == 0 {
				continue
			}

			// Honor common prefixes/suffixes/acronyms
			switch strings.ToUpper(p) {
			case "NE", "NW", "SE", "SW", "N", "S", "E", "W", "US":
				parts[i] = strings.ToUpper(p)
			default:
				runes := []rune(p)
				runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
				parts[i] = string(runes)
			}
		}
		return strings.Join(parts, " ")
	}

	country := strings.ToUpper(clean(a.Country))

	if country == "" {
		country = "US"
	}

	state := strings.ToUpper(clean(a.State))
	postal := strings.ToUpper(clean(a.PostalCode))
	zipDigits := regexp.MustCompile(`[^0-9]`) // For US ZIP+4: keep 5 or 9 digits with hyphen format

	if country == "US" || country == "USA" || country == "UNITED STATES" {
		d := zipDigits.ReplaceAllString(postal, "")
	
		if len(d) >= 9 {
			postal = fmt.Sprintf("%s-%s", d[:5], d[5:9])
		} else if len(d) >= 5 {
			postal = d[:5]
		}
	}

	return Address{
		Street1:    capWords(a.Street1),
		Street2:    capWords(a.Street2),
		City:       capWords(a.City),
		State:      state,
		PostalCode: postal,
		Country:    country,
	}
}