package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apRepo "nae/internal/modules/ap/adapter/repository"
	apUseCase "nae/internal/modules/ap/usecase"
)

func TestAPHandlerLifecycle(t *testing.T) {
	uc := apUseCase.NewAPUseCase(apRepo.NewInMemoryAPRepository())
	h := NewAPHandler(uc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	batchReq := map[string]any{"supplier_id": "SUP-1001", "reference": "AP-2026-001", "total": 500.0}
	batchBody, _ := json.Marshal(batchReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ap/payment-batches", bytes.NewReader(batchBody))
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected payment batch creation status 201, got %d: %s", res.Code, res.Body.String())
	}

	var batch map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &batch); err != nil {
		t.Fatalf("failed to parse payment batch response: %v", err)
	}
	batchID := batch["id"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/v1/ap/payment-batches/"+batchID+"/approve", nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected approve status 200, got %d: %s", res.Code, res.Body.String())
	}

	executeReq := map[string]any{"bank_account_id": "bank-001"}
	executeBody, _ := json.Marshal(executeReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/ap/payment-batches/"+batchID+"/execute", bytes.NewReader(executeBody))
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected execute status 200, got %d: %s", res.Code, res.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/ap/payment-batches/"+batchID+"/reconcile?supplier_id=SUP-1001", nil)
	res = httptest.NewRecorder()
	mux.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected reconcile status 200, got %d: %s", res.Code, res.Body.String())
	}
}
