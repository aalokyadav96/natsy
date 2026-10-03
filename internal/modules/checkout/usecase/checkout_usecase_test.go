package usecase

import (
	"context"
	"testing"

	cartDomain "nae/internal/modules/cart/domain"
	inventoryDomain "nae/internal/modules/inventory/domain"
)

type stubCartRepository struct {
	cart *cartDomain.Cart
}

func (s *stubCartRepository) GetCart(ctx context.Context, userID string) (*cartDomain.Cart, error) {
	return s.cart, nil
}

func (s *stubCartRepository) AddItem(ctx context.Context, userID string, item cartDomain.CartItem) error {
	return nil
}

func (s *stubCartRepository) ClearCart(ctx context.Context, userID string) error {
	return nil
}

type stubInventoryRepository struct {
	reserved []string
}

func (s *stubInventoryRepository) GetByProductID(ctx context.Context, productID string) (*inventoryDomain.InventoryItem, error) {
	return &inventoryDomain.InventoryItem{ID: productID, ProductID: productID, Available: 10}, nil
}

func (s *stubInventoryRepository) Set(ctx context.Context, item *inventoryDomain.InventoryItem) error {
	return nil
}

func (s *stubInventoryRepository) Reserve(ctx context.Context, productID string, quantity int) error {
	s.reserved = append(s.reserved, productID)
	return nil
}

func (s *stubInventoryRepository) Release(ctx context.Context, productID string, quantity int) error {
	return nil
}

func TestCheckoutReservesInventoryBeforePublishingOrder(t *testing.T) {
	cartRepo := &stubCartRepository{cart: &cartDomain.Cart{Items: []cartDomain.CartItem{{ProductID: "sku-1", Quantity: 2, UnitPrice: 19.99}}}}
	inventoryRepo := &stubInventoryRepository{}

	uc := NewCheckoutUseCase(cartRepo, inventoryRepo, nil)
	orderID, err := uc.ProcessCheckout(context.Background(), "user-123")
	if err != nil {
		t.Fatalf("ProcessCheckout returned error: %v", err)
	}
	if orderID == "" {
		t.Fatal("expected order ID to be returned")
	}
	if len(inventoryRepo.reserved) != 1 || inventoryRepo.reserved[0] != "sku-1" {
		t.Fatalf("expected inventory to be reserved for sku-1, got %#v", inventoryRepo.reserved)
	}
}
