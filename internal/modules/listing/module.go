package listing

import (
	"nae/internal/app"
	listingRepo "nae/internal/modules/listing/adapter/repository"
	listingHTTP "nae/internal/modules/listing/port/http"
	listingUC "nae/internal/modules/listing/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *listingHTTP.ListingHandler
}

func NewModule(pool *pgxpool.Pool) *Module {
	repo := listingRepo.NewPostgresProductRepository(pool)
	uc := listingUC.NewListingUseCase(repo)
	return &Module{handler: listingHTTP.NewListingHandler(uc)}
}

func (m *Module) Name() string { return "listing" }

func (m *Module) Register(c *app.Container) error {
	c.RegisterRoutes(m.handler.RegisterRoutes)
	return nil
}
