package httpapi

import (
	context "context"
	"fmt"
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/geocode"
	"net/http"
	"os"
	"strings"
	"time"
)

func HandleGeocode(w http.ResponseWriter, r *http.Request) {
	addr, err := ParseAddress(r)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, errorObj(err.Error()))
		return
	}

	// Normalize first for a cleaner geocode request
	norm := address.NormalizeAddress(addr)

	// Select provider: set GEOCODER=google|nominatim|mock and relevant API keys
	providerName := strings.ToLower(strings.TrimSpace(os.Getenv("GEOCODER")))
	var provider geocode.Geocoder

	switch providerName {
	case "google":
		provider = &geocode.GoogleGeocoder{APIKey: os.Getenv("GOOGLE_MAPS_API_KEY")}
	case "nominatim":
		provider = &geocode.NominatimGeocoder{BaseURL: "https://nominatim.openstreetmap.org"}
	default:
		provider = &geocode.MockGeocoder{}
		providerName = "mock"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	lat, lng, acc, err := provider.Geocode(ctx, norm)

	if err != nil {
		WriteJSON(w, http.StatusBadGateway, errorObj(fmt.Sprintf("geocoding failed via %s: %v", providerName, err)))
		return
	}

	WriteJSON(w, http.StatusOK, geocode.GeocodeResponse{
		Latitude:   lat,
		Longitude:  lng,
		Accuracy:   acc,
		Provider:   providerName,
		Normalized: norm,
	})
}