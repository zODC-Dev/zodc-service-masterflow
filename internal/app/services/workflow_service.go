package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	natsModel "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type WorkflowService struct {
	DB             *sql.DB
	WorkflowRepo   *repositories.WorkflowRepository
	FormRepo       *repositories.FormRepository
	CategoryRepo   *repositories.CategoryRepository
	UserAPI        *externals.UserAPI
	RequestRepo    *repositories.RequestRepository
	ConnectionRepo *repositories.ConnectionRepository
	NodeRepo       *repositories.NodeRepository
	NodeService    *NodeService
	NatsClient     *nats.NATSClient
}

func NewWorkflowService(cfg WorkflowService) *WorkflowService {
	workflowService := WorkflowService{
		DB:             cfg.DB,
		WorkflowRepo:   cfg.WorkflowRepo,
		FormRepo:       cfg.FormRepo,
		CategoryRepo:   cfg.CategoryRepo,
		UserAPI:        cfg.UserAPI,
		RequestRepo:    cfg.RequestRepo,
		ConnectionRepo: cfg.ConnectionRepo,
		NodeRepo:       cfg.NodeRepo,
		NodeService:    cfg.NodeService,
		NatsClient:     cfg.NatsClient,
	}
	return &workflowService
}

func (s *WorkflowService) CreateWorkFlow(ctx context.Context, tx *sql.Tx, workflowData interface{}, projectKey *string, userId int32) (model.Workflows, error) {
	workflow := model.Workflows{
		CurrentVersion: 1,
		IsArchived:     false,
		UserID:         userId,
	}
	if err := utils.Mapper(workflowData, &workflow); err != nil {
		return workflow, fmt.Errorf("mapping workflow failed: %w", err)
	}

	workflow.ProjectKey = projectKey

	return s.WorkflowRepo.CreateWorkflow(ctx, tx, workflow)
}

func (s *WorkflowService) CreateWorkFlowVersion(ctx context.Context, tx *sql.Tx, workflowId int32, hasSubWorkflow bool, version int32) (model.WorkflowVersions, error) {
	workFlowVersion := model.WorkflowVersions{
		Version:        version,
		WorkflowID:     workflowId,
		HasSubWorkflow: hasSubWorkflow,
	}

	return s.WorkflowRepo.CreateWorkflowVersion(ctx, tx, workFlowVersion)
}

func (s *WorkflowService) MapToWorkflowNodeResponse(node model.Nodes) (responses.NodeResponse, error) {
	nodeDataResponse := responses.NodeDataResponse{}
	if err := utils.Mapper(node, &nodeDataResponse); err != nil {
		return responses.NodeResponse{}, err
	}

	if node.AssigneeID != nil {
		nodeDataResponse.Assignee.Id = *node.AssigneeID
	}

	nodeResponse := responses.NodeResponse{
		Id:   node.ID,
		Type: node.Type,
		Position: types.Position{
			X: node.X,
			Y: node.Y,
		},
		Size: types.Size{
			Width:  node.Width,
			Height: node.Height,
		},
		Data:     nodeDataResponse,
		ParentId: node.ParentID,

		JiraKey: node.JiraKey,

		Status: node.Status,

		StartedAt:   node.ActualStartTime,
		CompletedAt: node.ActualEndTime,
	}

	return nodeResponse, nil
}

func (s *WorkflowService) RunWorkflowIfItStoryOrSubWorkflow(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
		if node.SubRequestID == nil {
			return fmt.Errorf("sub request not found")
		}

		if err := s.RunWorkflow(ctx, tx, *node.SubRequestID); err != nil {
			return err
		}
	}

	return nil
}

func (s *WorkflowService) RunWorkflow(ctx context.Context, tx *sql.Tx, requestId int32) error {
	request, err := s.RequestRepo.FindOneRequestByRequestIdTx(ctx, tx, requestId)
	if err != nil {
		return fmt.Errorf("request not found")
	}

	// Update request status to in processing
	request.Status = string(constants.RequestStatusInProgress)

	requestModel := model.Requests{}
	if err := utils.Mapper(request, &requestModel); err != nil {
		return fmt.Errorf("map request fail: %w", err)
	}
	if err := s.RequestRepo.UpdateRequest(ctx, tx, requestModel); err != nil {
		return fmt.Errorf("update request fail: %w", err)
	}

	// Store Next Node For Update status to processing
	nextNodeIds := make(map[string]bool)
	for i := range request.Nodes {
		if request.Nodes[i].Type == string(constants.NodeTypeStart) {
			currentTime := time.Now()

			request.Nodes[i].Status = string(constants.NodeStatusCompleted)
			request.Nodes[i].IsCurrent = true
			request.Nodes[i].ActualEndTime = &currentTime
			request.Nodes[i].ActualStartTime = &currentTime

			nodeModel := model.Nodes{}
			if err := utils.Mapper(request.Nodes[i], &nodeModel); err != nil {
				return fmt.Errorf("map node fail: %w", err)
			}
			err = s.NodeRepo.UpdateNode(ctx, tx, nodeModel)
			if err != nil {
				return fmt.Errorf("update node status to completed fail: %w", err)
			}

			for j := range request.Connections {
				if request.Connections[j].FromNodeID == request.Nodes[i].ID {

					// Update connection
					request.Connections[j].IsCompleted = true
					if err := s.ConnectionRepo.UpdateConnection(ctx, tx, request.Connections[j]); err != nil {
						return fmt.Errorf("update connection fail: %w", err)
					}

					nextNodeIds[request.Connections[j].ToNodeID] = true
				}
			}
		}
	}

	for i := range request.Nodes {
		if nextNodeIds[request.Nodes[i].ID] {

			nodeModel := model.Nodes{}
			if err := utils.Mapper(request.Nodes[i], &nodeModel); err != nil {
				return fmt.Errorf("map node fail: %w", err)
			}

			nodeModel.IsCurrent = true

			if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
				return fmt.Errorf("update node status to in processing fail: %w", err)
			}

			// if err := s.NodeService.UpdateNodeStatusToInProcessing(ctx, tx, nodeModel); err != nil {
			// 	return fmt.Errorf("update node status to in processing fail: %w", err)
			// }

			// Send notification
			if request.Nodes[i].AssigneeID != nil {
				// notification := types.Notification{
				// 	ToUserIds: []string{strconv.Itoa(int(*request.Nodes[i].AssigneeID))},
				// 	Subject:   "New task assigned",
				// 	Body:      fmt.Sprintf("New task assigned: %s – You have been assigned a new task by %d.", request.Nodes[i].Title, request.UserID),
				// }

				// notificationBytes, err := json.Marshal(notification)
				// if err != nil {
				// 	return fmt.Errorf("marshal notification failed: %w", err)
				// }

				// err = s.NatsClient.Publish("notifications", notificationBytes)
				// if err != nil {
				// 	return fmt.Errorf("publish notification failed: %w", err)
				// }

			}

		}
	}

	// Notification
	// uniqueUsers := make(map[int32]struct{})
	// for _, node := range request.Nodes {
	// 	if node.AssigneeID != nil {
	// 		uniqueUsers[*node.AssigneeID] = struct{}{}
	// 	}
	// }

	// userIdsStr := make([]string, 0, len(uniqueUsers))
	// for id := range uniqueUsers {
	// 	userIdsStr = append(userIdsStr, strconv.Itoa(int(id)))
	// }

	// Send notification
	// notification := types.Notification{
	// 	ToUserIds: userIdsStr,
	// 	Subject:   "Workflow Started",
	// 	Body:      fmt.Sprintf("Workflow started with request ID: %d", requestId),
	// }

	// notificationBytes, err := json.Marshal(notification)
	// if err != nil {
	// 	return fmt.Errorf("marshal notification failed: %w", err)
	// }

	// err = s.NatsClient.Publish("notifications", notificationBytes)
	// if err != nil {
	// 	return fmt.Errorf("publish notification failed: %w", err)
	// }

	return nil
}

