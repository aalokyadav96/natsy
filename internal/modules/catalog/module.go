package catalog

import (
	"nae/internal/app"
	catalogRepo "nae/internal/modules/catalog/adapter/repository"
	catalogHTTP "nae/internal/modules/catalog/port/http"
	catalogUC "nae/internal/modules/catalog/usecase"
)

type Module struct {
	handler *catalogHTTP.CatalogHandler
}

func NewModule() *Module {
	repo := catalogRepo.NewInMemoryCatalogRepository()
	return &Module{handler: catalogHTTP.NewCatalogHandler(catalogUC.NewCatalogUseCase(repo))}
}

func (m *Module) Name() string { return "catalog" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
