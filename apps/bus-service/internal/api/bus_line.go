package api

import (
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"log/slog"
	"net/http"
)

type BusLineHandler struct {
	busLineStore store.BusLineStore
	logger       *slog.Logger
}

func NewBusLineHandler(busLineStore store.BusLineStore, logger *slog.Logger) *BusLineHandler {
	return &BusLineHandler{
		busLineStore: busLineStore,
		logger:       logger.With(slog.String("handler", "BusLineHandler")),
	}
}

// GetBusLines godoc
// @Summary Get bus lines
// @Description Retrieve a list of bus lines
// @Tags BusLines
// @Accept json
// @Produce json
// @Success 200 {array} store.BusLine "List of bus lines"
// @Router /api/bus-lines [get]
func (h *BusLineHandler) GetBusLines(w http.ResponseWriter, r *http.Request) error {
	lines, err := h.busLineStore.ListBusLines()

	if err != nil {
		h.logger.Error("failed to fetch bus-lines", slog.Any("error", err))
		return err
	}

	return WriteJSON(w, http.StatusOK, lines)
}
