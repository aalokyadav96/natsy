package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/cart/usecase"
)

type CartHandler struct {
	uc usecase.CartUseCase
}

func NewCartHandler(uc usecase.CartUseCase) *CartHandler {
	return &CartHandler{uc: uc}
}

func (h *CartHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/cart", h.GetCart)
	mux.HandleFunc("POST /api/v1/cart/items", h.AddItem)
	mux.HandleFunc("DELETE /api/v1/cart", h.ClearCart)
}

type addItemReq struct {
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id query parameter", http.StatusBadRequest)
		return
	}

	cart, err := h.uc.GetCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req addItemReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.uc.AddToCart(r.Context(), req.UserID, req.ProductID, req.Quantity, req.UnitPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id query parameter", http.StatusBadRequest)
		return
	}

	if err := h.uc.ClearCart(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
