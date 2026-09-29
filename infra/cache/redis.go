package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"nae/internal/domain"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url: %w", err)
	}

	rdb := redis.NewClient(opts)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &RedisCache{client: rdb}, nil
}

func (r *RedisCache) SetOrder(ctx context.Context, order *domain.Order, ttl time.Duration) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("order:%s", order.ID)
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisCache) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	key := fmt.Sprintf("order:%s", id)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, err
	}
	return &order, nil
}
