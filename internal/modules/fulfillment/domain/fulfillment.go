package domain

import "time"

type Fulfillment struct {
	ID        string
	OrderID   string
	Status    string
	CreatedAt time.Time
}