func (s *WorkflowService) CreateNodesConnectionsStories(ctx context.Context, tx *sql.Tx, req *requests.NodesConnectionsStories, requestId int32, projectKey *string, userId int32) error {
	formSystems, err := s.FormRepo.FindAllFormSystem(ctx, s.DB)
	if err != nil {
		return err
	}

	formSystemFieldMap := make(map[string]int32)
	formSystemTagMap := make(map[string]int32)

	for _, formSystem := range formSystems {
		formSystemTagMap[formSystem.Tag] = formSystem.Version.ID
		for _, field := range formSystem.Fields {
			formSystemFieldMap[field.FieldID] = field.ID
		}
	}

	//Create Stories
	for _, storyReq := range req.Stories {
		category, err := s.CategoryRepo.FindOneCategoryByKey(ctx, s.DB, storyReq.CategoryKey)
		if err != nil {
			return fmt.Errorf("category key not found: %w", err)
		}
		storyReq.CategoryId = category.ID

		// Create Story Workflow
		storyWorkflow, err := s.CreateWorkFlow(ctx, tx, storyReq, projectKey, userId)
		if err != nil {
			return fmt.Errorf("create Story Workflow Fail: %w", err)
		}

		storyWorkflowVersion, err := s.CreateWorkFlowVersion(ctx, tx, storyWorkflow.ID, false, 1)
		if err != nil {
			return fmt.Errorf("create Story Workflow Version Fail: %w", err)
		}

		storyRequestModel := model.Requests{
			Title:             storyReq.Title,
			WorkflowVersionID: storyWorkflowVersion.ID,
			IsTemplate:        true,
			Status:            string(constants.RequestStatusTodo),
			ParentID:          &requestId,
			UserID:            userId,
			LastUpdateUserID:  userId,
		}
		storyRequest, err := s.RequestRepo.CreateRequest(ctx, tx, storyRequestModel)
		if err != nil {
			return fmt.Errorf("create Story Request Fail: %w", err)
		}

		// Create Story Node
		storyNode := model.Nodes{
			ID:        storyReq.Node.Id,
			RequestID: requestId,

			X:      storyReq.Node.Position.X,
			Y:      storyReq.Node.Position.Y,
			Width:  storyReq.Node.Size.Width,
			Height: storyReq.Node.Size.Height,

			Type: storyReq.Node.Type,

			// Data
			Title:      storyReq.Node.Data.Title,
			AssigneeID: &storyReq.Node.Data.Assignee.Id,

			SubRequestID: &storyRequest.ID,

			Status: string(constants.NodeStatusTodo),
		}

		if err := s.NodeRepo.CreateNodes(ctx, tx, []model.Nodes{storyNode}); err != nil {
			return fmt.Errorf("create Story Workflow MAIN Node Fail: %w", err)
		}

		// Create Story Nodes
		storyNodes := []model.Nodes{}

		i := 0
		for _, storyNodeReq := range req.Nodes {

			if storyNodeReq.ParentId != storyReq.Node.Id {
				req.Nodes[i] = storyNodeReq
				i++
				continue
			}

			storyNode := model.Nodes{
				ID:        storyNodeReq.Id,
				RequestID: storyRequest.ID,

				X:      storyNodeReq.Position.X,
				Y:      storyNodeReq.Position.Y,
				Width:  storyNodeReq.Size.Width,
				Height: storyNodeReq.Size.Height,

				Type: storyNodeReq.Type,

				// ParentID: &storyNodeReq.ParentId,

				Title: storyNodeReq.Data.Title,

				Status: string(constants.NodeStatusTodo),

				EstimatePoint: storyNodeReq.EstimatePoint,
			}

			if storyNodeReq.ParentId != "" {
				storyNode.ParentID = &storyNodeReq.ParentId
			}
			if storyNodeReq.Data.Assignee.Id != 0 {
				storyNode.AssigneeID = &storyNodeReq.Data.Assignee.Id
			}
			if storyNodeReq.Data.EndType != "" {
				storyNode.EndType = &storyNodeReq.Data.EndType
			}

			// Form Type System Tag Story // Create Form Data
			for _, formSystem := range formSystems {
				if formSystem.Tag == storyNodeReq.Type {
					// Create Form Data
					formData := model.FormData{
						FormTemplateVersionID: formSystem.Version.ID,
					}

					formData, err = s.FormRepo.CreateFormData(ctx, tx, formData)
					if err != nil {
						return fmt.Errorf("create form data system Fail: %w", err)
					}

					formFieldDatas := []model.FormFieldData{}
					for _, form := range storyNodeReq.Form {
						for _, field := range formSystem.Fields {
							if field.FieldID == form.FieldId {
								formFieldData := model.FormFieldData{
									Value:               form.Value,
									FormDataID:          formData.ID,
									FormTemplateFieldID: field.ID,
								}

								formFieldDatas = append(formFieldDatas, formFieldData)
							}
						}
					}

					if len(formFieldDatas) > 0 {
						err := s.FormRepo.CreateFormFieldDatas(ctx, tx, formFieldDatas)
						if err != nil {
							return fmt.Errorf("create form field datas  Fail: %w", err)
						}
					}

					storyNode.FormDataID = &formData.ID

					break
				}
			}

			storyNodes = append(storyNodes, storyNode)
		}
		req.Nodes = req.Nodes[:i]

		if len(storyNodes) > 0 {
			err = s.NodeRepo.CreateNodes(ctx, tx, storyNodes)
			if err != nil {
				return fmt.Errorf("create Story Node Fail: %w", err)
			}
		}

		// Create Story Connections
		storyConnections := []model.Connections{}

		i = 0
		for _, connReq := range req.Connections {
			shouldKeepConnection := true

			for _, storyNode := range storyNodes {

				if storyNode.ID == connReq.From {
					shouldKeepConnection = false

					storyConnection := model.Connections{
						ID:         connReq.Id,
						FromNodeID: connReq.From,
						ToNodeID:   connReq.To,
						RequestID:  storyRequest.ID,
					}

					if connReq.Text != "" {
						storyConnection.Text = &connReq.Text
					}

					storyConnections = append(storyConnections, storyConnection)
				}
			}

			if shouldKeepConnection {
				req.Connections[i] = connReq
				i++
			}
		}
		req.Connections = req.Connections[:i]

		if len(storyConnections) > 0 {
			err = s.ConnectionRepo.CreateConnections(ctx, tx, storyConnections)
			if err != nil {
				return fmt.Errorf("create Story Connection Fail: %w", err)
			}
		}

	}

	// Create Workflow Nodes
	workflowNodes := []model.Nodes{}
	formAttachedModels := []model.NodeForms{}
	nodeConditionDestinations := []model.NodeConditionDestinations{}

	for _, workflowNodeReq := range req.Nodes {
		workflowNode := model.Nodes{
			ID:        workflowNodeReq.Id,
			RequestID: requestId,

			X:      workflowNodeReq.Position.X,
			Y:      workflowNodeReq.Position.Y,
			Width:  workflowNodeReq.Size.Width,
			Height: workflowNodeReq.Size.Height,

			Type: workflowNodeReq.Type,

			AssigneeID: &workflowNodeReq.Data.Assignee.Id,

			SubRequestID: workflowNodeReq.Data.SubRequestID,

			// Data
			Title:   workflowNodeReq.Data.Title,
			EndType: &workflowNodeReq.Data.EndType,

			Status: string(constants.NodeStatusTodo),
		}

		for _, formSystem := range formSystems {
			if formSystem.Tag == workflowNodeReq.Type {
				// Create Form Data
				formData := model.FormData{
					FormTemplateVersionID: formSystem.Version.ID,
				}

				formData, err = s.FormRepo.CreateFormData(ctx, tx, formData)
				if err != nil {
					return fmt.Errorf("create form data fail: %w", err)
				}

				formFieldDatas := []model.FormFieldData{}
				for _, form := range workflowNodeReq.Form {
					if fieldID, exists := formSystemFieldMap[form.FieldId]; exists {
						formFieldData := model.FormFieldData{
							Value:               form.Value,
							FormDataID:          formData.ID,
							FormTemplateFieldID: fieldID,
						}
						formFieldDatas = append(formFieldDatas, formFieldData)
					}
				}

				if len(formFieldDatas) > 0 {
					err := s.FormRepo.CreateFormFieldDatas(ctx, tx, formFieldDatas)
					if err != nil {
						return fmt.Errorf("create form fields data fail: %w", err)
					}
				}

				workflowNode.FormDataID = &formData.ID

				break
			}
		}

		workflowNodes = append(workflowNodes, workflowNode)

		// Condition Node

		if workflowNodeReq.Type == string(constants.NodeTypeCondition) {

			for _, destination := range workflowNodeReq.Data.Condition.TrueDestinations {
				nodeConditionDestinations = append(nodeConditionDestinations, model.NodeConditionDestinations{
					IsTrue:            true,
					DestinationNodeID: destination,
					NodeID:            workflowNodeReq.Id,
				})
			}

			for _, destination := range workflowNodeReq.Data.Condition.FalseDestinations {
				nodeConditionDestinations = append(nodeConditionDestinations, model.NodeConditionDestinations{
					IsTrue:            false,
					DestinationNodeID: destination,
					NodeID:            workflowNodeReq.Id,
				})
			}

		}

		// Form Attached

		for _, formAttached := range workflowNodeReq.Data.FormAttached {
			formAttachedModel := model.NodeForms{
				Key:                      formAttached.Key,
				FromUserID:               formAttached.FromUserId,
				DataID:                   formAttached.DataId,
				OptionKey:                formAttached.OptionId,
				FromFormAttachedPosition: formAttached.FromFormAttachedPosition,
				Permission:               formAttached.Permission,
				IsOriginal:               formAttached.IsOriginal,
				TemplateID:               formAttached.FormTemplateId,
				NodeID:                   workflowNodeReq.Id,
			}

			formAttachedModels = append(formAttachedModels, formAttachedModel)
		}

	}
	// Create Workflow Nodes
	if len(workflowNodes) > 0 {
		err = s.NodeRepo.CreateNodes(ctx, tx, workflowNodes)
		if err != nil {
			return fmt.Errorf("create Workflow Nodes Fail: %w", err)
		}
	}

	// Create Node Forms
	if len(formAttachedModels) > 0 {
		err = s.NodeRepo.CreateNodeForms(ctx, tx, formAttachedModels)
		if err != nil {
			return fmt.Errorf("create node forms fail: %w", err)
		}
	}

	// Create Node Condition Destinations
	if len(nodeConditionDestinations) > 0 {
		err = s.NodeRepo.CreateNodeConditionDestinations(ctx, tx, nodeConditionDestinations)
		if err != nil {
			return fmt.Errorf("create node condition destinations fail: %w", err)
		}
	}

	// Create Workflow Connections
	workflowConnections := []model.Connections{}
	for _, workflowConnectionReq := range req.Connections {
		workflowConnection := model.Connections{
			ID: workflowConnectionReq.Id,

			FromNodeID: workflowConnectionReq.From,
			ToNodeID:   workflowConnectionReq.To,

			RequestID: requestId,
		}

		workflowConnections = append(workflowConnections, workflowConnection)
	}
	if len(workflowConnections) > 0 {
		err = s.ConnectionRepo.CreateConnections(ctx, tx, workflowConnections)
		if err != nil {
			return fmt.Errorf("create Workflow Connections Fail: %w", err)
		}
	}

	return nil
}

