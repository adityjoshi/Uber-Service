package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/adityjoshi/Uber-Service/services/matching-service/internal/dto"
)

type LocationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewLocationClient() *LocationClient {
	return &LocationClient{
		baseURL: getenv("LOCATION_SERVICE_URL", "http://location-service:8082"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *LocationClient) GetNearbyDrivers(ctx context.Context, latitude, longitude, radius float64) ([]dto.NearByDriverResponse, error) {
	url := fmt.Sprintf(
		"%s/api/v1/locations/driver/nearby?latitude=%f&longitude=%f&radius=%f",
		c.baseURL, latitude, longitude, radius,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("location client: build request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("location client: Get nearby drivers: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("location clien: unexpected status %d", resp.StatusCode)
	}

	var drivers []dto.NearByDriverResponse
	if err := json.NewDecoder(resp.Body).Decode(&drivers); err != nil {
		return nil, fmt.Errorf("location client: decode response: %w", err)
	}
	return drivers, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
