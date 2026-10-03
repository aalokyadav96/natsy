package domain

import (
	"context"
	"fmt"
	"time"
)

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

const (
	StatusCreated   = "created"
	StatusValidated = "validated"
	StatusReserved  = "reserved"
	StatusPaid      = "paid"
	StatusFulfilled = "fulfilled"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewOrderID(userID string) string {
	if userID == "" {
		return fmt.Sprintf("ord_%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("ord_%s_%d", userID, time.Now().UnixNano())
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	ListByUser(ctx context.Context, userID string) ([]*Order, error)
}
