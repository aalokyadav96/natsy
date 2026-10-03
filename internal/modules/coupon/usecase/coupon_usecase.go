package usecase

import (
	"context"
	"fmt"
	"time"

	"nae/internal/modules/coupon/domain"
	sharedDomain "nae/internal/shared/domain"
)

type CouponUseCase interface {
	Validate(ctx context.Context, code string) (*domain.Coupon, error)
}

type couponUseCase struct {
	repo domain.CouponRepository
}

func NewCouponUseCase(repo domain.CouponRepository) CouponUseCase { return &couponUseCase{repo: repo} }

func (u *couponUseCase) Validate(ctx context.Context, code string) (*domain.Coupon, error) {
	if code == "" {
		return nil, fmt.Errorf("%w: coupon code is required", sharedDomain.ErrInvalidInput)
	}
	if u.repo == nil {
		return nil, fmt.Errorf("%w: coupon repository is not configured", sharedDomain.ErrInternalError)
	}

	coupon, err := u.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if coupon == nil {
		return nil, fmt.Errorf("%w: coupon %s was not found", sharedDomain.ErrNotFound, code)
	}
	if time.Now().Before(coupon.ValidFrom) || time.Now().After(coupon.ValidTo) {
		return nil, fmt.Errorf("%w: coupon %s is not valid at this time", sharedDomain.ErrInvalidInput, code)
	}
	return coupon, nil
}
