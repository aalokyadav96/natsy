package ap

import (
	"nae/internal/app"
	apRepo "nae/internal/modules/ap/adapter/repository"
	apHTTP "nae/internal/modules/ap/port/http"
	apUC "nae/internal/modules/ap/usecase"
	"nae/internal/shared/infra/mq"
)

type Module struct {
	handler *apHTTP.APHandler
}

func NewModule(nats ...*mq.NATSClient) *Module {
	repo := apRepo.NewInMemoryAPRepository()
	var client *mq.NATSClient
	if len(nats) > 0 {
		client = nats[0]
	}
	uc := apUC.NewAPUseCase(repo, client)
	return &Module{handler: apHTTP.NewAPHandler(uc)}
}

func (m *Module) Name() string { return "ap" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
