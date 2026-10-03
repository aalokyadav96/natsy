package usecase

import (
	"context"
	"fmt"
	"time"

	"nae/internal/modules/audit/domain"
	sharedDomain "nae/internal/shared/domain"
)

type AuditUseCase interface {
	Log(ctx context.Context, module, action, entityID string) (*domain.AuditEntry, error)
}

type auditUseCase struct {
	repo domain.AuditRepository
}

func NewAuditUseCase(repo domain.AuditRepository) AuditUseCase { return &auditUseCase{repo: repo} }

func (u *auditUseCase) Log(ctx context.Context, module, action, entityID string) (*domain.AuditEntry, error) {
	if module == "" || action == "" || entityID == "" {
		return nil, fmt.Errorf("%w: module, action and entity_id are required", sharedDomain.ErrInvalidInput)
	}
	if u.repo == nil {
		return nil, fmt.Errorf("%w: audit repository is not configured", sharedDomain.ErrInternalError)
	}

	entry := &domain.AuditEntry{
		ID:        fmt.Sprintf("aud_%d", time.Now().UnixNano()),
		Module:    module,
		Action:    action,
		EntityID:  entityID,
		CreatedAt: time.Now().UTC(),
	}
	if err := u.repo.Save(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}
