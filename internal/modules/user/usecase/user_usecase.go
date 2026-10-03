package usecase

import (
	"context"
	"fmt"
	"time"

	"nae/internal/modules/user/domain"
	sharedDomain "nae/internal/shared/domain"
)

type UserUseCase interface {
	Register(ctx context.Context, name, email string) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
}

type userUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) UserUseCase {
	return &userUseCase{repo: repo}
}

func (u *userUseCase) Register(ctx context.Context, name, email string) (*domain.User, error) {
	if name == "" || email == "" {
		return nil, fmt.Errorf("%w: user name and email are required", sharedDomain.ErrInvalidInput)
	}
	if u.repo == nil {
		return nil, fmt.Errorf("%w: user repository is not configured", sharedDomain.ErrInternalError)
	}

	user := &domain.User{
		ID:        fmt.Sprintf("usr_%d", time.Now().UnixNano()),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := u.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user id is required", sharedDomain.ErrInvalidInput)
	}
	if u.repo == nil {
		return nil, fmt.Errorf("%w: user repository is not configured", sharedDomain.ErrInternalError)
	}
	return u.repo.GetByID(ctx, id)
}
