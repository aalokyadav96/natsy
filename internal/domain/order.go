package domain

import (
	"context"
	"time"
)

type Order struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Currency  string    `json:"currency" db:"currency"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, id string) (*Order, error)
}

type OrderCache interface {
	SetOrder(ctx context.Context, order *Order, ttl time.Duration) error
	GetOrder(ctx context.Context, id string) (*Order, error)
}
