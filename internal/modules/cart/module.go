package cart

import (
	"time"

	"nae/internal/app"
	cartRepo "nae/internal/modules/cart/adapter/repository"
	cartHTTP "nae/internal/modules/cart/port/http"
	cartUC "nae/internal/modules/cart/usecase"

	"github.com/redis/go-redis/v9"
)

type Module struct {
	handler *cartHTTP.CartHandler
}

func NewModule(client *redis.Client, ttl time.Duration) *Module {
	repo := cartRepo.NewRedisCartRepository(client, ttl)
	uc := cartUC.NewCartUseCase(repo)
	return &Module{handler: cartHTTP.NewCartHandler(uc)}
}

func (m *Module) Name() string { return "cart" }

func (m *Module) Register(c *app.Container) error {
	c.RegisterRoutes(m.handler.RegisterRoutes)
	return nil
}
