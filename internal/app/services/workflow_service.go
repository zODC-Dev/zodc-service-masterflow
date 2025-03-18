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
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type WorkflowService struct {
	db           *sql.DB
	workflowRepo *repositories.WorkflowRepository
	formRepo     *repositories.FormRepository
	categoryRepo *repositories.CategoryRepository
	userAPI      *externals.UserAPI
}

func NewWorkflowService(db *sql.DB, workflowRepo *repositories.WorkflowRepository, formRepo *repositories.FormRepository, categoryRepo *repositories.CategoryRepository, userAPI *externals.UserAPI) *WorkflowService {
	return &WorkflowService{
		db:           db,
		workflowRepo: workflowRepo,
		formRepo:     formRepo,
		categoryRepo: categoryRepo,
		userAPI:      userAPI,
	}
}

func (s *WorkflowService) CreateWorkFlow(ctx context.Context, tx *sql.Tx, workflowData interface{}) (model.Workflows, error) {
	workflow := model.Workflows{
		Currentversion: 1,
		IsArchived:     false,
	}
	if err := utils.Mapper(workflowData, &workflow); err != nil {
		return workflow, fmt.Errorf("mapping workflow failed: %w", err)
	}

	return s.workflowRepo.CreateWorkflow(ctx, tx, workflow)
}

func (s *WorkflowService) CreateWorkFlowVersion(ctx context.Context, tx *sql.Tx, workflowId int32, hasSubWorkflow bool) (model.WorkflowVersions, error) {
	workFlowVersion := model.WorkflowVersions{
		Version:        1,
		WorkflowID:     workflowId,
		HasSubWorkflow: hasSubWorkflow,
	}

	return s.workflowRepo.CreateWorkflowVersion(ctx, tx, workFlowVersion)
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
	request, err := s.workflowRepo.FindOneRequestByRequestId(ctx, s.db, requestId)
	if err != nil {
		return fmt.Errorf("request not found")
	}

	// Update request status to in processing
	request.Status = string(constants.RequestStatusInProcessing)

	requestModel := model.Requests{}
	if err := utils.Mapper(request, &requestModel); err != nil {
		return err
	}
	if err := s.workflowRepo.UpdateRequest(ctx, tx, requestModel); err != nil {
		return err
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
					if err := s.workflowRepo.UpdateConnection(ctx, tx, request.Connections[j]); err != nil {
						return err
					}

					nextNodeIds[request.Connections[j].ToNodeID] = true
				}
			}
		}
	}

	for i := range request.Nodes {
		if nextNodeIds[request.Nodes[i].ID] {
			if err := s.UpdateNodeStatusToInProcessing(ctx, tx, request.Nodes[i]); err != nil {
				return err
			}
		}
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *WorkflowService) UpdateNodeStatusToInProcessing(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	node.Status = string(constants.NodeStatusInProccessing)

	// Update Node
	if err := s.workflowRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	// If Story Or Sub Workflow
	if err := s.RunWorkflowIfItStoryOrSubWorkflow(ctx, tx, node); err != nil {
		return err
	}

	return nil
}

