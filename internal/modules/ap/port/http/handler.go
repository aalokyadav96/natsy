package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/ap/usecase"
)

type APHandler struct {
	uc usecase.APUseCase
}

func NewAPHandler(uc usecase.APUseCase) *APHandler { return &APHandler{uc: uc} }

func (h *APHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/ap/payment-batches", h.CreatePaymentBatch)
	mux.HandleFunc("POST /api/v1/ap/payment-batches/{id}/approve", h.ApprovePaymentBatch)
	mux.HandleFunc("POST /api/v1/ap/payment-batches/{id}/execute", h.ExecutePayment)
	mux.HandleFunc("POST /api/v1/ap/payment-batches/{id}/reconcile", h.ReconcilePayment)
}

type createPaymentBatchRequest struct {
	SupplierID string  `json:"supplier_id"`
	Reference  string  `json:"reference"`
	Total      float64 `json:"total"`
}

func (h *APHandler) CreatePaymentBatch(w http.ResponseWriter, r *http.Request) {
	var req createPaymentBatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	batch, err := h.uc.CreatePaymentBatch(r.Context(), req.SupplierID, req.Reference, req.Total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(batch)
}

func (h *APHandler) ApprovePaymentBatch(w http.ResponseWriter, r *http.Request) {
	batchID := r.PathValue("id")
	if batchID == "" {
		http.Error(w, "missing payment batch id", http.StatusBadRequest)
		return
	}
	if err := h.uc.ApprovePaymentBatch(r.Context(), batchID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"batch_id": batchID, "status": "approved"})
}

type executePaymentRequest struct {
	BankAccountID string `json:"bank_account_id"`
}

func (h *APHandler) ExecutePayment(w http.ResponseWriter, r *http.Request) {
	batchID := r.PathValue("id")
	if batchID == "" {
		http.Error(w, "missing payment batch id", http.StatusBadRequest)
		return
	}
	var req executePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	payment, err := h.uc.ExecutePayment(r.Context(), batchID, req.BankAccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payment)
}

func (h *APHandler) ReconcilePayment(w http.ResponseWriter, r *http.Request) {
	batchID := r.PathValue("id")
	if batchID == "" {
		http.Error(w, "missing payment batch id", http.StatusBadRequest)
		return
	}
	payments, err := h.uc.GetOutstandingBySupplier(r.Context(), r.URL.Query().Get("supplier_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, payment := range payments {
		if payment.BatchID == batchID {
			if err := h.uc.ReconcilePayment(r.Context(), payment.ID); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"batch_id": batchID, "status": "reconciled"})
}
