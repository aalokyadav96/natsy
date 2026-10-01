package repository

import (
	"context"
	"fmt"

	"nae/internal/modules/payment/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresPaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPaymentRepository(pool *pgxpool.Pool) domain.PaymentRepository {
	return &postgresPaymentRepository{pool: pool}
}

func (r *postgresPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	if payment == nil {
		return fmt.Errorf("payment is nil")
	}
	query := `
		INSERT INTO payments (id, order_id, status, amount, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, payment.ID, payment.OrderID, payment.Status, payment.Amount, payment.CreatedAt)
	return err
}

func (r *postgresPaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, order_id, status, amount, created_at FROM payments WHERE order_id = $1 LIMIT 1`, orderID)
	var payment domain.Payment
	if err := row.Scan(&payment.ID, &payment.OrderID, &payment.Status, &payment.Amount, &payment.CreatedAt); err != nil {
		return nil, err
	}
	return &payment, nil
}