// Handlers
func (s *WorkflowService) CreateWorkflowHandler(ctx context.Context, req *requests.CreateWorkflow, userId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("create workflow handler fail: %w", err)
	}
	defer tx.Rollback()

	// Clone req to reqClone for sending to Jira
	reqClone := *req

	// Create Workflow
	workflow, err := s.CreateWorkFlow(ctx, tx, req, &req.ProjectKey, userId)
	if err != nil {
		return fmt.Errorf("create Main workflow Fail: %w", err)
	}

	//Check If Workflow Has SubWorkflow
	hasSubWorkflow := len(req.Stories) > 0
	for i := range req.Nodes {
		if req.Nodes[i].Type == "SUB_WORKFLOW" {
			hasSubWorkflow = true
			break
		}
	}

	// Create Workflow Version
	workflowVersion, err := s.CreateWorkFlowVersion(ctx, tx, workflow.ID, hasSubWorkflow, 1)
	if err != nil {
		return fmt.Errorf("create Main Workflow Version Fail: %w", err)
	}

	// Create Request
	requestModel := model.Requests{
		Title:             req.Title,
		WorkflowVersionID: workflowVersion.ID,
		IsTemplate:        true,
		Status:            string(constants.RequestStatusTodo),
		UserID:            userId,
		LastUpdateUserID:  userId,
	}
	request, err := s.RequestRepo.CreateRequest(ctx, tx, requestModel)
	if err != nil {
		return fmt.Errorf("create Main Request Fail: %w", err)
	}

	nodeConnectionStoryReq := requests.NodesConnectionsStories{}
	if err := utils.Mapper(req.NodesConnectionsStories, &nodeConnectionStoryReq); err != nil {
		return fmt.Errorf("create Main Node Connection Story Fail: %w", err)
	}
	if err := s.CreateNodesConnectionsStories(ctx, tx, &req.NodesConnectionsStories, request.ID, &req.ProjectKey, userId); err != nil {
		return err
	}

	// Tính toán Gantt Chart nếu có project key
	if reqClone.ProjectKey != "" {
		// Đồng bộ với Jira trước (đối với các nodes mới được tạo)
		needJiraSync := false
		for _, node := range reqClone.Nodes {
			if (node.Type == string(constants.NodeTypeTask) ||
				node.Type == string(constants.NodeTypeBug) ||
				node.Type == string(constants.NodeTypeStory)) &&
				node.JiraKey == "" {
				needJiraSync = true
				break
			}
		}

		// Tạo bản đồ NodeId -> JiraKey để theo dõi các JiraKey
		jiraKeyMap := make(map[string]string)

		// Nếu có nodes mới cần đồng bộ
		if needJiraSync {
			slog.Info("Synchronizing new nodes with Jira before Gantt Chart calculation")
			jiraResponse, err := s.publishWorkflowToJira(ctx, tx, reqClone.Nodes, reqClone.Stories, reqClone.Connections, reqClone.ProjectKey, *reqClone.SprintId)
			if err != nil {
				slog.Error("Failed to sync new nodes with Jira", "error", err)
				// Tiếp tục xử lý, không return error
			} else {
				// Cập nhật jiraKeyMap từ response
				for _, issue := range jiraResponse.Data.Data.Issues {
					jiraKeyMap[issue.NodeId] = issue.JiraKey
					slog.Info("JiraKey mapping from Jira response", "nodeId", issue.NodeId, "jiraKey", issue.JiraKey)
				}
			}
		}

		// Dùng lại bản đồ ID -> JiraKey đã có hoặc JiraKey từ database
		updatedNodes := make([]requests.Node, len(reqClone.Nodes))
		for i, node := range reqClone.Nodes {
			updatedNode := node
			// Ưu tiên JiraKey từ Jira response
			if jiraKey, exists := jiraKeyMap[node.Id]; exists && jiraKey != "" {
				updatedNode.JiraKey = jiraKey
				slog.Info("Node JiraKey updated for Gantt Chart", "nodeId", node.Id, "jiraKey", jiraKey)
			}
			updatedNodes[i] = updatedNode
		}

		// Cập nhật JiraKey cho stories
		updatedStories := make([]requests.Story, len(reqClone.Stories))
		for i, story := range reqClone.Stories {
			updatedStory := story
			// Ưu tiên JiraKey từ Jira response
			if jiraKey, exists := jiraKeyMap[story.Node.Id]; exists && jiraKey != "" {
				updatedStory.Node.JiraKey = jiraKey
				slog.Info("Story JiraKey updated for Gantt Chart", "nodeId", story.Node.Id, "jiraKey", jiraKey)
			}
			updatedStories[i] = updatedStory
		}

		// Tính toán Gantt Chart với JiraKey đã cập nhật
		if err := s.publishWorkflowToGanttChart(ctx, tx, updatedNodes, updatedStories, reqClone.Connections, reqClone.ProjectKey, *reqClone.SprintId, workflow.ID); err != nil {
			slog.Error("Failed to calculate Gantt Chart", "error", err)
			// Không return error ở đây để không làm fail luồng chính nếu tính toán Gantt Chart lỗi
		}
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *WorkflowService) FindAllWorkflowHandler(ctx context.Context, workflowTemplateQueryParams queryparams.WorkflowQueryParam, userId int32) ([]responses.WorkflowResponse, error) {
	workflowResponses := []responses.WorkflowResponse{}

	userIds := []int32{int32(userId)}

	projects := []string{}
	if workflowTemplateQueryParams.Type != string(constants.WorkflowTypeGeneral) {
		users, err := s.UserAPI.FindUsersByUserIds(userIds)
		if err != nil {
			return workflowResponses, err
		}

		// Get User Project
		for _, projectRole := range users.Data[0].ProjectRoles {
			projects = append(projects, projectRole.ProjectKey)
		}
	}

	workflows, err := s.WorkflowRepo.FindAllWorkflowTemplates(ctx, s.DB, workflowTemplateQueryParams, projects)
	if err != nil {
		return workflowResponses, err
	}

	for _, workflow := range workflows {

		//Mapping workflow response
		workflowResponse := responses.WorkflowResponse{
			IsArchived:     workflow.IsArchived,
			RequestId:      workflow.Request.ID,
			CurrentVersion: workflow.CurrentVersion,
		}
		if err := utils.Mapper(workflow, &workflowResponse); err != nil {
			return workflowResponses, err
		}

		workflowResponses = append(workflowResponses, workflowResponse)
	}

	return workflowResponses, nil
}

