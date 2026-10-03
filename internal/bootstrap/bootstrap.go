package bootstrap

import (
	"fmt"
	"time"

	"nae/internal/app"
	"nae/internal/modules/audit"
	"nae/internal/modules/cart"
	"nae/internal/modules/catalog"
	"nae/internal/modules/checkout"
	"nae/internal/modules/coupon"
	"nae/internal/modules/fulfillment"
	"nae/internal/modules/inventory"
	"nae/internal/modules/invoice"
	"nae/internal/modules/listing"
	"nae/internal/modules/notification"
	"nae/internal/modules/order"
	"nae/internal/modules/payment"
	"nae/internal/modules/user"
	"nae/internal/shared/infra"
	"nae/internal/shared/infra/mq"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func BuildModules(clients *infra.Clients) ([]app.Module, error) {
	if clients == nil {
		return nil, fmt.Errorf("infra clients are required")
	}

	var pgPool *pgxpool.Pool
	if clients.Postgres != nil {
		pgPool = clients.Postgres.Pool
	}

	var redisClient *redis.Client
	if clients.Redis != nil {
		redisClient = clients.Redis.Client
	}

	var natsClient *mq.NATSClient
	if clients.NATS != nil {
		natsClient = clients.NATS
	}

	return []app.Module{
		user.NewModule(),
		catalog.NewModule(),
		listing.NewModule(pgPool),
		inventory.NewModule(redisClient),
		cart.NewModule(redisClient, 24*time.Hour),
		coupon.NewModule(),
		order.NewModule(pgPool, natsClient),
		payment.NewModule(pgPool, natsClient),
		fulfillment.NewModule(),
		notification.NewModule(natsClient),
		audit.NewModule(),
		checkout.NewModule(redisClient, natsClient),
		invoice.NewModule(natsClient),
	}, nil
}

func NewContainer(clients *infra.Clients) (*app.Container, error) {
	modules, err := BuildModules(clients)
	if err != nil {
		return nil, err
	}

	container := app.NewContainer()
	if err := container.RegisterModules(modules...); err != nil {
		return nil, err
	}
	return container, nil
}
