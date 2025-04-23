package services

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

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
	ConnectionRepo  *repositories.ConnectionRepository
	WorkflowService *WorkflowService
	NatsService     *NatsService
	NodeService     *NodeService
	FormService     *FormService
	FormRepo        *repositories.FormRepository
	HistoryRepo     *repositories.HistoryRepository
	HistoryService  *HistoryService
}

func NewRequestService(cfg RequestService) *RequestService {
	return &RequestService{
		DB:              cfg.DB,
		RequestRepo:     cfg.RequestRepo,
		UserAPI:         cfg.UserAPI,
		NodeRepo:        cfg.NodeRepo,
		ConnectionRepo:  cfg.ConnectionRepo,
		WorkflowService: cfg.WorkflowService,
		NatsService:     cfg.NatsService,
		NodeService:     cfg.NodeService,
		FormService:     cfg.FormService,
		FormRepo:        cfg.FormRepo,
		HistoryRepo:     cfg.HistoryRepo,
		HistoryService:  cfg.HistoryService,
	}
}

func (s *RequestService) UpdateCalculateRequestProgress(ctx context.Context, tx *sql.Tx, requestId int32) error {
	request, err := s.RequestRepo.FindOneRequestByRequestIdTx(ctx, tx, requestId)
	if err != nil {
		return fmt.Errorf("find request by request id fail: %w", err)
	}

	totalCompletedNode := 0
	totalNode := len(request.Nodes)
	for _, requestNode := range request.Nodes {
		if requestNode.Type == string(constants.NodeTypeStart) || requestNode.Type == string(constants.NodeTypeEnd) || requestNode.Type == string(constants.NodeTypeCondition) {
			totalNode--
		} else if requestNode.Status == string(constants.NodeStatusCompleted) {
			totalCompletedNode++
		}
	}

	if totalNode == 0 {
		request.Progress = 100
	} else {
		request.Progress = float32(float64(totalCompletedNode) / float64(totalNode) * 100)
	}
	requestModel := model.Requests{}
	if err := utils.Mapper(request, &requestModel); err != nil {
		return err
	}

	if err := s.RequestRepo.UpdateRequest(ctx, tx, requestModel); err != nil {
		return err
	}

	// Update Parent Request Count Progress
	if request.ParentID != nil {
		if err := s.UpdateCalculateRequestProgress(ctx, tx, *request.ParentID); err != nil {
			return err
		}
	}

	return nil
}

//

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

		// Project Key
		if request.Workflow.ProjectKey != nil {
			requestResponse.ProjectKey = request.Workflow.ProjectKey
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
		if node.Type == string(constants.NodeTypeStart) || node.Type == string(constants.NodeTypeEnd) || node.Type == string(constants.NodeTypeCondition) {
			continue
		}
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
		if node.Type == string(constants.NodeTypeStart) || node.Type == string(constants.NodeTypeEnd) || node.Type == string(constants.NodeTypeCondition) {
			continue
		}

		// If Node is Story or SubWorkflow, get all nodes from subRequest
		if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
			subRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *node.SubRequestID)
			if err != nil {
				return paginatedResponse, err
			}

			for _, subNode := range subRequest.Nodes {
				if subNode.Type == string(constants.NodeTypeStart) || subNode.Type == string(constants.NodeTypeEnd) {
					continue
				}
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
			IsApproved:       node.IsApproved,
			IsRejected:       node.IsRejected,
			ProjectKey:       request.Workflow.ProjectKey,
			JiraLinkUrl:      node.JiraLinkURL,
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

	total, tasks, err := s.RequestRepo.FindAllTasksByProject(ctx, s.DB, userId, requestTaskProjectQueryParam)
	if err != nil {
		return paginatedResponse, err
	}

	// Only unique existingUserIds
	nodes := []results.NodeResult{}
	existingUserIds := make(map[int32]bool)
	userIds := []int32{}

	// Get all nodes from request
	for _, node := range tasks {

		// If Node is Story or SubWorkflow, get all nodes from subRequest
		if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
			if node.Status == string(constants.NodeStatusTodo) {
				total--
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
			RequestID:        node.RequestID,
			IsApproved:       node.IsApproved,
			IsRejected:       node.IsRejected,
			ProjectKey:       node.Workflows.ProjectKey,
			JiraLinkUrl:      node.JiraLinkURL,
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

	inProcessCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, constants.NodeStatusInProgress, queryParams)
	if err != nil {
		return taskCountResponse, err
	}
	taskCountResponse.InProgressCount = int32(inProcessCount)

	taskCountResponse.TotalCount = taskCountResponse.CompletedCount + taskCountResponse.OverdueCount + taskCountResponse.TodoCount + int32(inProcessCount)

	todayCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, constants.NodeStatusToday, queryParams)
	if err != nil {
		return taskCountResponse, err
	}
	taskCountResponse.TodayCount = int32(todayCount)

	inComingCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, constants.NodeStatusInComing, queryParams)
	if err != nil {
		return taskCountResponse, err
	}
	taskCountResponse.InComingCount = int32(inComingCount)

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
	requestOverviewResponse.Category = responses.CategoryFindAll{
		ID:       request.Category.ID,
		Name:     request.Category.Name,
		Type:     request.Category.Type,
		Key:      request.Category.Key,
		IsActive: request.Category.IsActive,
	}

	return requestOverviewResponse, nil
}

