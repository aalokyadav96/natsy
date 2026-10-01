package http

import (
	"encoding/json"
	"net/http"

	"nae/internal/modules/catalog/usecase"
)

type CatalogHandler struct {
	uc usecase.CatalogUseCase
}

func NewCatalogHandler(uc usecase.CatalogUseCase) *CatalogHandler {
	return &CatalogHandler{uc: uc}
}

func (h *CatalogHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/catalog", h.ListItems)
}

func (h *CatalogHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.ListItems(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}
