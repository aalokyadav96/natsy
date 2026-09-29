package api

import (
	"log"
	"net/http"
	"time"

	"nae/infra/cache"
	"nae/infra/db"
	myNats "nae/infra/mq"
)

// NewRouter wires up the handlers, dependencies, and middleware
func NewRouter(
	pgRepo *db.PostgresRepo,
	redisCache *cache.RedisCache,
	pub *myNats.Publisher,
) http.Handler {
	mux := http.NewServeMux()

	// Initialize handler with all dependencies (DB, Cache, NATS Publisher)
	orderHandler := NewOrderHandler(pgRepo, redisCache, pub)

	// Routes (Go 1.22+ syntax matching HTTP methods)
	mux.HandleFunc("POST /api/v1/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /api/v1/orders", orderHandler.GetOrder)

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
