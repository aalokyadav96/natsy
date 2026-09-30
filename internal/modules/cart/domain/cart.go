package domain

import (
	"context"
)

type CartItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartRepository interface {
	GetCart(ctx context.Context, userID string) (*Cart, error)
	AddItem(ctx context.Context, userID string, item CartItem) error
	ClearCart(ctx context.Context, userID string) error
}
