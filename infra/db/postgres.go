package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"nae/internal/domain"
)

type PostgresRepo struct {
	Pool *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, connString string) (*PostgresRepo, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create db pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	repo := &PostgresRepo{Pool: pool}
	if err := repo.initTables(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return repo, nil
}

func (r *PostgresRepo) initTables(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(64) PRIMARY KEY,
		user_id VARCHAR(64) NOT NULL,
		amount NUMERIC(12, 2) NOT NULL,
		currency VARCHAR(3) NOT NULL,
		status VARCHAR(32) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := r.Pool.Exec(ctx, query)
	return err
}

func (r *PostgresRepo) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, user_id, amount, currency, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.Pool.Exec(ctx, query, order.ID, order.UserID, order.Amount, order.Currency, order.Status, order.CreatedAt)
	return err
}

func (r *PostgresRepo) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT id, user_id, amount, currency, status, created_at FROM orders WHERE id = $1`
	row := r.Pool.QueryRow(ctx, query, id)

	var o domain.Order
	err := row.Scan(&o.ID, &o.UserID, &o.Amount, &o.Currency, &o.Status, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}
