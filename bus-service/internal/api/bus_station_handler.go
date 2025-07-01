package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/perkzen/mbus/bus-service/internal/store"
	"net/http"
	"strconv"
)

type BusStationHandler struct {
	busStationStore store.BusStationStore
}

func NewBusStationHandler(busStationStore store.BusStationStore) *BusStationHandler {
	return &BusStationHandler{
		busStationStore: busStationStore,
	}
}

func (h *BusStationHandler) GetBusStations(w http.ResponseWriter, r *http.Request) error {
	limit := QueryInt(r, "limit", 10)
	offset := QueryInt(r, "offset", 0)

	name, _ := QueryStr(r, "name")
	line, _ := QueryStr(r, "line")

	busStations, err := h.busStationStore.ListBusStations(limit, offset, &store.BusStationFilterOptions{
		Name: name,
		Line: line,
	})
	if err != nil {
		fmt.Println("Error fetching bus stations:", err.Error())
		return err
	}

	return WriteJSON(w, http.StatusOK, busStations)
}

func (h *BusStationHandler) GetBusStationByCode(w http.ResponseWriter, r *http.Request) error {

	code := chi.URLParam(r, "code")
	if code == "" {
		return BadRequestError("Bus station code is required")
	}

	stationCode, err := strconv.Atoi(code)
	if err != nil {
		return BadRequestError("Invalid bus station code format")
	}

	busStation, err := h.busStationStore.FindBusStationByCode(stationCode)
	if err != nil || busStation == nil {
		fmt.Println("Error fetching bus station:", err.Error())
		return NotFoundError("Bus station not found")
	}

	return WriteJSON(w, http.StatusOK, busStation)
}
