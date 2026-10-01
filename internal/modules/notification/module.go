package notification

import (
	"nae/internal/app"
	notificationHTTP "nae/internal/modules/notification/port/http"
	notificationUC "nae/internal/modules/notification/usecase"
	"nae/internal/shared/infra/mq"
)

type Module struct {
	handler *notificationHTTP.NotificationHandler
}

func NewModule(nats *mq.NATSClient) *Module {
	return &Module{handler: notificationHTTP.NewNotificationHandler(notificationUC.NewNotificationUseCase(nats))}
}

func (m *Module) Name() string { return "notification" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	return nil
}
