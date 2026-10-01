package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nae/internal/modules/payment/domain"
	"nae/internal/shared/infra/mq"
)

type PaymentUseCase interface {
	Create(ctx context.Context, orderID string, amount float64) (*domain.Payment, error)
}

type paymentUseCase struct {
	repo domain.PaymentRepository
	nats *mq.NATSClient
}

func NewPaymentUseCase(repo domain.PaymentRepository, nats *mq.NATSClient) PaymentUseCase {
	return &paymentUseCase{repo: repo, nats: nats}
}

func (u *paymentUseCase) Create(ctx context.Context, orderID string, amount float64) (*domain.Payment, error) {
	if orderID == "" {
		return nil, fmt.Errorf("order_id is required")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	payment := &domain.Payment{
		ID:        fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		OrderID:   orderID,
		Amount:    amount,
		Status:    "authorized",
		CreatedAt: time.Now().UTC(),
	}

	if err := u.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	if u.nats != nil {
		payload, err := json.Marshal(map[string]any{
			"payment_id": payment.ID,
			"order_id":   payment.OrderID,
			"amount":     payment.Amount,
			"status":     payment.Status,
		})
		if err == nil {
			_, _ = u.nats.JS.Publish(ctx, "PAYMENTS.created", payload)
		}
	}

	return payment, nil
}
