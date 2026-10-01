package usecase

import (
	"context"

	"nae/internal/modules/catalog/domain"
)

type CatalogUseCase interface {
	ListItems(ctx context.Context) ([]*domain.CatalogItem, error)
	GetItem(ctx context.Context, id string) (*domain.CatalogItem, error)
}

type catalogUseCase struct{}

func NewCatalogUseCase() CatalogUseCase { return &catalogUseCase{} }

func (u *catalogUseCase) ListItems(ctx context.Context) ([]*domain.CatalogItem, error) {
	return []*domain.CatalogItem{}, nil
}

func (u *catalogUseCase) GetItem(ctx context.Context, id string) (*domain.CatalogItem, error) {
	return &domain.CatalogItem{ID: id, Name: "placeholder"}, nil
}
