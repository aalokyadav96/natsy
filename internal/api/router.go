package api

import (
	"log"
	"net/http"
	"time"

	myNats "nae/infra/nats"
)

func NewRouter(pub *myNats.Publisher) http.Handler {
	mux := http.NewServeMux()

	orderHandler := NewOrderHandler(pub)

	// Routes
	mux.HandleFunc("/api/v1/orders", orderHandler.CreateOrder)

	// Attach middleware chain
	return loggingMiddleware(mux)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s - %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}
