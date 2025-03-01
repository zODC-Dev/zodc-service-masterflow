package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type WorkflowService struct {
	db           *sql.DB
	workflowRepo *repositories.WorkflowRepository
	formRepo     *repositories.FormRepository
}

func NewWorkflowService(db *sql.DB, workflowRepo *repositories.WorkflowRepository, formRepo *repositories.FormRepository) *WorkflowService {
	return &WorkflowService{
		db:           db,
		workflowRepo: workflowRepo,
		formRepo:     formRepo,
	}
}

func (s *WorkflowService) CreateWorkFlow(ctx context.Context, tx *sql.Tx, workflowData interface{}) (model.Workflows, error) {
	workflow := model.Workflows{}
	if err := utils.Mapper(workflowData, &workflow); err != nil {
		return workflow, fmt.Errorf("mapping workflow failed: %w", err)
	}

	workflow, err := s.workflowRepo.CreateWorkflow(ctx, tx, workflow)
	if err != nil {
		return workflow, fmt.Errorf("create workflow failed: %w", err)
	}

	return workflow, nil
}

func (s *WorkflowService) CreateWorkFlowVersion(ctx context.Context, tx *sql.Tx, workflowId int32) (model.WorkflowVersions, error) {
	workFlowVersion := model.WorkflowVersions{
		Version:    1,
		WorkflowID: workflowId,
	}

	workFlowVersion, err := s.workflowRepo.CreateWorkflowVersion(ctx, tx, workFlowVersion)
	if err != nil {
		return workFlowVersion, fmt.Errorf("create workflow template version failed: %w", err)
	}

	return workFlowVersion, nil

}

