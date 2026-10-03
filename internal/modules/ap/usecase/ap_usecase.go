package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	apRepo "nae/internal/modules/ap/adapter/repository"
	apDomain "nae/internal/modules/ap/domain"
	sharedDomain "nae/internal/shared/domain"
	"nae/internal/shared/infra/mq"
)

type APUseCase interface {
	CreatePaymentBatch(ctx context.Context, supplierID, reference string, total float64) (*apDomain.PaymentBatch, error)
	ApprovePaymentBatch(ctx context.Context, batchID string) error
	ExecutePayment(ctx context.Context, batchID, bankAccountID string) (*apDomain.VendorPayment, error)
	ReconcilePayment(ctx context.Context, paymentID string) error
	GetOutstandingBySupplier(ctx context.Context, supplierID string) ([]*apDomain.VendorPayment, error)
}

type apUseCase struct {
	repo apDomain.APRepository
	nats *mq.NATSClient
}

func NewAPUseCase(repo apDomain.APRepository, nats ...*mq.NATSClient) APUseCase {
	if repo == nil {
		repo = apRepo.NewInMemoryAPRepository()
	}
	var client *mq.NATSClient
	if len(nats) > 0 {
		client = nats[0]
	}
	return &apUseCase{repo: repo, nats: client}
}

func (u *apUseCase) CreatePaymentBatch(ctx context.Context, supplierID, reference string, total float64) (*apDomain.PaymentBatch, error) {
	if supplierID == "" {
		return nil, fmt.Errorf("%w: supplier id is required", sharedDomain.ErrInvalidInput)
	}
	if total <= 0 {
		return nil, fmt.Errorf("%w: batch total must be greater than zero", sharedDomain.ErrInvalidInput)
	}

	batch := &apDomain.PaymentBatch{
		ID:         fmt.Sprintf("batch_%d", time.Now().UnixNano()),
		SupplierID: supplierID,
		Reference:  reference,
		Total:      total,
		Status:     sharedDomain.PaymentBatchStatusDraft,
		CreatedAt:  time.Now().UTC(),
	}
	if batch.Reference == "" {
		batch.Reference = fmt.Sprintf("AP-%d", time.Now().UnixNano())
	}
	if err := u.repo.CreatePaymentBatch(ctx, batch); err != nil {
		return nil, err
	}
	if u.nats != nil {
		payload, err := json.Marshal(map[string]any{"batch_id": batch.ID, "supplier_id": supplierID, "reference": batch.Reference, "total": total, "status": batch.Status})
		if err == nil {
			_, _ = u.nats.JS.Publish(ctx, "AP.payment_batch.created", payload)
		}
	}
	return batch, nil
}

func (u *apUseCase) ApprovePaymentBatch(ctx context.Context, batchID string) error {
	if batchID == "" {
		return fmt.Errorf("%w: batch id is required", sharedDomain.ErrInvalidInput)
	}
	batch, err := u.repo.GetPaymentBatch(ctx, batchID)
	if err != nil {
		return err
	}
	if batch.Status == sharedDomain.PaymentBatchStatusApproved {
		return nil
	}
	batch.Status = sharedDomain.PaymentBatchStatusApproved
	batch.ApprovedAt = time.Now().UTC()
	if err := u.repo.UpdatePaymentBatch(ctx, batch); err != nil {
		return err
	}
	if u.nats != nil {
		payload, marshalErr := json.Marshal(map[string]any{"batch_id": batch.ID, "supplier_id": batch.SupplierID, "status": batch.Status})
		if marshalErr == nil {
			_, _ = u.nats.JS.Publish(ctx, "AP.payment_batch.approved", payload)
		}
	}
	return nil
}

