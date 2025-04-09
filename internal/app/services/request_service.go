package services

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type RequestService struct {
	DB              *sql.DB
	RequestRepo     *repositories.RequestRepository
	UserAPI         *externals.UserAPI
	NodeRepo        *repositories.NodeRepository
	WorkflowService *WorkflowService
}

func NewRequestService(cfg RequestService) *RequestService {
	return &RequestService{
		DB:              cfg.DB,
		RequestRepo:     cfg.RequestRepo,
		UserAPI:         cfg.UserAPI,
		NodeRepo:        cfg.NodeRepo,
		WorkflowService: cfg.WorkflowService,
	}
}

func (s *RequestService) FindAllRequestHandler(ctx context.Context, requestQueryParam queryparams.RequestQueryParam, userId int32) (responses.Paginate[[]responses.RequestResponse], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestResponse]{}

	count, requests, err := s.RequestRepo.FindAllRequest(ctx, s.DB, requestQueryParam, userId)
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
			parentRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *request.ParentID)
			if err != nil {
				return paginatedResponse, err
			}
			requestResponse.ParentKey = parentRequest.Key
		}

		// Tasks and Process - Current Tasks
		requestResponse.CurrentTasks = []responses.CurrentTaskResponse{}
		for _, node := range request.Nodes {
			if node.Type == string(constants.NodeTypeEnd) || node.Type == string(constants.NodeTypeStart) {
				continue
			}

			userIdsTask := []int32{}
			// Participants
			if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
				subRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *node.SubRequestID)
				if err != nil {
					return paginatedResponse, err
				}

				for _, subNode := range subRequest.Nodes {
					userIdsTask = append(userIdsTask, *subNode.AssigneeID)
				}

			} else {
				userIdsTask = append(userIdsTask, *node.AssigneeID)
			}

			taskUsers, err := s.UserAPI.FindUsersByUserIds(userIdsTask)
			if err != nil {
				return paginatedResponse, err
			}

			participants := []types.Assignee{}
			existingUserIds := make(map[int32]bool)

			for _, user := range taskUsers.Data {
				if _, exists := existingUserIds[user.ID]; exists {
					continue
				}

				participant := types.Assignee{}
				if err := utils.Mapper(user, &participant); err != nil {
					return paginatedResponse, err
				}

				participants = append(participants, participant)
				existingUserIds[user.ID] = true
			}

			currentTask := responses.CurrentTaskResponse{
				Title:        node.Title,
				UpdatedAt:    node.UpdatedAt,
				Participants: participants,
			}

			requestResponse.CurrentTasks = append(requestResponse.CurrentTasks, currentTask)

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

func (s *RequestService) GetRequestCountHandler(ctx context.Context, userId int32) (responses.RequestCountResponse, error) {
	requestOverviewResponse := responses.RequestCountResponse{}
	var err error

	count, err := s.RequestRepo.CountRequestByStatusAndUserId(ctx, s.DB, userId, "")
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.MyRequests = int32(count)

	count, err = s.RequestRepo.CountRequestByStatusAndUserId(ctx, s.DB, userId, constants.RequestStatusInProgress)
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.InProcess = int32(count)

	count, err = s.RequestRepo.CountRequestByStatusAndUserId(ctx, s.DB, userId, constants.RequestStatusCompleted)
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.Completed = int32(count)

	count, err = s.RequestRepo.CountUserAppendInRequestAndNodeUserId(ctx, s.DB, userId)
	if err != nil {
		return requestOverviewResponse, err
	}
	requestOverviewResponse.AllRequests = int32(count)

	return requestOverviewResponse, nil
}

