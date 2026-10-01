package payment

import (
	"nae/internal/app"
	paymentRepo "nae/internal/modules/payment/adapter/repository"
	paymentHTTP "nae/internal/modules/payment/port/http"
	paymentUC "nae/internal/modules/payment/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *paymentHTTP.PaymentHandler
}

func NewModule(pool *pgxpool.Pool, nats *mq.NATSClient) *Module {
	repo := paymentRepo.NewPostgresPaymentRepository(pool)
	uc := paymentUC.NewPaymentUseCase(repo, nats)
	return &Module{handler: paymentHTTP.NewPaymentHandler(uc)}
}

func (m *Module) Name() string { return "payment" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
