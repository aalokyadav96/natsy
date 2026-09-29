package mq

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

// Handler signature for processing domain messages
type Handler func(ctx context.Context, msg jetstream.Msg) error

type ConsumerConfig struct {
	StreamName    string
	DurableName   string
	FilterSubject string
}

type ConsumerRunner struct {
	js jetstream.JetStream
}

func NewConsumerRunner(js jetstream.JetStream) *ConsumerRunner {
	return &ConsumerRunner{js: js}
}

func (r *ConsumerRunner) RegisterWorker(ctx context.Context, cfg ConsumerConfig, handler Handler) (jetstream.ConsumeContext, error) {
	cons, err := r.js.CreateOrUpdateConsumer(ctx, cfg.StreamName, jetstream.ConsumerConfig{
		Durable:       cfg.DurableName,
		FilterSubject: cfg.FilterSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer %s: %w", cfg.DurableName, err)
	}

	consumeCtx, err := cons.Consume(func(msg jetstream.Msg) {
		// Invoke domain logic
		if err := handler(ctx, msg); err != nil {
			// NAK message so it can be retried according to consumer config
			fmt.Printf("[Worker Error] %s, Nak'ing message\n", err)
			_ = msg.Nak()
			return
		}

		// Acknowledge on success
		_ = msg.Ack()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start consume loop: %w", err)
	}

	return consumeCtx, nil
}