func (s *RequestService) FindAllSubRequestByRequestId(ctx context.Context, requestId int32, requestSubRequestQueryParam queryparams.RequestSubRequestQueryParam) (responses.Paginate[[]responses.RequestResponse], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestResponse]{}
	requestsResponse := []responses.RequestResponse{}
	paginatedResponse.Items = requestsResponse

	total, requests, err := s.RequestRepo.FindAllSubRequestByParentId(ctx, s.DB, requestId, requestSubRequestQueryParam)
	if err != nil {
		errStr := err.Error()
		if errStr == "qrm: no rows in result set" {
			return paginatedResponse, nil
		}
		return paginatedResponse, fmt.Errorf("find all children request by request id fail: %w", err)
	}

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

	totalPages := (int(total) + requestSubRequestQueryParam.PageSize - 1) / requestSubRequestQueryParam.PageSize

	paginatedResponse = responses.Paginate[[]responses.RequestResponse]{
		Items:      requestsResponse,
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

	originalRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return fmt.Errorf("find one original request by request id fail: %w", err)
	}

	// Get original nodes and connections for Jira sync
	origNodes := originalRequest.Nodes
	origConnections := originalRequest.Connections

	// Remove existing nodes/connections for the main request and its original sub-requests first
	for _, node := range originalRequest.Nodes {
		if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
			if node.SubRequestID != nil {
				// Remove Nodes, Connections and Stories for the original sub-request
				err = s.RequestRepo.RemoveNodesConnectionsStoriesByRequestId(ctx, tx, *node.SubRequestID)
				if err != nil {
					// Consider if failing to remove an old sub-request should halt the whole update
					// For now, we return the error.
					return fmt.Errorf("remove nodes connections stories for original subrequest %d fail: %w", *node.SubRequestID, err)
				}
			}
		}
	}

	// Remove Nodes, Connections and Stories for the main request
	err = s.RequestRepo.RemoveNodesConnectionsStoriesByRequestId(ctx, tx, requestId)
	if err != nil {
		return fmt.Errorf("remove nodes connections stories for main request %d fail: %w", requestId, err)
	}

	// --- Start Modification ---
	// Prepare slices to accumulate nodes and connections from sub-requests
	subNodesToAdd := []requests.Node{}
	subConnectionsToAdd := []requests.Connection{}

	// Iterate through the incoming nodes to find sub-workflows
	for _, node := range req.Nodes {
		if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
			// Ensure Data and SubRequestID are present
			if node.Data.SubRequestID != nil {
				subRequestID := *node.Data.SubRequestID

				subRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, subRequestID)
				if err != nil {
					return fmt.Errorf("find one sub-request by request id fail: %w", err)
				}

				// Fetch nodes for the sub-request
				subReqNodesModel := subRequest.Nodes

				// Fetch connections for the sub-request
				subReqConnectionsModel := subRequest.Connections

				// Transform and append nodes
				for _, subNodeModel := range subReqNodesModel {
					var subNodeReq requests.Node
					// Assuming utils.Mapper can handle this conversion. Adjust if manual mapping is needed.
					if err := utils.Mapper(subNodeModel, &subNodeReq); err != nil {
						return fmt.Errorf("map sub-request node %s fail: %w", subNodeModel.ID, err)
					}
					// Important: Associate these nodes with the main request ID for creation
					// subNodeReq.RequestID = requestId // Ensure CreateNodesConnectionsStories handles this or set it here if needed.
					subNodesToAdd = append(subNodesToAdd, subNodeReq)
				}

				// Transform and append connections
				for _, subConnModel := range subReqConnectionsModel {
					var subConnReq requests.Connection
					// Assuming utils.Mapper can handle this conversion. Adjust if manual mapping is needed.
					if err := utils.Mapper(subConnModel, &subConnReq); err != nil {
						return fmt.Errorf("map sub-request connection fail: %w", err) // Connections might not have an ID, use index or other identifier if needed for error msg
					}
					// Important: Associate these connections with the main request ID for creation
					// subConnReq.RequestID = requestId // Ensure CreateNodesConnectionsStories handles this or set it here if needed.
					subConnectionsToAdd = append(subConnectionsToAdd, subConnReq)
				}

				// Note: We removed the redundant `CreateRequest` logic here.
				// The original node of type Story/SubWorkflow in `req.Nodes` might still be needed
				// by CreateNodesConnectionsStories to establish linkage or context.
				// We are *adding* the content of the sub-workflow, not replacing the node itself.
			} else {
				// Handle cases where a Story/SubWorkflow node is missing Data or SubRequestID if necessary
				// For now, we just skip it.
				fmt.Printf("Warning: Node type %s is missing Data or SubRequestID.\n", node.Type)
			}
		}
	}

	// Combine original nodes/connections with those from sub-requests
	finalNodes := append(req.Nodes, subNodesToAdd...)
	finalConnections := append(req.Connections, subConnectionsToAdd...)

	nodesConnectionsStories := requests.NodesConnectionsStories{
		Nodes:       finalNodes,       // Use combined nodes
		Connections: finalConnections, // Use combined connections
		Stories:     req.Stories,      // Assuming stories are not nested this way
	}
	// --- End Modification ---

	// Create the new structure using the combined nodes and connections
	// Pass the originalRequest.WorkflowVersionID, not the potentially different one from the fetched subRequest
	err = s.WorkflowService.CreateNodesConnectionsStories(ctx, tx, &nodesConnectionsStories, requestId, originalRequest.Workflow.ProjectKey, userId, originalRequest.SprintID, false)
	if err != nil {
		return fmt.Errorf("create nodes connections stories fail: %w", err)
	}

	// ================================ SYNC WITH JIRA ================================
	// Sync with Jira if this is a project workflow with project key
	if originalRequest.Workflow.Type == string(constants.WorkflowTypeProject) && originalRequest.Workflow.ProjectKey != nil {
		// Get the NatsService from WorkflowService

		// Need to convert the original nodes and connections to the proper types
		var modelNodes []model.Nodes
		var modelConnections []model.Connections

		// Convert original nodes to model.Nodes
		for _, node := range origNodes {
			var modelNode model.Nodes
			if err := utils.Mapper(node, &modelNode); err != nil {
				slog.Error("Failed to map original node", "error", err)
				continue
			}
			modelNodes = append(modelNodes, modelNode)
		}

		// Convert original connections to model.Connections
		for _, conn := range origConnections {
			var modelConn model.Connections
			if err := utils.Mapper(conn, &modelConn); err != nil {
				slog.Error("Failed to map original connection", "error", err)
				continue
			}
			modelConnections = append(modelConnections, modelConn)
		}

		// Sync the updated workflow with Jira using edit mode
		_, err := s.NatsService.PublishWorkflowEditToJira(ctx, tx, requestId, finalNodes, modelNodes, req.Stories,
			finalConnections, modelConnections, *originalRequest.Workflow.ProjectKey, originalRequest.SprintID)
		if err != nil {
			// Log the error but continue - we don't want to fail the update if Jira sync fails
			slog.Error("Failed to sync workflow edit with Jira", "error", err)
		}
	}

	// ================================ END SYNC WITH JIRA ================================

	// history
	if originalRequest.Status == string(constants.RequestStatusInProgress) {
		for _, node := range originalRequest.Nodes {
			if node.Type == string(constants.NodeTypeStart) {
				err = s.HistoryService.HistoryEditRequest(ctx, requestId, node.ID)
				if err != nil {
					return fmt.Errorf("history edit request fail: %w", err)
				}
				break
			}
		}
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *RequestService) GetRequestCompletedFormHandler(ctx context.Context, requestId int32, queryParams queryparams.RequestTaskQueryParam) (responses.Paginate[[]responses.RequestCompletedFormInputResponse], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestCompletedFormInputResponse]{}
	requestCompletedFormResponse := []responses.RequestCompletedFormInputResponse{}
	paginatedResponse.Items = requestCompletedFormResponse

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return paginatedResponse, err
	}

	total := 0

	if request.Workflow.Type == string(constants.WorkflowTypeProject) {
		for _, node := range request.Nodes {
			if node.Type == string(constants.NodeTypeTask) || node.Type == string(constants.NodeTypeStory) {
				formData, err := s.NodeRepo.FindOneFormDataByNodeId(ctx, s.DB, node.ID)
				if err != nil {
					return paginatedResponse, err
				}

				total = len(formData)

				users, err := s.UserAPI.FindUsersByUserIds([]int32{*node.AssigneeID})
				if err != nil {
					return paginatedResponse, err
				}

				assignee := types.Assignee{}
				if err := utils.Mapper(users.Data[0], &assignee); err != nil {
					return paginatedResponse, err
				}

				formTemplateSystem, err := s.FormRepo.FindOneFormTemplateByFormTemplateId(ctx, s.DB, constants.FormTemplateIDJiraSystemForm)
				if err != nil {
					return paginatedResponse, err
				}

				formFieldMap := make(map[int32]string)
				for _, formTemplateField := range formTemplateSystem.Fields {
					formFieldMap[formTemplateField.ID] = formTemplateField.FieldID
				}

				formDataResponse := []responses.RequestCompletedFormDataResponse{}
				for _, form := range formData {
					for _, formFieldData := range form.FormFieldData {
						formDataResponse = append(formDataResponse, responses.RequestCompletedFormDataResponse{
							FieldID: formFieldMap[formFieldData.FormTemplateFieldID],
							Value:   formFieldData.Value,
						})
					}
				}

				formTemplate, err := s.FormService.FindOneFormTemplateDetailByFormTemplateId(ctx, constants.FormTemplateIDJiraSystemForm)
				if err != nil {
					return paginatedResponse, err
				}

				requestCompletedFormRes := responses.RequestCompletedFormInputResponse{
					SubmittedAt: node.UpdatedAt,
					Type:        node.Type,
					FormData:    formDataResponse,
					Approval:    []responses.RequestCompletedFormApprovalResponse{},
					Template:    formTemplate,
					Submitter:   assignee,
					LastUpdate:  assignee,
					DataId:      node.FormDataID,
				}

				if node.JiraKey != nil {
					requestCompletedFormRes.Key = *node.JiraKey
				} else {
					requestCompletedFormRes.Key = strconv.Itoa(int(node.Key))
				}

				requestCompletedFormResponse = append(requestCompletedFormResponse, requestCompletedFormRes)
			}

			if node.Type == string(constants.NodeTypeSubWorkflow) || node.Type == string(constants.NodeTypeStory) {
				subRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *node.SubRequestID)
				if err != nil {
					return paginatedResponse, err
				}

				for _, subNode := range subRequest.Nodes {
					if subNode.Type == string(constants.NodeTypeTask) || subNode.Type == string(constants.NodeTypeStory) {
						formData, err := s.NodeRepo.FindOneFormDataByNodeId(ctx, s.DB, subNode.ID)
						if err != nil {
							return paginatedResponse, err
						}

						total = len(formData)

						users, err := s.UserAPI.FindUsersByUserIds([]int32{*subNode.AssigneeID})
						if err != nil {
							return paginatedResponse, err
						}

						assignee := types.Assignee{}
						if err := utils.Mapper(users.Data[0], &assignee); err != nil {
							return paginatedResponse, err
						}

						formTemplateSystem, err := s.FormRepo.FindOneFormTemplateByFormTemplateId(ctx, s.DB, constants.FormTemplateIDJiraSystemForm)
						if err != nil {
							return paginatedResponse, err
						}

						formFieldMap := make(map[int32]string)
						for _, formTemplateField := range formTemplateSystem.Fields {
							formFieldMap[formTemplateField.ID] = formTemplateField.FieldID
						}

						formDataResponse := []responses.RequestCompletedFormDataResponse{}
						for _, form := range formData {
							for _, formFieldData := range form.FormFieldData {
								formDataResponse = append(formDataResponse, responses.RequestCompletedFormDataResponse{
									FieldID: formFieldMap[formFieldData.FormTemplateFieldID],
									Value:   formFieldData.Value,
								})
							}
						}

						formTemplate, err := s.FormService.FindOneFormTemplateDetailByFormTemplateId(ctx, constants.FormTemplateIDJiraSystemForm)
						if err != nil {
							return paginatedResponse, err
						}

						requestCompletedFormRes := responses.RequestCompletedFormInputResponse{
							SubmittedAt: subNode.UpdatedAt,
							Type:        subNode.Type,
							FormData:    formDataResponse,
							Approval:    []responses.RequestCompletedFormApprovalResponse{},
							Template:    formTemplate,
							Submitter:   assignee,
							LastUpdate:  assignee,
							DataId:      subNode.FormDataID,
						}

						if subNode.JiraKey != nil {
							requestCompletedFormRes.Key = *subNode.JiraKey
						} else {
							requestCompletedFormRes.Key = strconv.Itoa(int(subNode.Key))
						}

						requestCompletedFormResponse = append(requestCompletedFormResponse, requestCompletedFormRes)
					}
				}
			}
		}
	} else {
		count, nodeForms, err := s.RequestRepo.FindAllRequestCompletedFormByRequestId(ctx, s.DB, requestId, queryParams.Page, queryParams.PageSize)
		total = int(count)
		if err != nil {
			return paginatedResponse, err
		}

		userIds := []int32{}
		existUserIds := make(map[int32]bool)

		for _, nodeForm := range nodeForms {
			if nodeForm.Permission != string(constants.NodeFormPermissionInput) {
				continue
			}

			if nodeForm.SubmittedByUserID != nil && !existUserIds[*nodeForm.SubmittedByUserID] {
				userIds = append(userIds, *nodeForm.SubmittedByUserID)
				existUserIds[*nodeForm.SubmittedByUserID] = true
			}
			if nodeForm.LastUpdateUserID != nil && !existUserIds[*nodeForm.LastUpdateUserID] {
				userIds = append(userIds, *nodeForm.LastUpdateUserID)
				existUserIds[*nodeForm.LastUpdateUserID] = true
			}
			for _, approveOrRejectUser := range nodeForm.ApproveOrRejectUsers {
				if approveOrRejectUser.UserID != 0 && !existUserIds[approveOrRejectUser.UserID] {
					userIds = append(userIds, approveOrRejectUser.UserID)
					existUserIds[approveOrRejectUser.UserID] = true
				}
			}
		}

		userApiMap := map[int32]results.UserApiDataResult{}
		if len(userIds) > 0 {
			assigneeResult, err := s.UserAPI.FindUsersByUserIds(userIds)
			if err != nil {
				return paginatedResponse, err
			}
			for _, userApi := range assigneeResult.Data {
				userApiMap[userApi.ID] = userApi
			}
		}

		mapUser := func(id *int32) types.Assignee {
			assignee := types.Assignee{}
			if id != nil {
				if user, ok := userApiMap[*id]; ok {
					assignee.Id = user.ID
					assignee.Name = user.Name
					assignee.Email = user.Email
					assignee.AvatarUrl = user.AvatarUrl
					assignee.IsSystemUser = user.IsSystemUser
				}
			}
			return assignee
		}

		//
		for _, nodeForm := range nodeForms {
			requestCompletedFormRes := responses.RequestCompletedFormInputResponse{}
			if nodeForm.SubmittedAt != nil {
				requestCompletedFormRes.Submitter = mapUser(nodeForm.SubmittedByUserID)
			}
			if nodeForm.LastUpdateUserID != nil {
				requestCompletedFormRes.LastUpdate = mapUser(nodeForm.LastUpdateUserID)
			}

			requestCompletedFormRes.Type = nodeForm.Node.Type

			if nodeForm.Node.JiraKey != nil {
				requestCompletedFormRes.Key = *nodeForm.Node.JiraKey
			} else {
				requestCompletedFormRes.Key = strconv.Itoa(int(nodeForm.Node.Key))
			}

			//
			approval := []responses.RequestCompletedFormApprovalResponse{}
			for _, approveOrRejectUser := range nodeForm.ApproveOrRejectUsers {
				approval = append(approval, responses.RequestCompletedFormApprovalResponse{
					Assignee:   mapUser(&approveOrRejectUser.UserID),
					IsApproved: approveOrRejectUser.IsApproved,
				})
			}
			requestCompletedFormRes.Approval = approval

			formTemplateNodeForm, err := s.FormRepo.FindOneFormTemplateByFormTemplateId(ctx, s.DB, nodeForm.FormData.FormTemplate.ID)
			if err != nil {
				return paginatedResponse, err
			}

			fieldMap := make(map[int32]string)
			for _, formTemplateField := range formTemplateNodeForm.Fields {
				fieldMap[formTemplateField.ID] = formTemplateField.FieldID
			}

			//
			formData := []responses.RequestCompletedFormDataResponse{}
			for _, formFieldData := range nodeForm.FormData.FormFieldData {
				formData = append(formData, responses.RequestCompletedFormDataResponse{
					FieldID: fieldMap[formFieldData.FormTemplateFieldID],
					Value:   formFieldData.Value,
				})
			}
			requestCompletedFormRes.FormData = formData

			//
			formTemplate, err := s.FormService.FindOneFormTemplateDetailByFormTemplateId(ctx, nodeForm.TemplateID)
			if err != nil {
				return paginatedResponse, err
			}
			requestCompletedFormRes.Template = formTemplate

			requestCompletedFormRes.DataId = nodeForm.DataID

			requestCompletedFormResponse = append(requestCompletedFormResponse, requestCompletedFormRes)
		}

	}

	totalPages := (int(total) + queryParams.PageSize - 1) / queryParams.PageSize

	paginatedResponse = responses.Paginate[[]responses.RequestCompletedFormInputResponse]{
		Items:      requestCompletedFormResponse,
		Total:      int(total),
		Page:       queryParams.Page,
		PageSize:   queryParams.PageSize,
		TotalPages: totalPages,
	}
	return paginatedResponse, nil
}

