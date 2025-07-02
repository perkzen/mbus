package api

import (
	"github.com/perkzen/mbus/bus-service/internal/store"
	"log/slog"
	"net/http"
)

type BusLineHandler struct {
	busLineStore store.BusLineStore
	logger       *slog.Logger
}

// NewBusLineHandler creates a new BusLineHandler with the provided bus line store and logger, attaching a handler-specific context to the logger.
func NewBusLineHandler(busLineStore store.BusLineStore, logger *slog.Logger) *BusLineHandler {
	return &BusLineHandler{
		busLineStore: busLineStore,
		logger:       logger.With(slog.String("handler", "BusLineHandler")),
	}
}

func (h *BusLineHandler) GetBusLines(w http.ResponseWriter, r *http.Request) error {
	lines, err := h.busLineStore.ListBusLines()

	if err != nil {
		h.logger.Error("failed to fetch station", slog.Any("error", err))
		return err
	}

	return WriteJSON(w, http.StatusOK, lines)
}
