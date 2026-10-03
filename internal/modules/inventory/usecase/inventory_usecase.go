package usecase

import (
	"context"

	"nae/internal/modules/inventory/domain"
)

type InventoryUseCase interface {
	GetInventory(ctx context.Context, productID string) (*domain.InventoryItem, error)
	Reserve(ctx context.Context, productID string, quantity int) error
	Release(ctx context.Context, productID string, quantity int) error
}

type inventoryUseCase struct {
	repo domain.InventoryRepository
}

func NewInventoryUseCase(repo domain.InventoryRepository) InventoryUseCase {
	return &inventoryUseCase{repo: repo}
}

func (u *inventoryUseCase) GetInventory(ctx context.Context, productID string) (*domain.InventoryItem, error) {
	return u.repo.GetByProductID(ctx, productID)
}

func (u *inventoryUseCase) Reserve(ctx context.Context, productID string, quantity int) error {
	return u.repo.Reserve(ctx, productID, quantity)
}

func (u *inventoryUseCase) Release(ctx context.Context, productID string, quantity int) error {
	return u.repo.Release(ctx, productID, quantity)
}