func (s *RequestService) GetRequestDetailHandler(ctx context.Context, userId int32, requestId int32) (responses.RequestDetailResponse, error) {
	requestDetailResponse := responses.RequestDetailResponse{}

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return requestDetailResponse, err
	}

	if err := utils.Mapper(request, &requestDetailResponse); err != nil {
		return requestDetailResponse, err
	}

	// Parent Request
	if request.ParentID != nil {
		parentRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *request.ParentID)
		if err != nil {
			return requestDetailResponse, err
		}

		requestDetailResponse.ParentRequest = &responses.RequestResponse{}
		if err := utils.Mapper(parentRequest, &requestDetailResponse.ParentRequest); err != nil {
			return requestDetailResponse, err
		}

		if request.SprintID != nil {
			requestDetailResponse.ParentRequest.SprintId = *request.SprintID
		}

		requestDetailResponse.ParentKey = parentRequest.Key
	}

	// Participants
	userIds := []int32{}
	for _, node := range request.Nodes {
		userIds = append(userIds, *node.AssigneeID)
	}

	users, err := s.UserAPI.FindUsersByUserIds(userIds)
	if err != nil {
		return requestDetailResponse, err
	}

	participants := []types.Assignee{}
	existingUserIds := make(map[int32]bool)

	for _, user := range users.Data {
		if _, exists := existingUserIds[user.ID]; exists {
			continue
		}

		participant := types.Assignee{}
		if err := utils.Mapper(user, &participant); err != nil {
			return requestDetailResponse, err
		}

		participants = append(participants, participant)
		existingUserIds[user.ID] = true
	}
	requestDetailResponse.Participants = participants

	// Workflow
	workflowResponse := responses.WorkflowResponse{}
	if err := utils.Mapper(request.Workflow, &workflowResponse); err != nil {
		return requestDetailResponse, err
	}
	requestDetailResponse.Workflow = workflowResponse

	categoryResponse := responses.CategoryResponse{}
	if err := utils.Mapper(request.Category, &categoryResponse); err != nil {
		return requestDetailResponse, err
	}
	requestDetailResponse.Workflow.Category = categoryResponse

	// Childen Request
	childrenRequests, err := s.RequestRepo.FindAllChildrenRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return requestDetailResponse, err
	}

	childrenRequestsResponse := []responses.RequestResponse{}
	for _, childRequest := range childrenRequests {
		childRequestResponse := responses.RequestResponse{}
		if err := utils.Mapper(childRequest, &childRequestResponse); err != nil {
			return requestDetailResponse, err
		}

		childrenRequestsResponse = append(childrenRequestsResponse, childRequestResponse)
	}
	requestDetailResponse.ChildRequests = childrenRequestsResponse

	// RequestedBy
	requestedBy, err := s.UserAPI.FindUsersByUserIds([]int32{request.UserID})
	if err != nil {
		return requestDetailResponse, err
	}
	requestedByResponse := types.Assignee{}
	if err := utils.Mapper(requestedBy.Data[0], &requestedByResponse); err != nil {
		return requestDetailResponse, err
	}
	requestDetailResponse.RequestedBy = requestedByResponse

	return requestDetailResponse, nil
}

func (s *RequestService) GetRequestTasksHandler(ctx context.Context, requestId int32, requestTaskQueryParam queryparams.RequestTaskQueryParam) (responses.Paginate[[]responses.RequestTaskResponse], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestTaskResponse]{}
	requestTaskResponse := []responses.RequestTaskResponse{}

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return paginatedResponse, err
	}

	total, err := s.NodeRepo.CountAllNodeByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return paginatedResponse, err
	}

	nodes := []model.Nodes{}

	// Only unique existingUserIds
	existingUserIds := make(map[int32]bool)
	userIds := []int32{}

	// Get all nodes from request
	for _, node := range request.Nodes {

		// Skip Start and End Node
		if node.Type == string(constants.NodeTypeStart) || node.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// If Node is Story or SubWorkflow, get all nodes from subRequest
		if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
			subRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *node.SubRequestID)
			if err != nil {
				return paginatedResponse, err
			}

			for _, subNode := range subRequest.Nodes {
				// Only unique userIds
				if _, exists := existingUserIds[*subNode.AssigneeID]; !exists {
					existingUserIds[*subNode.AssigneeID] = true
					userIds = append(userIds, *subNode.AssigneeID)
				}

				// Append node
				subNodeModel := model.Nodes{}
				if err := utils.Mapper(subNode, &subNodeModel); err != nil {
					return paginatedResponse, err
				}
				nodes = append(nodes, subNodeModel)
			}
		} else {
			// Only unique userIds
			if _, exists := existingUserIds[*node.AssigneeID]; !exists {
				existingUserIds[*node.AssigneeID] = true
				userIds = append(userIds, *node.AssigneeID)
			}

			// Append node
			nodeModel := model.Nodes{}
			if err := utils.Mapper(node, &nodeModel); err != nil {
				return paginatedResponse, err
			}
			nodes = append(nodes, nodeModel)
		}
	}

	// Get all users from userIds
	users, err := s.UserAPI.FindUsersByUserIds(userIds)
	if err != nil {
		return paginatedResponse, err
	}

	// Create a map of assignees
	assignees := make(map[int32]types.Assignee)
	for _, user := range users.Data {
		assignees[user.ID] = types.Assignee{
			Id:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			AvatarUrl:    user.AvatarUrl,
			IsSystemUser: user.IsSystemUser,
		}
	}

	for _, node := range nodes {
		requestTask := responses.RequestTaskResponse{
			Id:               node.ID,
			Title:            node.Title,
			Status:           node.Status,
			PlannedStartTime: node.PlannedStartTime,
			PlannedEndTime:   node.PlannedEndTime,
			ActualStartTime:  node.ActualStartTime,
			ActualEndTime:    node.ActualEndTime,
			EstimatePoint:    node.EstimatePoint,
			RequestProgress:  request.Progress,
			RequestTitle:     request.Title,
			RequestID:        request.ID,
			Assignee:         assignees[*node.AssigneeID],
		}

		// Task Key
		if node.JiraKey != nil {
			requestTask.Key = *node.JiraKey
		} else {
			requestTask.Key = strconv.Itoa(int(node.Key))
		}

		requestTaskResponse = append(requestTaskResponse, requestTask)
	}

	totalPages := (int(total.Count) + requestTaskQueryParam.PageSize - 1) / requestTaskQueryParam.PageSize

	paginatedResponse = responses.Paginate[[]responses.RequestTaskResponse]{
		Items:      requestTaskResponse,
		Total:      int(total.Count),
		Page:       requestTaskQueryParam.Page,
		PageSize:   requestTaskQueryParam.PageSize,
		TotalPages: totalPages,
	}
	return paginatedResponse, nil
}

