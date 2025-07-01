package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/perkzen/mbus/bus-service/internal/integrations/marprom"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type DepartureHandler struct {
	marpromClient *marprom.APIClient
	cache         *redis.Client
}

func NewDepartureHandler(cache *redis.Client) *DepartureHandler {
	return &DepartureHandler{
		marpromClient: marprom.NewAPIClient(),
		cache:         cache,
	}
}

func (h *DepartureHandler) GetDeparturesByStation(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	code := chi.URLParam(r, "code")
	if code == "" {
		return BadRequestError("Bus station code is required")
	}

	line, _ := QueryStr(r, "line")

	cacheKey := buildCacheKey(code, line)
	cached, err := h.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		var departures []marprom.Departure
		if err := json.Unmarshal([]byte(cached), &departures); err == nil {
			return WriteJSON(w, http.StatusOK, departures)
		}
	}

	details, err := h.marpromClient.GetDeparturesByBusStation(&marprom.DepartureFilterOptions{
		Code: code,
		Line: line,
	})

	if err != nil {
		return ServiceUnavailableError("Marprom service is unavailable")
	}

	data, _ := json.Marshal(details)
	_ = h.cache.Set(ctx, cacheKey, data, 24*time.Hour).Err()
	return WriteJSON(w, http.StatusOK, details)
}

func buildCacheKey(code, line string) string {
	if line != "" {
		return fmt.Sprintf("departures_%s_%s", code, line)
	}
	return fmt.Sprintf("departures_%s", code)
}
