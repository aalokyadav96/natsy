package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nae/internal/modules/procurement/adapter/repository"
	"nae/internal/modules/procurement/domain"
	sharedDomain "nae/internal/shared/domain"
	"nae/internal/shared/infra/mq"
)

type ProcurementUseCase interface {
	CreateSupplier(ctx context.Context, supplierID, name string) (*domain.Supplier, error)
	CreatePurchaseOrder(ctx context.Context, supplierID, purchaseOrderNumber string, amount float64) (*domain.PurchaseOrder, error)
	ApprovePurchaseOrder(ctx context.Context, purchaseOrderID string) error
	ReceiveGoods(ctx context.Context, purchaseOrderID string, quantity int) (*domain.GoodsReceipt, error)
	CreateInvoice(ctx context.Context, purchaseOrderID string, amount float64) (*domain.Invoice, error)
	PostAccountsPayable(ctx context.Context, purchaseOrderID string) (*domain.AccountsPayableEntry, error)
	GetAccountsPayableBySupplier(ctx context.Context, supplierID string) ([]*domain.AccountsPayableEntry, error)
}

type procurementUseCase struct {
	repo domain.ProcurementRepository
	nats *mq.NATSClient
}

func NewProcurementUseCase(repo domain.ProcurementRepository, nats ...*mq.NATSClient) ProcurementUseCase {
	if repo == nil {
		repo = repository.NewInMemoryProcurementRepository()
	}
	var client *mq.NATSClient
	if len(nats) > 0 {
		client = nats[0]
	}
	return &procurementUseCase{repo: repo, nats: client}
}

func (u *procurementUseCase) CreateSupplier(ctx context.Context, supplierID, name string) (*domain.Supplier, error) {
	if supplierID == "" {
		return nil, fmt.Errorf("%w: supplier id is required", sharedDomain.ErrInvalidInput)
	}
	if name == "" {
		return nil, fmt.Errorf("%w: supplier name is required", sharedDomain.ErrInvalidInput)
	}

	supplier := &domain.Supplier{ID: supplierID, Name: name}
	if err := u.repo.CreateSupplier(ctx, supplier); err != nil {
		return nil, err
	}
	return supplier, nil
}

func (u *procurementUseCase) CreatePurchaseOrder(ctx context.Context, supplierID, purchaseOrderNumber string, amount float64) (*domain.PurchaseOrder, error) {
	if supplierID == "" {
		return nil, fmt.Errorf("%w: supplier id is required", sharedDomain.ErrInvalidInput)
	}
	if _, err := u.repo.GetSupplier(ctx, supplierID); err != nil {
		return nil, fmt.Errorf("%w: supplier %s is not registered", sharedDomain.ErrNotFound, supplierID)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("%w: purchase amount must be greater than zero", sharedDomain.ErrInvalidInput)
	}

	po := &domain.PurchaseOrder{
		ID:         fmt.Sprintf("po_%d", time.Now().UnixNano()),
		SupplierID: supplierID,
		Number:     purchaseOrderNumber,
		Amount:     amount,
		Status:     sharedDomain.PurchaseOrderStatusDraft,
		CreatedAt:  time.Now().UTC(),
	}
	if po.Number == "" {
		po.Number = fmt.Sprintf("PO-%d", time.Now().UnixNano())
	}
	if err := u.repo.CreatePurchaseOrder(ctx, po); err != nil {
		return nil, err
	}
	return po, nil
}

func (u *procurementUseCase) ApprovePurchaseOrder(ctx context.Context, purchaseOrderID string) error {
	if purchaseOrderID == "" {
		return fmt.Errorf("%w: purchase order id is required", sharedDomain.ErrInvalidInput)
	}
	po, err := u.repo.GetPurchaseOrder(ctx, purchaseOrderID)
	if err != nil {
		return err
	}
	if po.Status == sharedDomain.PurchaseOrderStatusApproved {
		return nil
	}
	po.Status = sharedDomain.PurchaseOrderStatusApproved
	po.ApprovedAt = time.Now().UTC()
	if err := u.repo.UpdatePurchaseOrder(ctx, po); err != nil {
		return err
	}
	if u.nats != nil {
		payload, marshalErr := json.Marshal(map[string]any{
			"purchase_order_id": purchaseOrderID,
			"supplier_id":      po.SupplierID,
			"status":           po.Status,
		})
		if marshalErr == nil {
			_, _ = u.nats.JS.Publish(ctx, "PURCHASE_ORDERS.approved", payload)
		}
	}
	return nil
}