func (s *WorkflowService) MapToWorkflowNodeResponse(node model.WorkflowNodes) (responses.NodeResponse, error) {
	nodeDataResponse := responses.NodeDataResponse{}
	if err := utils.Mapper(node, &nodeDataResponse); err != nil {
		return responses.NodeResponse{}, err
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

// Handlers
func (s *WorkflowService) CreateWorkFlowHandler(ctx context.Context, req *requests.WorkflowRequest) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare Form System
	formSystems, err := s.formRepo.FindAllFormSystem(ctx, s.db)
	if err != nil {
		return err
	}

	workflow, err := s.CreateWorkFlow(ctx, tx, req)
	if err != nil {
		return fmt.Errorf("create Main Workflow Fail: %w", err)
	}

	workflowVersion, err := s.CreateWorkFlowVersion(ctx, tx, workflow.ID)
	if err != nil {
		return fmt.Errorf("create Main Workflow Version Fail: %w", err)
	}

	//Create Stories
	for _, storyReq := range req.Stories {

		storyReq.Decoration = workflow.Decoration
		storyReq.Description = workflow.Description

		storyWorkflow, err := s.CreateWorkFlow(ctx, tx, storyReq)
		if err != nil {
			return fmt.Errorf("create Story Workflow Fail: %w", err)
		}

		storyWorkflowVersion, err := s.CreateWorkFlowVersion(ctx, tx, storyWorkflow.ID)
		if err != nil {
			return fmt.Errorf("create Story Workflow Version Fail: %w", err)
		}

		storyWorkflowNode := model.WorkflowNodes{
			ID:                storyReq.Node.Id,
			WorkflowVersionID: workflowVersion.ID,

			X:      storyReq.Node.Position.X,
			Y:      storyReq.Node.Position.Y,
			Width:  storyReq.Node.Size.Width,
			Height: storyReq.Node.Size.Height,

			Type: storyReq.Node.Type,

			// Data
			Title: &storyReq.Node.Data.Title,

			//subworkflow ??? // can delete if it wrong
			SubWorkflowID: &storyWorkflow.ID,
		}

		err = s.workflowRepo.CreateWorkflowNodes(ctx, tx, []model.WorkflowNodes{storyWorkflowNode})
		if err != nil {
			return fmt.Errorf("create Story Workflow MAIN Node Fail: %w", err)
		}

		// Create Story Nodes
		storyNodes := []model.WorkflowNodes{}

		i := 0
		for _, storyNodeReq := range req.Nodes {
			if storyNodeReq.ParentId != storyReq.Node.Id {
				req.Nodes[i] = storyNodeReq
				i++
				continue
			}

			storyNode := model.WorkflowNodes{
				ID:                storyNodeReq.Id,
				WorkflowVersionID: storyWorkflowVersion.ID,

				X:      storyNodeReq.Position.X,
				Y:      storyNodeReq.Position.Y,
				Width:  storyNodeReq.Size.Width,
				Height: storyNodeReq.Size.Height,

				Type: storyNodeReq.Type,

				ParentID: &storyNodeReq.ParentId,

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

		err = s.workflowRepo.CreateWorkflowNodes(ctx, tx, storyNodes)
		if err != nil {
			return fmt.Errorf("create Story Node Fail: %w", err)
		}

		// Create Story Connections
		storyConnections := []model.WorkflowConnections{}

		i = 0
		for _, connReq := range req.Connections {
			shouldKeepConnection := true

			for _, storyNode := range storyNodes {

				if storyNode.ID == connReq.From {
					shouldKeepConnection = false

					storyConnection := model.WorkflowConnections{
						ID:                 connReq.Id,
						FromWorkflowNodeID: connReq.From,
						ToWorkflowNodeID:   connReq.To,
						Type:               connReq.Type,
						WorkflowVersionID:  storyWorkflowVersion.ID,
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

		err = s.workflowRepo.CreateWorkflowConnections(ctx, tx, storyConnections)
		if err != nil {
			return fmt.Errorf("create Story Connection Fail: %w", err)
		}
	}

	// Create workflow node
	workflowNodes := []model.WorkflowNodes{}

	for _, workflowNodeReq := range req.Nodes {
		workflowNode := model.WorkflowNodes{
			ID:                workflowNodeReq.Id,
			WorkflowVersionID: workflowVersion.ID,

			X:      workflowNodeReq.Position.X,
			Y:      workflowNodeReq.Position.Y,
			Width:  workflowNodeReq.Size.Width,
			Height: workflowNodeReq.Size.Height,

			Type: workflowNodeReq.Type,

			// ParentID: &workflowNodeReq.ParentId,

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

	err = s.workflowRepo.CreateWorkflowNodes(ctx, tx, workflowNodes)
	if err != nil {
		return fmt.Errorf("create Workflow Nodes Fail: %w", err)
	}

	// Create workflow connection
	workflowConnections := []model.WorkflowConnections{}

	for _, workflowConnectionReq := range req.Connections {
		workflowConnection := model.WorkflowConnections{
			ID: workflowConnectionReq.Id,

			Type: workflowConnectionReq.Type,

			FromWorkflowNodeID: workflowConnectionReq.From,
			ToWorkflowNodeID:   workflowConnectionReq.To,

			WorkflowVersionID: workflowVersion.ID,
		}

		workflowConnections = append(workflowConnections, workflowConnection)
	}

	err = s.workflowRepo.CreateWorkflowConnections(ctx, tx, workflowConnections)
	if err != nil {
		return fmt.Errorf("create Workflow Connections Fail: %w", err)
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *WorkflowService) FindAllWorkflowHandler(ctx context.Context) ([]responses.WorkflowResponse, error) {
	workflowResponses := []responses.WorkflowResponse{}

	workflows, err := s.workflowRepo.FindAllWorkflowTemplates(ctx, s.db)
	if err != nil {
		return workflowResponses, err
	}

	for _, workflow := range workflows {

		//Mapping workflow response
		workflowResponse := responses.WorkflowResponse{}
		if err := utils.Mapper(workflow, &workflowResponse); err != nil {
			return workflowResponses, err
		}

		workflowResponse.Version = workflow.Version.Version

		workflowResponses = append(workflowResponses, workflowResponse)
	}

	return workflowResponses, nil
}

func (s *WorkflowService) FindOneWorkflowDetailHandler(ctx context.Context, workflowVersionId int32) (responses.WorkflowDetailResponse, error) {
	workflowResponse := responses.WorkflowDetailResponse{}

	workflow, err := s.workflowRepo.FindOneWorkflowDetailByWorkflowVersionId(ctx, s.db, workflowVersionId)
	if err != nil {
		return workflowResponse, err
	}

	//Mapping workflow response
	if err := utils.Mapper(workflow, &workflowResponse); err != nil {
		return workflowResponse, err
	}

	workflowResponse.Version = workflow.Version.Version

	// Stories
	storiesResponse := []responses.StoryResponse{}

	i := 0
	for _, node := range workflow.Nodes {
		if node.Type != "STORY" {
			workflow.Nodes[i] = node
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

			Type:        workflow.Type,
			Decoration:  workflow.Decoration,
			Description: workflow.Description,
			Title:       workflow.Title,
		}

		storiesResponse = append(storiesResponse, story)
	}
	workflow.Nodes = workflow.Nodes[:i]
	workflowResponse.Stories = storiesResponse

	// Nodes
	nodesResponse := []responses.NodeResponse{}
	for _, node := range workflow.Nodes {

		nodeResponse, err := s.MapToWorkflowNodeResponse(node)
		if err != nil {
			return workflowResponse, fmt.Errorf("map workflow node response fail: %w", err)
		}

		nodesResponse = append(nodesResponse, nodeResponse)
	}
	workflowResponse.Nodes = nodesResponse

	// Connections
	connectionsResponse := []responses.ConnectionResponse{}
	for _, connection := range workflow.Connections {
		connectionResponse := responses.ConnectionResponse{
			Id:   connection.ID,
			To:   connection.ToWorkflowNodeID,
			From: connection.FromWorkflowNodeID,
			Type: connection.Type,
		}

		connectionsResponse = append(connectionsResponse, connectionResponse)
	}
	workflowResponse.Connections = connectionsResponse

	return workflowResponse, nil
}