func (s *RequestService) GetRequestFileManagerHandler(ctx context.Context, requestId int32, queryParams queryparams.RequestTaskQueryParam) (responses.Paginate[[]responses.RequestFileManagerResponse], error) {
	paginatedResponse := responses.Paginate[[]responses.RequestFileManagerResponse]{}
	requestFileManagerResponse := []responses.RequestFileManagerResponse{}
	paginatedResponse.Items = requestFileManagerResponse

	total, nodeForms, err := s.RequestRepo.FindAllRequestFileManagerByRequestId(ctx, s.DB, requestId, queryParams.Page, queryParams.PageSize)
	if err != nil {
		return paginatedResponse, err
	}

	userIds := []int32{}
	existUserIds := make(map[int32]bool)

	for _, nodeForm := range nodeForms {
		if nodeForm.SubmittedByUserID != nil && !existUserIds[*nodeForm.SubmittedByUserID] {
			userIds = append(userIds, *nodeForm.SubmittedByUserID)
			existUserIds[*nodeForm.SubmittedByUserID] = true
		}
	}

	userApiMap := map[int32]results.UserApiDataResult{}
	if len(userIds) > 0 {
		assigneeResult, err := s.UserAPI.FindUsersByUserIds(userIds)
		if err != nil {
			return paginatedResponse, err
		}
		for _, userApi := range assigneeResult.Data {
			userApiMap[userApi.ID] = userApi
		}
	}

	mapUser := func(id *int32) types.Assignee {
		assignee := types.Assignee{}
		if id != nil {
			if user, ok := userApiMap[*id]; ok {
				assignee.Id = user.ID
				assignee.Name = user.Name
				assignee.Email = user.Email
				assignee.AvatarUrl = user.AvatarUrl
				assignee.IsSystemUser = user.IsSystemUser
			}
		}
		return assignee
	}

	for _, nodeForm := range nodeForms {
		requestFileManageRes := responses.RequestFileManagerResponse{
			SubmittedAt: *nodeForm.SubmittedAt,
			Assignee:    mapUser(nodeForm.SubmittedByUserID),
		}

		for _, formFieldData := range nodeForm.FormFieldData {
			requestFileManageRes.Data = append(requestFileManageRes.Data, formFieldData.Value)
		}

		requestFileManagerResponse = append(requestFileManagerResponse, requestFileManageRes)
	}

	totalPages := (int(total) + queryParams.PageSize - 1) / queryParams.PageSize

	paginatedResponse = responses.Paginate[[]responses.RequestFileManagerResponse]{
		Items:      requestFileManagerResponse,
		Total:      int(total),
		Page:       queryParams.Page,
		PageSize:   queryParams.PageSize,
		TotalPages: totalPages,
	}

	return paginatedResponse, nil
}

