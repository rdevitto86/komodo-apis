package geocode

import (
	"context"
	"errors"
	"komodo-address-api/internal/address"
	"strings"
)

type NominatimGeocoder struct{ BaseURL string }

func (n *NominatimGeocoder) Geocode(ctx context.Context, a address.Address) (float64, float64, string, error) {
	if strings.TrimSpace(n.BaseURL) == "" {
		return 0, 0, "", errors.New("nominatim base URL not set")
	}
	return 0, 0, "", errors.New("nominatim geocoding not implemented in this template")
}