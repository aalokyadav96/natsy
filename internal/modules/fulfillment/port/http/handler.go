package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/fulfillment/usecase"
)

type FulfillmentHandler struct {
	uc usecase.FulfillmentUseCase
}

func NewFulfillmentHandler(uc usecase.FulfillmentUseCase) *FulfillmentHandler {
	return &FulfillmentHandler{uc: uc}
}

func (h *FulfillmentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/fulfillment", h.Create)
}

func (h *FulfillmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID string `json:"order_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	item, err := h.uc.Create(r.Context(), req.OrderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(item)
}
