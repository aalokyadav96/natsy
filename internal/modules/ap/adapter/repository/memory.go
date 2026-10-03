package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nae/internal/modules/ap/domain"
)

type inMemoryAPRepository struct {
	mu                 sync.RWMutex
	batches            map[string]*domain.PaymentBatch
	payments           map[string]*domain.VendorPayment
	paymentsByBatch    map[string][]*domain.VendorPayment
	paymentsBySupplier map[string][]*domain.VendorPayment
}

func NewInMemoryAPRepository() domain.APRepository {
	return &inMemoryAPRepository{
		batches:            make(map[string]*domain.PaymentBatch),
		payments:           make(map[string]*domain.VendorPayment),
		paymentsByBatch:    make(map[string][]*domain.VendorPayment),
		paymentsBySupplier: make(map[string][]*domain.VendorPayment),
	}
}

func (r *inMemoryAPRepository) CreatePaymentBatch(ctx context.Context, batch *domain.PaymentBatch) error {
	if batch == nil {
		return fmt.Errorf("payment batch is required")
	}
	if batch.ID == "" {
		batch.ID = fmt.Sprintf("batch_%d", time.Now().UnixNano())
	}
	if batch.Status == "" {
		batch.Status = "draft"
	}
	if batch.CreatedAt.IsZero() {
		batch.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.batches[batch.ID] = batch
	return nil
}

func (r *inMemoryAPRepository) GetPaymentBatch(ctx context.Context, id string) (*domain.PaymentBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	batch, ok := r.batches[id]
	if !ok {
		return nil, fmt.Errorf("payment batch %s not found", id)
	}
	clone := *batch
	return &clone, nil
}

func (r *inMemoryAPRepository) UpdatePaymentBatch(ctx context.Context, batch *domain.PaymentBatch) error {
	if batch == nil {
		return fmt.Errorf("payment batch is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.batches[batch.ID]; !ok {
		return fmt.Errorf("payment batch %s not found", batch.ID)
	}
	r.batches[batch.ID] = batch
	return nil
}

func (r *inMemoryAPRepository) CreateVendorPayment(ctx context.Context, payment *domain.VendorPayment) error {
	if payment == nil {
		return fmt.Errorf("vendor payment is required")
	}
	if payment.ID == "" {
		payment.ID = fmt.Sprintf("vp_%d", time.Now().UnixNano())
	}
	if payment.Status == "" {
		payment.Status = "paid"
	}
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.payments[payment.ID] = payment
	r.paymentsByBatch[payment.BatchID] = append(r.paymentsByBatch[payment.BatchID], payment)
	r.paymentsBySupplier[payment.SupplierID] = append(r.paymentsBySupplier[payment.SupplierID], payment)
	return nil
}

func (r *inMemoryAPRepository) GetVendorPayment(ctx context.Context, id string) (*domain.VendorPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	payment, ok := r.payments[id]
	if !ok {
		return nil, fmt.Errorf("vendor payment %s not found", id)
	}
	clone := *payment
	return &clone, nil
}

func (r *inMemoryAPRepository) UpdateVendorPayment(ctx context.Context, payment *domain.VendorPayment) error {
	if payment == nil {
		return fmt.Errorf("vendor payment is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.payments[payment.ID]; !ok {
		return fmt.Errorf("vendor payment %s not found", payment.ID)
	}
	r.payments[payment.ID] = payment
	for i, existing := range r.paymentsByBatch[payment.BatchID] {
		if existing.ID == payment.ID {
			r.paymentsByBatch[payment.BatchID][i] = payment
			break
		}
	}
	for i, existing := range r.paymentsBySupplier[payment.SupplierID] {
		if existing.ID == payment.ID {
			r.paymentsBySupplier[payment.SupplierID][i] = payment
			break
		}
	}
	return nil
}

func (r *inMemoryAPRepository) ListPaymentsByBatch(ctx context.Context, batchID string) ([]*domain.VendorPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.paymentsByBatch[batchID]
	result := make([]*domain.VendorPayment, len(items))
	for i, payment := range items {
		clone := *payment
		result[i] = &clone
	}
	return result, nil
}

func (r *inMemoryAPRepository) ListPaymentsBySupplier(ctx context.Context, supplierID string) ([]*domain.VendorPayment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.paymentsBySupplier[supplierID]
	result := make([]*domain.VendorPayment, len(items))
	for i, payment := range items {
		clone := *payment
		result[i] = &clone
	}
	return result, nil
}
