package event

import (
	"context"
	"encoding/json"
	"log"

	"nae/internal/modules/payment/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/nats-io/nats.go/jetstream"
)

type PaymentConsumer struct {
	uc   usecase.PaymentUseCase
	nats *mq.NATSClient
}

func NewPaymentConsumer(uc usecase.PaymentUseCase, nats *mq.NATSClient) *PaymentConsumer {
	return &PaymentConsumer{uc: uc, nats: nats}
}

func (c *PaymentConsumer) Start(ctx context.Context) error {
	if c.nats == nil {
		return nil
	}

	cons, err := c.nats.JS.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Durable:       "payment_service",
		FilterSubject: "ORDERS.created",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		var evt struct {
			OrderID string  `json:"order_id"`
			Total   float64 `json:"total"`
		}
		if err := json.Unmarshal(msg.Data(), &evt); err != nil {
			msg.Nak()
			return
		}
		if evt.OrderID == "" {
			msg.Ack()
			return
		}
		if _, err := c.uc.Create(ctx, evt.OrderID, evt.Total); err != nil {
			log.Printf("[Payment Module] failed to create payment for order %s: %v", evt.OrderID, err)
			msg.Nak()
			return
		}
		_ = msg.Ack()
	})
	return err
}
