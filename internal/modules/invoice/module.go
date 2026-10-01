package invoice

import (
	"nae/internal/app"
	invoiceEvent "nae/internal/modules/invoice/port/event"
	"nae/internal/shared/infra/mq"
)

type Module struct {
	consumer *invoiceEvent.InvoiceConsumer
}

func NewModule(nats *mq.NATSClient) *Module {
	return &Module{consumer: invoiceEvent.NewInvoiceConsumer(nats)}
}

func (m *Module) Name() string { return "invoice" }

func (m *Module) Register(c *app.Container) error {
	c.RegisterStartup(m.consumer.Start)
	return nil
}
