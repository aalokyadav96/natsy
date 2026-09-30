package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nae/internal/modules/cart/domain"

	"github.com/redis/go-redis/v9"
)

type redisCartRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCartRepository(rdb *redis.Client, ttl time.Duration) domain.CartRepository {
	return &redisCartRepository{
		client: rdb,
		ttl:    ttl,
	}
}

func (r *redisCartRepository) getCartKey(userID string) string {
	return fmt.Sprintf("cart:%s", userID)
}

func (r *redisCartRepository) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	key := r.getCartKey(userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Return empty cart on miss
		return &domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve cart from redis: %w", err)
	}

	var cart domain.Cart
	if err := json.Unmarshal([]byte(val), &cart); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cart data: %w", err)
	}

	return &cart, nil
}

func (r *redisCartRepository) AddItem(ctx context.Context, userID string, item domain.CartItem) error {
	cart, err := r.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// Check if product already exists in cart; update quantity if found
	itemFound := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			cart.Items[i].UnitPrice = item.UnitPrice // update to latest price
			itemFound = true
			break
		}
	}

	if !itemFound {
		cart.Items = append(cart.Items, item)
	}

	data, err := json.Marshal(cart)
	if err != nil {
		return fmt.Errorf("failed to marshal cart payload: %w", err)
	}

	key := r.getCartKey(userID)
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *redisCartRepository) ClearCart(ctx context.Context, userID string) error {
	key := r.getCartKey(userID)
	return r.client.Del(ctx, key).Err()
}
