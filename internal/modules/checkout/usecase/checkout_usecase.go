package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	cartDomain "nae/internal/modules/cart/domain"
	inventoryDomain "nae/internal/modules/inventory/domain"
	orderDomain "nae/internal/modules/order/domain"
	"nae/internal/shared/infra/mq"
)

type CheckoutUseCase interface {
	ProcessCheckout(ctx context.Context, userID string) (string, error)
}

type checkoutUseCase struct {
	cartRepo      cartDomain.CartRepository
	inventoryRepo inventoryDomain.InventoryRepository
	nats          *mq.NATSClient
}

func NewCheckoutUseCase(cartRepo cartDomain.CartRepository, inventoryRepo inventoryDomain.InventoryRepository, nats *mq.NATSClient) CheckoutUseCase {
	return &checkoutUseCase{cartRepo: cartRepo, inventoryRepo: inventoryRepo, nats: nats}
}

type OrderCreatedEvent struct {
	OrderID string                `json:"order_id"`
	UserID  string                `json:"user_id"`
	Items   []cartDomain.CartItem `json:"items"`
	Total   float64               `json:"total"`
	Status  string                `json:"status"`
}

func (u *checkoutUseCase) ProcessCheckout(ctx context.Context, userID string) (string, error) {
	cart, err := u.cartRepo.GetCart(ctx, userID)
	if err != nil || len(cart.Items) == 0 {
		return "", fmt.Errorf("cart is empty")
	}

	var total float64
	for _, item := range cart.Items {
		total += item.UnitPrice * float64(item.Quantity)
		if u.inventoryRepo != nil {
			if err := u.inventoryRepo.Reserve(ctx, item.ProductID, item.Quantity); err != nil {
				return "", fmt.Errorf("reserve inventory for %s: %w", item.ProductID, err)
			}
		}
	}

	orderID := orderDomain.NewOrderID(userID)

	evt := OrderCreatedEvent{
		OrderID: orderID,
		UserID:  userID,
		Items:   cart.Items,
		Total:   total,
		Status:  orderDomain.StatusReserved,
	}

	payload, _ := json.Marshal(evt)
	if u.nats != nil {
		_, err = u.nats.JS.Publish(ctx, "ORDERS.created", payload)
		if err != nil {
			return "", fmt.Errorf("failed to emit event: %w", err)
		}
	}

	_ = u.cartRepo.ClearCart(ctx, userID)
	return orderID, nil
}
