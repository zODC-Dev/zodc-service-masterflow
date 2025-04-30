package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/externals"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type WorkflowService struct {
	DB                  *sql.DB
	WorkflowRepo        *repositories.WorkflowRepository
	FormRepo            *repositories.FormRepository
	CategoryRepo        *repositories.CategoryRepository
	UserAPI             *externals.UserAPI
	RequestRepo         *repositories.RequestRepository
	ConnectionRepo      *repositories.ConnectionRepository
	NodeRepo            *repositories.NodeRepository
	NodeService         *NodeService
	NatsClient          *nats.NATSClient
	NatsService         *NatsService
	NotificationService *NotificationService
	HistoryService      *HistoryService
}

func NewWorkflowService(cfg WorkflowService) *WorkflowService {
	workflowService := WorkflowService{
		DB:                  cfg.DB,
		WorkflowRepo:        cfg.WorkflowRepo,
		FormRepo:            cfg.FormRepo,
		CategoryRepo:        cfg.CategoryRepo,
		UserAPI:             cfg.UserAPI,
		RequestRepo:         cfg.RequestRepo,
		ConnectionRepo:      cfg.ConnectionRepo,
		NodeRepo:            cfg.NodeRepo,
		NodeService:         cfg.NodeService,
		NatsClient:          cfg.NatsClient,
		NatsService:         cfg.NatsService,
		NotificationService: cfg.NotificationService,
		HistoryService:      cfg.HistoryService,
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
	if node.CcEmails != nil {
		err := json.Unmarshal([]byte(*node.CcEmails), &cc)
		if err != nil {
			return responses.NodeResponse{}, err
		}
	}

	var to []string
	if node.ToEmails != nil {
		err := json.Unmarshal([]byte(*node.ToEmails), &to)
		if err != nil {
			return responses.NodeResponse{}, err
		}
	}

	var bcc []string
	if node.BccEmails != nil {
		err := json.Unmarshal([]byte(*node.BccEmails), &bcc)
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

		Level: node.Level,

		IsCurrent: node.IsCurrent,
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

	currentTime := time.Now().UTC().Add(7 * time.Hour)

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
			currentTime := time.Now().UTC().Add(7 * time.Hour)

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

				currentTime := time.Now().UTC().Add(7 * time.Hour)
				nodeModel.ActualStartTime = &currentTime
			}

			if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
				return fmt.Errorf("update node status to in processing fail: %w", err)
			}

			if nodeModel.Type == string(constants.NodeTypeStory) || nodeModel.Type == string(constants.NodeTypeSubWorkflow) {

				currentTime := time.Now().UTC().Add(7 * time.Hour)
				nodeModel.ActualStartTime = &currentTime
				if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
					return fmt.Errorf("update node status to in processing fail: %w", err)
				}

				s.RunWorkflow(ctx, tx, *nodeModel.SubRequestID)
			}

		}
	}

	return nil
}

