package event

import (
	"context"
	"encoding/json"
	"log"

	"nae/internal/shared/infra/mq"

	"github.com/nats-io/nats.go/jetstream"
)

type InvoiceConsumer struct {
	nats *mq.NATSClient
}

func NewInvoiceConsumer(nats *mq.NATSClient) *InvoiceConsumer {
	return &InvoiceConsumer{nats: nats}
}

func (c *InvoiceConsumer) Start(ctx context.Context) error {
	cons, err := c.nats.JS.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Durable:       "invoice_service",
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

		log.Printf("[Invoice Module] Generating invoice for Order ID: %s, User: %s, Total: $%.2f", evt.OrderID, evt.UserID, evt.Total)
		if c.nats != nil {
			payload, marshalErr := json.Marshal(map[string]any{
				"order_id": evt.OrderID,
				"user_id":  evt.UserID,
				"total":    evt.Total,
				"status":   "issued",
			})
			if marshalErr == nil {
				_, _ = c.nats.JS.Publish(ctx, "INVOICES.created", payload)
			}
		}
		_ = msg.Ack()
	})

	return err
}
