package usecase

import (
	"context"
	"nae/internal/modules/cart/domain"
)

type CartUseCase interface {
	GetCart(ctx context.Context, userID string) (*domain.Cart, error)
	AddToCart(ctx context.Context, userID, productID string, quantity int, unitPrice float64) error
	ClearCart(ctx context.Context, userID string) error
}

type cartUseCase struct {
	repo domain.CartRepository
}

func NewCartUseCase(repo domain.CartRepository) CartUseCase {
	return &cartUseCase{repo: repo}
}

func (u *cartUseCase) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	return u.repo.GetCart(ctx, userID)
}

func (u *cartUseCase) AddToCart(ctx context.Context, userID, productID string, quantity int, unitPrice float64) error {
	item := domain.CartItem{
		ProductID: productID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
	return u.repo.AddItem(ctx, userID, item)
}

func (u *cartUseCase) ClearCart(ctx context.Context, userID string) error {
	return u.repo.ClearCart(ctx, userID)
}
