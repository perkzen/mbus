package api

import (
	"github.com/perkzen/mbus/bus-service/internal/store"
	"net/http"
)

type BusLineHandler struct {
	busLineStore store.BusLineStore
}

func NewBusLineHandler(busLineStore store.BusLineStore) *BusLineHandler {
	return &BusLineHandler{
		busLineStore: busLineStore,
	}
}

func (h *BusLineHandler) ListBusLines(w http.ResponseWriter, r *http.Request) error {
	lines, err := h.busLineStore.ListBusLines()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, lines)
}