func (s *RequestService) GetRequestCompletedFormApprovalHandler(ctx context.Context, requestId int32, dataId string) ([]responses.RequestCompletedFormApprovalOverviewResponse, error) {
	requestCompletedFormApprovalResponse := []responses.RequestCompletedFormApprovalOverviewResponse{}

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return requestCompletedFormApprovalResponse, err
	}

	userIds := []int32{}
	existUserIds := make(map[int32]bool)

	for _, node := range request.Nodes {
		if node.AssigneeID != nil && !existUserIds[*node.AssigneeID] {
			userIds = append(userIds, *node.AssigneeID)
			existUserIds[*node.AssigneeID] = true
		}
	}

	userApiMap := map[int32]results.UserApiDataResult{}
	if len(userIds) > 0 {
		assigneeResult, err := s.UserAPI.FindUsersByUserIds(userIds)
		if err != nil {
			return requestCompletedFormApprovalResponse, err
		}
		for _, userApi := range assigneeResult.Data {
			userApiMap[userApi.ID] = userApi
		}
	}

	mapUser := func(id *int32) types.Assignee {
		assignee := types.Assignee{}
		if id != nil {
			if user, ok := userApiMap[*id]; ok {
				assignee.Id = user.ID
				assignee.Name = user.Name
				assignee.Email = user.Email
				assignee.AvatarUrl = user.AvatarUrl
				assignee.IsSystemUser = user.IsSystemUser
			}
		}
		return assignee
	}

	for _, node := range request.Nodes {
		isApproveNodeInlcudeNodeFormDataId := false

		for _, nodeForm := range node.NodeForms {
			if nodeForm.Permission == string(constants.NodeFormPermissionView) && nodeForm.DataID != nil && *nodeForm.DataID == dataId {
				isApproveNodeInlcudeNodeFormDataId = true
				break
			}
		}

		if isApproveNodeInlcudeNodeFormDataId {
			requestCompletedFormApprovalResponse = append(requestCompletedFormApprovalResponse, responses.RequestCompletedFormApprovalOverviewResponse{
				Key:        node.Key,
				TaskTitle:  node.Title,
				IsApproved: node.IsApproved,
				IsRejected: node.IsRejected,
				Assignee:   mapUser(node.AssigneeID),
			})
		}
	}

	return requestCompletedFormApprovalResponse, nil
}

