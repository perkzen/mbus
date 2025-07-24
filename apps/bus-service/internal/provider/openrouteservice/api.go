package openrouteservice

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"github.com/redis/go-redis/v9"
)

type API interface {
	GetMatrix(locations [][]float64) (*MatrixResponse, error)
}

type APIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
	cache   *redis.Client
}

func NewAPIClient(apiKey string, cache *redis.Client) *APIClient {
	return &APIClient{
		apiKey:  apiKey,
		baseURL: "https://api.openrouteservice.org/v2",
		client:  &http.Client{},
		cache:   cache,
	}
}

func (c *APIClient) GetMatrix(locations [][]float64) (*MatrixResponse, error) {
	ctx := context.Background()

	locBytes, _ := json.Marshal(locations)
	hash := sha256.Sum256(locBytes)
	cacheKey := fmt.Sprintf("ors_matrix_%x", hash)

	loader := func() (MatrixResponse, error) {
		return c.fetchMatrix(locations)
	}

	result, err := utils.WithCache(ctx, c.cache, cacheKey, 24*time.Hour, loader)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *APIClient) fetchMatrix(locations [][]float64) (MatrixResponse, error) {
	reqBody := MatrixRequest{
		Locations:        locations,
		Metrics:          []string{"distance", "duration"},
		ResolveLocations: true,
		Units:            "km",
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return MatrixResponse{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/matrix/driving-car", bytes.NewBuffer(data))
	if err != nil {
		return MatrixResponse{}, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return MatrixResponse{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return MatrixResponse{}, fmt.Errorf("ORS API error: %s", resp.Status)
	}

	var result MatrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return MatrixResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert durations from seconds to minutes + 1min/km
	for i := range result.Durations {
		for j := range result.Durations[i] {
			durationMin := result.Durations[i][j] / 60
			extraMinutes := result.Distances[i][j]
			result.Durations[i][j] = math.Round(durationMin + extraMinutes)
		}
	}

	return result, nil
}
