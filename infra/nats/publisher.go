package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	js jetstream.JetStream
}

func NewPublisher(js jetstream.JetStream) *Publisher {
	return &Publisher{js: js}
}

// Event wraps body data and standard operational metadata
type Event struct {
	Subject string
	Payload []byte
	Headers map[string]string
}

func (p *Publisher) Publish(ctx context.Context, evt Event) (*jetstream.PubAck, error) {
	headers := nats.Header{}
	for k, v := range evt.Headers {
		headers.Set(k, v)
	}

	msg := &nats.Msg{
		Subject: evt.Subject,
		Header:  headers,
		Data:    evt.Payload,
	}

	ack, err := p.js.PublishMsg(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("publish failed for subject %s: %w", evt.Subject, err)
	}
	return ack, nil
}
