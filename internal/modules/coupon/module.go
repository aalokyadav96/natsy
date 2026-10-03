package coupon

import (
	"nae/internal/app"
	couponRepo "nae/internal/modules/coupon/adapter/repository"
	couponHTTP "nae/internal/modules/coupon/port/http"
	couponUC "nae/internal/modules/coupon/usecase"
)

type Module struct {
	handler *couponHTTP.CouponHandler
}

func NewModule() *Module {
	repo := couponRepo.NewInMemoryCouponRepository()
	return &Module{handler: couponHTTP.NewCouponHandler(couponUC.NewCouponUseCase(repo))}
}

func (m *Module) Name() string { return "coupon" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
