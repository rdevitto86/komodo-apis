package handlers

import (
	"context"
	"fmt"
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/geocode"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// HandleGeocodeGin processes geocoding requests for Gin.
func HandleGeocode(ctx *gin.Context) {
	// Parse the address from the request
	var addr address.Address

	if err := ctx.ShouldBindJSON(&addr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Normalize the address for a cleaner geocode request
	normalizedAddr := address.NormalizeAddress(addr)

	// Select the geocoder provider
	provider, providerName, err := selectGeocoder()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Set a timeout for the geocoding request
	timeout := getGeocodeTimeout()
	c, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
	defer cancel()

	// Perform the geocoding
	lat, lng, acc, err := provider.Geocode(c, normalizedAddr)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, err.Error())
		return
	}

	// Return the geocoding response
	ctx.JSON(http.StatusOK, geocode.GeocodeResponse{
		Latitude:   lat,
		Longitude:  lng,
		Accuracy:   acc,
		Provider:   providerName,
		Normalized: normalizedAddr,
	})
}

// selectGeocoder selects the geocoding provider based on environment variables.
func selectGeocoder() (geocode.Geocoder, string, error) {
	providerName := strings.ToLower(strings.TrimSpace(os.Getenv("GEOCODER_PROVIDER")))

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