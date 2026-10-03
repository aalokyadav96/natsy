package inventory

import (
	"nae/internal/app"
	inventoryRepo "nae/internal/modules/inventory/adapter/repository"
	inventoryEvent "nae/internal/modules/inventory/port/event"
	inventoryHTTP "nae/internal/modules/inventory/port/http"
	inventoryUC "nae/internal/modules/inventory/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/redis/go-redis/v9"
)

type Module struct {
	handler  *inventoryHTTP.InventoryHandler
	consumer *inventoryEvent.InventoryConsumer
}

func NewModule(client *redis.Client, nats *mq.NATSClient) *Module {
	repo := inventoryRepo.NewRedisInventoryRepository(client)
	uc := inventoryUC.NewInventoryUseCase(repo)
	return &Module{handler: inventoryHTTP.NewInventoryHandler(uc), consumer: inventoryEvent.NewInventoryConsumer(uc, nats)}
}

func (m *Module) Name() string { return "inventory" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	if m.consumer != nil {
		c.RegisterStartup(m.consumer.Start)
	}
	return nil
}
