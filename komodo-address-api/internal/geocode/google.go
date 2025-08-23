package geocode

import (
	"context"
	"errors"
	"komodo-address-api/internal/address"
	"strings"
)

type GoogleGeocoder struct{ APIKey string }

func (g *GoogleGeocoder) Geocode(ctx context.Context, a address.Address) (float64, float64, string, error) {
	if strings.TrimSpace(g.APIKey) == "" {
		return 0, 0, "", errors.New("GOOGLE_MAPS_API_KEY not set")
	}
	return 0, 0, "", errors.New("google geocoding not implemented in this template")
}