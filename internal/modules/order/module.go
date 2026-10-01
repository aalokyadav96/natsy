package order

import (
	"nae/internal/app"
	orderRepo "nae/internal/modules/order/adapter/repository"
	orderHTTP "nae/internal/modules/order/port/http"
	orderUC "nae/internal/modules/order/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *orderHTTP.OrderHandler
}

func NewModule(pool *pgxpool.Pool, nats *mq.NATSClient) *Module {
	repo := orderRepo.NewPostgresOrderRepository(pool)
	uc := orderUC.NewOrderUseCase(repo, nats)
	return &Module{handler: orderHTTP.NewOrderHandler(uc)}
}

func (m *Module) Name() string { return "order" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
