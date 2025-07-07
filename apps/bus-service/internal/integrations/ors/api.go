package ors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type API interface {
	GetMatrix(locations [][]float64) (*MatrixResponse, error)
}

type APIClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewAPIClient(apiKey string) *APIClient {

	return &APIClient{
		apiKey:  apiKey,
		baseURL: "https://api.openrouteservice.org/v2",
		client:  &http.Client{},
	}
}

func (c *APIClient) GetMatrix(locations [][]float64) (*MatrixResponse, error) {
	reqBody := MatrixRequest{
		Locations:        locations,
		Metrics:          []string{"distance", "duration"},
		ResolveLocations: true,
		Units:            "km",
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/matrix/driving-car", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ORS API error: %s", resp.Status)
	}

	var result MatrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}