func (s *WorkflowService) FindOneWorkflowDetailHandler(ctx context.Context, requestId int32) (responses.WorkflowDetailResponse, error) {
	workflowResponse := responses.WorkflowDetailResponse{}

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, requestId)
	if err != nil {
		return workflowResponse, err
	}

	//Mapping workflow response
	if err := utils.Mapper(request.Workflow, &workflowResponse); err != nil {
		return workflowResponse, err
	}

	categoryResponse := responses.CategoryResponse{}
	if err := utils.Mapper(request.Category, &categoryResponse); err != nil {
		return workflowResponse, err
	}
	workflowResponse.Category = categoryResponse

	workflowResponse.CurrentVersion = request.Version.Version
	workflowResponse.IsArchived = request.Workflow.IsArchived

	workflowResponse.RequestId = requestId

	if request.Workflow.ProjectKey != nil {
		workflowResponse.ProjectKey = *request.Workflow.ProjectKey
	}

	workflowResponse.Connections = []responses.ConnectionResponse{}
	workflowResponse.Nodes = []responses.NodeResponse{}

	// Stories
	storiesResponse := []responses.StoryResponse{}
	i := 0
	for _, node := range request.Nodes {
		if node.Type != "STORY" {
			request.Nodes[i] = node
			i++
			continue
		}

		storyRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *node.SubRequestID)
		if err != nil {
			return workflowResponse, fmt.Errorf("find story request fail: %w", err)
		}

		// Map Response
		nodeModel := model.Nodes{}
		if err := utils.Mapper(node, &nodeModel); err != nil {
			return workflowResponse, fmt.Errorf("map node fail: %w", err)
		}

		nodeResponse, err := s.MapToWorkflowNodeResponse(nodeModel)
		nodeResponse.ParentId = nil
		if err != nil {
			return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
		}

		story := responses.StoryResponse{
			Node: nodeResponse,

			Type:        request.Workflow.Type,
			Decoration:  request.Workflow.Decoration,
			Description: request.Workflow.Description,
			Title:       request.Workflow.Title,
			CategoryKey: storyRequest.Category.Key,
		}

		storiesResponse = append(storiesResponse, story)
	}
	request.Nodes = request.Nodes[:i]
	workflowResponse.Stories = storiesResponse

	// Nodes

	for _, node := range request.Nodes {
		nodeModel := model.Nodes{}
		if err := utils.Mapper(node, &nodeModel); err != nil {
			return workflowResponse, fmt.Errorf("map node fail: %w", err)
		}

		nodeResponse, err := s.MapToWorkflowNodeResponse(nodeModel)
		nodeResponse.ParentId = nil
		if err != nil {
			return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
		}

		// Form Attached
		for _, nodeForm := range node.NodeForms {

			approveUserIds := []int32{}
			for _, approveUser := range nodeForm.NodeFormApproveUsers {
				approveUserIds = append(approveUserIds, approveUser.UserID)
			}

			formAttachedResponse := responses.NodeFormResponse{
				Key:                      nodeForm.Key,
				FromUserId:               nodeForm.FromUserID,
				DataId:                   nodeForm.DataID,
				OptionKey:                nodeForm.OptionKey,
				FromFormAttachedPosition: nodeForm.FromFormAttachedPosition,
				Permission:               nodeForm.Permission,
				IsOriginal:               nodeForm.IsOriginal,
				FormTemplateId:           nodeForm.TemplateID,
				ApproveUserIds:           approveUserIds,
			}

			nodeResponse.Data.FormAttached = append(nodeResponse.Data.FormAttached, formAttachedResponse)
		}

		// Form Data
		for _, formFieldData := range node.FormData.FormFieldData {
			nodeResponse.Form = append(nodeResponse.Form, responses.NodeFormDataResponse{
				FieldId: formFieldData.FormTemplateField.FieldID,
				Value:   formFieldData.Value,
			})
		}

		workflowResponse.Nodes = append(workflowResponse.Nodes, nodeResponse)
	}

	// Connections
	for _, connection := range request.Connections {
		connectionResponse := responses.ConnectionResponse{
			Id:          connection.ID,
			To:          connection.ToNodeID,
			From:        connection.FromNodeID,
			IsCompleted: connection.IsCompleted,
		}

		workflowResponse.Connections = append(workflowResponse.Connections, connectionResponse)
	}

	// Add assignee
	userIds := make([]int32, 0, len(workflowResponse.Nodes))
	for _, node := range workflowResponse.Nodes {
		userIds = append(userIds, node.Data.Assignee.Id)
	}

	results, err := s.UserAPI.FindUsersByUserIds(userIds)
	if err != nil {
		return workflowResponse, err
	}

	userMap := make(map[int32]struct {
		Name         string
		Email        string
		AvatarUrl    string
		IsSystemUser bool
	})
	for _, user := range results.Data {
		userMap[user.ID] = struct {
			Name         string
			Email        string
			AvatarUrl    string
			IsSystemUser bool
		}{
			Name:         user.Name,
			Email:        user.Email,
			AvatarUrl:    user.AvatarUrl,
			IsSystemUser: user.IsSystemUser,
		}
	}

	for i, node := range workflowResponse.Nodes {
		if user, exists := userMap[node.Data.Assignee.Id]; exists {
			workflowResponse.Nodes[i].Data.Assignee.Name = user.Name
			workflowResponse.Nodes[i].Data.Assignee.Email = user.Email
			workflowResponse.Nodes[i].Data.Assignee.AvatarUrl = user.AvatarUrl
			workflowResponse.Nodes[i].Data.Assignee.IsSystemUser = user.IsSystemUser

		}
	}

	return workflowResponse, nil
}