func (s *RequestService) FindAllHistoryByRequestId(ctx context.Context, requestId int32) ([]responses.HistoryResponse, error) {

	historiesResponse := []responses.HistoryResponse{}

	histories, err := s.HistoryRepo.FindAllHistoryByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return historiesResponse, err
	}

	userIds := []int32{}
	existingUserIds := make(map[int32]bool)
	for _, history := range histories {
		if history.UserID != nil {
			if !existingUserIds[*history.UserID] {
				existingUserIds[*history.UserID] = true
				userIds = append(userIds, *history.UserID)
			}
		}

		if history.TypeAction == constants.HistoryTypeAssignee {
			if history.FromValue != nil {
				fromValueInt32, err := strconv.Atoi(*history.FromValue)
				if err != nil {
					return nil, err
				}
				if !existingUserIds[int32(fromValueInt32)] {
					existingUserIds[int32(fromValueInt32)] = true
					userIds = append(userIds, int32(fromValueInt32))
				}
			}

			toValueInt32, err := strconv.Atoi(*history.ToValue)
			if err != nil {
				return nil, err
			}
			if !existingUserIds[int32(toValueInt32)] {
				existingUserIds[int32(toValueInt32)] = true
				userIds = append(userIds, int32(toValueInt32))
			}
		} else if history.TypeAction == constants.HistoryTypeNewTask {
			toValueInt32, err := strconv.Atoi(*history.ToValue)
			if err != nil {
				return nil, err
			}
			if !existingUserIds[int32(toValueInt32)] {
				existingUserIds[int32(toValueInt32)] = true
				userIds = append(userIds, int32(toValueInt32))
			}
		}
	}

	userApiMap := map[int32]results.UserApiDataResult{}
	if len(userIds) > 0 {
		assigneeResult, err := s.UserAPI.FindUsersByUserIds(userIds)
		if err != nil {
			return historiesResponse, err
		}
		for _, userApi := range assigneeResult.Data {
			userApiMap[userApi.ID] = userApi
		}
	}

	mapUser := func(id *int32) types.Assignee {
		assignee := types.Assignee{}
		if id != nil {
			if user, ok := userApiMap[*id]; ok {
				assignee.Id = user.ID
				assignee.Name = user.Name
				assignee.Email = user.Email
				assignee.AvatarUrl = user.AvatarUrl
				assignee.IsSystemUser = user.IsSystemUser
			}
		}
		return assignee
	}

	for _, history := range histories {

		historyNodeResponse := responses.HistoryNodeResponse{
			ID:    history.Node.ID,
			Title: history.Node.Title,
		}

		if history.Node.JiraKey != nil {
			historyNodeResponse.Key = *history.Node.JiraKey
		} else {
			historyNodeResponse.Key = strconv.Itoa(int(history.Node.Key))
		}

		historyResponse := responses.HistoryResponse{
			ID:        history.ID,
			CreatedAt: history.CreatedAt,
			Assignee:  mapUser(history.UserID),
			Node:      historyNodeResponse,
			Type:      history.TypeAction,
		}

		if history.TypeAction == constants.HistoryTypeAssignee {
			if history.FromValue != nil {
				fromValueInt, err := strconv.Atoi(*history.FromValue)
				if err != nil {
					return nil, err
				}
				fromValueInt32 := int32(fromValueInt)

				historyResponse.From = mapUser(&fromValueInt32)
			}

			toValueInt, err := strconv.Atoi(*history.ToValue)
			if err != nil {
				return nil, err
			}
			toValueInt32 := int32(toValueInt)

			historyResponse.To = mapUser(&toValueInt32)

		} else if history.TypeAction == constants.HistoryTypeStatus {
			historyResponse.From = history.FromValue
			historyResponse.To = history.ToValue

		} else if history.TypeAction == constants.HistoryTypeApprove || history.TypeAction == constants.HistoryTypeReject {
			historyResponse.From = history.FromValue
			historyResponse.To = history.ToValue
		} else if history.TypeAction == constants.HistoryTypeNewTask {
			toValueInt, err := strconv.Atoi(*history.ToValue)
			if err != nil {
				return nil, err
			}
			toValueInt32 := int32(toValueInt)

			historyResponse.To = mapUser(&toValueInt32)
		}
		historiesResponse = append(historiesResponse, historyResponse)
	}

	return historiesResponse, nil
}

