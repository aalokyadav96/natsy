package user

import (
	"nae/internal/app"
	userHTTP "nae/internal/modules/user/port/http"
	userUC "nae/internal/modules/user/usecase"
)

type Module struct {
	handler *userHTTP.UserHandler
}

func NewModule() *Module {
	return &Module{handler: userHTTP.NewUserHandler(userUC.NewUserUseCase(nil))}
}

func (m *Module) Name() string { return "user" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
