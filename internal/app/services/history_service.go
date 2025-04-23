package services

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
)

type HistoryService struct {
	DB          *sql.DB
	HistoryRepo *repositories.HistoryRepository
	UserAPI     *externals.UserAPI
}

func NewHistoryService(db *sql.DB, historyRepo *repositories.HistoryRepository, userAPI *externals.UserAPI) *HistoryService {
	return &HistoryService{
		DB:          db,
		HistoryRepo: historyRepo,
		UserAPI:     userAPI,
	}
}

func (s *HistoryService) HistoryChangeNodeStatus(ctx context.Context, userId int32, requestId int32, nodeId string, fromStatus *string, toStatus string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	history := model.Histories{
		UserID:     &userId,
		RequestID:  requestId,
		NodeID:     nodeId,
		TypeAction: constants.HistoryTypeStatus,
		FromValue:  fromStatus,
		ToValue:    &toStatus,
	}

	err = s.HistoryRepo.CreateHistory(ctx, tx, history)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *HistoryService) HistoryApproveOrRejectNode(ctx context.Context, userId int32, requestId int32, nodeId string, fromStatus *string, toStatus string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	history := model.Histories{
		UserID:     &userId,
		RequestID:  requestId,
		NodeID:     nodeId,
		TypeAction: constants.HistoryTypeApproveReject,
		ToValue:    &toStatus,
		FromValue:  fromStatus,
	}

	err = s.HistoryRepo.CreateHistory(ctx, tx, history)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *HistoryService) HistoryChangeNodeAssignee(ctx context.Context, userId int32, requestId int32, nodeId string, fromAssigneeId *int32, toAssigneeId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromAssigneeIdStr *string = nil
	if fromAssigneeId != nil {
		tempStr := strconv.Itoa(int(*fromAssigneeId))
		fromAssigneeIdStr = &tempStr
	}

	toAssigneeIdStr := strconv.Itoa(int(toAssigneeId))

	history := model.Histories{
		UserID:     &userId,
		RequestID:  requestId,
		NodeID:     nodeId,
		TypeAction: constants.HistoryTypeAssignee,
		FromValue:  fromAssigneeIdStr,
		ToValue:    &toAssigneeIdStr,
	}

	err = s.HistoryRepo.CreateHistory(ctx, tx, history)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *HistoryService) HistoryNewTask(ctx context.Context, requestId int32, nodeId string, toUserId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	toUserIdStr := strconv.Itoa(int(toUserId))
	history := model.Histories{
		UserID:     nil,
		RequestID:  requestId,
		NodeID:     nodeId,
		TypeAction: constants.HistoryTypeNewTask,
		ToValue:    &toUserIdStr,
		FromValue:  nil,
	}

	err = s.HistoryRepo.CreateHistory(ctx, tx, history)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *HistoryService) HistoryStartRequest(ctx context.Context, requestId int32, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	history := model.Histories{
		UserID:     nil,
		RequestID:  requestId,
		NodeID:     nodeId,
		TypeAction: constants.HistoryTypeStartRequest,
		FromValue:  nil,
		ToValue:    nil,
	}

	err = s.HistoryRepo.CreateHistory(ctx, tx, history)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *HistoryService) HistoryEndRequest(ctx context.Context, requestId int32, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	history := model.Histories{
		UserID:     nil,
		RequestID:  requestId,
		NodeID:     nodeId,
		TypeAction: constants.HistoryTypeEndRequest,
		FromValue:  nil,
		ToValue:    nil,
	}

	err = s.HistoryRepo.CreateHistory(ctx, tx, history)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
