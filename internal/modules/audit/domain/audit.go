package domain

import "time"

type AuditEntry struct {
	ID        string
	Module    string
	Action    string
	EntityID  string
	CreatedAt time.Time
}
