package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

type NotificationService struct {
	DB          *sql.DB
	NatsClient  *nats.NATSClient
	UserAPI     *externals.UserAPI
	RequestRepo *repositories.RequestRepository
}

func NewNotificationService(db *sql.DB, natsClient *nats.NATSClient, userAPI *externals.UserAPI, requestRepo *repositories.RequestRepository) *NotificationService {
	return &NotificationService{
		DB:          db,
		NatsClient:  natsClient,
		UserAPI:     userAPI,
		RequestRepo: requestRepo,
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

func (s *NotificationService) NotifyTaskCompleted(ctx context.Context, node model.Nodes) error {

	userIds := []string{}
	existingUserIds := map[int32]bool{}
	isSendNotification := false

	if node.TaskCompletedAssignee {
		isSendNotification = true
		userIds = append(userIds, strconv.Itoa(int(*node.AssigneeID)))
		existingUserIds[*node.AssigneeID] = true
	}

	if node.TaskCompletedRequester || node.TaskCompletedParticipants {
		isSendNotification = true

		request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, node.RequestID)
		if err != nil {
			return err
		}

		if node.TaskCompletedRequester {
			userIds = append(userIds, strconv.Itoa(int(request.UserID)))
			existingUserIds[request.UserID] = true
		}

		if node.TaskCompletedParticipants {
			for _, nodeRequest := range request.Nodes {
				if !existingUserIds[*nodeRequest.AssigneeID] {
					userIds = append(userIds, strconv.Itoa(int(*nodeRequest.AssigneeID)))
					existingUserIds[*nodeRequest.AssigneeID] = true
				}
			}

		}
	}

	if isSendNotification {

		users, err := s.UserAPI.FindUsersByUserIds([]int32{*node.AssigneeID})
		if err != nil {
			return err
		}

		notification := types.Notification{
			ToUserIds: userIds,
			Subject:   fmt.Sprintf("Task completed: %s", node.Title),
			Body:      fmt.Sprintf("%s has marked this task as done.", users.Data[0].Name),
		}
		return s.SendNotification(ctx, notification)
	}

	return nil
}

func (s *NotificationService) NotifyTaskStarted(ctx context.Context, node model.Nodes) error {

	userIds := []string{}
	existingUserIds := map[int32]bool{}
	isSendNotification := false

	if node.TaskCompletedAssignee {
		isSendNotification = true
		userIds = append(userIds, strconv.Itoa(int(*node.AssigneeID)))
		existingUserIds[*node.AssigneeID] = true
	}

	if node.TaskCompletedRequester || node.TaskCompletedParticipants {
		isSendNotification = true

		request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, node.RequestID)
		if err != nil {
			return err
		}

		if node.TaskCompletedRequester {
			userIds = append(userIds, strconv.Itoa(int(request.UserID)))
			existingUserIds[request.UserID] = true
		}

		if node.TaskCompletedParticipants {
			for _, nodeRequest := range request.Nodes {
				if !existingUserIds[*nodeRequest.AssigneeID] {
					userIds = append(userIds, strconv.Itoa(int(*nodeRequest.AssigneeID)))
					existingUserIds[*nodeRequest.AssigneeID] = true
				}
			}

		}
	}

	if isSendNotification {

		users, err := s.UserAPI.FindUsersByUserIds([]int32{*node.AssigneeID})
		if err != nil {
			return err
		}

		notification := types.Notification{
			ToUserIds: userIds,
			Subject:   fmt.Sprintf("Task started: %s", node.Title),
			Body:      fmt.Sprintf("%s has started this task.", users.Data[0].Name),
		}
		return s.SendNotification(ctx, notification)
	}

	return nil
}

func (s *NotificationService) NotifyTaskAvailable(ctx context.Context, taskTitle string, userId int32) error {
	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   "You Have New Task Available Today",
		Body:      fmt.Sprintf("%s is ready to start.", taskTitle),
	}
	return s.SendNotification(ctx, notification)
}

func (s *NotificationService) NotifyStartRequestWithDetail(ctx context.Context, userId int32, subject string, body string) error {
	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   subject,
		Body:      body,
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
