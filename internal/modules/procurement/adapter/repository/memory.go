package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nae/internal/modules/procurement/domain"
)

type inMemoryProcurementRepository struct {
	mu                 sync.RWMutex
	suppliers          map[string]*domain.Supplier
	purchaseOrders     map[string]*domain.PurchaseOrder
	receipts           map[string]*domain.GoodsReceipt
	invoices           map[string]*domain.Invoice
	apEntriesByPO      map[string]*domain.AccountsPayableEntry
	apEntriesBySupplier map[string][]*domain.AccountsPayableEntry
}

func NewInMemoryProcurementRepository() domain.ProcurementRepository {
	return &inMemoryProcurementRepository{
		suppliers:           make(map[string]*domain.Supplier),
		purchaseOrders:      make(map[string]*domain.PurchaseOrder),
		receipts:            make(map[string]*domain.GoodsReceipt),
		invoices:            make(map[string]*domain.Invoice),
		apEntriesByPO:       make(map[string]*domain.AccountsPayableEntry),
		apEntriesBySupplier: make(map[string][]*domain.AccountsPayableEntry),
	}
}

func (r *inMemoryProcurementRepository) CreateSupplier(ctx context.Context, supplier *domain.Supplier) error {
	if supplier == nil {
		return fmt.Errorf("supplier is required")
	}
	if supplier.ID == "" {
		supplier.ID = fmt.Sprintf("sup_%d", time.Now().UnixNano())
	}
	if supplier.CreatedAt.IsZero() {
		supplier.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.suppliers[supplier.ID] = supplier
	return nil
}

func (r *inMemoryProcurementRepository) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	supplier, ok := r.suppliers[id]
	if !ok {
		return nil, fmt.Errorf("supplier %s not found", id)
	}
	clone := *supplier
	return &clone, nil
}

func (r *inMemoryProcurementRepository) CreatePurchaseOrder(ctx context.Context, purchaseOrder *domain.PurchaseOrder) error {
	if purchaseOrder == nil {
		return fmt.Errorf("purchase order is required")
	}
	if purchaseOrder.ID == "" {
		purchaseOrder.ID = fmt.Sprintf("po_%d", time.Now().UnixNano())
	}
	if purchaseOrder.Status == "" {
		purchaseOrder.Status = "draft"
	}
	if purchaseOrder.CreatedAt.IsZero() {
		purchaseOrder.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.purchaseOrders[purchaseOrder.ID] = purchaseOrder
	return nil
}

func (r *inMemoryProcurementRepository) GetPurchaseOrder(ctx context.Context, id string) (*domain.PurchaseOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	po, ok := r.purchaseOrders[id]
	if !ok {
		return nil, fmt.Errorf("purchase order %s not found", id)
	}
	clone := *po
	return &clone, nil
}

func (r *inMemoryProcurementRepository) UpdatePurchaseOrder(ctx context.Context, purchaseOrder *domain.PurchaseOrder) error {
	if purchaseOrder == nil {
		return fmt.Errorf("purchase order is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.purchaseOrders[purchaseOrder.ID]; !ok {
		return fmt.Errorf("purchase order %s not found", purchaseOrder.ID)
	}
	r.purchaseOrders[purchaseOrder.ID] = purchaseOrder
	return nil
}

func (r *inMemoryProcurementRepository) CreateGoodsReceipt(ctx context.Context, receipt *domain.GoodsReceipt) error {
	if receipt == nil {
		return fmt.Errorf("goods receipt is required")
	}
	if receipt.ID == "" {
		receipt.ID = fmt.Sprintf("gr_%d", time.Now().UnixNano())
	}
	if receipt.Status == "" {
		receipt.Status = "received"
	}
	if receipt.ReceivedAt.IsZero() {
		receipt.ReceivedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.receipts[receipt.PurchaseOrderID] = receipt
	return nil
}

func (r *inMemoryProcurementRepository) GetGoodsReceiptByPurchaseOrder(ctx context.Context, purchaseOrderID string) (*domain.GoodsReceipt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	receipt, ok := r.receipts[purchaseOrderID]
	if !ok {
		return nil, fmt.Errorf("goods receipt for purchase order %s not found", purchaseOrderID)
	}
	clone := *receipt
	return &clone, nil
}

func (r *inMemoryProcurementRepository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	if invoice == nil {
		return fmt.Errorf("invoice is required")
	}
	if invoice.ID == "" {
		invoice.ID = fmt.Sprintf("inv_%d", time.Now().UnixNano())
	}
	if invoice.Status == "" {
		invoice.Status = "pending"
	}
	if invoice.CreatedAt.IsZero() {
		invoice.CreatedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoices[invoice.PurchaseOrderID] = invoice
	return nil
}

func (r *inMemoryProcurementRepository) GetInvoiceByPurchaseOrder(ctx context.Context, purchaseOrderID string) (*domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	invoice, ok := r.invoices[purchaseOrderID]
	if !ok {
		return nil, fmt.Errorf("invoice for purchase order %s not found", purchaseOrderID)
	}
	clone := *invoice
	return &clone, nil
}

func (r *inMemoryProcurementRepository) CreateAccountsPayableEntry(ctx context.Context, entry *domain.AccountsPayableEntry) error {
	if entry == nil {
		return fmt.Errorf("accounts payable entry is required")
	}
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("ap_%d", time.Now().UnixNano())
	}
	if entry.Status == "" {
		entry.Status = "posted"
	}
	if entry.PostedAt.IsZero() {
		entry.PostedAt = time.Now().UTC()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.apEntriesByPO[entry.PurchaseOrderID] = entry
	r.apEntriesBySupplier[entry.SupplierID] = append(r.apEntriesBySupplier[entry.SupplierID], entry)
	return nil
}

func (r *inMemoryProcurementRepository) GetAccountsPayableBySupplier(ctx context.Context, supplierID string) ([]*domain.AccountsPayableEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entries, ok := r.apEntriesBySupplier[supplierID]
	if !ok {
		return nil, fmt.Errorf("no accounts payable entries for supplier %s", supplierID)
	}
	result := make([]*domain.AccountsPayableEntry, len(entries))
	for i, entry := range entries {
		clone := *entry
		result[i] = &clone
	}
	return result, nil
}

func (r *inMemoryProcurementRepository) GetAccountsPayableByPurchaseOrder(ctx context.Context, purchaseOrderID string) (*domain.AccountsPayableEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.apEntriesByPO[purchaseOrderID]
	if !ok {
		return nil, fmt.Errorf("accounts payable entry for purchase order %s not found", purchaseOrderID)
	}
	clone := *entry
	return &clone, nil
}
