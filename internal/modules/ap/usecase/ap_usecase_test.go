package usecase

import (
	"context"
	"testing"
)

func TestPaymentBatchLifecycle(t *testing.T) {
	uc := NewAPUseCase(nil)

	batch, err := uc.CreatePaymentBatch(context.Background(), "SUP-1001", "AP-2026-001", 500.0)
	if err != nil {
		t.Fatalf("CreatePaymentBatch returned error: %v", err)
	}
	if batch == nil || batch.Status != "draft" {
		t.Fatalf("expected batch to start in draft state, got %#v", batch)
	}

	if err := uc.ApprovePaymentBatch(context.Background(), batch.ID); err != nil {
		t.Fatalf("ApprovePaymentBatch returned error: %v", err)
	}

	payment, err := uc.ExecutePayment(context.Background(), batch.ID, "bank-001")
	if err != nil {
		t.Fatalf("ExecutePayment returned error: %v", err)
	}
	if payment == nil || payment.Status != "paid" {
		t.Fatalf("expected payment to be marked paid, got %#v", payment)
	}

	if err := uc.ReconcilePayment(context.Background(), payment.ID); err != nil {
		t.Fatalf("ReconcilePayment returned error: %v", err)
	}

	outstanding, err := uc.GetOutstandingBySupplier(context.Background(), "SUP-1001")
	if err != nil {
		t.Fatalf("GetOutstandingBySupplier returned error: %v", err)
	}
	if len(outstanding) != 0 {
		t.Fatalf("expected no outstanding payments after reconciliation, got %#v", outstanding)
	}
}
