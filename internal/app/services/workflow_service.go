package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
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
	NatsService    *NatsService
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
		NatsService:    cfg.NatsService,
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

	var cc []string
	if node.CcEmail != nil {
		err := json.Unmarshal([]byte(*node.CcEmail), &cc)
		if err != nil {
			return responses.NodeResponse{}, err
		}
	}

	var to []string
	if node.ToEmail != nil {
		err := json.Unmarshal([]byte(*node.ToEmail), &to)
		if err != nil {
			return responses.NodeResponse{}, err
		}
	}

	var bcc []string
	if node.BccEmail != nil {
		err := json.Unmarshal([]byte(*node.BccEmail), &bcc)
		if err != nil {
			return responses.NodeResponse{}, err
		}
	}
	nodeDataResponse.EditorContent = responses.NodeDataResponseEditorContent{
		Subject: node.Subject,
		Body:    node.Body,
		Cc:      &cc,
		To:      &to,
		Bcc:     &bcc,
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

	currentTime := time.Now()

	// Update request status to in processing
	request.Status = string(constants.RequestStatusInProgress)
	request.StartedAt = &currentTime

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
			if request.Workflow.Type == string(constants.WorkflowTypeGeneral) {
				nodeModel.Status = string(constants.NodeStatusInProgress)

				currentTime := time.Now()
				nodeModel.ActualStartTime = &currentTime
			}

			if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
				return fmt.Errorf("update node status to in processing fail: %w", err)
			}

			if nodeModel.Type == string(constants.NodeTypeStory) || nodeModel.Type == string(constants.NodeTypeSubWorkflow) {

				currentTime := time.Now()
				nodeModel.ActualStartTime = &currentTime
				if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
					return fmt.Errorf("update node status to in processing fail: %w", err)
				}

				s.RunWorkflow(ctx, tx, *nodeModel.SubRequestID)
			}

			// Send notification
			if request.Nodes[i].AssigneeID != nil {
				notification := types.Notification{
					ToUserIds: []string{strconv.Itoa(int(*request.Nodes[i].AssigneeID))},
					Subject:   "New task assigned",
					Body:      fmt.Sprintf("New task assigned: %s – You have been assigned a new task by %d.", request.Nodes[i].Title, request.UserID),
				}

				notificationBytes, err := json.Marshal(notification)
				if err != nil {
					return fmt.Errorf("marshal notification failed: %w", err)
				}

				err = s.NatsClient.Publish("notifications", notificationBytes)
				if err != nil {
					return fmt.Errorf("publish notification failed: %w", err)
				}

			}

		}
	}

	// Notification
	uniqueUsers := make(map[int32]struct{})
	for _, node := range request.Nodes {
		if node.AssigneeID != nil {
			uniqueUsers[*node.AssigneeID] = struct{}{}
		}
	}

	userIdsStr := make([]string, 0, len(uniqueUsers))
	for id := range uniqueUsers {
		userIdsStr = append(userIdsStr, strconv.Itoa(int(id)))
	}

	// Send notification
	notification := types.Notification{
		ToUserIds: userIdsStr,
		Subject:   "Workflow Started",
		Body:      fmt.Sprintf("Workflow started with request ID: %d", requestId),
	}

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

