package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

type NotificationService struct {
	NatsClient *nats.NATSClient
	UserAPI    *externals.UserAPI
}

func NewNotificationService(natsClient *nats.NATSClient, userAPI *externals.UserAPI) *NotificationService {
	return &NotificationService{
		NatsClient: natsClient,
		UserAPI:    userAPI,
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

func (s *NotificationService) NotifyRequestTerminated(ctx context.Context, requestTitle string, userIds []string) error {
	notification := types.Notification{
		ToUserIds: userIds,
		Subject:   "Request Terminated",
		Body:      fmt.Sprintf("The request “%s” has been terminated and will no longer proceed.", requestTitle),
	}
	return s.SendNotification(ctx, notification)
}

func (s *NotificationService) NotifyRequestCompleted(ctx context.Context, requestTitle string, userIds []string) error {
	notification := types.Notification{
		ToUserIds: userIds,
		Subject:   "Request Completed",
		Body:      fmt.Sprintf("The request “%s” has been completed.", requestTitle),
	}
	return s.SendNotification(ctx, notification)
}

func (s *NotificationService) NotifyTaskCompleted(ctx context.Context, taskTitle string, userId int32) error {
	users, err := s.UserAPI.FindUsersByUserIds([]int32{userId})
	if err != nil {
		return err
	}

	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   fmt.Sprintf("Task completed: %s", taskTitle),
		Body:      fmt.Sprintf("%s has marked this task as done.", users.Data[0].Name),
	}
	return s.SendNotification(ctx, notification)
}

func (s *NotificationService) NotifyTaskReassigned(ctx context.Context, taskTitle string, userId int32, userName string) error {
	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   fmt.Sprintf("Task reassigned: %s", taskTitle),
		Body:      fmt.Sprintf("The task “%s” has been reassigned to you.", taskTitle),
	}
	return s.SendNotification(ctx, notification)
}

func (s *NotificationService) NotifyTaskAvailable(ctx context.Context, taskTitle string, userId int32) error {
	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   "You Have New Task Available Today",
		Body:      fmt.Sprintf("%s is ready to start.", taskTitle),
	}
	return s.SendNotification(ctx, notification)
}

func (s *NotificationService) NotifyStartWorkflow(ctx context.Context, workflowTitle string, userIds []string) error {
	notification := types.Notification{
		ToUserIds: userIds,
		Subject:   "New Request Started",
		Body:      fmt.Sprintf("A new request “%s” has been initiated.", workflowTitle),
	}
	return s.SendNotification(ctx, notification)
}

func (s *NotificationService) NotifyNodeApproveNeeded(ctx context.Context, requestTitle string, userId int32) error {
	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   "Approval Needed",
		Body:      fmt.Sprintf("The request “%s” is pending your approval.", requestTitle),
	}
	return s.SendNotification(ctx, notification)
}