func (u *apUseCase) ExecutePayment(ctx context.Context, batchID, bankAccountID string) (*apDomain.VendorPayment, error) {
	if batchID == "" {
		return nil, fmt.Errorf("%w: batch id is required", sharedDomain.ErrInvalidInput)
	}
	if bankAccountID == "" {
		return nil, fmt.Errorf("%w: bank account id is required", sharedDomain.ErrInvalidInput)
	}
	batch, err := u.repo.GetPaymentBatch(ctx, batchID)
	if err != nil {
		return nil, err
	}
	if batch.Status != sharedDomain.PaymentBatchStatusApproved {
		return nil, fmt.Errorf("%w: payment batch %s is not approved", sharedDomain.ErrConflict, batchID)
	}

	payment := &apDomain.VendorPayment{
		ID:            fmt.Sprintf("vp_%d", time.Now().UnixNano()),
		BatchID:       batch.ID,
		SupplierID:    batch.SupplierID,
		BankAccountID: bankAccountID,
		Amount:        batch.Total,
		Status:        sharedDomain.PaymentStatusPaid,
		CreatedAt:     time.Now().UTC(),
		PaidAt:        time.Now().UTC(),
	}
	if err := u.repo.CreateVendorPayment(ctx, payment); err != nil {
		return nil, err
	}
	batch.Status = sharedDomain.PaymentBatchStatusPaid
	batch.PaidAt = time.Now().UTC()
	if err := u.repo.UpdatePaymentBatch(ctx, batch); err != nil {
		return nil, err
	}
	if u.nats != nil {
		payload, marshalErr := json.Marshal(map[string]any{"payment_id": payment.ID, "batch_id": payment.BatchID, "supplier_id": payment.SupplierID, "amount": payment.Amount, "status": payment.Status})
		if marshalErr == nil {
			_, _ = u.nats.JS.Publish(ctx, "AP.payment.paid", payload)
		}
	}
	return payment, nil
}

func (u *apUseCase) ReconcilePayment(ctx context.Context, paymentID string) error {
	if paymentID == "" {
		return fmt.Errorf("%w: payment id is required", sharedDomain.ErrInvalidInput)
	}
	payment, err := u.repo.GetVendorPayment(ctx, paymentID)
	if err != nil {
		return err
	}
	if payment.Status == sharedDomain.PaymentStatusReconciled {
		return nil
	}
	payment.Status = sharedDomain.PaymentStatusReconciled
	payment.ReconciledAt = time.Now().UTC()
	if err := u.repo.UpdateVendorPayment(ctx, payment); err != nil {
		return err
	}

	payments, err := u.repo.ListPaymentsByBatch(ctx, payment.BatchID)
	if err != nil {
		return nil
	}
	allReconciled := true
	for _, item := range payments {
		if item.Status != sharedDomain.PaymentStatusReconciled {
			allReconciled = false
			break
		}
	}
	if allReconciled {
		batch, getErr := u.repo.GetPaymentBatch(ctx, payment.BatchID)
		if getErr == nil {
			batch.Status = sharedDomain.PaymentBatchStatusReconciled
			_ = u.repo.UpdatePaymentBatch(ctx, batch)
		}
	}
	if u.nats != nil {
		payload, marshalErr := json.Marshal(map[string]any{"payment_id": payment.ID, "batch_id": payment.BatchID, "supplier_id": payment.SupplierID, "status": payment.Status})
		if marshalErr == nil {
			_, _ = u.nats.JS.Publish(ctx, "AP.payment.reconciled", payload)
		}
	}
	return nil
}

func (u *apUseCase) GetOutstandingBySupplier(ctx context.Context, supplierID string) ([]*apDomain.VendorPayment, error) {
	if supplierID == "" {
		return nil, fmt.Errorf("%w: supplier id is required", sharedDomain.ErrInvalidInput)
	}
	payments, err := u.repo.ListPaymentsBySupplier(ctx, supplierID)
	if err != nil {
		return nil, err
	}
	var outstanding []*apDomain.VendorPayment
	for _, payment := range payments {
		if payment.Status != sharedDomain.PaymentStatusReconciled {
			outstanding = append(outstanding, payment)
		}
	}
	return outstanding, nil
}