func (s *WorkflowService) CreateNodesConnectionsStories(ctx context.Context, tx *sql.Tx, req *requests.NodesConnectionsStories, requestId int32, projectKey *string, userId int32, isStoryIsTemplate bool) error {
	formSystems, err := s.FormRepo.FindAllFormSystem(ctx, s.DB)
	if err != nil {
		return err
	}

	parentId := ""
	node, err := s.NodeRepo.FindOneNodeBySubRequestID(ctx, tx, requestId)
	if err == nil {
		if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
			parentId = node.ID
		}
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
			IsTemplate:        isStoryIsTemplate,
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

			SubRequestID:  &storyRequest.ID,
			EstimatePoint: storyReq.Node.Data.EstimatePoint,

			Status: string(constants.NodeStatusTodo),

			JiraKey: storyReq.Node.JiraKey,

			JiraLinkURL: storyReq.Node.Data.JiraLinkURL,
		}

		if formSystemVersionId, exists := formSystemTagMap["TASK"]; exists {
			// Create Form Data
			uuid := uuid.New()

			formData := model.FormData{
				ID:                    uuid.String(),
				FormTemplateVersionID: formSystemVersionId,
			}

			formData, err = s.FormRepo.CreateFormData(ctx, tx, formData)
			if err != nil {
				return fmt.Errorf("create form data system Fail: %w", err)
			}

			formFieldDatas := []model.FormFieldData{}
			for _, form := range storyReq.Node.Form {
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
					return fmt.Errorf("create form field datas  Fail: %w", err)
				}
			}

			storyNode.FormDataID = &formData.ID

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

				EstimatePoint: storyNodeReq.Data.EstimatePoint,

				JiraKey: storyNodeReq.JiraKey,

				JiraLinkURL: storyNodeReq.Data.JiraLinkURL,
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

			if formSystemVersionId, exists := formSystemTagMap[storyNodeReq.Type]; exists {
				// Create Form Data
				uuid := uuid.New()

				formData := model.FormData{
					ID:                    uuid.String(),
					FormTemplateVersionID: formSystemVersionId,
				}

				formData, err = s.FormRepo.CreateFormData(ctx, tx, formData)
				if err != nil {
					return fmt.Errorf("create form data system Fail: %w", err)
				}

				formFieldDatas := []model.FormFieldData{}
				for _, form := range storyNodeReq.Form {
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
						return fmt.Errorf("create form field datas  Fail: %w", err)
					}
				}

				storyNode.FormDataID = &formData.ID

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

			Status:        string(constants.NodeStatusTodo),
			EstimatePoint: workflowNodeReq.Data.EstimatePoint,

			JiraKey: workflowNodeReq.JiraKey,

			Subject: workflowNodeReq.Data.EditorContent.Subject,
			Body:    workflowNodeReq.Data.EditorContent.Body,
		}

		if workflowNodeReq.Data.EditorContent.Cc != nil {
			ccEmail, err := json.Marshal(workflowNodeReq.Data.EditorContent.Cc)
			if err != nil {
				return fmt.Errorf("marshal cc email fail: %w", err)
			}

			ccEmailString := string(ccEmail)
			workflowNode.CcEmail = &ccEmailString
		}

		if workflowNodeReq.Data.EditorContent.To != nil {
			toEmail, err := json.Marshal(workflowNodeReq.Data.EditorContent.To)
			if err != nil {
				return fmt.Errorf("marshal to email fail: %w", err)
			}

			toEmailString := string(toEmail)
			workflowNode.ToEmail = &toEmailString
		}

		if workflowNodeReq.Data.EditorContent.Bcc != nil {
			bccEmail, err := json.Marshal(workflowNodeReq.Data.EditorContent.Bcc)
			if err != nil {
				return fmt.Errorf("marshal bcc email fail: %w", err)
			}

			bccEmailString := string(bccEmail)
			workflowNode.BccEmail = &bccEmailString
		}

		if parentId != "" {
			workflowNode.ParentID = &parentId
		}

		if workflowNodeReq.Type == string(constants.NodeTypeTask) {
			formTemplateId := int32(1)
			workflowNode.FormTemplateID = &formTemplateId
		}

		for _, formSystem := range formSystems {
			if formSystem.Tag == workflowNodeReq.Type {
				// Create Form Data
				uuid := uuid.New()
				formData := model.FormData{
					FormTemplateVersionID: formSystem.Version.ID,
					ID:                    uuid.String(),
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
				OptionKey:                formAttached.OptionKey,
				FromFormAttachedPosition: formAttached.FromFormAttachedPosition,
				Permission:               formAttached.Permission,
				IsOriginal:               formAttached.IsOriginal,
				TemplateID:               formAttached.FormTemplateId,
				NodeID:                   workflowNodeReq.Id,
			}
			if formAttached.DataId != "" {
				formAttachedModel.DataID = &formAttached.DataId
			}

			formAttachedModels = append(formAttachedModels, formAttachedModel)

			// Form Data

			if formAttached.Permission == string("INPUT") && formAttached.DataId != "" {
				formTemplate, err := s.FormRepo.FindOneFormTemplateByFormTemplateId(ctx, s.DB, formAttached.FormTemplateId)
				if err != nil {
					return fmt.Errorf("find form template fail: %w", err)
				}

				formData := model.FormData{
					FormTemplateVersionID: formTemplate.Version.ID,
					ID:                    formAttached.DataId,
				}

				_, err = s.FormRepo.CreateFormData(ctx, tx, formData)
				if err != nil {
					return fmt.Errorf("create form data fail: %w", err)
				}
			}
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
	if err := s.CreateNodesConnectionsStories(ctx, tx, &req.NodesConnectionsStories, request.ID, &req.ProjectKey, userId, true); err != nil {
		return err
	}

	// Tính toán Gantt Chart nếu có project key
	if reqClone.ProjectKey != "" && reqClone.SprintId != nil {
		// Tạo bản đồ NodeId -> JiraKey để theo dõi các JiraKey
		jiraKeyMap := make(map[string]string)

		// Luôn đồng bộ với Jira để thiết lập mối quan hệ giữa các tasks
		slog.Info("Synchronizing with Jira before Gantt Chart calculation")
		jiraResponse, err := s.NatsService.publishWorkflowToJira(ctx, tx, reqClone.Nodes, reqClone.Stories, reqClone.Connections, reqClone.ProjectKey, *reqClone.SprintId)
		if err != nil {
			slog.Error("Failed to sync with Jira", "error", err)
			// Tiếp tục xử lý, không return error
		} else {
			// Cập nhật jiraKeyMap từ response
			for _, issue := range jiraResponse.Data.Data.Issues {
				jiraKeyMap[issue.NodeId] = issue.JiraKey
				slog.Info("JiraKey mapping from Jira response", "nodeId", issue.NodeId, "jiraKey", issue.JiraKey)
			}
		}

		// Dùng lại bản đồ ID -> JiraKey đã có hoặc JiraKey từ database
		updatedNodes := make([]requests.Node, len(reqClone.Nodes))
		for i, node := range reqClone.Nodes {
			updatedNode := node
			// Ưu tiên JiraKey từ Jira response
			if jiraKey, exists := jiraKeyMap[node.Id]; exists && jiraKey != "" {
				updatedNode.JiraKey = &jiraKey
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
				updatedStory.Node.JiraKey = &jiraKey
				slog.Info("Story JiraKey updated for Gantt Chart", "nodeId", story.Node.Id, "jiraKey", jiraKey)
			}
			updatedStories[i] = updatedStory
		}

		// Tính toán Gantt Chart với JiraKey đã cập nhật
		if err := s.NatsService.publishWorkflowToGanttChart(ctx, tx, updatedNodes, updatedStories, reqClone.Connections, reqClone.ProjectKey, *reqClone.SprintId, workflow.ID); err != nil {
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

	users, err := s.UserAPI.FindUsersByUserIds([]int32{userId})
	if err != nil {
		return workflowResponses, err
	}

	// Get User Project
	projects := []string{}
	for _, projectRole := range users.Data[0].ProjectRoles {
		projects = append(projects, projectRole.ProjectKey)
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

	users, err := s.UserAPI.FindUsersByUserIds([]int32{request.LastUpdateUserID})
	if err != nil {
		return workflowResponse, err
	}

	lastAssignee := types.Assignee{
		Id:           users.Data[0].ID,
		Name:         users.Data[0].Name,
		Email:        users.Data[0].Email,
		AvatarUrl:    users.Data[0].AvatarUrl,
		IsSystemUser: users.Data[0].IsSystemUser,
	}

	workflowResponse.LastAssignee = lastAssignee

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
		nodeResponse.ParentId = node.ParentID
		nodeResponse.JiraKey = node.JiraKey

		if node.AssigneeID != nil {
			user, err := s.UserAPI.FindUsersByUserIds([]int32{*node.AssigneeID})
			if err != nil {
				return workflowResponse, fmt.Errorf("find user fail: %w", err)
			}

			nodeResponse.Data.Assignee.Id = user.Data[0].ID
			nodeResponse.Data.Assignee.Name = user.Data[0].Name
			nodeResponse.Data.Assignee.Email = user.Data[0].Email
			nodeResponse.Data.Assignee.AvatarUrl = user.Data[0].AvatarUrl
			nodeResponse.Data.Assignee.IsSystemUser = user.Data[0].IsSystemUser
		}

		for _, formFieldData := range node.FormData.FormFieldData {
			nodeResponse.Form = append(nodeResponse.Form, responses.NodeFormDataResponse{
				FieldId: formFieldData.FormTemplateField.FieldID,
				Value:   formFieldData.Value,
			})
		}

		if err != nil {
			return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
		}

		story := responses.StoryResponse{
			Node: nodeResponse,

			Type:        request.Workflow.Type,
			Decoration:  request.Workflow.Decoration,
			Description: request.Workflow.Description,
			Title:       node.Title,
			CategoryKey: storyRequest.Category.Key,

			Progress: storyRequest.Progress,
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
		nodeResponse.ParentId = node.ParentID
		if err != nil {
			return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
		}

		// Form Attached
		for _, nodeForm := range node.NodeForms {

			approveUserIds := []responses.NodeFormApprovalOrRejectUsersResponse{}
			for _, approveUser := range nodeForm.NodeFormApproveOrRejectUsers {
				approveUserIds = append(approveUserIds, responses.NodeFormApprovalOrRejectUsersResponse{
					IsApproved: approveUser.IsApproved,
					Assignee: types.Assignee{
						Id:           approveUser.UserID,
						Name:         "",
						Email:        "",
						AvatarUrl:    "",
						IsSystemUser: false,
					},
				})
			}

			formAttachedResponse := responses.NodeFormResponse{
				Key:                           nodeForm.Key,
				FromUserId:                    nodeForm.FromUserID,
				OptionKey:                     nodeForm.OptionKey,
				FromFormAttachedPosition:      nodeForm.FromFormAttachedPosition,
				Permission:                    nodeForm.Permission,
				IsOriginal:                    nodeForm.IsOriginal,
				FormTemplateId:                nodeForm.TemplateID,
				NodeFormApprovalOrRejectUsers: approveUserIds,
			}
			if nodeForm.DataID != nil {
				formAttachedResponse.DataId = *nodeForm.DataID
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

		if node.Type == string(constants.NodeTypeCondition) {

			falseNodeDestinationIds := []string{}
			trueNodeDestinationIds := []string{}

			falseNodeDestinations, err := s.NodeRepo.FindAllNodeConditionDestinationByNodeId(ctx, s.DB, node.ID, false)
			if err != nil {
				return workflowResponse, fmt.Errorf("find node condition destination fail: %w", err)
			}
			for _, destination := range falseNodeDestinations {
				falseNodeDestinationIds = append(falseNodeDestinationIds, destination.DestinationNodeID)
			}

			trueNodeDestinations, err := s.NodeRepo.FindAllNodeConditionDestinationByNodeId(ctx, s.DB, node.ID, true)
			if err != nil {
				return workflowResponse, fmt.Errorf("find node condition destination fail: %w", err)
			}
			for _, destination := range trueNodeDestinations {
				trueNodeDestinationIds = append(trueNodeDestinationIds, destination.DestinationNodeID)
			}

			nodeResponse.Data.Condition = responses.NodeDataConditionResponse{
				TrueDestinations:  trueNodeDestinationIds,
				FalseDestinations: falseNodeDestinationIds,
			}
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
	existedUserIds := make(map[int32]bool)
	for _, node := range workflowResponse.Nodes {
		if existedUserIds[node.Data.Assignee.Id] {
			continue
		}
		userIds = append(userIds, node.Data.Assignee.Id)
		existedUserIds[node.Data.Assignee.Id] = true
	}

	results, err := s.UserAPI.FindUsersByUserIds(userIds)
	if err != nil {
		return workflowResponse, err
	}

	userMap := make(map[int32]types.Assignee)
	for _, user := range results.Data {
		userMap[user.ID] = types.Assignee{
			Id:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			AvatarUrl:    user.AvatarUrl,
			IsSystemUser: user.IsSystemUser,
		}
	}

	for i, node := range workflowResponse.Nodes {
		if user, exists := userMap[node.Data.Assignee.Id]; exists {
			workflowResponse.Nodes[i].Data.Assignee.Id = user.Id
			workflowResponse.Nodes[i].Data.Assignee.Name = user.Name
			workflowResponse.Nodes[i].Data.Assignee.Email = user.Email
			workflowResponse.Nodes[i].Data.Assignee.AvatarUrl = user.AvatarUrl
			workflowResponse.Nodes[i].Data.Assignee.IsSystemUser = user.IsSystemUser

		}
	}

	for i := range workflowResponse.Stories {
		workflowResponse.Stories[i].IsSystemLinked = true
		if user, exists := userMap[workflowResponse.Stories[i].Node.Data.Assignee.Id]; exists {
			workflowResponse.Stories[i].Node.Data.Assignee.Id = user.Id
			workflowResponse.Stories[i].Node.Data.Assignee.Name = user.Name
			workflowResponse.Stories[i].Node.Data.Assignee.Email = user.Email
			workflowResponse.Stories[i].Node.Data.Assignee.AvatarUrl = user.AvatarUrl
			workflowResponse.Stories[i].Node.Data.Assignee.IsSystemUser = user.IsSystemUser

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

	// Create Clone Request for Sync Jira
	reqClone := req
	reqDetailClone := request

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

	if err := s.CreateNodesConnectionsStories(ctx, tx, &nodeConnectionStoryReq, newRequest.ID, request.Workflow.ProjectKey, userId, false); err != nil {
		return 0, err
	}

	// =========================== SYNC JIRA ===========================
	// Tính toán Gantt Chart nếu có project key
	if reqDetailClone.Workflow.ProjectKey != nil && reqClone.SprintID != nil && s.NatsService != nil {
		// Tạo bản đồ NodeId -> JiraKey để theo dõi các JiraKey
		jiraKeyMap := make(map[string]string)

		// Luôn đồng bộ với Jira để thiết lập mối quan hệ giữa các tasks
		slog.Info("Synchronizing with Jira before Gantt Chart calculation")

		if s.NatsService == nil {
			slog.Error("Nats service is nil")
			return 0, fmt.Errorf("nats service is nil")
		}

		jiraResponse, err := s.NatsService.publishWorkflowToJira(ctx, tx, reqClone.Nodes, reqClone.Stories, reqClone.Connections, *reqDetailClone.Workflow.ProjectKey, *reqClone.SprintID)
		if err != nil {
			slog.Error("Failed to sync with Jira", "error", err)
			// Tiếp tục xử lý, không return error
		} else {
			// Cập nhật jiraKeyMap từ response
			for _, issue := range jiraResponse.Data.Data.Issues {
				jiraKeyMap[issue.NodeId] = issue.JiraKey
				slog.Info("JiraKey mapping from Jira response", "nodeId", issue.NodeId, "jiraKey", issue.JiraKey)
			}
		}

		// Dùng lại bản đồ ID -> JiraKey đã có hoặc JiraKey từ database
		updatedNodes := make([]requests.Node, len(reqClone.Nodes))
		for i, node := range reqClone.Nodes {
			updatedNode := node
			// Ưu tiên JiraKey từ Jira response
			if jiraKey, exists := jiraKeyMap[node.Id]; exists && jiraKey != "" {
				updatedNode.JiraKey = &jiraKey
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
				updatedStory.Node.JiraKey = &jiraKey
				slog.Info("Story JiraKey updated for Gantt Chart", "nodeId", story.Node.Id, "jiraKey", jiraKey)
			}
			updatedStories[i] = updatedStory
		}

		// Tính toán Gantt Chart với JiraKey đã cập nhật
		if err := s.NatsService.publishWorkflowToGanttChart(ctx, tx, updatedNodes, updatedStories, reqClone.Connections, *reqDetailClone.Workflow.ProjectKey, *reqClone.SprintID, request.Workflow.ID); err != nil {
			slog.Error("Failed to calculate Gantt Chart", "error", err)
			// Không return error ở đây để không làm fail luồng chính nếu tính toán Gantt Chart lỗi
		}
	}

	// =========================== END SYNC JIRA ===========================

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
			if err := s.CreateNodesConnectionsStories(ctx, tx, &nodeConnectionStoryReq, copyRequest.ID, subRequest.Workflow.ProjectKey, userId, false); err != nil {
				return 0, fmt.Errorf("create copy request nodes connections stories fail: %w", err)
			}
		}

		if node.Type == string(constants.NodeTypeTask) {
			uuid := uuid.New()

			formData := model.FormData{
				FormTemplateVersionID: int32(constants.FormTemplateIDJiraSystemForm),
				ID:                    uuid.String(),
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

	// Archive Workflow If Workflow Is Project
	if request.Workflow.Type == string(constants.WorkflowTypeProject) {
		if err := s.ArchiveWorkflowHandler(ctx, request.Workflow.ID); err != nil {
			return 0, fmt.Errorf("archive workflow fail: %w", err)
		}
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

	workflowModel := model.Workflows{}
	if err := utils.Mapper(workflow, &workflowModel); err != nil {
		return fmt.Errorf("map workflow fail: %w", err)
	}

	if err := s.WorkflowRepo.UpdateWorkflow(ctx, tx, workflowModel); err != nil {
		return fmt.Errorf("update workflow fail: %w", err)
	}

	for _, request := range workflow.Requests {
		subWorkflowModel := model.Workflows{}

		// Its Workflow when request parentID is nil, its story when request parentID is not nil
		if request.ParentID == nil {
			subRequests, err := s.RequestRepo.FindAllSubRequestByParentIdWithoutPagination(ctx, s.DB, request.ID)
			if err != nil {
				return fmt.Errorf("find sub request fail: %w", err)
			}

			for _, subRequest := range subRequests {
				if err := utils.Mapper(subRequest.Workflows, &subWorkflowModel); err != nil {
					return fmt.Errorf("map sub workflow fail: %w", err)
				}

				subWorkflowModel.IsArchived = true
				if err := s.WorkflowRepo.UpdateWorkflow(ctx, tx, subWorkflowModel); err != nil {
					return fmt.Errorf("update sub workflow fail: %w", err)
				}
			}
		} else {
			subMainRequest, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, *request.ParentID)
			if err != nil {
				return fmt.Errorf("find sub request fail: %w", err)
			}

			workflowModel := model.Workflows{}
			if err := utils.Mapper(subMainRequest.Workflow, &workflowModel); err != nil {
				return fmt.Errorf("map workflow fail: %w", err)
			}

			workflowModel.IsArchived = true
			if err := s.WorkflowRepo.UpdateWorkflow(ctx, tx, workflowModel); err != nil {
				return fmt.Errorf("update workflow fail: %w", err)
			}

			subRequests, err := s.RequestRepo.FindAllSubRequestByParentIdWithoutPagination(ctx, s.DB, subMainRequest.ID)
			if err != nil {
				return fmt.Errorf("find sub request fail: %w", err)
			}

			for _, subRequest := range subRequests {
				if err := utils.Mapper(subRequest.Workflows, &subWorkflowModel); err != nil {
					return fmt.Errorf("map sub workflow fail: %w", err)
				}

				subWorkflowModel.IsArchived = true
				if err := s.WorkflowRepo.UpdateWorkflow(ctx, tx, subWorkflowModel); err != nil {
					return fmt.Errorf("update sub workflow fail: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}
