package repository

import (
	"context"
	"fmt"
	"sync"

	"nae/internal/modules/fulfillment/domain"
	sharedDomain "nae/internal/shared/domain"
)

type inMemoryFulfillmentRepository struct {
	mu          sync.RWMutex
	fulfillments map[string]*domain.Fulfillment
}

func NewInMemoryFulfillmentRepository() domain.FulfillmentRepository {
	return &inMemoryFulfillmentRepository{fulfillments: make(map[string]*domain.Fulfillment)}
}

func (r *inMemoryFulfillmentRepository) Create(ctx context.Context, fulfillment *domain.Fulfillment) error {
	if fulfillment == nil {
		return fmt.Errorf("%w: fulfillment is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.fulfillments == nil {
		r.fulfillments = make(map[string]*domain.Fulfillment)
	}
	r.fulfillments[fulfillment.OrderID] = fulfillment
	return nil
}

func (r *inMemoryFulfillmentRepository) GetByOrderID(ctx context.Context, orderID string) (*domain.Fulfillment, error) {
	if orderID == "" {
		return nil, fmt.Errorf("%w: order id is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()
	fulfillment, ok := r.fulfillments[orderID]
	if !ok {
		return nil, fmt.Errorf("%w: fulfillment for order %s was not found", sharedDomain.ErrNotFound, orderID)
	}
	clone := *fulfillment
	return &clone, nil
}
