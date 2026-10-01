package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nae/internal/modules/order/domain"
	"nae/internal/shared/infra/mq"
)

type OrderUseCase interface {
	Create(ctx context.Context, userID string, items []domain.OrderItem) (*domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
}

type orderUseCase struct {
	repo domain.OrderRepository
	nats *mq.NATSClient
}

func NewOrderUseCase(repo domain.OrderRepository, nats *mq.NATSClient) OrderUseCase {
	return &orderUseCase{repo: repo, nats: nats}
}

func (u *orderUseCase) Create(ctx context.Context, userID string, items []domain.OrderItem) (*domain.Order, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("order must contain at least one item")
	}

	var total float64
	for _, item := range items {
		if item.ProductID == "" {
			return nil, fmt.Errorf("product_id is required")
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("quantity must be positive for product %s", item.ProductID)
		}
		total += item.UnitPrice * float64(item.Quantity)
	}

	order := &domain.Order{
		ID:        fmt.Sprintf("ord_%s_%d", userID, time.Now().UnixNano()),
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    "created",
		CreatedAt: time.Now().UTC(),
	}

	if err := u.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	if u.nats != nil {
		payload, err := json.Marshal(map[string]any{
			"order_id": order.ID,
			"user_id":  order.UserID,
			"total":    order.Total,
			"items":    order.Items,
		})
		if err == nil {
			_, _ = u.nats.JS.Publish(ctx, "ORDERS.created", payload)
		}
	}

	return order, nil
}

func (u *orderUseCase) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order id is required")
	}
	return u.repo.GetByID(ctx, id)
}
