package usecase

import (
	"context"
	"fmt"
	"time"

	"nae/internal/modules/fulfillment/domain"
	sharedDomain "nae/internal/shared/domain"
)

type FulfillmentUseCase interface {
	Create(ctx context.Context, orderID string) (*domain.Fulfillment, error)
}

type fulfillmentUseCase struct {
	repo domain.FulfillmentRepository
}

func NewFulfillmentUseCase(repo domain.FulfillmentRepository) FulfillmentUseCase {
	return &fulfillmentUseCase{repo: repo}
}

func (u *fulfillmentUseCase) Create(ctx context.Context, orderID string) (*domain.Fulfillment, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: order id is required", sharedDomain.ErrInvalidInput)
	}
	if u.repo == nil {
		return nil, fmt.Errorf("%w: fulfillment repository is not configured", sharedDomain.ErrInternalError)
	}

	fulfillment := &domain.Fulfillment{
		ID:        fmt.Sprintf("ful_%d", time.Now().UnixNano()),
		OrderID:   orderID,
		Status:    "pending",
		CreatedAt: time.Now().UTC(),
	}
	if err := u.repo.Create(ctx, fulfillment); err != nil {
		return nil, err
	}
	return fulfillment, nil
}
