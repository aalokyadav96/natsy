package domain

import "time"

type Coupon struct {
	ID        string
	Code      string
	Discount  float64
	ValidFrom time.Time
	ValidTo   time.Time
}
