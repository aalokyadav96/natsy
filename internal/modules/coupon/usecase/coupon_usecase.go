package usecase

import (
	"context"

	"nae/internal/modules/coupon/domain"
)

type CouponUseCase interface {
	Validate(ctx context.Context, code string) (*domain.Coupon, error)
}

type couponUseCase struct{}

func NewCouponUseCase() CouponUseCase { return &couponUseCase{} }

func (u *couponUseCase) Validate(ctx context.Context, code string) (*domain.Coupon, error) {
	return &domain.Coupon{ID: "cp_1", Code: code, Discount: 10}, nil
}
