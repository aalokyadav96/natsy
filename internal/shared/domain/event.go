package domain

import "time"

// Event represents a standard domain event interface emitted across modules
type Event interface {
	EventName() string
	OccurredAt() time.Time
}

// BaseEvent provides standard metadata fields for all domain events
type BaseEvent struct {
	EventID   string    `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}
