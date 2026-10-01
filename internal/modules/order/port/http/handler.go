package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/order/domain"
	"nae/internal/modules/order/usecase"
)

type OrderHandler struct {
	uc usecase.OrderUseCase
}

func NewOrderHandler(uc usecase.OrderUseCase) *OrderHandler { return &OrderHandler{uc: uc} }

func (h *OrderHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/orders", h.Create)
	mux.HandleFunc("GET /api/v1/orders/{id}", h.GetByID)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string             `json:"user_id"`
		Items  []domain.OrderItem `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.uc.Create(r.Context(), req.UserID, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}
