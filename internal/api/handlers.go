package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	myNats "nae/infra/nats"
)

type OrderHandler struct {
	publisher *myNats.Publisher
}

func NewOrderHandler(pub *myNats.Publisher) *OrderHandler {
	return &OrderHandler{publisher: pub}
}

// Request payload definition
type CreateOrderRequest struct {
	OrderID  string  `json:"order_id"`
	UserID   string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// Response payload definition
type CreateOrderResponse struct {
	Status    string `json:"status"`
	OrderID   string `json:"order_id"`
	StreamSeq uint64 `json:"stream_sequence"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// 1. Serialize payload for NATS
	payload, err := json.Marshal(req)
	if err != nil {
		http.Error(w, "Failed to serialize event", http.StatusInternalServerError)
		return
	}

	// 2. Extract or generate metadata headers
	traceID := r.Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = fmt.Sprintf("tr_%d", time.Now().UnixNano()) // fallback trace ID generator
	}

	event := myNats.Event{
		Subject: "ORDERS.created",
		Payload: payload,
		Headers: map[string]string{
			"X-Trace-ID":   traceID,
			"X-User-ID":    req.UserID,
			"User-Agent":   r.UserAgent(),
			"Content-Type": "application/json",
		},
	}

	// 3. Publish to NATS JetStream synchronously during request execution
	ack, err := h.publisher.Publish(r.Context(), event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to publish event: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Return REST HTTP response with JetStream Ack Sequence
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(CreateOrderResponse{
		Status:    "ORDER_RECEIVED",
		OrderID:   req.OrderID,
		StreamSeq: ack.Sequence,
	})
}
