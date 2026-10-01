package usecase

import (
	"context"

	"nae/internal/modules/user/domain"
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
	user := &domain.User{
		ID:    "usr_" + name,
		Email: email,
		Name:  name,
	}
	if err := u.repo.Save(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return u.repo.GetByID(ctx, id)
}
