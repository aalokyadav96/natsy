package usecase

import (
	"context"

	"nae/internal/modules/fulfillment/domain"
)

type FulfillmentUseCase interface {
	Create(ctx context.Context, orderID string) (*domain.Fulfillment, error)
}

type fulfillmentUseCase struct{}

func NewFulfillmentUseCase() FulfillmentUseCase { return &fulfillmentUseCase{} }

func (u *fulfillmentUseCase) Create(ctx context.Context, orderID string) (*domain.Fulfillment, error) {
	return &domain.Fulfillment{ID: "ful_1", OrderID: orderID, Status: "pending"}, nil
}
