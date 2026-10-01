package infra

import (
	"context"
	"fmt"

	"nae/internal/shared/infra/cache"
	"nae/internal/shared/infra/db"
	"nae/internal/shared/infra/mq"
)

type Clients struct {
	Postgres *db.PostgresClient
	Redis    *cache.RedisClient
	NATS     *mq.NATSClient
}

func NewClients(ctx context.Context, cfg Config) (*Clients, error) {
	pgClient, err := db.NewPostgresClient(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	redisClient, err := cache.NewRedisClient(cfg.RedisURL)
	if err != nil {
		pgClient.Close()
		return nil, fmt.Errorf("redis: %w", err)
	}

	natsClient, err := mq.NewNATSClient(cfg.NATSURL)
	if err != nil {
		pgClient.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("nats: %w", err)
	}

	if err := natsClient.EnsureStream(ctx, "ORDERS", []string{"ORDERS.*"}); err != nil {
		natsClient.Close()
		pgClient.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("nats stream: %w", err)
	}

	return &Clients{
		Postgres: pgClient,
		Redis:    redisClient,
		NATS:     natsClient,
	}, nil
}

func (c *Clients) Close() {
	if c == nil {
		return
	}
	if c.NATS != nil {
		c.NATS.Close()
	}
	if c.Redis != nil {
		_ = c.Redis.Close()
	}
	if c.Postgres != nil {
		c.Postgres.Close()
	}
}
