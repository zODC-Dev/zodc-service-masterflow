package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type NodeService struct {
	DB              *sql.DB
	NodeRepo        *repositories.NodeRepository
	ConnectionRepo  *repositories.ConnectionRepository
	RequestRepo     *repositories.RequestRepository
	WorkflowService *WorkflowService
	NatsClient      *nats.NATSClient
}

func NewNodeService(cfg NodeService) *NodeService {
	nodeService := NodeService{
		DB:              cfg.DB,
		NodeRepo:        cfg.NodeRepo,
		ConnectionRepo:  cfg.ConnectionRepo,
		RequestRepo:     cfg.RequestRepo,
		WorkflowService: cfg.WorkflowService,
		NatsClient:      cfg.NatsClient,
	}
	return &nodeService
}

// Function
func (s *NodeService) UpdateNodeStatusToInProcessing(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	node.Status = string(constants.NodeStatusInProgress)

	// Update Node
	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	// If Story Or Sub Workflow
	if err := s.WorkflowService.RunWorkflowIfItStoryOrSubWorkflow(ctx, tx, node); err != nil {
		return err
	}

	return nil
}

//

func (s *NodeService) StartNodeHandler(ctx context.Context, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("start node handler fail: %w", err)
	}
	defer tx.Rollback()

	//
	node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return fmt.Errorf("find node by node id fail: %w", err)
	}

	// Update Current Node Status To In Process
	node.Status = string(constants.NodeStatusInProgress)

	// Set actual start time
	now := time.Now()
	node.ActualStartTime = &now

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return fmt.Errorf("update node fail: %w", err)
	}

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) CompleteNodeHandler(ctx context.Context, nodeId string, userId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
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
	node.ActualEndTime = &now

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	connectionsToNode, err := s.ConnectionRepo.FindConnectionsWithToNodesByFromNodeId(ctx, s.DB, node.ID)
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
		if err := s.ConnectionRepo.UpdateConnection(ctx, tx, connectionModel); err != nil {
			return err
		}

		// If Prevous Nodes not finish yet // If More than one node not completed then next node dont need to update status
		isUpdateNodeStatus := true

		connections, err := s.ConnectionRepo.FindConnectionsByToNodeId(ctx, s.DB, connectionsToNode[i].Node.ID)
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
				if err := s.NodeRepo.UpdateNode(ctx, tx, connectionsToNode[i].Node); err != nil {
					return err
				}

				// Mark request completed
				request, err := s.RequestRepo.FindRequestByNodeId(ctx, s.DB, connectionsToNode[i].Node.ID)
				if err != nil {
					return err
				}
				request.Status = string(constants.RequestStatusCompleted)
				if err := s.RequestRepo.UpdateRequest(ctx, tx, request); err != nil {
					return err
				}
			} else {
				connectionsToNode[i].Node.IsCurrent = true
				if err := s.UpdateNodeStatusToInProcessing(ctx, tx, connectionsToNode[i].Node); err != nil {
					return err
				}

			}
		}
	}

	// send notification
	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   "Task completed",
		Body:      fmt.Sprintf("Task completed: %s – %s has marked this task as done.", node.Title, userId),
	}
	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("marshal notification failed: %w", err)
	}
	s.NatsClient.Publish("notifications", notificationBytes)

	// Calculate Request Process
	request, _ := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, node.RequestID)
	totalCompletedNode := 0
	totalNode := len(request.Nodes)
	for _, requestNode := range request.Nodes {
		if requestNode.Status == string(constants.NodeStatusCompleted) {
			totalCompletedNode++
		}

		if requestNode.Type == string(constants.NodeTypeStart) || requestNode.Type == string(constants.NodeTypeEnd) {
			totalNode--
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

	// Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) ApproveNodeHandler(ctx context.Context, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	if node.Type != string(constants.NodeTypeCondition) {
		return fmt.Errorf("node is not a condition node")
	}

	node.IsApproved = true

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	//

	// Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) GetNodeFormWithPermission(ctx context.Context, nodeId string, permission string) ([]responses.NodeFormDetailResponse, error) {

	nodeForm, err := s.NodeRepo.FindAllNodeFormByNodeIdAndPermission(ctx, s.DB, nodeId, permission)
	if err != nil {
		return nil, err
	}

	response := []responses.NodeFormDetailResponse{}

	for _, nodeForm := range nodeForm {
		fieldMap := map[int32]string{}
		for _, formTemplateField := range nodeForm.FormTemplateFields {
			fieldMap[formTemplateField.ID] = formTemplateField.FieldID
		}

		formTemplate := responses.FormTemplateFindAll{}
		utils.Mapper(nodeForm.FormTemplates, &formTemplate)

		formTemplateFields := []responses.FormTemplateFieldsFindAll{}
		utils.Mapper(nodeForm.FormTemplateFields, &formTemplateFields)

		formDatas := []responses.NodeFormDataResponse{}
		for _, formData := range nodeForm.FormFieldData {
			formDatas = append(formDatas, responses.NodeFormDataResponse{
				FieldId: fieldMap[formData.FormTemplateFieldID],
				Value:   formData.Value,
			})
		}

		response = append(response, responses.NodeFormDetailResponse{
			Template:    formTemplate,
			Fields:      formTemplateFields,
			Data:        formDatas,
			DataId:      nodeForm.NodeForms.DataID,
			IsSubmitted: nodeForm.NodeForms.IsSubmitted,
			IsApproved:  nodeForm.NodeForms.IsApproved,
		})
	}

	return response, nil
}

func (s *NodeService) GetNodeJiraForm(ctx context.Context, nodeId string) (responses.JiraFormDetailResponse, error) {
	nodeJiraForm, err := s.NodeRepo.FindJiraFormByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return responses.JiraFormDetailResponse{}, err
	}

	formTemplate := responses.FormTemplateFindAll{}
	utils.Mapper(nodeJiraForm.FormTemplates, &formTemplate)

	formTemplateFields := []responses.FormTemplateFieldsFindAll{}
	utils.Mapper(nodeJiraForm.FormTemplateFields, &formTemplateFields)

	fieldMap := map[int32]string{}
	for _, formTemplateField := range formTemplateFields {
		fieldMap[formTemplateField.ID] = formTemplateField.FieldID
	}

	formDatas := []responses.NodeFormDataResponse{}
	for _, formData := range nodeJiraForm.FormFieldData {
		formDatas = append(formDatas, responses.NodeFormDataResponse{
			FieldId: fieldMap[formData.FormTemplateFieldID],
			Value:   formData.Value,
		})
	}

	response := responses.JiraFormDetailResponse{
		Template: formTemplate,
		Fields:   formTemplateFields,
		Data:     formDatas,
	}

	return response, nil
}
