package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	myNats "nae/infra/nats"
	"nae/internal/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize NATS JetStream Client
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	client, err := myNats.NewClient(natsURL)
	if err != nil {
		log.Fatalf("NATS Connection Error: %v", err)
	}
	defer client.Close()

	// 2. Provision Streams idempotently
	_, err = client.EnsureStream(ctx, myNats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})
	if err != nil {
		log.Fatalf("Stream Setup Error: %v", err)
	}

	// 3. Initialize Publisher & Router
	publisher := myNats.NewPublisher(client.JS)
	router := api.NewRouter(publisher)

	server := &http.Server{
		Addr:         ":4000",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// 4. Start HTTP Server asynchronously
	go func() {
		log.Printf("HTTP Server listening on port %s...", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP Server failed: %v", err)
		}
	}()

	// 5. Graceful Shutdown Signal Interceptor
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	log.Println("Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP Shutdown Error: %v", err)
	}

	fmt.Println("Server exited successfully.")
}
