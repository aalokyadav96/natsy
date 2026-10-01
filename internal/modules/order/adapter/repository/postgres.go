package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"nae/internal/modules/order/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) domain.OrderRepository {
	return &postgresOrderRepository{pool: pool}
}

func (r *postgresOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return fmt.Errorf("marshal items: %w", err)
	}
	query := `
		INSERT INTO orders (id, user_id, items, total, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = r.pool.Exec(ctx, query, order.ID, order.UserID, itemsJSON, order.Total, order.Status, order.CreatedAt)
	return err
}

func (r *postgresOrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, user_id, items, total, status, created_at FROM orders WHERE id = $1`, id)
	var order domain.Order
	var itemsJSON []byte
	if err := row.Scan(&order.ID, &order.UserID, &itemsJSON, &order.Total, &order.Status, &order.CreatedAt); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, fmt.Errorf("unmarshal items: %w", err)
	}
	return &order, nil
}

func (r *postgresOrderRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Order, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, user_id, items, total, status, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*domain.Order{}
	for rows.Next() {
		var order domain.Order
		var itemsJSON []byte
		if err := rows.Scan(&order.ID, &order.UserID, &itemsJSON, &order.Total, &order.Status, &order.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, fmt.Errorf("unmarshal items: %w", err)
		}
		orders = append(orders, &order)
	}
	return orders, nil
}
