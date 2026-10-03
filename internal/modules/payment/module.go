package payment

import (
	"nae/internal/app"
	paymentRepo "nae/internal/modules/payment/adapter/repository"
	paymentEvent "nae/internal/modules/payment/port/event"
	paymentHTTP "nae/internal/modules/payment/port/http"
	paymentUC "nae/internal/modules/payment/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler  *paymentHTTP.PaymentHandler
	consumer *paymentEvent.PaymentConsumer
}

func NewModule(pool *pgxpool.Pool, nats *mq.NATSClient) *Module {
	repo := paymentRepo.NewPostgresPaymentRepository(pool)
	uc := paymentUC.NewPaymentUseCase(repo, nats)
	return &Module{handler: paymentHTTP.NewPaymentHandler(uc), consumer: paymentEvent.NewPaymentConsumer(uc, nats)}
}

func (m *Module) Name() string { return "payment" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	if m.consumer != nil {
		c.RegisterStartup(m.consumer.Start)
	}
	return nil
}
