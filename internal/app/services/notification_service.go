package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

type NotificationService struct {
	NatsClient *nats.NATSClient
}

func NewNotificationService(natsClient *nats.NATSClient) *NotificationService {
	return &NotificationService{
		NatsClient: natsClient,
	}
}

func (s *NotificationService) SendNotification(ctx context.Context, notification types.Notification) error {
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshal notification failed: %w", err)
	}

	err = s.NatsClient.Publish("notifications", notificationBytes)
	if err != nil {
		return fmt.Errorf("publish notification failed: %w", err)
	}

	return nil
}
