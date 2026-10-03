package domain

import (
	"context"
	"time"
)

type Coupon struct {
	ID        string
	Code      string
	Discount  float64
	ValidFrom time.Time
	ValidTo   time.Time
}

type CouponRepository interface {
	GetByCode(ctx context.Context, code string) (*Coupon, error)
}
