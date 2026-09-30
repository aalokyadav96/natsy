package repository

import (
	"context"

	"nae/internal/modules/listing/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresProductRepository(pool *pgxpool.Pool) domain.ProductRepository {
	return &postgresRepo{pool: pool}
}

func (r *postgresRepo) Create(ctx context.Context, p *domain.Product) error {
	query := `INSERT INTO products (id, title, description, price, stock, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, p.ID, p.Title, p.Description, p.Price, p.Stock, p.CreatedAt)
	return err
}

func (r *postgresRepo) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `SELECT id, title, description, price, stock, created_at FROM products WHERE id = $1`
	p := &domain.Product{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)
	if err != nil {
		return nil, domain.ErrProductNotFound
	}
	return p, nil
}

func (r *postgresRepo) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	query := `SELECT id, title, description, price, stock, created_at FROM products LIMIT $1 OFFSET $2`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