func (s *RequestService) ReportMidSprintTasks(ctx context.Context, queryParams queryparams.RequestMidSprintReportQueryParam) ([]responses.RequestTaskResponse, error) {
	requestTasksResponse := []responses.RequestTaskResponse{}

	requests, err := s.RequestRepo.FindAllTasksByMidSprintReport(ctx, s.DB, queryParams)
	if err != nil {
		return nil, err
	}

	userIds := []int32{}
	existUserIds := make(map[int32]bool)
	for _, request := range requests {
		for _, node := range request.Nodes {
			if node.AssigneeID != nil && !existUserIds[*node.AssigneeID] {
				existUserIds[*node.AssigneeID] = true
				userIds = append(userIds, *node.AssigneeID)
			}
		}
	}

	userApiMap := map[int32]results.UserApiDataResult{}
	if len(userIds) > 0 {
		assigneeResult, err := s.UserAPI.FindUsersByUserIds(userIds)
		if err != nil {
			return nil, err
		}
		for _, userApi := range assigneeResult.Data {
			userApiMap[userApi.ID] = userApi
		}
	}

	mapUser := func(id *int32) types.Assignee {
		assignee := types.Assignee{}
		if id != nil {
			if user, ok := userApiMap[*id]; ok {
				assignee.Id = user.ID
				assignee.Name = user.Name
				assignee.Email = user.Email
				assignee.AvatarUrl = user.AvatarUrl
				assignee.IsSystemUser = user.IsSystemUser
			}
		}
		return assignee
	}

	for _, request := range requests {
		for _, node := range request.Nodes {
			if node.Type == string(constants.NodeTypeTask) {
				taskRes := responses.RequestTaskResponse{
					Id:               node.ID,
					Title:            node.Title,
					Status:           node.Status,
					Type:             node.Type,
					RequestID:        request.ID,
					RequestTitle:     request.Title,
					RequestProgress:  request.Progress,
					ProjectKey:       request.Workflow.ProjectKey,
					JiraLinkUrl:      node.JiraLinkURL,
					Assignee:         mapUser(node.AssigneeID),
					PlannedStartTime: node.PlannedStartTime,
					PlannedEndTime:   node.PlannedEndTime,
					ActualStartTime:  node.ActualStartTime,
					ActualEndTime:    node.ActualEndTime,
					EstimatePoint:    node.EstimatePoint,
					IsCurrent:        node.IsCurrent,
					IsApproved:       node.IsApproved,
					IsRejected:       node.IsRejected,
				}

				if node.JiraKey != nil {
					taskRes.Key = *node.JiraKey
				} else {
					taskRes.Key = strconv.Itoa(int(node.Key))
				}

				requestTasksResponse = append(requestTasksResponse, taskRes)
			}

		}
	}

	return requestTasksResponse, nil
}

func (s *RequestService) CancelRequestHandler(ctx context.Context, requestId int32) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return err
	}

	request.Status = string(constants.RequestStatusCanceled)
	now := time.Now()
	request.CanceledAt = &now

	requestModel := model.Requests{}
	utils.Mapper(request, &requestModel)
	err = s.RequestRepo.UpdateRequest(ctx, tx, requestModel)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
