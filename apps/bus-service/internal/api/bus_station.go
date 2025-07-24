package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/perkzen/mbus/apps/bus-service/internal/errs"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"log/slog"
	"net/http"
	"strconv"
)

type BusStationHandler struct {
	busStationStore store.BusStationStore
	logger          *slog.Logger
}

func NewBusStationHandler(busStationStore store.BusStationStore, logger *slog.Logger) *BusStationHandler {
	return &BusStationHandler{
		busStationStore: busStationStore,
		logger:          logger.With(slog.String("handler", "BusStationHandler")),
	}
}

// GetBusStations godoc
// @Summary Get bus stations
// @Description Retrieve a list of bus stations with optional filters
// @Tags Bus Stations
// @Accept json
// @Produce json
// @Param limit query int false "Limit the number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Param name query string false "Filter by bus station name"
// @Param line query string false "Filter by bus line"
// @Success 200 {array} store.BusStation "List of bus stations"
// @Router /api/bus-stations [get]
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
		return err
	}

	return WriteJSON(w, http.StatusOK, busStations)
}

// GetBusStationByID godoc
// @Summary Get bus station by id
// @Description Retrieve a bus station by its id
// @Tags Bus Stations
// @Accept json
// @Produce json
// @Param id path int true "Bus station id"
// @Success 200 {object} store.BusStation "Bus station details"
// @Router /api/bus-stations/{id} [get]
func (h *BusStationHandler) GetBusStationByID(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	if id == "" {
		return errs.BadRequestError("Bus station id is required")
	}

	stationID, err := strconv.Atoi(id)
	if err != nil {
		return errs.BadRequestError("Invalid bus station id format")
	}

	busStation, err := h.busStationStore.FindBusStationByID(stationID)
	if err != nil {
		return errs.BusStationNotFoundError(stationID)
	}

	return WriteJSON(w, http.StatusOK, busStation)
}
