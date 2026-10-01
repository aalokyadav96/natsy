package domain

import "context"

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt string
}

type UserRepository interface {
	Save(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
}
