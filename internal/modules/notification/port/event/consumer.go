package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"nae/internal/modules/notification/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/nats-io/nats.go/jetstream"
)

type NotificationConsumer struct {
	uc   usecase.NotificationUseCase
	nats *mq.NATSClient
}

func NewNotificationConsumer(uc usecase.NotificationUseCase, nats *mq.NATSClient) *NotificationConsumer {
	return &NotificationConsumer{uc: uc, nats: nats}
}

func (c *NotificationConsumer) Start(ctx context.Context) error {
	if c.nats == nil {
		return nil
	}

	cons, err := c.nats.JS.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Durable:       "notification_service",
		FilterSubject: "ORDERS.created",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		var evt struct {
			OrderID string  `json:"order_id"`
			UserID  string  `json:"user_id"`
			Total   float64 `json:"total"`
		}
		if err := json.Unmarshal(msg.Data(), &evt); err != nil {
			msg.Nak()
			return
		}
		if evt.UserID == "" {
			msg.Ack()
			return
		}
		if _, err := c.uc.Send(ctx, evt.UserID, fmt.Sprintf("Order %s created successfully for $%.2f", evt.OrderID, evt.Total)); err != nil {
			log.Printf("[Notification Module] failed to send user %s notification: %v", evt.UserID, err)
			msg.Nak()
			return
		}
		_ = msg.Ack()
	})
	return err
}
