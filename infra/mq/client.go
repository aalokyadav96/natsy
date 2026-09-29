package mq

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	NC *nats.Conn
	JS jetstream.JetStream
}

func NewClient(url string) (*Client, error) {
	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			fmt.Println("[NATS Warning] Disconnected:", err)
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			fmt.Println("[NATS Info] Reconnected to", c.ConnectedUrl())
		}),
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream init failed: %w", err)
	}

	return &Client{NC: nc, JS: js}, nil
}

func (c *Client) Close() {
	if c.NC != nil {
		c.NC.Drain()
	}
}
