package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/procurement/usecase"
)

type ProcurementHandler struct {
	uc usecase.ProcurementUseCase
}

func NewProcurementHandler(uc usecase.ProcurementUseCase) *ProcurementHandler {
	return &ProcurementHandler{uc: uc}
}

func (h *ProcurementHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/procurement/suppliers", h.CreateSupplier)
	mux.HandleFunc("POST /api/v1/procurement/purchase-orders", h.CreatePurchaseOrder)
	mux.HandleFunc("POST /api/v1/procurement/purchase-orders/{id}/approve", h.ApprovePurchaseOrder)
	mux.HandleFunc("POST /api/v1/procurement/purchase-orders/{id}/receive", h.ReceiveGoods)
	mux.HandleFunc("POST /api/v1/procurement/purchase-orders/{id}/invoice", h.CreateInvoice)
}

type createSupplierRequest struct {
	SupplierID string `json:"supplier_id"`
	Name       string `json:"name"`
}

func (h *ProcurementHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	var req createSupplierRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	supplier, err := h.uc.CreateSupplier(r.Context(), req.SupplierID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(supplier)
}

type createPurchaseOrderRequest struct {
	SupplierID string  `json:"supplier_id"`
	Number     string  `json:"number"`
	Amount     float64 `json:"amount"`
}

func (h *ProcurementHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var req createPurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	purchaseOrder, err := h.uc.CreatePurchaseOrder(r.Context(), req.SupplierID, req.Number, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(purchaseOrder)
}

func (h *ProcurementHandler) ApprovePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	purchaseOrderID := r.PathValue("id")
	if purchaseOrderID == "" {
		http.Error(w, "missing purchase order id", http.StatusBadRequest)
		return
	}

	if err := h.uc.ApprovePurchaseOrder(r.Context(), purchaseOrderID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"purchase_order_id": purchaseOrderID, "status": "approved"})
}

type receiveGoodsRequest struct {
	Quantity int `json:"quantity"`
}

func (h *ProcurementHandler) ReceiveGoods(w http.ResponseWriter, r *http.Request) {
	purchaseOrderID := r.PathValue("id")
	if purchaseOrderID == "" {
		http.Error(w, "missing purchase order id", http.StatusBadRequest)
		return
	}

	var req receiveGoodsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	receipt, err := h.uc.ReceiveGoods(r.Context(), purchaseOrderID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(receipt)
}

type createInvoiceRequest struct {
	Amount float64 `json:"amount"`
}

func (h *ProcurementHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	purchaseOrderID := r.PathValue("id")
	if purchaseOrderID == "" {
		http.Error(w, "missing purchase order id", http.StatusBadRequest)
		return
	}

	var req createInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	invoice, err := h.uc.CreateInvoice(r.Context(), purchaseOrderID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(invoice)
}
