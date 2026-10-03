package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nae/internal/modules/coupon/domain"
	sharedDomain "nae/internal/shared/domain"
)

type inMemoryCouponRepository struct {
	mu      sync.RWMutex
	coupons map[string]*domain.Coupon
}

func NewInMemoryCouponRepository() domain.CouponRepository {
	now := time.Now().UTC()
	return &inMemoryCouponRepository{coupons: map[string]*domain.Coupon{
		"SAVE10":    {ID: "cp_1", Code: "SAVE10", Discount: 10, ValidFrom: now.Add(-time.Hour), ValidTo: now.Add(7 * 24 * time.Hour)},
		"WELCOME20": {ID: "cp_2", Code: "WELCOME20", Discount: 20, ValidFrom: now.Add(-time.Hour), ValidTo: now.Add(30 * 24 * time.Hour)},
	}}
}

func (r *inMemoryCouponRepository) GetByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	if code == "" {
		return nil, fmt.Errorf("%w: coupon code is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()
	coupon, ok := r.coupons[code]
	if !ok {
		return nil, fmt.Errorf("%w: coupon %s was not found", sharedDomain.ErrNotFound, code)
	}
	clone := *coupon
	return &clone, nil
}
