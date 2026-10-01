package usecase

import (
	"context"

	"nae/internal/modules/audit/domain"
)

type AuditUseCase interface {
	Log(ctx context.Context, module, action, entityID string) (*domain.AuditEntry, error)
}

type auditUseCase struct{}

func NewAuditUseCase() AuditUseCase { return &auditUseCase{} }

func (u *auditUseCase) Log(ctx context.Context, module, action, entityID string) (*domain.AuditEntry, error) {
	return &domain.AuditEntry{ID: "aud_1", Module: module, Action: action, EntityID: entityID}, nil
}
