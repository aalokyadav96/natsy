package usecase

import (
	"context"
	"fmt"

	"nae/internal/modules/catalog/domain"
	sharedDomain "nae/internal/shared/domain"
)

type CatalogUseCase interface {
	ListItems(ctx context.Context) ([]*domain.CatalogItem, error)
	GetItem(ctx context.Context, id string) (*domain.CatalogItem, error)
}

type catalogUseCase struct {
	repo domain.CatalogRepository
}

func NewCatalogUseCase(repo domain.CatalogRepository) CatalogUseCase {
	return &catalogUseCase{repo: repo}
}

func (u *catalogUseCase) ListItems(ctx context.Context) ([]*domain.CatalogItem, error) {
	if u.repo == nil {
		return nil, fmt.Errorf("%w: catalog repository is not configured", sharedDomain.ErrInternalError)
	}
	return u.repo.List(ctx)
}

func (u *catalogUseCase) GetItem(ctx context.Context, id string) (*domain.CatalogItem, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: catalog item id is required", sharedDomain.ErrInvalidInput)
	}
	if u.repo == nil {
		return nil, fmt.Errorf("%w: catalog repository is not configured", sharedDomain.ErrInternalError)
	}
	return u.repo.GetByID(ctx, id)
}