func (s *RequestService) GetRequestTasksByProjectHandler(ctx context.Context, requestTaskProjectQueryParam queryparams.RequestTaskProjectQueryParam, userId int32) (responses.Paginate[[]responses.RequestTaskResponse], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestTaskResponse]{}
	requestTaskResponse := []responses.RequestTaskResponse{}

	count, err := s.RequestRepo.CountRequestByStatusAndUserId(ctx, s.DB, userId, "")
	if err != nil {
		return paginatedResponse, err
	}

	total := count

	tasks, err := s.RequestRepo.FindAllTasksByProject(ctx, s.DB, userId, requestTaskProjectQueryParam)
	if err != nil {
		return paginatedResponse, err
	}

	// Only unique existingUserIds
	nodes := []results.NodeResult{}
	existingUserIds := make(map[int32]bool)
	userIds := []int32{}

	// Get all nodes from request
	for _, node := range tasks {

		// Skip Start and End Node
		if node.Type == string(constants.NodeTypeStart) || node.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// If Node is Story or SubWorkflow, get all nodes from subRequest
		if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
			subRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *node.SubRequestID)
			if err != nil {
				return paginatedResponse, err
			}

			for _, subNode := range subRequest.Nodes {
				// Only unique userIds
				if _, exists := existingUserIds[*subNode.AssigneeID]; !exists {
					existingUserIds[*subNode.AssigneeID] = true
					userIds = append(userIds, *subNode.AssigneeID)
				}

				// Append node
				subNodeModel := results.NodeResult{}
				if err := utils.Mapper(subNode, &subNodeModel); err != nil {
					return paginatedResponse, err
				}
				nodes = append(nodes, subNodeModel)
			}
		} else {
			// Only unique userIds
			if _, exists := existingUserIds[*node.AssigneeID]; !exists {
				existingUserIds[*node.AssigneeID] = true
				userIds = append(userIds, *node.AssigneeID)
			}

			// Append node
			nodeModel := results.NodeResult{}
			if err := utils.Mapper(node, &nodeModel); err != nil {
				return paginatedResponse, err
			}
			nodes = append(nodes, nodeModel)
		}
	}

	// Get all users from userIds
	users, err := s.UserAPI.FindUsersByUserIds(userIds)
	if err != nil {
		return paginatedResponse, err
	}

	// Create a map of assignees
	assignees := make(map[int32]types.Assignee)
	for _, user := range users.Data {
		assignees[user.ID] = types.Assignee{
			Id:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			AvatarUrl:    user.AvatarUrl,
			IsSystemUser: user.IsSystemUser,
		}
	}

	for _, node := range nodes {

		requestTask := responses.RequestTaskResponse{
			Id:               node.ID,
			Title:            node.Title,
			Status:           node.Status,
			Type:             node.Type,
			PlannedStartTime: node.PlannedStartTime,
			PlannedEndTime:   node.PlannedEndTime,
			ActualStartTime:  node.ActualStartTime,
			ActualEndTime:    node.ActualEndTime,
			EstimatePoint:    node.EstimatePoint,
			RequestProgress:  node.Request.Progress,
			RequestTitle:     node.Request.Title,
			Assignee:         assignees[*node.AssigneeID],
			IsCurrent:        node.IsCurrent,
		}

		if node.JiraKey != nil {
			requestTask.Key = *node.JiraKey
		} else {
			requestTask.Key = strconv.Itoa(int(node.Key))
		}

		requestTaskResponse = append(requestTaskResponse, requestTask)
	}

	totalPages := (int(total) + requestTaskProjectQueryParam.PageSize - 1) / requestTaskProjectQueryParam.PageSize

	paginatedResponse = responses.Paginate[[]responses.RequestTaskResponse]{
		Items:      requestTaskResponse,
		Total:      int(total),
		Page:       requestTaskProjectQueryParam.Page,
		PageSize:   requestTaskProjectQueryParam.PageSize,
		TotalPages: totalPages,
	}

	return paginatedResponse, nil
}

