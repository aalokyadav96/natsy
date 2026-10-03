package domain

import (
	"context"
	"time"
)

type CatalogItem struct {
	ID          string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type CatalogRepository interface {
	List(ctx context.Context) ([]*CatalogItem, error)
	GetByID(ctx context.Context, id string) (*CatalogItem, error)
}