func (s *WorkflowService) CreateNodesConnectionsStories(ctx context.Context, tx *sql.Tx, req *requests.NodesConnectionsStories, requestId int32, projectKey *string, userId int32, sprintId *int32, isStoryIsTemplate bool) error {
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
			SprintID:          sprintId,
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

			JiraLinkURL: storyReq.Node.Data.JiraLinkUrl,

			PlannedEndTime: setTimeToEndOfWorkday(storyReq.Node.Data.EndDate),
			// Index
			Level: storyReq.Node.Level,

			IsSendApprovedForm: storyReq.Node.Data.EditorContent.IsSendApprovedForm,
			IsSendRejectedForm: storyReq.Node.Data.EditorContent.IsSendRejectedForm,

			TaskStartedRequester:    storyReq.Node.Data.TaskStarted.Requester,
			TaskStartedAssignee:     storyReq.Node.Data.TaskStarted.Assignee,
			TaskStartedParticipants: storyReq.Node.Data.TaskStarted.Participants,

			TaskCompletedRequester:    storyReq.Node.Data.TaskCompleted.Requester,
			TaskCompletedAssignee:     storyReq.Node.Data.TaskCompleted.Assignee,
			TaskCompletedParticipants: storyReq.Node.Data.TaskCompleted.Participants,

			//
			Description: storyReq.Node.Data.Description,
			AttachFile:  storyReq.Node.Data.AttachFile,
		}

		// For Update
		if storyReq.Node.Status != string(constants.NodeStatusTodo) {
			storyNode.Status = storyReq.Node.Status
		}
		if storyReq.Node.Data.PlannedStartTime != nil {
			storyNode.PlannedStartTime = storyReq.Node.Data.PlannedStartTime
		}
		if storyReq.Node.Data.PlannedEndTime != nil {
			storyNode.PlannedEndTime = storyReq.Node.Data.PlannedEndTime
		}
		if storyReq.Node.Data.ActualStartTime != nil {
			storyNode.ActualStartTime = storyReq.Node.Data.ActualStartTime
		}
		if storyReq.Node.Data.ActualEndTime != nil {
			storyNode.ActualEndTime = storyReq.Node.Data.ActualEndTime
		}

		storyNode.IsCurrent = storyReq.Node.IsCurrent
		// End For Update

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

				JiraLinkURL: storyNodeReq.Data.JiraLinkUrl,

				PlannedEndTime: setTimeToEndOfWorkday(storyNodeReq.Data.EndDate),
				//
				Level: storyNodeReq.Level,

				IsSendApprovedForm: storyNodeReq.Data.EditorContent.IsSendApprovedForm,
				IsSendRejectedForm: storyNodeReq.Data.EditorContent.IsSendRejectedForm,

				TaskStartedRequester:    storyNodeReq.Data.TaskStarted.Requester,
				TaskStartedAssignee:     storyNodeReq.Data.TaskStarted.Assignee,
				TaskStartedParticipants: storyNodeReq.Data.TaskStarted.Participants,

				TaskCompletedRequester:    storyNodeReq.Data.TaskCompleted.Requester,
				TaskCompletedAssignee:     storyNodeReq.Data.TaskCompleted.Assignee,
				TaskCompletedParticipants: storyNodeReq.Data.TaskCompleted.Participants,

				//
				Description: storyNodeReq.Data.Description,
				AttachFile:  storyNodeReq.Data.AttachFile,
			}

			// For Update
			if storyNodeReq.Status != string(constants.NodeStatusTodo) {
				storyNode.Status = storyNodeReq.Status
			}
			if storyNodeReq.Data.PlannedStartTime != nil {
				storyNode.PlannedStartTime = storyNodeReq.Data.PlannedStartTime
			}
			if storyNodeReq.Data.PlannedEndTime != nil {
				storyNode.PlannedEndTime = storyNodeReq.Data.PlannedEndTime
			}
			if storyNodeReq.Data.ActualStartTime != nil {
				storyNode.ActualStartTime = storyNodeReq.Data.ActualStartTime
			}
			if storyNodeReq.Data.ActualEndTime != nil {
				storyNode.ActualEndTime = storyNodeReq.Data.ActualEndTime
			}

			storyNode.IsCurrent = storyNodeReq.IsCurrent
			// End For Update

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
						ID:          connReq.Id,
						FromNodeID:  connReq.From,
						ToNodeID:    connReq.To,
						RequestID:   storyRequest.ID,
						IsCompleted: connReq.IsCompleted,
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

			PlannedEndTime: setTimeToEndOfWorkday(workflowNodeReq.Data.EndDate),

			Subject: workflowNodeReq.Data.EditorContent.Subject,
			Body:    workflowNodeReq.Data.EditorContent.Body,

			// Index
			Level: workflowNodeReq.Level,

			IsSendApprovedForm: workflowNodeReq.Data.EditorContent.IsSendApprovedForm,
			IsSendRejectedForm: workflowNodeReq.Data.EditorContent.IsSendRejectedForm,

			TaskStartedRequester:    workflowNodeReq.Data.TaskStarted.Requester,
			TaskStartedAssignee:     workflowNodeReq.Data.TaskStarted.Assignee,
			TaskStartedParticipants: workflowNodeReq.Data.TaskStarted.Participants,

			TaskCompletedRequester:    workflowNodeReq.Data.TaskCompleted.Requester,
			TaskCompletedAssignee:     workflowNodeReq.Data.TaskCompleted.Assignee,
			TaskCompletedParticipants: workflowNodeReq.Data.TaskCompleted.Participants,

			//
			Description: workflowNodeReq.Data.Description,
			AttachFile:  workflowNodeReq.Data.AttachFile,
		}

		// For Update
		if workflowNodeReq.Status != string(constants.NodeStatusTodo) {
			workflowNode.Status = workflowNodeReq.Status
		}
		if workflowNodeReq.Data.PlannedStartTime != nil {
			workflowNode.PlannedStartTime = workflowNodeReq.Data.PlannedStartTime
		}
		if workflowNodeReq.Data.PlannedEndTime != nil {
			workflowNode.PlannedEndTime = workflowNodeReq.Data.PlannedEndTime
		}
		if workflowNodeReq.Data.ActualStartTime != nil {
			workflowNode.ActualStartTime = workflowNodeReq.Data.ActualStartTime
		}
		if workflowNodeReq.Data.ActualEndTime != nil {
			workflowNode.ActualEndTime = workflowNodeReq.Data.ActualEndTime
		}

		workflowNode.IsCurrent = workflowNodeReq.IsCurrent
		// End For Update

		if workflowNodeReq.Data.EditorContent.Cc != nil {
			ccEmail, err := json.Marshal(workflowNodeReq.Data.EditorContent.Cc)
			if err != nil {
				return fmt.Errorf("marshal cc email fail: %w", err)
			}

			ccEmailString := string(ccEmail)
			workflowNode.CcEmails = &ccEmailString
		}

		if workflowNodeReq.Data.EditorContent.To != nil {
			toEmail, err := json.Marshal(workflowNodeReq.Data.EditorContent.To)
			if err != nil {
				return fmt.Errorf("marshal to email fail: %w", err)
			}

			toEmailString := string(toEmail)
			workflowNode.ToEmails = &toEmailString
		}

		if workflowNodeReq.Data.EditorContent.Bcc != nil {
			bccEmail, err := json.Marshal(workflowNodeReq.Data.EditorContent.Bcc)
			if err != nil {
				return fmt.Errorf("marshal bcc email fail: %w", err)
			}

			bccEmailString := string(bccEmail)
			workflowNode.BccEmails = &bccEmailString
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
				TemplateVersionID:        formAttached.FormTemplateVersionId,
				NodeID:                   workflowNodeReq.Id,
				Level:                    formAttached.Level,
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

			IsCompleted: workflowConnectionReq.IsCompleted,

			RequestID: requestId,
		}

		if workflowConnectionReq.Text != "" {
			workflowConnection.Text = &workflowConnectionReq.Text
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
		SprintID:          req.SprintId,
	}
	request, err := s.RequestRepo.CreateRequest(ctx, tx, requestModel)
	if err != nil {
		return fmt.Errorf("create Main Request Fail: %w", err)
	}

	nodeConnectionStoryReq := requests.NodesConnectionsStories{}
	if err := utils.Mapper(req.NodesConnectionsStories, &nodeConnectionStoryReq); err != nil {
		return fmt.Errorf("create Main Node Connection Story Fail: %w", err)
	}
	if err := s.CreateNodesConnectionsStories(ctx, tx, &req.NodesConnectionsStories, request.ID, &req.ProjectKey, userId, req.SprintId, true); err != nil {
		return err
	}

	// Tính toán Gantt Chart nếu có project key
	if reqClone.ProjectKey != "" && reqClone.SprintId != nil {
		// Luôn đồng bộ với Jira để thiết lập mối quan hệ giữa các tasks
		_, err := s.NatsService.PublishWorkflowToJira(ctx, tx, reqClone.Nodes, reqClone.Stories, reqClone.Connections, reqClone.ProjectKey, *reqClone.SprintId)
		if err != nil {
			slog.Error("Failed to sync with Jira", "error", err)
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
	workflowResponse.Title = request.Title

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
				FormTemplateVersionId:         nodeForm.TemplateVersionID,
				NodeFormApprovalOrRejectUsers: approveUserIds,
				Level:                         nodeForm.Level,
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

		if connection.Text != nil {
			connectionResponse.Text = *connection.Text
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

		requestTemplate, err = s.RequestRepo.CreateRequest(ctx, tx, requestTemplate)
		if err != nil {
			return 0, err
		}

		nodeConnectionStoryReq := requests.NodesConnectionsStories{}
		if err := utils.Mapper(req.Template, &nodeConnectionStoryReq); err != nil {
			return 0, err
		}

		if err := s.CreateNodesConnectionsStories(ctx, tx, &nodeConnectionStoryReq, requestTemplate.ID, request.Workflow.ProjectKey, userId, req.SprintID, false); err != nil {
			return 0, err
		}

	}

	startedAt := time.Now().UTC().Add(7 * time.Hour)

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

	if err := s.CreateNodesConnectionsStories(ctx, tx, &nodeConnectionStoryReq, newRequest.ID, request.Workflow.ProjectKey, userId, req.SprintID, false); err != nil {
		return 0, err
	}

	// =========================== SYNC JIRA ===========================
	// Tính toán Gantt Chart nếu có project key
	if reqDetailClone.Workflow.ProjectKey != nil && reqClone.SprintID != nil && s.NatsService != nil {
		// Tạo bản đồ NodeId -> JiraKey để theo dõi các JiraKey
		jiraKeyMap := make(map[string]string)

		if s.NatsService == nil {
			slog.Error("Nats service is nil")
			return 0, fmt.Errorf("nats service is nil")
		}

		_, err := s.NatsService.PublishWorkflowToJira(ctx, tx, reqClone.Nodes, reqClone.Stories, reqClone.Connections, *reqDetailClone.Workflow.ProjectKey, *reqClone.SprintID)
		if err != nil {
			slog.Error("Failed to sync with Jira", "error", err)
			// Tiếp tục xử lý, không return error
		} else {
			// Dùng lại bản đồ ID -> JiraKey đã có hoặc JiraKey từ database
			updatedNodes := make([]requests.Node, len(reqClone.Nodes))
			for i, node := range reqClone.Nodes {
				updatedNode := node
				// Ưu tiên JiraKey từ Jira response
				if jiraKey, exists := jiraKeyMap[node.Id]; exists && jiraKey != "" {
					updatedNode.JiraKey = &jiraKey
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
				}
				updatedStories[i] = updatedStory
			}

			// Tính toán Gantt Chart với JiraKey đã cập nhật
			if err := s.NatsService.PublishWorkflowToGanttChart(ctx, tx, updatedNodes, updatedStories, reqClone.Connections, *reqDetailClone.Workflow.ProjectKey, *reqClone.SprintID, request.Workflow.ID); err != nil {
				slog.Error("Failed to calculate Gantt Chart", "error", err)
			}
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
			if err := s.CreateNodesConnectionsStories(ctx, tx, &nodeConnectionStoryReq, copyRequest.ID, subRequest.Workflow.ProjectKey, userId, req.SprintID, false); err != nil {
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

	for _, node := range req.Nodes {
		if node.Type == string(constants.NodeTypeStart) {
			err = s.HistoryService.HistoryStartRequest(ctx, tx, userId, newRequest.ID, node.Id)
			if err != nil {
				return 0, fmt.Errorf("history start request fail: %w", err)
			}
			break
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit fail: %w", err)
	}

	// Notify Start Request With Detail
	// sort nodes by level
	userIdsInt32 := []int32{}
	existingUserIds := make(map[int32]bool)
	for _, node := range req.Nodes {
		if node.Data.Assignee.Id != 0 {
			if !existingUserIds[node.Data.Assignee.Id] {
				userIdsInt32 = append(userIdsInt32, node.Data.Assignee.Id)
				existingUserIds[node.Data.Assignee.Id] = true
			}
		}
	}

	userApiMap := map[int32]results.UserApiDataResult{}
	if len(userIdsInt32) > 0 {
		assigneeResult, err := s.UserAPI.FindUsersByUserIds(userIdsInt32)
		if err != nil {
			return 0, err
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

	sort.Slice(req.Nodes, func(i, j int) bool {
		return *req.Nodes[i].Level < *req.Nodes[j].Level
	})

	userTasks := map[int32][]requests.Node{}
	for _, node := range req.Nodes {
		if node.Type == string(constants.NodeTypeTask) || node.Type == string(constants.NodeTypeInput) || node.Type == string(constants.NodeTypeApproval) {
			userTasks[node.Data.Assignee.Id] = append(userTasks[node.Data.Assignee.Id], node)
		}
	}

	for _, tasks := range userTasks {
		subject := fmt.Sprintf("[ZODC] You’ve Been Assigned to a New Request – “%s”", req.Title)
		bodyTasks := fmt.Sprintf("<p>Hi %s,</p>", mapUser(&tasks[0].Data.Assignee.Id).Name)
		bodyTasks += fmt.Sprintf("<p>You have been added as a participant in the request “%s”, which has just been started by the Product Owner.</p>", req.Title)
		bodyTasks += "<p>Below is a list of tasks assigned to you in this request:</p><br>"

		userId := int32(0)

		for i, node := range tasks {
			bodyTasks += fmt.Sprintf("<p> => Task %d: %s</p>", i+1, node.Data.Title)
			userId = node.Data.Assignee.Id
		}

		bodyTasks += "<br><p>Please log in to the ZODC system to review the workflow, complete your input forms, and track task dependencies.</p"
		bodyTasks += "<p>If you have any questions regarding this request, feel free to contact the PO or check the request details in ZODC.</p>"
		bodyTasks += "<p>Best regards,<br>ZODC System</p>"

		s.NotificationService.NotifyStartRequestWithDetail(ctx, userId, subject, bodyTasks)
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

// setTimeToEndOfWorkday sets the time part of a date to 17:30 (end of workday)
// while preserving the date part
func setTimeToEndOfWorkday(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}

	// Create new time with same date but time set to 17:30
	newTime := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		17, // Hour: 17
		30, // Minute: 30
		0,  // Second: 0
		0,  // Nanosecond: 0
		t.Location(),
	)

	return &newTime
}

func (s *WorkflowService) UpdateWorkflowHandler(ctx context.Context, req *requests.UpdateWorkflow, workflowId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("begin tx fail: %w", err)
	}
	defer tx.Rollback()

	workflow, err := s.WorkflowRepo.FindOneWorkflowByWorkflowId(ctx, s.DB, workflowId)
	if err != nil {
		return fmt.Errorf("find workflow fail: %w", err)
	}

	if req.Decoration != "" {
		workflow.Decoration = req.Decoration
	}
	if req.Description != "" {
		workflow.Description = req.Description
	}
	if req.Title != "" {
		workflow.Title = req.Title
	}
	if req.Type != "" {
		workflow.Type = req.Type
	}
	if req.CategoryId != 0 {
		workflow.CategoryID = req.CategoryId
	}

	workflowModel := model.Workflows{}
	if err := utils.Mapper(workflow, &workflowModel); err != nil {
		return fmt.Errorf("map workflow fail: %w", err)
	}

	if err := s.WorkflowRepo.UpdateWorkflow(ctx, tx, workflowModel); err != nil {
		return fmt.Errorf("update workflow fail: %w", err)
	}

	requestTemplate, err := s.RequestRepo.FindOneRequestTemplateByWorkflowId(ctx, s.DB, workflowId)
	if err != nil {
		return fmt.Errorf("find request template fail: %w", err)
	}

	requestTemplate.Title = req.Title

	if err := s.RequestRepo.UpdateRequest(ctx, tx, requestTemplate); err != nil {
		return fmt.Errorf("update request fail: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}
