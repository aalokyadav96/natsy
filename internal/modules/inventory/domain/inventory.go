package domain

import "context"

type InventoryItem struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Available int    `json:"available"`
}

type InventoryRepository interface {
	GetByProductID(ctx context.Context, productID string) (*InventoryItem, error)
	Set(ctx context.Context, item *InventoryItem) error
	Reserve(ctx context.Context, productID string, quantity int) error
	Release(ctx context.Context, productID string, quantity int) error
}
