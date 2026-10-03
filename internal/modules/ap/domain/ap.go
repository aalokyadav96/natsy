package domain

import (
	"context"
	"time"
)

type BankAccount struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
}

type PaymentBatch struct {
	ID         string    `json:"id"`
	SupplierID string    `json:"supplier_id"`
	Reference  string    `json:"reference"`
	Total      float64   `json:"total"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	ApprovedAt time.Time `json:"approved_at,omitempty"`
	PaidAt     time.Time `json:"paid_at,omitempty"`
}

type VendorPayment struct {
	ID            string    `json:"id"`
	BatchID       string    `json:"batch_id"`
	SupplierID    string    `json:"supplier_id"`
	BankAccountID string    `json:"bank_account_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	PaidAt        time.Time `json:"paid_at,omitempty"`
	ReconciledAt  time.Time `json:"reconciled_at,omitempty"`
}

type APRepository interface {
	CreatePaymentBatch(ctx context.Context, batch *PaymentBatch) error
	GetPaymentBatch(ctx context.Context, id string) (*PaymentBatch, error)
	UpdatePaymentBatch(ctx context.Context, batch *PaymentBatch) error
	CreateVendorPayment(ctx context.Context, payment *VendorPayment) error
	GetVendorPayment(ctx context.Context, id string) (*VendorPayment, error)
	UpdateVendorPayment(ctx context.Context, payment *VendorPayment) error
	ListPaymentsByBatch(ctx context.Context, batchID string) ([]*VendorPayment, error)
	ListPaymentsBySupplier(ctx context.Context, supplierID string) ([]*VendorPayment, error)
}
