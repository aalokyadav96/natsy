package domain

import "time"

type Notification struct {
	ID        string
	UserID    string
	Type      string
	Message   string
	CreatedAt time.Time
}
