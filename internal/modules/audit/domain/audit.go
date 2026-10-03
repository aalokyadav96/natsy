package domain

import (
	"context"
	"time"
)

type AuditEntry struct {
	ID        string
	Module    string
	Action    string
	EntityID  string
	CreatedAt time.Time
}

type AuditRepository interface {
	Save(ctx context.Context, entry *AuditEntry) error
	ListByEntity(ctx context.Context, entityID string) ([]*AuditEntry, error)
}
