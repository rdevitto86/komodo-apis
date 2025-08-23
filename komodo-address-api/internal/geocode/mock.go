package geocode

import (
	"context"
	"komodo-address-api/internal/address"
	"math"
	"strings"
)

type MockGeocoder struct{}

func (m *MockGeocoder) Geocode(ctx context.Context, a address.Address) (float64, float64, string, error) {
	// Very simple hash -> lat/lng within plausible bounds.
	key := strings.ToUpper(strings.Join([]string{a.Street1, a.City, a.State, a.PostalCode, a.Country}, ","))
	var h uint64 = 1469598103934665603 // FNV-1a 64-bit offset basis

	for i := 0; i < len(key); i++ {
		h ^= uint64(key[i])
		h *= 1099511628211
	}

	lat := float64(int64(h%180000)-90000) / 1000.0                 // -90..+90
	lng := float64(int64((h/180000)%360000)-180000) / 1000.0 // -180..+180
	// Snap to 6 decimals
	lat = math.Round(lat*1e6) / 1e6
	lng = math.Round(lng*1e6) / 1e6

	return lat, lng, "APPROXIMATE", nil
}