package event

import (
	"context"
	"encoding/json"
	"log"

	"nae/internal/modules/fulfillment/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/nats-io/nats.go/jetstream"
)

type FulfillmentConsumer struct {
	uc   usecase.FulfillmentUseCase
	nats *mq.NATSClient
}

func NewFulfillmentConsumer(uc usecase.FulfillmentUseCase, nats *mq.NATSClient) *FulfillmentConsumer {
	return &FulfillmentConsumer{uc: uc, nats: nats}
}

func (c *FulfillmentConsumer) Start(ctx context.Context) error {
	if c.nats == nil {
		return nil
	}

	cons, err := c.nats.JS.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Durable:       "fulfillment_service",
		FilterSubject: "ORDERS.created",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		var evt struct {
			OrderID string `json:"order_id"`
		}
		if err := json.Unmarshal(msg.Data(), &evt); err != nil {
			msg.Nak()
			return
		}
		if evt.OrderID == "" {
			msg.Ack()
			return
		}
		fulfillment, err := c.uc.Create(ctx, evt.OrderID)
		if err != nil {
			log.Printf("[Fulfillment Module] failed to create fulfillment for order %s: %v", evt.OrderID, err)
			msg.Nak()
			return
		}
		if c.nats != nil {
			payload, marshalErr := json.Marshal(map[string]any{
				"fulfillment_id": fulfillment.ID,
				"order_id":       evt.OrderID,
				"status":         fulfillment.Status,
			})
			if marshalErr == nil {
				_, _ = c.nats.JS.Publish(ctx, "FULFILLMENTS.created", payload)
			}
		}
		_ = msg.Ack()
	})
	return err
}