func (s *WorkflowService) CreateNodesConnectionsStories(ctx context.Context, tx *sql.Tx, req *requests.NodesConnectionsStories, requestId int32) error {
	formSystems, err := s.formRepo.FindAllFormSystem(ctx, s.db)
	if err != nil {
		return err
	}

	//Create Stories
	for _, storyReq := range req.Stories {
		category, err := s.categoryRepo.FindOneCategoryByKey(ctx, s.db, storyReq.CategoryKey)
		if err != nil {
			return fmt.Errorf("category key not found: %w", err)
		}
		storyReq.CategoryId = category.ID

		// Create Story Workflow
		storyWorkflow, err := s.CreateWorkFlow(ctx, tx, storyReq)
		if err != nil {
			return fmt.Errorf("create Story Workflow Fail: %w", err)
		}

		storyWorkflowVersion, err := s.CreateWorkFlowVersion(ctx, tx, storyWorkflow.ID, false)
		if err != nil {
			return fmt.Errorf("create Story Workflow Version Fail: %w", err)
		}

		storyRequestModel := model.Requests{
			WorkflowVersionID: storyWorkflowVersion.ID,
			IsTemplate:        false,
			Status:            "IN_ACTIVE",
		}
		storyRequest, err := s.workflowRepo.CreateRequest(ctx, tx, storyRequestModel)
		if err != nil {
			return fmt.Errorf("create Story Request Fail: %w", err)
		}

		// Create Story Node
		storyNode := model.Nodes{
			ID:        storyReq.Node.Id,
			RequestID: storyRequest.ID,

			X:      storyReq.Node.Position.X,
			Y:      storyReq.Node.Position.Y,
			Width:  storyReq.Node.Size.Width,
			Height: storyReq.Node.Size.Height,

			Type: storyReq.Node.Type,

			// Data
			Title:      &storyReq.Node.Data.Title,
			AssigneeID: &storyReq.Node.Data.Assignee.Id,

			SubRequestID: &storyRequest.ID,
		}

		if err := s.workflowRepo.CreateWorkflowNodes(ctx, tx, []model.Nodes{storyNode}); err != nil {
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

				ParentID: &storyNodeReq.ParentId,

				AssigneeID: &storyNodeReq.Data.Assignee.Id,

				// Data
				Title:   &storyNodeReq.Data.Title,
				EndType: &storyNodeReq.Data.EndType,
				DueIn:   &storyNodeReq.Data.DueIn,
			}

			// Form Type System Tag Story // Create Form Data
			for _, formSystem := range formSystems {
				if formSystem.Tag == storyNodeReq.Type {
					// Create Form Data
					formData := model.FormData{
						FormTemplateVersionID: formSystem.Version.ID,
					}

					formData, err = s.formRepo.CreateFormData(ctx, tx, formData)
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
						err := s.formRepo.CreateFormFieldDatas(ctx, tx, formFieldDatas)
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
			err = s.workflowRepo.CreateWorkflowNodes(ctx, tx, storyNodes)
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
						Type:       connReq.Type,
						RequestID:  storyRequest.ID,
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
			err = s.workflowRepo.CreateWorkflowConnections(ctx, tx, storyConnections)
			if err != nil {
				return fmt.Errorf("create Story Connection Fail: %w", err)
			}
		}

	}

	// Create Workflow Nodes
	workflowNodes := []model.Nodes{}
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

			SubRequestID: workflowNodeReq.Data.SubWorkflowVersionID,

			// Data
			Title:   &workflowNodeReq.Data.Title,
			EndType: &workflowNodeReq.Data.EndType,
			DueIn:   &workflowNodeReq.Data.DueIn,
		}

		for _, formSystem := range formSystems {
			if formSystem.Tag == workflowNodeReq.Type {
				// Create Form Data
				formData := model.FormData{
					FormTemplateVersionID: formSystem.Version.ID,
				}

				formData, err = s.formRepo.CreateFormData(ctx, tx, formData)
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
					err := s.formRepo.CreateFormFieldDatas(ctx, tx, formFieldDatas)
					if err != nil {
						return fmt.Errorf("create form fields data fail: %w", err)
					}
				}

				workflowNode.FormDataID = &formData.ID

				break
			}
		}

		workflowNodes = append(workflowNodes, workflowNode)
	}
	if len(workflowNodes) > 0 {
		err = s.workflowRepo.CreateWorkflowNodes(ctx, tx, workflowNodes)
		if err != nil {
			return fmt.Errorf("create Workflow Nodes Fail: %w", err)
		}
	}

	// Create Workflow Connections
	workflowConnections := []model.Connections{}
	for _, workflowConnectionReq := range req.Connections {
		workflowConnection := model.Connections{
			ID: workflowConnectionReq.Id,

			Type: workflowConnectionReq.Type,

			FromNodeID: workflowConnectionReq.From,
			ToNodeID:   workflowConnectionReq.To,

			RequestID: requestId,
		}

		workflowConnections = append(workflowConnections, workflowConnection)
	}
	if len(workflowConnections) > 0 {
		err = s.workflowRepo.CreateWorkflowConnections(ctx, tx, workflowConnections)
		if err != nil {
			return fmt.Errorf("create Workflow Connections Fail: %w", err)
		}
	}

	return nil
}

