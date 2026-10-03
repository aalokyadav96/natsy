package domain

import (
	"context"
	"time"
)

type Supplier struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type PurchaseOrder struct {
	ID         string    `json:"id"`
	SupplierID string    `json:"supplier_id"`
	Number     string    `json:"number"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
}

type GoodsReceipt struct {
	ID              string    `json:"id"`
	PurchaseOrderID string    `json:"purchase_order_id"`
	Quantity        int       `json:"quantity"`
	Status          string    `json:"status"`
	ReceivedAt      time.Time `json:"received_at"`
}

type Invoice struct {
	ID              string    `json:"id"`
	PurchaseOrderID string    `json:"purchase_order_id"`
	SupplierID      string    `json:"supplier_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	MatchedAt       time.Time `json:"matched_at,omitempty"`
}

type AccountsPayableEntry struct {
	ID              string    `json:"id"`
	SupplierID      string    `json:"supplier_id"`
	PurchaseOrderID string    `json:"purchase_order_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	PostedAt        time.Time `json:"posted_at"`
}

type ProcurementRepository interface {
	CreateSupplier(ctx context.Context, supplier *Supplier) error
	GetSupplier(ctx context.Context, id string) (*Supplier, error)
	CreatePurchaseOrder(ctx context.Context, purchaseOrder *PurchaseOrder) error
	GetPurchaseOrder(ctx context.Context, id string) (*PurchaseOrder, error)
	UpdatePurchaseOrder(ctx context.Context, purchaseOrder *PurchaseOrder) error
	CreateGoodsReceipt(ctx context.Context, receipt *GoodsReceipt) error
	GetGoodsReceiptByPurchaseOrder(ctx context.Context, purchaseOrderID string) (*GoodsReceipt, error)
	CreateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoiceByPurchaseOrder(ctx context.Context, purchaseOrderID string) (*Invoice, error)
	CreateAccountsPayableEntry(ctx context.Context, entry *AccountsPayableEntry) error
	GetAccountsPayableBySupplier(ctx context.Context, supplierID string) ([]*AccountsPayableEntry, error)
	GetAccountsPayableByPurchaseOrder(ctx context.Context, purchaseOrderID string) (*AccountsPayableEntry, error)
}
