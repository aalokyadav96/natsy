package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	repository "nae/internal/modules/procurement/adapter/repository"
	"nae/internal/modules/procurement/usecase"
)

func TestProcurementHandlerLifecycle(t *testing.T) {
	uc := usecase.NewProcurementUseCase(repository.NewInMemoryProcurementRepository())
	h := NewProcurementHandler(uc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	supplierReq := map[string]any{"supplier_id": "SUP-1001", "name": "Contoso Labs"}
	supplierBody, _ := json.Marshal(supplierReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/procurement/suppliers", bytes.NewReader(supplierBody))
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected supplier creation status 201, got %d: %s", res.Code, res.Body.String())
	}

	poReq := map[string]any{"supplier_id": "SUP-1001", "number": "PO-1001", "amount": 2500.0}
	poBody, _ := json.Marshal(poReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/procurement/purchase-orders", bytes.NewReader(poBody))
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected purchase order creation status 201, got %d: %s", res.Code, res.Body.String())
	}

	var po map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &po); err != nil {
		t.Fatalf("failed to parse purchase order response: %v", err)
	}
	poID := po["id"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/procurement/purchase-orders/"+poID+"/approve", nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected purchase order approval status 200, got %d: %s", res.Code, res.Body.String())
	}

	receiveReq := map[string]any{"quantity": 12}
	receiveBody, _ := json.Marshal(receiveReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/procurement/purchase-orders/"+poID+"/receive", bytes.NewReader(receiveBody))
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected goods receipt status 200, got %d: %s", res.Code, res.Body.String())
	}

	invoiceReq := map[string]any{"amount": 2500.0}
	invoiceBody, _ := json.Marshal(invoiceReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/procurement/purchase-orders/"+poID+"/invoice", bytes.NewReader(invoiceBody))
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected invoice creation status 201, got %d: %s", res.Code, res.Body.String())
	}
}