// Handlers
func (s *WorkflowService) CreateWorkflowHandler(ctx context.Context, req *requests.CreateWorkflow) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create Workflow
	workflow, err := s.CreateWorkFlow(ctx, tx, req)
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
	workflowVersion, err := s.CreateWorkFlowVersion(ctx, tx, workflow.ID, hasSubWorkflow)
	if err != nil {
		return fmt.Errorf("create Main Workflow Version Fail: %w", err)
	}

	// Create Request
	requestModel := model.Requests{
		WorkflowVersionID: workflowVersion.ID,
		IsTemplate:        true,
		Status:            string(constants.RequestStatusTodo),
	}
	request, err := s.workflowRepo.CreateRequest(ctx, tx, requestModel)
	if err != nil {
		return err
	}

	nodeConnectionStoryReq := requests.NodesConnectionsStories{}
	if err := utils.Mapper(req.NodesConnectionsStories, &nodeConnectionStoryReq); err != nil {
		return err
	}
	if err := s.CreateNodesConnectionsStories(ctx, tx, &req.NodesConnectionsStories, request.ID); err != nil {
		return err
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *WorkflowService) FindAllWorkflowHandler(ctx context.Context, workflowTemplateQueryParams queryparams.WorkflowQueryParam) ([]responses.WorkflowResponse, error) {
	workflowResponses := []responses.WorkflowResponse{}

	workflows, err := s.workflowRepo.FindAllWorkflowTemplates(ctx, s.db, workflowTemplateQueryParams)
	if err != nil {
		return workflowResponses, err
	}

	for _, workflow := range workflows {

		//Mapping workflow response
		workflowResponse := responses.WorkflowResponse{
			IsArchived: workflow.IsArchived,
		}
		if err := utils.Mapper(workflow, &workflowResponse); err != nil {
			return workflowResponses, err
		}

		workflowResponse.Id = workflow.Version.ID

		workflowResponse.Version = workflow.Version.Version

		workflowResponses = append(workflowResponses, workflowResponse)
	}

	return workflowResponses, nil
}

func (s *WorkflowService) FindOneWorkflowDetailHandler(ctx context.Context, requestId int32) (responses.WorkflowDetailResponse, error) {
	workflowResponse := responses.WorkflowDetailResponse{}

	request, err := s.workflowRepo.FindOneRequestByRequestId(ctx, s.db, requestId)
	if err != nil {
		return workflowResponse, err
	}

	//Mapping workflow response
	if err := utils.Mapper(request.Workflow, &workflowResponse); err != nil {
		return workflowResponse, err
	}

	workflowResponse.Version = request.Version.Version
	workflowResponse.IsArchived = request.Workflow.IsArchived

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

		// Map Response
		nodeResponse, err := s.MapToWorkflowNodeResponse(node)
		if err != nil {
			return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
		}

		story := responses.StoryResponse{
			Node: nodeResponse,

			Type:        request.Workflow.Type,
			Decoration:  request.Workflow.Decoration,
			Description: request.Workflow.Description,
			Title:       request.Workflow.Title,
		}

		storiesResponse = append(storiesResponse, story)

		// Story nodes
		if node.SubRequestID != nil {
			storyNodes, err := s.workflowRepo.FindAllNodeByRequestId(ctx, s.db, *node.SubRequestID)
			if err != nil {
				return workflowResponse, err
			}

			for _, storyNode := range storyNodes {
				nodeResponse, err := s.MapToWorkflowNodeResponse(storyNode)
				if err != nil {
					return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
				}

				workflowResponse.Nodes = append(workflowResponse.Nodes, nodeResponse)
			}

		}

		// Story connections
		if node.SubRequestID != nil {
			storyConnections, err := s.workflowRepo.FindAllConnectionByREquestId(ctx, s.db, *node.SubRequestID)

			if err != nil {
				return workflowResponse, err
			}

			for _, storyConnection := range storyConnections {
				connectionResponse := responses.ConnectionResponse{
					Id:   storyConnection.ID,
					To:   storyConnection.ToNodeID,
					From: storyConnection.FromNodeID,
					Type: storyConnection.Type,
				}

				workflowResponse.Connections = append(workflowResponse.Connections, connectionResponse)
			}

		}
	}
	request.Nodes = request.Nodes[:i]
	workflowResponse.Stories = storiesResponse

	// Nodes
	for _, node := range request.Nodes {

		nodeResponse, err := s.MapToWorkflowNodeResponse(node)
		if err != nil {
			return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
		}

		workflowResponse.Nodes = append(workflowResponse.Nodes, nodeResponse)
	}

	// Connections
	for _, connection := range request.Connections {
		connectionResponse := responses.ConnectionResponse{
			Id:   connection.ID,
			To:   connection.ToNodeID,
			From: connection.FromNodeID,
			Type: connection.Type,
		}

		workflowResponse.Connections = append(workflowResponse.Connections, connectionResponse)
	}

	// Add assignee
	userIds := make([]int32, 0, len(workflowResponse.Nodes))
	for _, node := range workflowResponse.Nodes {
		userIds = append(userIds, node.Data.Assignee.Id)
	}

	results, err := s.userAPI.FindUsersByUserIds(userIds)
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

