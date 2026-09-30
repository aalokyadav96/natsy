package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nae/internal/shared/infra/cache"
	"nae/internal/shared/infra/db"
	"nae/internal/shared/infra/mq"

	// Listing Module
	listingRepo "nae/internal/modules/listing/adapter/repository"
	listingHTTP "nae/internal/modules/listing/port/http"
	listingUC "nae/internal/modules/listing/usecase"

	// Cart Module
	cartRepo "nae/internal/modules/cart/adapter/repository"
	cartHTTP "nae/internal/modules/cart/port/http"
	cartUC "nae/internal/modules/cart/usecase"

	// Checkout Module
	checkoutHTTP "nae/internal/modules/checkout/port/http"
	checkoutUC "nae/internal/modules/checkout/usecase"

	// Invoice Module
	invoiceEvent "nae/internal/modules/invoice/port/event"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// -------------------------------------------------------------------------
	// 1. Shared Infrastructure Setup
	// -------------------------------------------------------------------------
	pgClient, err := db.NewPostgresClient(
		ctx,
		getEnv("POSTGRES_URL", "postgres://apeman:ningning@localhost:5432/eventdb?sslmode=disable"),
	)
	if err != nil {
		log.Fatalf("Postgres connection failed: %v", err)
	}
	defer pgClient.Close()

	redisClient, err := cache.NewRedisClient(
		getEnv("REDIS_URL", "redis://localhost:6379/0"),
	)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer redisClient.Close()

	natsClient, err := mq.NewNATSClient(
		getEnv("NATS_URL", "nats://localhost:4222"),
	)
	if err != nil {
		log.Fatalf("NATS connection failed: %v", err)
	}
	defer natsClient.Close()

	if err := natsClient.EnsureStream(ctx, "ORDERS", []string{"ORDERS.*"}); err != nil {
		log.Fatalf("NATS stream creation failed: %v", err)
	}

	// -------------------------------------------------------------------------
	// 2. Module Wireups
	// -------------------------------------------------------------------------

	// A. Listing Module
	lRepo := listingRepo.NewPostgresProductRepository(pgClient.Pool)
	lUC := listingUC.NewListingUseCase(lRepo)
	listingHandler := listingHTTP.NewListingHandler(lUC)

	// B. Cart Module
	cRepo := cartRepo.NewRedisCartRepository(redisClient.Client, 24*time.Hour)
	cUC := cartUC.NewCartUseCase(cRepo)
	cartHandler := cartHTTP.NewCartHandler(cUC)

	// C. Checkout Module
	coUC := checkoutUC.NewCheckoutUseCase(cRepo, natsClient)
	checkoutHandler := checkoutHTTP.NewCheckoutHandler(coUC)

	// D. Invoice Module (Background Event Consumer)
	invoiceConsumer := invoiceEvent.NewInvoiceConsumer(natsClient)
	if err := invoiceConsumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start invoice consumer: %v", err)
	}

	// -------------------------------------------------------------------------
	// 3. HTTP Router Registration
	// -------------------------------------------------------------------------
	mux := http.NewServeMux()

	listingHandler.RegisterRoutes(mux)
	cartHandler.RegisterRoutes(mux)
	checkoutHandler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":4000",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// -------------------------------------------------------------------------
	// 4. Server Execution & Graceful Shutdown
	// -------------------------------------------------------------------------
	go func() {
		log.Printf("Modular Monolith running on %s...", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server crash: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP Shutdown Error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
