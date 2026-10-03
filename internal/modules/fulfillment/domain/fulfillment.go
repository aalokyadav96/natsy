package domain

import (
	"context"
	"time"
)

type Fulfillment struct {
	ID        string
	OrderID   string
	Status    string
	CreatedAt time.Time
}

type FulfillmentRepository interface {
	Create(ctx context.Context, fulfillment *Fulfillment) error
	GetByOrderID(ctx context.Context, orderID string) (*Fulfillment, error)
}
