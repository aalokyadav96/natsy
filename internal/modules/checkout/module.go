package checkout

import (
	"time"

	"nae/internal/app"
	cartRepo "nae/internal/modules/cart/adapter/repository"
	checkoutHTTP "nae/internal/modules/checkout/port/http"
	checkoutUC "nae/internal/modules/checkout/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/redis/go-redis/v9"
)

type Module struct {
	handler *checkoutHTTP.CheckoutHandler
}

func NewModule(client *redis.Client, nats *mq.NATSClient) *Module {
	repo := cartRepo.NewRedisCartRepository(client, 24*time.Hour)
	uc := checkoutUC.NewCheckoutUseCase(repo, nats)
	return &Module{handler: checkoutHTTP.NewCheckoutHandler(uc)}
}

func (m *Module) Name() string { return "checkout" }

func (m *Module) Register(c *app.Container) error {
	c.RegisterRoutes(m.handler.RegisterRoutes)
	return nil
}
