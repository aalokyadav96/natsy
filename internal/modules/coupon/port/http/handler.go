package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/coupon/usecase"
)

type CouponHandler struct {
	uc usecase.CouponUseCase
}

func NewCouponHandler(uc usecase.CouponUseCase) *CouponHandler { return &CouponHandler{uc: uc} }

func (h *CouponHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/coupons/validate", h.Validate)
}

func (h *CouponHandler) Validate(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	coupon, err := h.uc.Validate(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(coupon)
}
