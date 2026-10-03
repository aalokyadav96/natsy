package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nae/internal/bootstrap"
	"nae/internal/shared/infra"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// -------------------------------------------------------------------------
	// 1. Shared Infrastructure Setup
	// -------------------------------------------------------------------------
	cfg := infra.LoadConfig()
	clients, err := infra.NewClients(ctx, cfg)
	if err != nil {
		log.Fatalf("Infrastructure initialization failed: %v", err)
	}
	defer clients.Close()

	// -------------------------------------------------------------------------
	// 2. Module registration via composition root
	// -------------------------------------------------------------------------
	moduleContainer, err := bootstrap.NewContainer(clients)
	if err != nil {
		log.Fatalf("Module composition failed: %v", err)
	}

	if err := moduleContainer.Start(ctx); err != nil {
		log.Fatalf("Module startup failed: %v", err)
	}

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      moduleContainer.Mux(),
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