func (s *RequestService) GetRequestTaskCount(ctx context.Context, userId int32, queryParams queryparams.RequestTaskCount) (responses.RequestTaskCountResponse, error) {
	taskCountResponse := responses.RequestTaskCountResponse{}

	totalCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, "", queryParams)
	if err != nil {
		return taskCountResponse, err
	}
	taskCountResponse.TotalCount = int32(totalCount)

	completedCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, constants.NodeStatusCompleted, queryParams)
	if err != nil {
		return taskCountResponse, err
	}
	taskCountResponse.CompletedCount = int32(completedCount)

	overdueCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, constants.NodeStatusOverDue, queryParams)
	if err != nil {
		return taskCountResponse, err
	}
	taskCountResponse.OverdueCount = int32(overdueCount)

	todoCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, constants.NodeStatusTodo, queryParams)
	if err != nil {
		return taskCountResponse, err
	}
	taskCountResponse.TodoCount = int32(todoCount)

	return taskCountResponse, nil
}

func (s *RequestService) GetRequestOverviewHandler(ctx context.Context, requestId int32) (responses.RequestOverviewResponse, error) {
	requestOverviewResponse := responses.RequestOverviewResponse{}
	workflowRequest, err := s.WorkflowService.FindOneWorkflowDetailHandler(ctx, requestId)
	if err != nil {
		return requestOverviewResponse, err
	}

	if err := utils.Mapper(workflowRequest, &requestOverviewResponse); err != nil {
		return requestOverviewResponse, err
	}

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return requestOverviewResponse, err
	}

	requestOverviewResponse.Progress = request.Progress

	return requestOverviewResponse, nil
}

func (s *RequestService) FindAllSubRequestByRequestId(ctx context.Context, requestId int32, requestSubRequestQueryParam queryparams.RequestSubRequestQueryParam) (responses.Paginate[[]responses.RequestSubRequest], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestSubRequest]{}

	total, requests, err := s.RequestRepo.FindAllSubRequestByParentId(ctx, s.DB, requestId, requestSubRequestQueryParam)
	if err != nil {
		return paginatedResponse, fmt.Errorf("find all children request by request id fail: %w", err)
	}

	subRequests := []responses.RequestSubRequest{}
	for _, request := range requests {

		assignee := types.Assignee{}
		if request.Nodes.AssigneeID != nil {
			users, err := s.UserAPI.FindUsersByUserIds([]int32{*request.Nodes.AssigneeID})
			if err != nil {
				return paginatedResponse, fmt.Errorf("find users by user ids fail: %w", err)
			}

			if err := utils.Mapper(users.Data[0], &assignee); err != nil {
				return paginatedResponse, fmt.Errorf("mapper assignee fail: %w", err)
			}
		}

		subRequests = append(subRequests, responses.RequestSubRequest{
			Key:           request.ID,
			WorkflowTitle: request.Workflows.Title,
			TaskTitle:     request.Title,
			Assignee:      assignee,
			Status:        request.Status,
			StartedAt:     request.StartedAt,
			CompletedAt:   request.CompletedAt,
			CanceledAt:    request.CanceledAt,
			TerminatedAt:  request.TerminatedAt,
		})
	}

	totalPages := (int(total) + requestSubRequestQueryParam.PageSize - 1) / requestSubRequestQueryParam.PageSize

	paginatedResponse = responses.Paginate[[]responses.RequestSubRequest]{
		Items:      subRequests,
		Total:      int(total),
		Page:       requestSubRequestQueryParam.Page,
		PageSize:   requestSubRequestQueryParam.PageSize,
		TotalPages: totalPages,
	}

	return paginatedResponse, nil
}

func (s *RequestService) UpdateRequestHandler(ctx context.Context, requestId int32, req *requests.RequestUpdateRequest, userId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx fail: %w", err)
	}
	defer tx.Rollback()

	// Remove Nodes, Connections and Stories
	err = s.RequestRepo.RemoveNodesConnectionsStoriesByRequestId(ctx, tx, requestId)
	if err != nil {
		return fmt.Errorf("remove nodes connections stories fail: %w", err)
	}

	// Create Nodes, Connections and Stories
	nodesConnectionsStories := requests.NodesConnectionsStories{
		Nodes:       req.Nodes,
		Connections: req.Connections,
		Stories:     req.Stories,
	}
	err = s.WorkflowService.CreateNodesConnectionsStories(ctx, tx, &nodesConnectionsStories, requestId, nil, userId)
	if err != nil {
		return fmt.Errorf("create nodes connections stories fail: %w", err)
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}
