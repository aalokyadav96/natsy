package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nae/internal/modules/notification/domain"
	"nae/internal/shared/infra/mq"
)

type NotificationUseCase interface {
	Send(ctx context.Context, userID, msg string) (*domain.Notification, error)
}

type notificationUseCase struct {
	nats *mq.NATSClient
}

func NewNotificationUseCase(nats *mq.NATSClient) NotificationUseCase {
	return &notificationUseCase{nats: nats}
}

func (u *notificationUseCase) Send(ctx context.Context, userID, msg string) (*domain.Notification, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if msg == "" {
		return nil, fmt.Errorf("message is required")
	}

	notification := &domain.Notification{
		ID:        fmt.Sprintf("nt_%d", time.Now().UnixNano()),
		UserID:    userID,
		Type:      "info",
		Message:   msg,
		CreatedAt: time.Now().UTC(),
	}

	if u.nats != nil {
		payload, err := json.Marshal(map[string]any{
			"notification_id": notification.ID,
			"user_id":         notification.UserID,
			"message":         notification.Message,
		})
		if err == nil {
			_, _ = u.nats.JS.Publish(ctx, "NOTIFICATIONS.created", payload)
		}
	}

	return notification, nil
}
