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

func (h *DepartureHandler) GetDepartures(w http.ResponseWriter, r *http.Request) error {
	fromCode := QueryInt(r, "from", -1)
	toCode := QueryInt(r, "to", -1)
	if fromCode == -1 || toCode == -1 {
		h.logger.Error("Invalid request parameters", slog.String("from", strconv.Itoa(fromCode)), slog.String("to", strconv.Itoa(toCode)))
		return BadRequestError("Both 'from' and 'to' parameters are required")
	}

	date := QueryDateStr(r, "date", utils.Today())

	data, err := h.departureService.GenerateTimetable(fromCode, toCode, date)
	if err != nil {
		h.logger.Error("Failed to generate timetable", slog.Any("error", err))
		return InternalServerError()
	}

	return WriteJSON(w, http.StatusOK, data)
}