func (s *WorkflowService) StartWorkflowHandler(ctx context.Context, req requests.StartWorkflow) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requestModel := model.Requests{
		Title:             req.Title,
		IsTemplate:        false,
		WorkflowVersionID: req.WorkflowVersionId,
		Status:            string(constants.RequestStatusInProcessing),
	}
	request, err := s.workflowRepo.CreateRequest(ctx, tx, requestModel)
	if err != nil {
		return err
	}

	nodeConnectionStoryReq := requests.NodesConnectionsStories{}
	if err := utils.Mapper(req.NodesConnectionsStories, &nodeConnectionStoryReq); err != nil {
		return err
	}
	if err := s.CreateNodesConnectionsStories(ctx, tx, &req.NodesConnectionsStories, request.ID); err != nil {
		return err
	}

	//Logic Start Workflow
	if err := s.RunWorkflow(ctx, tx, request.ID); err != nil {
		return err
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil

}

func (s *WorkflowService) StartNodeHandler(ctx context.Context, nodeId string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//
	node, err := s.workflowRepo.FindOneNodeByNodeId(ctx, s.db, nodeId)
	if err != nil {
		return err
	}

	// Update Current Node Status To In Process
	node.Status = string(constants.NodeStatusInProccessing)

	// Set actual start time
	now := time.Now()
	node.ActualStartTime = &now

	if err := s.workflowRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *WorkflowService) CompleteNodeHandler(ctx context.Context, nodeId string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	node, err := s.workflowRepo.FindOneNodeByNodeId(ctx, s.db, nodeId)
	if err != nil {
		return err
	}

	if node.Type == string(constants.NodeTypeStory) || node.Type == string(constants.NodeTypeSubWorkflow) {
		return fmt.Errorf("story or sub workflow is auto complete by system, cant mark as complete by user")
	}

	// Update Current Node Status To Completed
	node.Status = string(constants.NodeStatusCompleted)

	// Set actual finish time
	now := time.Now()
	node.ActualFinishTime = &now

	if err := s.workflowRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	connectionsToNode, err := s.workflowRepo.FindConnectionsWithToNodesByFromNodeId(ctx, s.db, node.ID)
	if err != nil {
		return err
	}

	// Update Next Node To In Processing
	for i := range connectionsToNode {

		// Update connection to completeed
		connectionsToNode[i].IsCompleted = true
		connectionModel := model.Connections{}
		if err := utils.Mapper(connectionsToNode[i], &connectionModel); err != nil {
			return err
		}
		if err := s.workflowRepo.UpdateConnection(ctx, tx, connectionModel); err != nil {
			return err
		}

		// If Prevous Nodes not finish yet // If More than one node not completed then next node dont need to update status
		isUpdateNodeStatus := true

		connections, err := s.workflowRepo.FindConnectionsByToNodeId(ctx, s.db, connectionsToNode[i].Node.ID)
		if err != nil {
			return err
		}
		for j := range connections {
			if !connections[j].IsCompleted {
				isUpdateNodeStatus = false
				break
			}
		}

		if isUpdateNodeStatus {
			// If Node is End Node
			if connectionsToNode[i].Node.Type == string(constants.NodeTypeEnd) {
				// Update end node to completed
				connectionsToNode[i].Node.Status = string(constants.NodeStatusCompleted)
				if err := s.workflowRepo.UpdateNode(ctx, tx, connectionsToNode[i].Node); err != nil {
					return err
				}

				// Mark request completed
				request, err := s.workflowRepo.FindRequestByNodeId(ctx, s.db, connectionsToNode[i].Node.ID)
				if err != nil {
					return err
				}
				request.Status = string(constants.RequestStatusCompleted)
				if err := s.workflowRepo.UpdateRequest(ctx, tx, request); err != nil {
					return err
				}
			} else {
				if err := s.UpdateNodeStatusToInProcessing(ctx, tx, connectionsToNode[i].Node); err != nil {
					return err
				}
			}
		}
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}
