package inventory

import (
	"nae/internal/app"
	inventoryRepo "nae/internal/modules/inventory/adapter/repository"
	inventoryHTTP "nae/internal/modules/inventory/port/http"
	inventoryUC "nae/internal/modules/inventory/usecase"

	"github.com/redis/go-redis/v9"
)

type Module struct {
	handler *inventoryHTTP.InventoryHandler
}

func NewModule(client *redis.Client) *Module {
	repo := inventoryRepo.NewRedisInventoryRepository(client)
	uc := inventoryUC.NewInventoryUseCase(repo)
	return &Module{handler: inventoryHTTP.NewInventoryHandler(uc)}
}

func (m *Module) Name() string { return "inventory" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
