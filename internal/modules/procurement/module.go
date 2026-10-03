package procurement

import (
	"nae/internal/app"
	repository "nae/internal/modules/procurement/adapter/repository"
	procurementHTTP "nae/internal/modules/procurement/port/http"
	procurementUC "nae/internal/modules/procurement/usecase"
)

type Module struct {
	handler *procurementHTTP.ProcurementHandler
}

func NewModule() *Module {
	repo := repository.NewInMemoryProcurementRepository()
	uc := procurementUC.NewProcurementUseCase(repo)
	return &Module{handler: procurementHTTP.NewProcurementHandler(uc)}
}

func (m *Module) Name() string { return "procurement" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
