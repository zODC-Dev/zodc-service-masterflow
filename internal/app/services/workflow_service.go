package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
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

		StartedAt: node.ActualStartTime,
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
			request.Nodes[i].Status = string(constants.NodeStatusCompleted)

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
			WorkflowVersionID: storyWorkflowVersion.ID,
			IsTemplate:        false,
			Status:            "IN_ACTIVE",
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

				// shuold remove len if check ?
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
			Id:   connection.ID,
			To:   connection.ToNodeID,
			From: connection.FromNodeID,
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
				FormTemplateVersionID: constants.FormTemplateIDJiraSystemForm,
			}
			formData, err := s.FormRepo.CreateFormData(ctx, tx, formData)
			if err != nil {
				return 0, fmt.Errorf("create form data fail: %w", err)
			}

			formFieldDatas := []model.FormFieldData{}
			for _, form := range node.Form {
				formTemplateField, _ := s.FormRepo.FindOneFormTemplateFieldByFieldId(ctx, tx, form.FieldId, constants.FormTemplateIDJiraSystemForm)

				formFieldData := model.FormFieldData{
					FormTemplateFieldID: formTemplateField.ID,
					Value:               form.Value,
					FormDataID:          formData.ID,
				}

				formFieldDatas = append(formFieldDatas, formFieldData)

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
