package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/inventory/usecase"
)

type InventoryHandler struct {
	uc usecase.InventoryUseCase
}

func NewInventoryHandler(uc usecase.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{uc: uc}
}

func (h *InventoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/inventory/{product_id}", h.GetInventory)
}

func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("product_id")
	item, err := h.uc.GetInventory(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(item)
}