func (s *WorkflowService) StartWorkflowHandler(ctx context.Context, req requests.StartWorkflow, userId int32) (int32, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return 0, fmt.Errorf("start workflow handler fail: %w", err)
	}
	defer tx.Rollback()

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, req.RequestID)
	if err != nil {
		return 0, fmt.Errorf("find request fail: %w", err)
	}

	workflowVersionId := request.Version.ID

	// If Request Is Template
	if req.IsTemplate {
		currentVersion := request.Workflow.CurrentVersion + 1
		request.Workflow.CurrentVersion = currentVersion

		if err := s.WorkflowRepo.UpdateWorkflow(ctx, tx, request.Workflow); err != nil {
			return 0, fmt.Errorf("update workflow fail: %w", err)
		}

		workflowVersion, err := s.CreateWorkFlowVersion(ctx, tx, request.Workflow.ID, request.Version.HasSubWorkflow, currentVersion)
		if err != nil {
			return 0, err
		}

		// Set workflowVersionId to new one
		workflowVersionId = workflowVersion.ID

		requestTemplate := model.Requests{
			IsTemplate:        true,
			WorkflowVersionID: workflowVersion.ID,
			Status:            string(constants.RequestStatusTodo),
			Title:             req.Title,
			UserID:            userId,
			LastUpdateUserID:  userId,
			Progress:          0,
		}

		_, err = s.RequestRepo.CreateRequest(ctx, tx, requestTemplate)
		if err != nil {
			return 0, err
		}

	}

	startedAt := time.Now()

	requestModel := model.Requests{
		Title:             req.Title,
		IsTemplate:        false,
		WorkflowVersionID: workflowVersionId,
		Status:            string(constants.RequestStatusInProgress),
		UserID:            userId,
		LastUpdateUserID:  userId,
		SprintID:          req.SprintID,
		StartedAt:         &startedAt,
	}
	newRequest, err := s.RequestRepo.CreateRequest(ctx, tx, requestModel)
	if err != nil {
		return 0, err
	}

	nodeConnectionStoryReq := requests.NodesConnectionsStories{}
	if err := utils.Mapper(req, &nodeConnectionStoryReq); err != nil {
		return 0, err
	}

	if err := s.CreateNodesConnectionsStories(ctx, tx, &nodeConnectionStoryReq, newRequest.ID, request.Workflow.ProjectKey, userId); err != nil {
		return 0, err
	}

	// Create Sub Workflow and Stories
	for _, node := range req.Nodes {
		// Only create sub workflow or stories
		if node.Type == string(constants.NodeTypeSubWorkflow) || node.Type == string(constants.NodeTypeStory) {
			subRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *node.Data.SubRequestID)
			if err != nil {
				return 0, fmt.Errorf("find story request fail: %w", err)
			}

			// Copy request
			requestModel := model.Requests{}
			if err := utils.Mapper(subRequest, &requestModel); err != nil {
				return 0, fmt.Errorf("map request fail: %w", err)
			}
			requestModel.SprintID = req.SprintID
			requestModel.Status = string(constants.RequestStatusTodo)
			copyRequest, err := s.RequestRepo.CreateRequest(ctx, tx, requestModel)
			if err != nil {
				return 0, fmt.Errorf("create copy request fail: %w", err)
			}

			// Copy Request Nodes Connections SubWorkflow
			nodeConnectionStoryReq := requests.NodesConnectionsStories{}
			if err := utils.Mapper(subRequest, &nodeConnectionStoryReq); err != nil {
				return 0, fmt.Errorf("map node connection story request fail: %w", err)
			}
			if err := s.CreateNodesConnectionsStories(ctx, tx, &nodeConnectionStoryReq, copyRequest.ID, subRequest.Workflow.ProjectKey, userId); err != nil {
				return 0, fmt.Errorf("create copy request nodes connections stories fail: %w", err)
			}
		}

		if node.Type == string(constants.NodeTypeTask) {
			formData := model.FormData{
				FormTemplateVersionID: int32(constants.FormTemplateIDJiraSystemForm),
			}
			formData, err := s.FormRepo.CreateFormData(ctx, tx, formData)
			if err != nil {
				return 0, fmt.Errorf("create form data fail: %w", err)
			}

			formFieldDatas := []model.FormFieldData{}

			formTemplate, err := s.FormRepo.FindOneFormTemplateByFormTemplateId(ctx, s.DB, constants.FormTemplateIDJiraSystemForm)
			if err != nil {
				return 0, fmt.Errorf("find form template by form template id fail: %w", err)
			}

			formFieldMap := make(map[string]int32)
			for _, formTemplateField := range formTemplate.Fields {
				formFieldMap[formTemplateField.FieldID] = formTemplateField.ID
			}

			for _, form := range node.Form {
				if formTemplateFieldId, exists := formFieldMap[form.FieldId]; exists {
					formFieldData := model.FormFieldData{
						FormTemplateFieldID: formTemplateFieldId,
						Value:               form.Value,
						FormDataID:          formData.ID,
					}

					formFieldDatas = append(formFieldDatas, formFieldData)
				}
			}
			if len(formFieldDatas) > 0 {
				s.FormRepo.CreateFormFieldDatas(ctx, tx, formFieldDatas)
			}
		}

	}

	// Run Workflow
	if err := s.RunWorkflow(ctx, tx, newRequest.ID); err != nil {
		return 0, fmt.Errorf("run workflow fail: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit fail: %w", err)
	}

	return newRequest.ID, nil
}

