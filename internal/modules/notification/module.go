package notification

import (
	"nae/internal/app"
	notificationEvent "nae/internal/modules/notification/port/event"
	notificationHTTP "nae/internal/modules/notification/port/http"
	notificationUC "nae/internal/modules/notification/usecase"
	"nae/internal/shared/infra/mq"
)

type Module struct {
	handler  *notificationHTTP.NotificationHandler
	consumer *notificationEvent.NotificationConsumer
}

func NewModule(nats *mq.NATSClient) *Module {
	uc := notificationUC.NewNotificationUseCase(nats)
	return &Module{handler: notificationHTTP.NewNotificationHandler(uc), consumer: notificationEvent.NewNotificationConsumer(uc, nats)}
}

func (m *Module) Name() string { return "notification" }

func (m *Module) Register(c *app.Container) error {
	if m.handler != nil {
		c.RegisterRoutes(m.handler.RegisterRoutes)
	}
	if m.consumer != nil {
		c.RegisterStartup(m.consumer.Start)
	}
	return nil
}
