package api

import (
	"github.com/perkzen/mbus/apps/bus-service/internal/service/departure"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type DepartureHandler struct {
	departureService *departure.Service
	logger           *slog.Logger
}

func NewDepartureHandler(departureService *departure.Service, logger *slog.Logger) *DepartureHandler {
	return &DepartureHandler{
		departureService: departureService,
		logger:           logger.With(slog.String("handler", "DepartureHandler")),
	}
}

// GetDepartures godoc
// @Summary Get departures
// @Description Retrieve departures between two bus stations on a specific date
// @Tags Departures
// @Accept json
// @Produce json
// @Param from query int true "Departure station code"
// @Param to query int true "Arrival station code"
// @Param date query string false "Date in YYYY-MM-DD format"
// @Success 200 {array} departure.TimetableRow "List of departures"
// @Router /api/departures [get]
func (h *DepartureHandler) GetDepartures(w http.ResponseWriter, r *http.Request) error {
	fromID := QueryInt(r, "from", -1)
	toID := QueryInt(r, "to", -1)
	if fromID == -1 || toID == -1 {
		h.logger.Error("Invalid request parameters", slog.String("from", strconv.Itoa(fromID)), slog.String("to", strconv.Itoa(toID)))
		return BadRequestError("Both 'from' and 'to' parameters are required")
	}

	date := QueryDateStr(r, "date", utils.Today())

	data, err := h.departureService.GenerateTimetable(fromID, toID, date)
	if err != nil {
		h.logger.Error("Failed to generate timetable", slog.Any("error", err))
		return NotFoundError(err.Error())
	}

	return WriteJSON(w, http.StatusOK, data)
}
