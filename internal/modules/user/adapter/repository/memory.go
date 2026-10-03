package repository

import (
	"context"
	"fmt"
	"sync"

	"nae/internal/modules/user/domain"
	sharedDomain "nae/internal/shared/domain"
)

type inMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserRepository() domain.UserRepository {
	return &inMemoryUserRepository{users: make(map[string]*domain.User)}
}

func (r *inMemoryUserRepository) Save(ctx context.Context, user *domain.User) error {
	if user == nil {
		return fmt.Errorf("%w: user is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.users == nil {
		r.users = make(map[string]*domain.User)
	}
	clone := *user
	r.users[user.ID] = &clone
	return nil
}

func (r *inMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: user id is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, fmt.Errorf("%w: user %s was not found", sharedDomain.ErrNotFound, id)
	}
	clone := *user
	return &clone, nil
}
