package event

import (
	"context"
	"encoding/json"
	"log"

	"nae/internal/modules/inventory/usecase"
	"nae/internal/shared/infra/mq"

	"github.com/nats-io/nats.go/jetstream"
)

type InventoryConsumer struct {
	uc   usecase.InventoryUseCase
	nats *mq.NATSClient
}

func NewInventoryConsumer(uc usecase.InventoryUseCase, nats *mq.NATSClient) *InventoryConsumer {
	return &InventoryConsumer{uc: uc, nats: nats}
}

func (c *InventoryConsumer) Start(ctx context.Context) error {
	if c.nats == nil {
		return nil
	}

	cons, err := c.nats.JS.CreateOrUpdateConsumer(ctx, "ORDERS", jetstream.ConsumerConfig{
		Durable:       "inventory_service",
		FilterSubject: "ORDERS.created",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	_, err = cons.Consume(func(msg jetstream.Msg) {
		var evt struct {
			OrderID string `json:"order_id"`
			Items   []struct {
				ProductID string `json:"product_id"`
				Quantity  int    `json:"quantity"`
			} `json:"items"`
		}

		if err := json.Unmarshal(msg.Data(), &evt); err != nil {
			msg.Nak()
			return
		}

		for _, item := range evt.Items {
			if item.ProductID == "" {
				continue
			}
			if err := c.uc.Reserve(ctx, item.ProductID, item.Quantity); err != nil {
				log.Printf("[Inventory Module] failed to reserve stock for product %s: %v", item.ProductID, err)
				msg.Nak()
				return
			}
		}
		_ = msg.Ack()
	})
	return err
}
