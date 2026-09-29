package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"

	"nae/infra/cache"
	"nae/infra/db"
	myNats "nae/infra/mq"
	"nae/internal/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Connect Postgres
	pgConn := os.Getenv("POSTGRES_URL")
	if pgConn == "" {
		pgConn = "postgres://apeman:ningning@localhost:5432/eventdb?sslmode=disable"
	}
	pgRepo, err := db.NewPostgresRepo(ctx, pgConn)
	if err != nil {
		log.Fatalf("Postgres connection failed: %v", err)
	}
	defer pgRepo.Pool.Close()

	// 2. Connect Redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}
	redisCache, err := cache.NewRedisCache(redisURL)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	// 3. Connect NATS JetStream
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}
	natsClient, err := myNats.NewClient(natsURL)
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}
	defer natsClient.Close()

	// Provision stream
	_, err = natsClient.EnsureStream(ctx, myNats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})
	if err != nil {
		log.Fatalf("NATS Stream setup failed: %v", err)
	}

	// 4. Setup Routes
	pub := myNats.NewPublisher(natsClient.JS)
	orderHandler := api.NewOrderHandler(pgRepo, redisCache, pub)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /api/v1/orders", orderHandler.GetOrder)

	server := &http.Server{
		Addr:         ":4000",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 5. Run Server
	go func() {
		log.Printf("Server listening on %s...", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server stopped: %v", err)
		}
	}()

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
	}
}
