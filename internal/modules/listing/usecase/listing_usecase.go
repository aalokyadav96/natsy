package usecase

import (
	"context"
	"time"

	"nae/internal/modules/listing/domain"
)

type ListingUseCase interface {
	CreateProduct(ctx context.Context, title, description string, price float64, stock int) (*domain.Product, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error)
}

type listingUseCase struct {
	repo domain.ProductRepository
}

func NewListingUseCase(repo domain.ProductRepository) ListingUseCase {
	return &listingUseCase{repo: repo}
}

func (u *listingUseCase) CreateProduct(ctx context.Context, title, description string, price float64, stock int) (*domain.Product, error) {
	p := &domain.Product{
		ID:          "prod_" + time.Now().Format("20060102150405"),
		Title:       title,
		Description: description,
		Price:       price,
		Stock:       stock,
		CreatedAt:   time.Now().UTC(),
	}

	if err := u.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (u *listingUseCase) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *listingUseCase) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	if limit <= 0 {
		limit = 10
	}
	return u.repo.List(ctx, limit, offset)
}
