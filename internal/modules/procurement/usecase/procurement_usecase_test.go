package usecase

import (
	"context"
	"testing"
)

func TestProcurementLifecycle(t *testing.T) {
	uc := NewProcurementUseCase(nil)

	supplier, err := uc.CreateSupplier(context.Background(), "SUP-1001", "Contoso Labs")
	if err != nil {
		t.Fatalf("CreateSupplier returned error: %v", err)
	}
	if supplier == nil || supplier.ID == "" {
		t.Fatal("expected supplier to be created")
	}

	po, err := uc.CreatePurchaseOrder(context.Background(), supplier.ID, "PO-1001", 2500.00)
	if err != nil {
		t.Fatalf("CreatePurchaseOrder returned error: %v", err)
	}
	if po == nil || po.Status != "draft" {
		t.Fatalf("expected purchase order to start in draft state, got %#v", po)
	}

	if err := uc.ApprovePurchaseOrder(context.Background(), po.ID); err != nil {
		t.Fatalf("ApprovePurchaseOrder returned error: %v", err)
	}

	receipt, err := uc.ReceiveGoods(context.Background(), po.ID, 12)
	if err != nil {
		t.Fatalf("ReceiveGoods returned error: %v", err)
	}
	if receipt == nil || receipt.Status != "received" {
		t.Fatalf("expected receipt status to be received, got %#v", receipt)
	}

	invoice, err := uc.CreateInvoice(context.Background(), po.ID, 2500.00)
	if err != nil {
		t.Fatalf("CreateInvoice returned error: %v", err)
	}
	if invoice == nil || invoice.Status != "matched" {
		t.Fatalf("expected invoice status to be matched after goods receipt, got %#v", invoice)
	}

	entries, err := uc.GetAccountsPayableBySupplier(context.Background(), supplier.ID)
	if err != nil {
		t.Fatalf("GetAccountsPayableBySupplier returned error: %v", err)
	}
	if len(entries) != 1 || entries[0].Status != "posted" {
		t.Fatalf("expected one posted AP entry for supplier %s, got %#v", supplier.ID, entries)
	}
}
