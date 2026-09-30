package mq

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NATSClient struct {
	NC *nats.Conn
	JS jetstream.JetStream
}

func NewNATSClient(url string) (*NATSClient, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream init failed: %w", err)
	}

	return &NATSClient{NC: nc, JS: js}, nil
}

func (c *NATSClient) EnsureStream(ctx context.Context, name string, subjects []string) error {
	_, err := c.JS.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     name,
		Subjects: subjects,
		Storage:  jetstream.FileStorage,
	})
	return err
}

func (c *NATSClient) Close() {
	if c.NC != nil {
		c.NC.Drain()
	}
}
