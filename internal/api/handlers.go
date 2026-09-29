package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nae/infra/cache"
	"nae/infra/db"
	myNats "nae/infra/mq"
	"nae/internal/domain"
)

type OrderHandler struct {
	db        *db.PostgresRepo
	cache     *cache.RedisCache
	publisher *myNats.Publisher
}

func NewOrderHandler(db *db.PostgresRepo, cache *cache.RedisCache, pub *myNats.Publisher) *OrderHandler {
	return &OrderHandler{
		db:        db,
		cache:     cache,
		publisher: pub,
	}
}

type CreateOrderRequest struct {
	OrderID  string  `json:"order_id"`
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	order := &domain.Order{
		ID:        req.OrderID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    "CREATED",
		CreatedAt: time.Now().UTC(),
	}

	// 1. Write to PostgreSQL DB
	if err := h.db.CreateOrder(r.Context(), order); err != nil {
		http.Error(w, fmt.Sprintf("DB Write Error: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Cache in Redis (Cache-Aside pattern)
	_ = h.cache.SetOrder(r.Context(), order, 10*time.Minute)

	// 3. Publish Event to NATS JetStream
	payload, _ := json.Marshal(order)
	traceID := r.Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = fmt.Sprintf("tr_%d", time.Now().UnixNano())
	}

	ack, err := h.publisher.Publish(r.Context(), myNats.Event{
		Subject: "ORDERS.created",
		Payload: payload,
		Headers: map[string]string{
			"X-Trace-ID": traceID,
			"X-User-ID":  order.UserID,
		},
	})
	if err != nil {
		// Log warning: Record was saved to DB, but message queue publish failed
		fmt.Printf("[Warning] Event publish failed: %v\n", err)
	}

	// 4. Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":          "SUCCESS",
		"order":           order,
		"stream_sequence": ack.Sequence,
	})
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Missing id query param", http.StatusBadRequest)
		return
	}

	// 1. Try reading from Redis Cache
	cachedOrder, err := h.cache.GetOrder(r.Context(), orderID)
	if err == nil && cachedOrder != nil {
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(cachedOrder)
		return
	}

	// 2. Cache Miss: Fall back to PostgreSQL DB
	order, err := h.db.GetOrderByID(r.Context(), orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// 3. Populate Redis Cache asynchronously
	go func() {
		_ = h.cache.SetOrder(context.Background(), order, 10*time.Minute)
	}()

	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}
