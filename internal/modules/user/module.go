package user

import (
	"nae/internal/app"
	userRepo "nae/internal/modules/user/adapter/repository"
	userHTTP "nae/internal/modules/user/port/http"
	userUC "nae/internal/modules/user/usecase"
)

type Module struct {
	handler *userHTTP.UserHandler
}

func NewModule() *Module {
	repo := userRepo.NewInMemoryUserRepository()
	return &Module{handler: userHTTP.NewUserHandler(userUC.NewUserUseCase(repo))}
}

func (m *Module) Name() string { return "user" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
