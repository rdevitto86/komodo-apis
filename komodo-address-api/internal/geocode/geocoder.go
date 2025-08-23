package geocode

import (
	"context"
	"komodo-address-api/internal/address"
)

type GeocodeResponse struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Accuracy   string  `json:"accuracy"` // e.g., ROOFTOP, APPROXIMATE
	Provider   string  `json:"provider"`
	Normalized address.Address `json:"normalized"`
}

type Geocoder interface {
	Geocode(ctx context.Context, a address.Address) (lat float64, lng float64, accuracy string, err error)
}
