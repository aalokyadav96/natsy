package audit

import (
	"nae/internal/app"
	auditHTTP "nae/internal/modules/audit/port/http"
	auditUC "nae/internal/modules/audit/usecase"
)

type Module struct {
	handler *auditHTTP.AuditHandler
}

func NewModule() *Module {
	return &Module{handler: auditHTTP.NewAuditHandler(auditUC.NewAuditUseCase())}
}

func (m *Module) Name() string { return "audit" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
