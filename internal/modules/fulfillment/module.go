package fulfillment

import (
	"nae/internal/app"
	fulfillmentRepo "nae/internal/modules/fulfillment/adapter/repository"
	fulfillmentEvent "nae/internal/modules/fulfillment/port/event"
	fulfillmentHTTP "nae/internal/modules/fulfillment/port/http"
	fulfillmentUC "nae/internal/modules/fulfillment/usecase"
	"nae/internal/shared/infra/mq"
)

type Module struct {
	handler  *fulfillmentHTTP.FulfillmentHandler
	consumer *fulfillmentEvent.FulfillmentConsumer
}

func NewModule(nats *mq.NATSClient) *Module {
	repo := fulfillmentRepo.NewInMemoryFulfillmentRepository()
	uc := fulfillmentUC.NewFulfillmentUseCase(repo)
	return &Module{handler: fulfillmentHTTP.NewFulfillmentHandler(uc), consumer: fulfillmentEvent.NewFulfillmentConsumer(uc, nats)}
}

func (m *Module) Name() string { return "fulfillment" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	if m.consumer != nil {
		c.RegisterStartup(m.consumer.Start)
	}
	return nil
}
