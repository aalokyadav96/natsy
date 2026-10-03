package repository

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"nae/internal/modules/catalog/domain"
	sharedDomain "nae/internal/shared/domain"
)

type inMemoryCatalogRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.CatalogItem
}

func NewInMemoryCatalogRepository() domain.CatalogRepository {
	items := map[string]*domain.CatalogItem{
		"prod_1": {ID: "prod_1", Name: "Laptop", Description: "14-inch laptop", Price: 999.99, CreatedAt: time.Now().UTC()},
		"prod_2": {ID: "prod_2", Name: "Keyboard", Description: "Mechanical keyboard", Price: 129.99, CreatedAt: time.Now().UTC()},
		"prod_3": {ID: "prod_3", Name: "Mouse", Description: "Wireless mouse", Price: 49.99, CreatedAt: time.Now().UTC()},
	}
	return &inMemoryCatalogRepository{items: items}
}

func (r *inMemoryCatalogRepository) List(ctx context.Context) ([]*domain.CatalogItem, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*domain.CatalogItem, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

func (r *inMemoryCatalogRepository) GetByID(ctx context.Context, id string) (*domain.CatalogItem, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: catalog item id is required", sharedDomain.ErrInvalidInput)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	if !ok {
		return nil, fmt.Errorf("%w: catalog item %s was not found", sharedDomain.ErrNotFound, id)
	}
	clone := *item
	return &clone, nil
}
