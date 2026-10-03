package repository

import (
	"context"
	"fmt"
	"sync"

	"nae/internal/modules/audit/domain"
	sharedDomain "nae/internal/shared/domain"
)

type inMemoryAuditRepository struct {
	mu      sync.RWMutex
	entries map[string][]*domain.AuditEntry
}

func NewInMemoryAuditRepository() domain.AuditRepository {
	return &inMemoryAuditRepository{entries: make(map[string][]*domain.AuditEntry)}
}

func (r *inMemoryAuditRepository) Save(ctx context.Context, entry *domain.AuditEntry) error {
	if entry == nil {
		return fmt.Errorf("%w: audit entry is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.entries == nil {
		r.entries = make(map[string][]*domain.AuditEntry)
	}
	r.entries[entry.EntityID] = append(r.entries[entry.EntityID], entry)
	return nil
}

func (r *inMemoryAuditRepository) ListByEntity(ctx context.Context, entityID string) ([]*domain.AuditEntry, error) {
	if entityID == "" {
		return nil, fmt.Errorf("%w: entity id is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()
	entries := append([]*domain.AuditEntry(nil), r.entries[entityID]...)
	return entries, nil
}
