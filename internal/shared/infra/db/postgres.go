package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	Pool *pgxpool.Pool
}

func NewPostgresClient(ctx context.Context, connString string) (*PostgresClient, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create db pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresClient{Pool: pool}, nil
}

func (p *PostgresClient) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
