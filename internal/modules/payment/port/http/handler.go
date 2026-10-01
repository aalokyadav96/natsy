package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/payment/usecase"
)

type PaymentHandler struct {
	uc usecase.PaymentUseCase
}

func NewPaymentHandler(uc usecase.PaymentUseCase) *PaymentHandler { return &PaymentHandler{uc: uc} }

func (h *PaymentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/payments", h.Create)
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OrderID string  `json:"order_id"`
		Amount  float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	payment, err := h.uc.Create(r.Context(), req.OrderID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payment)
}
