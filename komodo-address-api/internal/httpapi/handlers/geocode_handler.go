package handlers

import (
	"context"
	"fmt"
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/geocode"
	"komodo-address-api/internal/httpapi"
	"net/http"
	"os"
	"strings"
	"time"
)

// HandleGeocode processes geocoding requests.
func HandleGeocode(w http.ResponseWriter, r *http.Request) {
	// Parse the address from the request
	addr, err := httpapi.ParseAddress(r)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorObj("invalid address: "+err.Error()))
		return
	}

	// Normalize the address for a cleaner geocode request
	norm := address.NormalizeAddress(addr)

	// Select the geocoder provider
	provider, providerName, err := selectGeocoder()
	if err != nil {
		httpapi.WriteJSON(w, http.StatusInternalServerError, httpapi.ErrorObj(err.Error()))
		return
	}

	// Set a timeout for the geocoding request
	timeout := getGeocodeTimeout()
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	// Perform the geocoding
	lat, lng, acc, err := provider.Geocode(ctx, norm)
	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadGateway, httpapi.ErrorObj(fmt.Sprintf("geocoding failed via %s: %v", providerName, err)))
		return
	}

	// Return the geocoding response
	httpapi.WriteJSON(w, http.StatusOK, geocode.GeocodeResponse{
		Latitude:   lat,
		Longitude:  lng,
		Accuracy:   acc,
		Provider:   providerName,
		Normalized: norm,
	})
}

// selectGeocoder selects the geocoding provider based on environment variables.
func selectGeocoder() (geocode.Geocoder, string, error) {
	providerName := strings.ToLower(strings.TrimSpace(os.Getenv("GEOCODER")))
	switch providerName {
	case "google":
		apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
		if strings.TrimSpace(apiKey) == "" {
			return nil, "", fmt.Errorf("missing GOOGLE_MAPS_API_KEY for Google geocoder")
		}
		return &geocode.GoogleGeocoder{APIKey: apiKey}, "google", nil
	case "nominatim":
		baseURL := os.Getenv("NOMINATIM_BASE_URL")
		if strings.TrimSpace(baseURL) == "" {
			baseURL = "https://nominatim.openstreetmap.org" // Default Nominatim URL
		}
		return &geocode.NominatimGeocoder{BaseURL: baseURL}, "nominatim", nil
	default:
		return &geocode.MockGeocoder{}, "mock", nil
	}
}

// getGeocodeTimeout retrieves the geocoding timeout from an environment variable or defaults to 5 seconds.
func getGeocodeTimeout() time.Duration {
	timeoutStr := os.Getenv("GEOCODE_TIMEOUT")
	if timeoutStr == "" {
		return 5 * time.Second // Default timeout
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return 5 * time.Second // Fallback to default if parsing fails
	}
	return timeout
}