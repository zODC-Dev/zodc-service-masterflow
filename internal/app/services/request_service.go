package services

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type RequestService struct {
	db          *sql.DB
	requestRepo *repositories.RequestRepository
	userApi     *externals.UserAPI
}

func NewRequestService(db *sql.DB, requestRepo *repositories.RequestRepository, userApi *externals.UserAPI) *RequestService {
	return &RequestService{
		db:          db,
		requestRepo: requestRepo,
		userApi:     userApi,
	}
}

func (s *RequestService) FindAllRequestHandler(ctx context.Context, requestQueryParam queryparams.RequestQueryParam, userId int32) (responses.Paginate[[]responses.RequestResponse], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestResponse]{}

	count, requests, err := s.requestRepo.FindAllRequest(ctx, s.db, requestQueryParam, userId)
	if err != nil {
		return paginatedResponse, err
	}

	total := count.Count

	requestsResponse := []responses.RequestResponse{}
	for _, request := range requests {
		requestResponse := responses.RequestResponse{}
		if err := utils.Mapper(request, &requestResponse); err != nil {
			return paginatedResponse, err
		}

		// Parent Key
		if request.ParentID != nil {
			parentRequest, err := s.requestRepo.FindOneRequestByRequestId(ctx, s.db, *request.ParentID)
			if err != nil {
				return paginatedResponse, err
			}
			requestResponse.ParentKey = parentRequest.Key
		}

		// Tasks and Process - Current Tasks
		completedNodes := 0
		totalNodes := len(request.Nodes)

		requestResponse.CurrentTasks = []responses.CurrentTaskResponse{}
		for _, node := range request.Nodes {

			userIdsTask := []int32{}
			// Participants
			if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
				subRequest, err := s.requestRepo.FindOneRequestByRequestId(ctx, s.db, *node.SubRequestID)
				if err != nil {
					return paginatedResponse, err
				}

				for _, subNode := range subRequest.Nodes {
					userIdsTask = append(userIdsTask, *subNode.AssigneeID)
				}

			} else {
				userIdsTask = append(userIdsTask, *node.AssigneeID)
			}

			taskUsers, err := s.userApi.FindUsersByUserIds(userIdsTask)
			if err != nil {
				return paginatedResponse, err
			}

			participants := []types.Assignee{}
			for _, user := range taskUsers.Data {
				participant := types.Assignee{}
				if err := utils.Mapper(user, &participant); err != nil {
					return paginatedResponse, err
				}

				participants = append(participants, participant)
			}

			currentTask := responses.CurrentTaskResponse{
				Title:        *node.Title,
				UpdatedAt:    node.UpdatedAt,
				Participants: participants,
			}

			requestResponse.CurrentTasks = append(requestResponse.CurrentTasks, currentTask)

			if node.Status == string(constants.NodeStatusCompleted) {
				completedNodes++
			}

			if node.Type == string(constants.NodeTypeStart) || node.Type == string(constants.NodeTypeEnd) {
				totalNodes--
			}
		}

		// Set CompletedAt
		if request.Status == string(constants.RequestStatusCompleted) {
			requestResponse.CompletedAt = &request.UpdatedAt
		}

		requestsResponse = append(requestsResponse, requestResponse)
	}

	totalPages := (int(total) + requestQueryParam.PageSize - 1) / requestQueryParam.PageSize

	paginatedResponse = responses.Paginate[[]responses.RequestResponse]{
		Items:      requestsResponse,
		Total:      int(total),
		Page:       requestQueryParam.Page,
		PageSize:   requestQueryParam.PageSize,
		TotalPages: totalPages,
	}

	return paginatedResponse, nil
}

func (s *RequestService) GetRequestOverviewHandler(ctx context.Context, userId int32) (responses.RequestOverviewResponse, error) {
	requestOverviewResponse := responses.RequestOverviewResponse{}
	var err error

	count, err := s.requestRepo.CountRequestByStatusAndUserId(ctx, s.db, userId, "")
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.MyRequests = int32(count)

	count, err = s.requestRepo.CountRequestByStatusAndUserId(ctx, s.db, userId, constants.RequestStatusInProcessing)
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.InProcess = int32(count)

	count, err = s.requestRepo.CountRequestByStatusAndUserId(ctx, s.db, userId, constants.RequestStatusCompleted)
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.Completed = int32(count)

	count, err = s.requestRepo.CountUserAppendInRequestAndNodeUserId(ctx, s.db, userId)
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.AllRequests = int32(count)

	return requestOverviewResponse, nil
}