func (u *procurementUseCase) ReceiveGoods(ctx context.Context, purchaseOrderID string, quantity int) (*domain.GoodsReceipt, error) {
	if purchaseOrderID == "" {
		return nil, fmt.Errorf("%w: purchase order id is required", sharedDomain.ErrInvalidInput)
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("%w: quantity must be positive", sharedDomain.ErrInvalidInput)
	}
	po, err := u.repo.GetPurchaseOrder(ctx, purchaseOrderID)
	if err != nil {
		return nil, err
	}
	if po.Status != sharedDomain.PurchaseOrderStatusApproved {
		return nil, fmt.Errorf("%w: purchase order %s is not approved", sharedDomain.ErrConflict, purchaseOrderID)
	}

	receipt := &domain.GoodsReceipt{
		ID:              fmt.Sprintf("gr_%d", time.Now().UnixNano()),
		PurchaseOrderID: purchaseOrderID,
		Quantity:        quantity,
		Status:          sharedDomain.PurchaseOrderStatusReceived,
		ReceivedAt:      time.Now().UTC(),
	}
	if err := u.repo.CreateGoodsReceipt(ctx, receipt); err != nil {
		return nil, err
	}
	po.Status = sharedDomain.PurchaseOrderStatusReceived
	if err := u.repo.UpdatePurchaseOrder(ctx, po); err != nil {
		return nil, err
	}
	if u.nats != nil {
		payload, marshalErr := json.Marshal(map[string]any{
			"purchase_order_id": purchaseOrderID,
			"supplier_id":      po.SupplierID,
			"quantity":         quantity,
			"status":           po.Status,
		})
		if marshalErr == nil {
			_, _ = u.nats.JS.Publish(ctx, "PURCHASE_ORDERS.received", payload)
		}
	}
	return receipt, nil
}

func (u *procurementUseCase) CreateInvoice(ctx context.Context, purchaseOrderID string, amount float64) (*domain.Invoice, error) {
	if purchaseOrderID == "" {
		return nil, fmt.Errorf("%w: purchase order id is required", sharedDomain.ErrInvalidInput)
	}
	po, err := u.repo.GetPurchaseOrder(ctx, purchaseOrderID)
	if err != nil {
		return nil, err
	}
	if po.Status != sharedDomain.PurchaseOrderStatusReceived {
		return nil, fmt.Errorf("%w: purchase order %s must be received before invoicing", sharedDomain.ErrConflict, purchaseOrderID)
	}
	if amount <= 0 {
		amount = po.Amount
	}
	invoice := &domain.Invoice{
		ID:              fmt.Sprintf("inv_%d", time.Now().UnixNano()),
		PurchaseOrderID: purchaseOrderID,
		SupplierID:      po.SupplierID,
		Amount:          amount,
		Status:          sharedDomain.PurchaseOrderStatusMatched,
		CreatedAt:       time.Now().UTC(),
		MatchedAt:       time.Now().UTC(),
	}
	if err := u.repo.CreateInvoice(ctx, invoice); err != nil {
		return nil, err
	}
	po.Status = sharedDomain.PurchaseOrderStatusMatched
	if err := u.repo.UpdatePurchaseOrder(ctx, po); err != nil {
		return nil, err
	}
	if _, err := u.PostAccountsPayable(ctx, purchaseOrderID); err != nil {
		return nil, err
	}
	if u.nats != nil {
		payload, marshalErr := json.Marshal(map[string]any{
			"purchase_order_id": purchaseOrderID,
			"supplier_id":      po.SupplierID,
			"invoice_id":       invoice.ID,
			"amount":           invoice.Amount,
			"status":           invoice.Status,
		})
		if marshalErr == nil {
			_, _ = u.nats.JS.Publish(ctx, "INVOICES.matched", payload)
		}
	}
	return invoice, nil
}

func (u *procurementUseCase) PostAccountsPayable(ctx context.Context, purchaseOrderID string) (*domain.AccountsPayableEntry, error) {
	if purchaseOrderID == "" {
		return nil, fmt.Errorf("%w: purchase order id is required", sharedDomain.ErrInvalidInput)
	}
	po, err := u.repo.GetPurchaseOrder(ctx, purchaseOrderID)
	if err != nil {
		return nil, err
	}
	if po.Status != sharedDomain.PurchaseOrderStatusMatched {
		return nil, fmt.Errorf("%w: purchase order %s must be matched before AP posting", sharedDomain.ErrConflict, purchaseOrderID)
	}
	if existing, err := u.repo.GetAccountsPayableByPurchaseOrder(ctx, purchaseOrderID); err == nil && existing != nil {
		return existing, nil
	}

	entry := &domain.AccountsPayableEntry{
		ID:              fmt.Sprintf("ap_%d", time.Now().UnixNano()),
		SupplierID:      po.SupplierID,
		PurchaseOrderID: purchaseOrderID,
		Amount:          po.Amount,
		Status:          sharedDomain.AccountsPayableStatusPosted,
		PostedAt:        time.Now().UTC(),
	}
	if err := u.repo.CreateAccountsPayableEntry(ctx, entry); err != nil {
		return nil, err
	}
	if u.nats != nil {
		payload, marshalErr := json.Marshal(map[string]any{
			"purchase_order_id": purchaseOrderID,
			"supplier_id":      po.SupplierID,
			"amount":           entry.Amount,
			"status":           entry.Status,
		})
		if marshalErr == nil {
			_, _ = u.nats.JS.Publish(ctx, "AP.posted", payload)
		}
	}
	return entry, nil
}

func (u *procurementUseCase) GetAccountsPayableBySupplier(ctx context.Context, supplierID string) ([]*domain.AccountsPayableEntry, error) {
	if supplierID == "" {
		return nil, fmt.Errorf("%w: supplier id is required", sharedDomain.ErrInvalidInput)
	}
	return u.repo.GetAccountsPayableBySupplier(ctx, supplierID)
}
