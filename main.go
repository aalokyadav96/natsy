package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Connect to NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// 2. Initialize JetStream management client
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Failed to create JetStream client: %v", err)
	}

	streamName := "ORDERS"
	subjectName := "ORDERS.created"

	// 3. Create or update the Stream
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"ORDERS.*"},
	})
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}
	fmt.Printf("Stream '%s' ready.\n", streamName)

	// 4. Create or update a durable Consumer
	cons, err := js.CreateOrUpdateConsumer(ctx, streamName, jetstream.ConsumerConfig{
		Durable:       "orders-processor",
		FilterSubject: subjectName,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	// 5. Subscribe / Consume messages asynchronously
	consumeContext, err := cons.Consume(func(msg jetstream.Msg) {
		// Extract message payload
		payload := string(msg.Data())

		// Extract metadata headers
		headers := msg.Headers()
		traceID := headers.Get("X-Trace-ID")
		userID := headers.Get("X-User-ID")

		// Extract JetStream message metadata (Sequence ID, Timestamp, etc.)
		meta, err := msg.Metadata()
		if err != nil {
			log.Printf("Error reading metadata: %v", err)
		} else {
			fmt.Printf("\n[RECEIVED] Stream Sequence: %d | Time: %s\n", meta.Sequence.Stream, meta.Timestamp.Format(time.RFC3339))
		}

		fmt.Printf(" Headers -> Trace-ID: %s, User-ID: %s\n", traceID, userID)
		fmt.Printf(" Payload -> %s\n", payload)

		// Acknowledge the message so JetStream knows it was processed
		if err := msg.Ack(); err != nil {
			log.Printf("Failed to ack message: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to start consume loop: %v", err)
	}
	defer consumeContext.Stop()

	// 6. Publish messages with metadata (Headers)
	headers := nats.Header{}
	headers.Set("X-Trace-ID", "trace-abc-12345")
	headers.Set("X-User-ID", "usr_9982")

	msgPayload := []byte(`{"order_id": 1001, "item": "Laptop", "amount": 1200.00}`)

	msg := &nats.Msg{
		Subject: subjectName,
		Header:  headers,
		Data:    msgPayload,
	}

	// Publish to JetStream
	pubAck, err := js.PublishMsg(ctx, msg)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	fmt.Printf("[PUBLISHED] Message sent to '%s' (Stream Seq: %d)\n", pubAck.Stream, pubAck.Sequence)

	// Keep process alive to allow consumer to receive the message
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down gracefully...")
}
