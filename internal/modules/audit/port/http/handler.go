package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/audit/usecase"
)

type AuditHandler struct {
	uc usecase.AuditUseCase
}

func NewAuditHandler(uc usecase.AuditUseCase) *AuditHandler { return &AuditHandler{uc: uc} }

func (h *AuditHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/audit", h.Log)
}

func (h *AuditHandler) Log(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Module   string `json:"module"`
		Action   string `json:"action"`
		EntityID string `json:"entity_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	entry, err := h.uc.Log(r.Context(), req.Module, req.Action, req.EntityID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entry)
}
