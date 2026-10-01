package domain

import "time"

type CatalogItem struct {
	ID          string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
}
