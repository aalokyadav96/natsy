package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/checkout/usecase"
)

type CheckoutHandler struct {
	uc usecase.CheckoutUseCase
}

func NewCheckoutHandler(uc usecase.CheckoutUseCase) *CheckoutHandler {
	return &CheckoutHandler{uc: uc}
}

func (h *CheckoutHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/checkout", h.ProcessCheckout)
}

type checkoutRequest struct {
	UserID string `json:"user_id"`
}

type checkoutResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func (h *CheckoutHandler) ProcessCheckout(w http.ResponseWriter, r *http.Request) {
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "missing user_id in request body", http.StatusBadRequest)
		return
	}

	orderID, err := h.uc.ProcessCheckout(r.Context(), req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := checkoutResponse{
		OrderID: orderID,
		Status:  "ORDER_CREATED",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