func (s *WorkflowService) ArchiveWorkflowHandler(ctx context.Context, workflowId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("archive workflow handler fail: %w", err)
	}
	defer tx.Rollback()

	workflow, err := s.WorkflowRepo.FindOneWorkflowByWorkflowId(ctx, s.DB, workflowId)
	if err != nil {
		return fmt.Errorf("find workflow fail: %w", err)
	}

	workflow.IsArchived = true

	if err := s.WorkflowRepo.UpdateWorkflow(ctx, tx, workflow); err != nil {
		return fmt.Errorf("update workflow fail: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

//////////////////// JIRA ////////////////////

// publishWorkflowToJira gửi dữ liệu workflow đến Jira và trả về phản hồi
func (s *WorkflowService) publishWorkflowToJira(ctx context.Context, tx *sql.Tx, nodes []requests.Node, stories []requests.Story, connections []requests.Connection, projectKey string, sprintId int32) (natsModel.WorkflowSyncResponse, error) {
	slog.Info("Starting Jira synchronization", "projectKey", projectKey, "sprintId", sprintId)
	slog.Info("Processing stories", "len stories", len(stories), "len nodes", len(nodes), "len connections", len(connections))
	slog.Info("Processing stories", "stories", stories)
	slog.Info("Processing nodes", "nodes", nodes)
	slog.Info("Processing connections", "connections", connections)

	// First, ensure story assignees have feature_leader role
	if err := s.assignFeatureLeaderRoles(stories, projectKey); err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to assign feature leader roles: %w", err)
	}

	syncRequest := natsModel.WorkflowSyncRequest{
		TransactionId: uuid.New().String(),
		ProjectKey:    projectKey,
		SprintId:      sprintId,
		Issues:        make([]natsModel.WorkflowSyncIssue, 0),
		Connections:   make([]natsModel.WorkflowSyncConnection, 0),
	}

	// Process Stories
	for _, story := range stories {
		slog.Info("Processing story",
			"id", story.Node.Id,
			"title", story.Title,
			"jiraKey", story.Node.JiraKey)

		issue := natsModel.WorkflowSyncIssue{
			NodeId:     story.Node.Id,
			Type:       "Story",
			Title:      story.Title,
			AssigneeId: &story.Node.Data.Assignee.Id,
			Action:     "create",
		}

		if story.Node.JiraKey != "" {
			issue.Action = "update"
			issue.JiraKey = story.Node.JiraKey
		}

		syncRequest.Issues = append(syncRequest.Issues, issue)
	}

	// Process Tasks and Bugs
	for _, node := range nodes {
		slog.Info("Processing node details",
			"id", node.Id,
			"type", node.Type,
			"title", node.Data.Title,
			"jiraKey", node.JiraKey,
			"assigneeId", node.Data.Assignee.Id,
			"parentId", node.ParentId)

		if node.Type != string(constants.NodeTypeTask) && node.Type != string(constants.NodeTypeBug) {
			slog.Info("Skipping node - not a task or bug",
				"id", node.Id,
				"type", node.Type)
			continue
		}

		slog.Info("Processing node",
			"id", node.Id,
			"type", node.Type,
			"title", node.Data.Title,
			"jiraKey", node.JiraKey)

		issue := natsModel.WorkflowSyncIssue{
			NodeId:     node.Id,
			Type:       node.Type,
			Title:      node.Data.Title,
			AssigneeId: &node.Data.Assignee.Id,
			Action:     "create",
		}

		if node.JiraKey != "" {
			issue.Action = "update"
			issue.JiraKey = node.JiraKey
		}

		syncRequest.Issues = append(syncRequest.Issues, issue)
	}

	// Process Connections
	processedConnections := make(map[string]bool) // Để tránh duplicate connections

	// 1. Xử lý connections từ connections hiện có
	for _, conn := range connections {
		fromNode := findNodeByIdFromRequest(nodes, stories, conn.From)
		toNode := findNodeByIdFromRequest(nodes, stories, conn.To)

		if fromNode == nil || toNode == nil {
			slog.Info("Skipping connection - node not found",
				"fromId", conn.From,
				"toId", conn.To)
			continue
		}

		// Skip START/END connections
		if fromNode.Type == string(constants.NodeTypeStart) ||
			toNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Sử dụng node ID để tạo connectionKey
		connectionKey := fmt.Sprintf("%s-%s", fromNode.Id, toNode.Id)
		if processedConnections[connectionKey] {
			continue
		}

		connection := natsModel.WorkflowSyncConnection{
			FromIssueKey: fromNode.Id, // Luôn sử dụng node ID
			ToIssueKey:   toNode.Id,   // Luôn sử dụng node ID
			Type:         "relates to",
		}

		syncRequest.Connections = append(syncRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 2. Xử lý connections từ parent-child relationships
	for _, node := range nodes {
		if node.ParentId == "" {
			continue
		}

		// Tìm parent node (có thể là story hoặc node khác)
		parentNode := findNodeByIdFromRequest(nodes, stories, node.ParentId)
		if parentNode == nil {
			continue
		}

		// Skip nếu parent là START/END
		if parentNode.Type == string(constants.NodeTypeStart) ||
			parentNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Sử dụng node ID để tạo connectionKey
		connectionKey := fmt.Sprintf("%s-%s", parentNode.Id, node.Id)
		if processedConnections[connectionKey] {
			continue
		}

		slog.Info("Creating parent-child connection",
			"parent", parentNode.Data.Title,
			"child", node.Data.Title,
			"parentType", parentNode.Type,
			"childType", node.Type)

		connection := natsModel.WorkflowSyncConnection{
			FromIssueKey: parentNode.Id, // Luôn sử dụng node ID
			ToIssueKey:   node.Id,       // Luôn sử dụng node ID
			Type:         "contains",
		}

		syncRequest.Connections = append(syncRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 3. Xử lý connections giữa stories và nodes
	for _, story := range stories {
		// Tìm tất cả nodes liên quan đến story này
		for _, node := range nodes {
			// Skip nếu node là START/END
			if node.Type == string(constants.NodeTypeStart) ||
				node.Type == string(constants.NodeTypeEnd) {
				continue
			}

			// Sử dụng node ID để tạo connectionKey
			connectionKey := fmt.Sprintf("%s-%s", story.Node.Id, node.Id)
			if processedConnections[connectionKey] {
				continue
			}

			// Kiểm tra mối quan hệ giữa story và node
			if node.ParentId == story.Node.Id {
				slog.Info("Creating story-node connection",
					"story", story.Title,
					"node", node.Data.Title)

				connection := natsModel.WorkflowSyncConnection{
					FromIssueKey: story.Node.Id, // Luôn sử dụng node ID
					ToIssueKey:   node.Id,       // Luôn sử dụng node ID
					Type:         "contains",
				}

				syncRequest.Connections = append(syncRequest.Connections, connection)
				processedConnections[connectionKey] = true
			}
		}
	}

	// Send to NATS
	slog.Info("Sending sync request to NATS", "request", syncRequest)
	requestBytes, err := json.Marshal(syncRequest)
	if err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to marshal sync request: %w", err)
	}

	response, err := s.NatsClient.Request(constants.NatsTopicWorkflowSyncRequest, requestBytes, 30*time.Second)
	if err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to sync with Jira: %w", err)
	}

	// Log raw response để debug
	slog.Info("Received raw response from Jira sync service", "response", string(response.Data))

	// Process response
	var syncResponse natsModel.WorkflowSyncResponse
	if err := json.Unmarshal(response.Data, &syncResponse); err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to unmarshal Jira response: %w", err)
	}

	// Kiểm tra response success
	if !syncResponse.Success || !syncResponse.Data.Success {
		slog.Error("Jira synchronization failed",
			"outerSuccess", syncResponse.Success,
			"innerSuccess", syncResponse.Data.Success)
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("Jira synchronization failed")
	}

	// Update JiraKeys in database from nested structure
	for _, issue := range syncResponse.Data.Data.Issues {
		slog.Info("Updating JiraKey",
			"nodeId", issue.NodeId,
			"jiraKey", issue.JiraKey)

		if err := s.NodeRepo.UpdateJiraKey(ctx, tx, issue.NodeId, issue.JiraKey); err != nil {
			return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to update JiraKey: %w", err)
		}
	}

	slog.Info("Completed Jira synchronization", "projectKey", projectKey)
	return syncResponse, nil
}

// New function to assign feature_leader roles to story assignees
func (s *WorkflowService) assignFeatureLeaderRoles(stories []requests.Story, projectKey string) error {
	// Keep track of users we've already assigned the role to avoid duplicate requests
	assignedUsers := make(map[int32]bool)

	for _, story := range stories {
		// Skip if no assignee or already processed
		if story.Node.Data.Assignee.Id == 0 || assignedUsers[story.Node.Data.Assignee.Id] {
			continue
		}

		// Mark this user as processed
		assignedUsers[story.Node.Data.Assignee.Id] = true

		// Create role assignment request
		roleRequest := natsModel.RoleAssignmentRequest{
			UserID:     story.Node.Data.Assignee.Id,
			ProjectKey: projectKey,
			RoleName:   constants.RoleFeatureLeader,
		}

		slog.Info("Assigning feature_leader role",
			"userId", roleRequest.UserID,
			"projectKey", roleRequest.ProjectKey)

		// Convert to JSON
		requestBytes, err := json.Marshal(roleRequest)
		if err != nil {
			return fmt.Errorf("failed to marshal role assignment request: %w", err)
		}

		// Send request via NATS with a 10-second timeout
		response, err := s.NatsClient.Request(constants.NatsTopicRoleAssignmentRequest, requestBytes, 10*time.Second)
		if err != nil {
			return fmt.Errorf("failed to assign feature_leader role to user %d: %w",
				story.Node.Data.Assignee.Id, err)
		}

		// Process response
		var roleResponse natsModel.RoleAssignmentResponse
		if err := json.Unmarshal(response.Data, &roleResponse); err != nil {
			return fmt.Errorf("failed to unmarshal role assignment response: %w", err)
		}

		// Check if successful
		if !roleResponse.Success {
			return fmt.Errorf("role assignment failed for user %d: %s",
				story.Node.Data.Assignee.Id, roleResponse.Message)
		}

		slog.Info("Successfully assigned feature_leader role",
			"userId", roleRequest.UserID,
			"projectKey", roleRequest.ProjectKey)
	}

	return nil
}

// Helper function to find node by ID from request structs
func findNodeByIdFromRequest(nodes []requests.Node, stories []requests.Story, nodeId string) *requests.Node {
	// Tìm trong nodes
	for i := range nodes {
		if nodes[i].Id == nodeId {
			return &requests.Node{
				Id:       nodes[i].Id,
				Type:     nodes[i].Type,
				Data:     nodes[i].Data,
				ParentId: nodes[i].ParentId,
				JiraKey:  nodes[i].JiraKey,
			}
		}
	}

	// Tìm trong stories
	for i := range stories {
		if stories[i].Node.Id == nodeId {
			return &requests.Node{
				Id:       stories[i].Node.Id,
				Type:     stories[i].Node.Type,
				Data:     stories[i].Node.Data,
				ParentId: "",
				JiraKey:  stories[i].Node.JiraKey,
			}
		}
	}

	return nil
}

// publishWorkflowToGanttChart gửi dữ liệu workflow đến service tính toán Gantt Chart
func (s *WorkflowService) publishWorkflowToGanttChart(ctx context.Context, tx *sql.Tx, nodes []requests.Node, stories []requests.Story, connections []requests.Connection, projectKey string, sprintId int32, workflowId int32) error {
	slog.Info("Starting Gantt Chart calculation",
		"projectKey", projectKey,
		"sprintId", sprintId,
		"workflowId", workflowId,
		"nodes", len(nodes),
		"stories", len(stories))

	// Log chi tiết nodes IDs và JiraKeys từ request
	for i, node := range nodes {
		slog.Info(fmt.Sprintf("Request node #%d details", i),
			"nodeId", node.Id,
			"title", node.Data.Title,
			"type", node.Type,
			"jiraKey", node.JiraKey)
	}

	// Lấy thông tin từ database để so sánh
	nodeMap := make(map[string]string) // map nodeId -> jiraKey

	// Trước tiên sử dụng JiraKeys từ parameters
	for _, node := range nodes {
		if node.JiraKey != "" {
			nodeMap[node.Id] = node.JiraKey
			slog.Info("Using JiraKey from request", "nodeId", node.Id, "jiraKey", node.JiraKey)
		}
	}

	for _, story := range stories {
		if story.Node.JiraKey != "" {
			nodeMap[story.Node.Id] = story.Node.JiraKey
			slog.Info("Using JiraKey from request story", "nodeId", story.Node.Id, "jiraKey", story.Node.JiraKey)
		}
	}

	// Chuẩn bị request
	ganttRequest := natsModel.GanttChartCalculationRequest{
		WorkflowId:  workflowId,
		SprintId:    sprintId,
		ProjectKey:  projectKey,
		Issues:      []natsModel.GanttChartJiraIssue{},
		Connections: []natsModel.GanttChartConnection{},
	}

	// Processing Stories - Add to issues
	for _, story := range stories {
		jiraKey := story.Node.JiraKey
		// Ưu tiên JiraKey từ database
		if dbJiraKey, exists := nodeMap[story.Node.Id]; exists && dbJiraKey != "" {
			jiraKey = dbJiraKey
		}

		slog.Info("Processing story for Gantt Chart",
			"id", story.Node.Id,
			"title", story.Title,
			"type", "STORY",
			"requestJiraKey", story.Node.JiraKey,
			"dbJiraKey", nodeMap[story.Node.Id],
			"finalJiraKey", jiraKey)

		issue := natsModel.GanttChartJiraIssue{
			NodeId:  story.Node.Id,
			Type:    "STORY",
			JiraKey: jiraKey,
		}

		ganttRequest.Issues = append(ganttRequest.Issues, issue)
	}

	// Processing Tasks and Bugs - Add to issues
	for _, node := range nodes {
		if node.Type != string(constants.NodeTypeTask) &&
			node.Type != string(constants.NodeTypeBug) &&
			node.Type != string(constants.NodeTypeStory) {
			continue
		}

		jiraKey := node.JiraKey
		// Ưu tiên JiraKey từ database
		if dbJiraKey, exists := nodeMap[node.Id]; exists && dbJiraKey != "" {
			jiraKey = dbJiraKey
		}

		slog.Info("Processing node for Gantt Chart",
			"id", node.Id,
			"title", node.Data.Title,
			"type", node.Type,
			"requestJiraKey", node.JiraKey,
			"dbJiraKey", nodeMap[node.Id],
			"finalJiraKey", jiraKey)

		// Cảnh báo nếu không có JiraKey
		if jiraKey == "" {
			slog.Warn("Node missing JiraKey for Gantt Chart calculation",
				"nodeId", node.Id,
				"type", node.Type,
				"title", node.Data.Title)
		}

		issue := natsModel.GanttChartJiraIssue{
			NodeId:  node.Id,
			Type:    node.Type,
			JiraKey: jiraKey,
		}

		ganttRequest.Issues = append(ganttRequest.Issues, issue)
	}

	// Processing Connections
	processedConnections := make(map[string]bool)

	// 1. Connections từ connections hiện có
	for _, conn := range connections {
		fromNode := findNodeByIdFromRequest(nodes, stories, conn.From)
		toNode := findNodeByIdFromRequest(nodes, stories, conn.To)

		if fromNode == nil || toNode == nil {
			continue
		}

		// Skip START/END connections
		if fromNode.Type == string(constants.NodeTypeStart) ||
			toNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Skip if not story/task/bug nodes
		if (fromNode.Type != string(constants.NodeTypeStory) &&
			fromNode.Type != string(constants.NodeTypeTask) &&
			fromNode.Type != string(constants.NodeTypeBug)) ||
			(toNode.Type != string(constants.NodeTypeStory) &&
				toNode.Type != string(constants.NodeTypeTask) &&
				toNode.Type != string(constants.NodeTypeBug)) {
			continue
		}

		// Tạo connection key để tránh duplicate
		connectionKey := fmt.Sprintf("%s-%s", fromNode.Id, toNode.Id)
		if processedConnections[connectionKey] {
			continue
		}

		connection := natsModel.GanttChartConnection{
			FromNodeId: fromNode.Id,
			ToNodeId:   toNode.Id,
			Type:       "relates to",
		}

		ganttRequest.Connections = append(ganttRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 2. Mối quan hệ parent-child
	for _, node := range nodes {
		if node.ParentId == "" {
			continue
		}

		// Skip if not task/bug
		if node.Type != string(constants.NodeTypeTask) &&
			node.Type != string(constants.NodeTypeBug) {
			continue
		}

		// Tìm parent node
		parentNode := findNodeByIdFromRequest(nodes, stories, node.ParentId)
		if parentNode == nil {
			continue
		}

		// Skip nếu parent không phải story/task/bug
		if parentNode.Type != string(constants.NodeTypeStory) &&
			parentNode.Type != string(constants.NodeTypeTask) &&
			parentNode.Type != string(constants.NodeTypeBug) {
			continue
		}

		// Tạo connection key để tránh duplicate
		connectionKey := fmt.Sprintf("%s-%s", parentNode.Id, node.Id)
		if processedConnections[connectionKey] {
			continue
		}

		// Tạo connection type "contains" giữa parent và node
		connection := natsModel.GanttChartConnection{
			FromNodeId: parentNode.Id,
			ToNodeId:   node.Id,
			Type:       "contains",
		}

		ganttRequest.Connections = append(ganttRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 3. Mối quan hệ story-tasks
	for _, story := range stories {
		for _, node := range nodes {
			// Skip if not task/bug
			if node.Type != string(constants.NodeTypeTask) &&
				node.Type != string(constants.NodeTypeBug) {
				continue
			}

			// Tạo connection key để tránh duplicate
			connectionKey := fmt.Sprintf("%s-%s", story.Node.Id, node.Id)
			if processedConnections[connectionKey] {
				continue
			}

			// Kiểm tra nếu node thuộc story
			if node.ParentId == story.Node.Id {
				// Tạo connection type "contains" giữa story và node
				connection := natsModel.GanttChartConnection{
					FromNodeId: story.Node.Id,
					ToNodeId:   node.Id,
					Type:       "contains",
				}

				ganttRequest.Connections = append(ganttRequest.Connections, connection)
				processedConnections[connectionKey] = true
			}
		}
	}

	// Gửi request đến NATS
	slog.Info("Sending Gantt Chart calculation request to NATS", "request", ganttRequest)
	requestBytes, err := json.Marshal(ganttRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal Gantt Chart request: %w", err)
	}

	response, err := s.NatsClient.Request(constants.NatsTopicGanttChartCalculationRequest, requestBytes, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to calculate Gantt Chart: %w", err)
	}

	// Log raw response để debug
	slog.Info("Received raw response from Gantt Chart service", "response", string(response.Data))

	// Xử lý response có cấu trúc lồng ghép
	var ganttResponse natsModel.GanttChartCalculationResponse
	if err := json.Unmarshal(response.Data, &ganttResponse); err != nil {
		slog.Error("Failed to unmarshal Gantt Chart response", "error", err)
		return fmt.Errorf("failed to unmarshal Gantt Chart response: %w", err)
	}

	// Kiểm tra response success
	if !ganttResponse.Success || !ganttResponse.Data.Success {
		slog.Error("Gantt Chart calculation failed",
			"outerSuccess", ganttResponse.Success,
			"innerSuccess", ganttResponse.Data.Success)
		return fmt.Errorf("Gantt Chart calculation failed")
	}

	// Log cụ thể các node được cập nhật
	slog.Info("Issues returned from Gantt Chart service", "count", len(ganttResponse.Data.Data.Issues))
	for i, issue := range ganttResponse.Data.Data.Issues {
		slog.Info(fmt.Sprintf("Issue #%d details", i),
			"nodeId", issue.NodeId,
			"plannedStart", issue.PlannedStartTime,
			"plannedEnd", issue.PlannedEndTime)
	}

	// Cập nhật PlannedStartTime và PlannedEndTime vào database
	nodeUpdates := make([]struct {
		NodeId           string
		PlannedStartTime time.Time
		PlannedEndTime   time.Time
	}, len(ganttResponse.Data.Data.Issues))

	for i, issue := range ganttResponse.Data.Data.Issues {
		// Cập nhật PlannedStartTime và PlannedEndTime
		nodeUpdates[i] = struct {
			NodeId           string
			PlannedStartTime time.Time
			PlannedEndTime   time.Time
		}{
			NodeId:           issue.NodeId,
			PlannedStartTime: issue.PlannedStartTime,
			PlannedEndTime:   issue.PlannedEndTime,
		}
	}

	if len(nodeUpdates) > 0 {
		if err := s.NodeRepo.UpdateNodePlannedTimes(ctx, tx, nodeUpdates); err != nil {
			slog.Error("Failed to update node planned times", "error", err)
			return fmt.Errorf("failed to update node planned times: %w", err)
		}

		// Log thành công sau khi đã lưu
		slog.Info("Successfully updated planned times for nodes", "count", len(nodeUpdates))
	} else {
		slog.Warn("No planned times received from Gantt Chart service")
	}

	slog.Info("Completed Gantt Chart calculation", "projectKey", projectKey)
	return nil
}
