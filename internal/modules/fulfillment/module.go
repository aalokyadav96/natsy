package fulfillment

import (
	"nae/internal/app"
	fulfillmentHTTP "nae/internal/modules/fulfillment/port/http"
	fulfillmentUC "nae/internal/modules/fulfillment/usecase"
)

type Module struct {
	handler *fulfillmentHTTP.FulfillmentHandler
}

func NewModule() *Module {
	return &Module{handler: fulfillmentHTTP.NewFulfillmentHandler(fulfillmentUC.NewFulfillmentUseCase())}
}

func (m *Module) Name() string { return "fulfillment" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
