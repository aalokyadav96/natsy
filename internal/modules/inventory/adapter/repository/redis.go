package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"nae/internal/modules/inventory/domain"

	"github.com/redis/go-redis/v9"
)

type redisInventoryRepository struct {
	client *redis.Client
}

func NewRedisInventoryRepository(client *redis.Client) domain.InventoryRepository {
	return &redisInventoryRepository{client: client}
}

func (r *redisInventoryRepository) key(productID string) string {
	return fmt.Sprintf("inventory:%s", productID)
}

func (r *redisInventoryRepository) GetByProductID(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	data, err := r.client.Get(ctx, r.key(productID)).Result()
	if err == redis.Nil {
		return &domain.InventoryItem{ID: productID, ProductID: productID, Available: 0}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get inventory: %w", err)
	}

	var item domain.InventoryItem
	if err := json.Unmarshal([]byte(data), &item); err != nil {
		return nil, fmt.Errorf("decode inventory: %w", err)
	}
	return &item, nil
}

func (r *redisInventoryRepository) Set(ctx context.Context, item *domain.InventoryItem) error {
	if item == nil {
		return fmt.Errorf("inventory item is nil")
	}
	payload, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(item.ProductID), payload, 0).Err()
}

func (r *redisInventoryRepository) Reserve(ctx context.Context, productID string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	item, err := r.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}
	if item.Available < quantity {
		return fmt.Errorf("insufficient inventory for product %s", productID)
	}
	item.Available -= quantity
	return r.Set(ctx, item)
}

func (r *redisInventoryRepository) Release(ctx context.Context, productID string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	item, err := r.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}
	item.Available += quantity
	return r.Set(ctx, item)
}
