package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/configs"
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

type NodeService struct {
	DB                  *sql.DB
	NodeRepo            *repositories.NodeRepository
	ConnectionRepo      *repositories.ConnectionRepository
	RequestRepo         *repositories.RequestRepository
	WorkflowService     *WorkflowService
	NatsService         *NatsService
	NatsClient          *nats.NATSClient
	FormRepo            *repositories.FormRepository
	FormService         *FormService
	UserAPI             *externals.UserAPI
	RequestService      *RequestService
	NotificationService *NotificationService
	HistoryService      *HistoryService
	CommentRepository   *repositories.CommentRepository
}

func NewNodeService(cfg NodeService) *NodeService {
	nodeService := NodeService{
		DB:                  cfg.DB,
		NodeRepo:            cfg.NodeRepo,
		ConnectionRepo:      cfg.ConnectionRepo,
		RequestRepo:         cfg.RequestRepo,
		WorkflowService:     cfg.WorkflowService,
		NatsService:         cfg.NatsService,
		NatsClient:          cfg.NatsClient,
		FormRepo:            cfg.FormRepo,
		FormService:         cfg.FormService,
		UserAPI:             cfg.UserAPI,
		RequestService:      cfg.RequestService,
		NotificationService: cfg.NotificationService,
		HistoryService:      cfg.HistoryService,
		CommentRepository:   cfg.CommentRepository,
	}
	return &nodeService
}

// Function
func (s *NodeService) CompleteNodeSwitchCaseLogic(ctx context.Context, tx *sql.Tx, node results.NodeResult, nextNode model.Nodes, userId int32, users results.UserApiResult) error {
	nextNode.IsCurrent = true
	if err := s.NodeRepo.UpdateNode(ctx, tx, nextNode); err != nil {
		return err
	}

	now := time.Now().UTC().Add(7 * time.Hour)

	switch nextNode.Type {
	case string(constants.NodeTypeEnd):
		nextNode.Status = string(constants.NodeStatusCompleted)
		nextNode.ActualStartTime = &now
		nextNode.ActualEndTime = &now
		if err := s.NodeRepo.UpdateNode(ctx, tx, nextNode); err != nil {
			return err
		}

		// Update End Node Request
		endNodeRequest, err := s.RequestRepo.FindOneRequestByNodeId(ctx, s.DB, nextNode.ID)
		if err != nil {
			return err
		}

		// Get All UserId
		userIds := []string{}
		existingUserIds := map[string]bool{}
		for _, node := range endNodeRequest.Nodes {
			if node.AssigneeID != nil {
				if !existingUserIds[strconv.Itoa(int(*node.AssigneeID))] {
					userIds = append(userIds, strconv.Itoa(int(*node.AssigneeID)))
					existingUserIds[strconv.Itoa(int(*node.AssigneeID))] = true
				}
			}
		}

		// Check End Node Type For Status Request
		if nextNode.EndType != nil && *nextNode.EndType == string(constants.NodeEndTypeTerminate) {
			endNodeRequest.TerminatedAt = &now
			endNodeRequest.Status = string(constants.RequestStatusTerminated)

			// History
			if err := s.HistoryService.HistoryTerminateRequest(ctx, tx, endNodeRequest.ID, nextNode.ID); err != nil {
				return err
			}

			// Notify
			if err := s.NotificationService.NotifyRequestTerminated(ctx, endNodeRequest.Title, userIds); err != nil {
				return err
			}
		} else {
			endNodeRequest.CompletedAt = &now
			endNodeRequest.Status = string(constants.RequestStatusCompleted)

			// History
			if err := s.HistoryService.HistoryEndRequest(ctx, tx, endNodeRequest.ID, nextNode.ID); err != nil {
				return err
			}

			// Notify
			if err := s.NotificationService.NotifyRequestCompleted(ctx, endNodeRequest.Title, userIds); err != nil {
				return err
			}

		}

		endNodeRequest.Progress = 100

		endNodeRequestModel := model.Requests{}
		utils.Mapper(endNodeRequest, &endNodeRequestModel)
		if err := s.RequestRepo.UpdateRequest(ctx, tx, endNodeRequestModel); err != nil {
			return err
		}

		// If This Request is Sub Request so need to check to complete main node in main request
		nodeSubRequest, err := s.NodeRepo.FindOneNodeBySubRequestID(ctx, tx, endNodeRequest.ID)
		if err != nil {
			errStr := err.Error()
			if errStr != "qrm: no rows in result set" {
				return err
			}
		} else {
			err = s.CompleteNodeLogic(ctx, tx, nodeSubRequest.ID, userId)
			if err != nil {
				return err
			}
		}
	case string(constants.NodeTypeCondition):
		isTrue := true
		// If Node Aprroval is Reject = true => isTrue = false, alway isTrue=true for another node
		if node.Type == string(constants.NodeTypeApproval) && node.IsRejected {
			isTrue = false
		}
		for _, nodeForm := range node.NodeForms {
			if !nodeForm.IsApproved && !nodeForm.IsRejected {
				if isTrue {
					nodeForm.IsApproved = true
				} else {
					nodeForm.IsRejected = true
				}

				if err := s.NodeRepo.UpdateNodeForm(ctx, tx, nodeForm); err != nil {
					return err
				}

				approveOrRejectUser := model.NodeFormApproveOrRejectUsers{
					UserID:     userId,
					NodeFormID: nodeForm.ID,
					IsApproved: isTrue,
				}
				if err := s.NodeRepo.CreateApproveOrRejectUser(ctx, tx, approveOrRejectUser); err != nil {
					return err
				}
			}
		}

		// This node is auto complete
		conditionNode, err := s.NodeRepo.FindOneNodeByNodeIdTx(ctx, tx, nextNode.ID)
		if err != nil {
			return err
		}

		conditionNode.Status = string(constants.NodeStatusCompleted)
		conditionNode.ActualStartTime = &now
		conditionNode.ActualEndTime = &now

		conditionNodeModel := model.Nodes{}
		utils.Mapper(conditionNode, &conditionNodeModel)
		if err := s.NodeRepo.UpdateNode(ctx, tx, conditionNodeModel); err != nil {
			return err
		}

		// Check next node with condition
		nodeConditionDestinations, err := s.NodeRepo.FindAllNodeConditionDestinationByNodeId(ctx, s.DB, nextNode.ID, isTrue)
		if err != nil {
			return err
		}

		for _, nodeConditionDestination := range nodeConditionDestinations {

			connectionConditionDestinations, err := s.ConnectionRepo.FindConnectionsByToNodeIdTx(ctx, tx, nodeConditionDestination.DestinationNodeID)
			if err != nil {
				return fmt.Errorf("find connections by to node id fail: %w", err)
			}

			for _, connectionConditionDestination := range connectionConditionDestinations {
				conditionNode.Status = string(constants.NodeStatusCompleted)
				if connectionConditionDestination.FromNodeID == conditionNode.ID {
					connectionConditionDestination.IsCompleted = true
					if err := s.ConnectionRepo.UpdateConnection(ctx, tx, connectionConditionDestination); err != nil {
						return err
					}
					break
				}
			}

			isNodeIsCurrentNode := true
			for _, connectionConditionDestination := range connectionConditionDestinations {
				if connectionConditionDestination.FromNodeID == node.ID && !connectionConditionDestination.IsCompleted {
					isNodeIsCurrentNode = false
					break
				}
			}

			if isNodeIsCurrentNode {
				destinationNode, err := s.NodeRepo.FindOneNodeByNodeIdTx(ctx, tx, nodeConditionDestination.DestinationNodeID)
				if err != nil {
					return err
				}

				destinationNodeModel := model.Nodes{}
				utils.Mapper(destinationNode, &destinationNodeModel)
				if err := s.CompleteNodeSwitchCaseLogic(ctx, tx, node, destinationNodeModel, userId, users); err != nil {
					return err
				}
			}
		}
	case string(constants.NodeTypeApproval):
		request, _ := s.RequestRepo.FindOneRequestByNodeId(ctx, s.DB, nextNode.ID)
		s.NotificationService.NotifyNodeApproveNeeded(ctx, request.Title, *nextNode.AssigneeID)
		fallthrough
	case string(constants.NodeTypeInput):
		nextNode.Status = string(constants.NodeStatusInProgress)
		nextNode.ActualStartTime = &now

		if err := s.NodeRepo.UpdateNode(ctx, tx, nextNode); err != nil {
			return err
		}
		fallthrough
	case string(constants.NodeTypeBug):
		fallthrough
	case string(constants.NodeTypeTask):
		// Notify
		if err := s.NotificationService.NotifyTaskAvailable(ctx, nextNode.Title, users.Data[0].ID); err != nil {
			return err
		}

		// History
		if err := s.HistoryService.HistoryNewTask(ctx, tx, nextNode.RequestID, nextNode.ID, *nextNode.AssigneeID); err != nil {
			return err
		}
	case string(constants.NodeTypeNotification):
		// Send Notification
		var cc []string
		if nextNode.CcEmails != nil {
			err := json.Unmarshal([]byte(*nextNode.CcEmails), &cc)
			if err != nil {
				return err
			}
		}
		var to []string
		if nextNode.ToEmails != nil {
			err := json.Unmarshal([]byte(*nextNode.ToEmails), &to)
			if err != nil {
				return err
			}
		}
		var bcc []string
		if nextNode.BccEmails != nil {
			err := json.Unmarshal([]byte(*nextNode.BccEmails), &bcc)
			if err != nil {
				return err
			}
		}
		notification := types.Notification{
			ToEmails:    to,
			ToCcEmails:  cc,
			ToBccEmails: bcc,
		}
		if nextNode.Subject != nil {
			notification.Subject = *nextNode.Subject
		}
		if nextNode.Body != nil {
			notification.Body = *nextNode.Body
		}
		if nextNode.Subject != nil {
			notification.Subject = *nextNode.Subject
		}

		// Is Send Form
		if nextNode.IsSendApprovedForm || nextNode.IsSendRejectedForm {
			request, err := s.RequestRepo.FindOneRequestByRequestIdTx(ctx, tx, nextNode.RequestID)
			if err != nil {
				return err
			}

			notification.Body += "<br><br>"
			for _, nextNodeRequest := range request.Nodes {
				for _, nextNodeForm := range nextNodeRequest.NodeForms {
					if nextNodeRequest.Type == string(constants.NodeTypeApproval) {
						formDataUrl := configs.Env.FE_HOST + "/form-management/review/" + *nextNodeForm.DataID
						if nextNodeForm.IsApproved && nextNode.IsSendApprovedForm {
							notification.Body += fmt.Sprintf("<br><a href=\"%s\">%s</a>", formDataUrl, formDataUrl)
						}

						if nextNodeForm.IsRejected && nextNode.IsSendRejectedForm {
							notification.Body += fmt.Sprintf("<br><a href=\"%s\">%s</a>", formDataUrl, formDataUrl)
						}
					}
				}
			}
		}

		notificationBytes, err := json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("marshal notification failed: %w", err)
		}
		s.NatsClient.Publish("notifications", notificationBytes)

		// Update Node
		if err := s.CompleteNodeLogic(ctx, tx, nextNode.ID, userId); err != nil {
			return err
		}
	case string(constants.NodeTypeSubWorkflow):
		fallthrough
	case string(constants.NodeTypeStory):
		nextNode.ActualStartTime = &now
		if err := s.NodeRepo.UpdateNode(ctx, tx, nextNode); err != nil {
			return fmt.Errorf("update node status to in processing fail: %w", err)
		}

		if err := s.WorkflowService.RunWorkflow(ctx, tx, *nextNode.SubRequestID, userId); err != nil {
			return err
		}
	default:
		return fmt.Errorf("next node type not valid")
	}

	return nil
}

func (s *NodeService) CompleteNodeLogic(ctx context.Context, tx *sql.Tx, nodeId string, userId int32) error {
	node, err := s.NodeRepo.FindOneNodeByNodeIdTx(ctx, tx, nodeId)
	if err != nil {
		return err
	}

	if node.Type != string(constants.NodeTypeStart) && !node.IsCurrent {
		return fmt.Errorf("this node is not eligible to complete the node")
	}

	// Get Current Time
	now := time.Now().UTC().Add(7 * time.Hour)

	//
	node.ActualEndTime = &now
	node.Status = string(constants.NodeStatusCompleted)

	//
	if node.ActualStartTime != nil {
		node.ActualStartTime = &now
	}

	// Update
	nodeModel := model.Nodes{}
	utils.Mapper(node, &nodeModel)
	if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
		return err
	}

	// History
	oldStatus := string(constants.NodeStatusInProgress)
	if err := s.HistoryService.HistoryChangeNodeStatus(ctx, tx, userId, node.RequestID, nodeId, &oldStatus, string(constants.NodeStatusCompleted)); err != nil {
		return err
	}

	// Notify
	if err := s.NotificationService.NotifyTaskCompleted(ctx, tx, nodeModel); err != nil {
		return err
	}

	// Connections
	connectionsToNode, err := s.ConnectionRepo.FindConnectionsWithToNodesByFromNodeIdTx(ctx, tx, node.ID)
	if err != nil {
		return err
	}

	for _, connectionToNode := range connectionsToNode {
		// Update next connection to completed
		connectionToNodeModel := model.Connections{}
		utils.Mapper(connectionToNode, &connectionToNodeModel)
		connectionToNodeModel.IsCompleted = true
		if err := s.ConnectionRepo.UpdateConnection(ctx, tx, connectionToNodeModel); err != nil {
			return err
		}

		// Update next node to is current and prepare for work
		nextNode := connectionToNode.Node

		isUpdateNodeStatus := true
		connections, err := s.ConnectionRepo.FindConnectionsByToNodeIdTx(ctx, tx, nextNode.ID)
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
			users, err := s.UserAPI.FindUsersByUserIds([]int32{*nextNode.AssigneeID})
			if err != nil {
				return err
			}

			if err := s.CompleteNodeSwitchCaseLogic(ctx, tx, node, nextNode, userId, users); err != nil {
				return err
			}
		}
	}

	// Recalculate Request
	if err := s.RequestService.UpdateCalculateRequestProgress(ctx, tx, node.RequestID); err != nil {
		return err
	}

	return nil
}

// handler
func (s *NodeService) StartNodeHandler(ctx context.Context, userId int32, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("start node handler fail: %w", err)
	}
	defer tx.Rollback()

	nodeResult, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return fmt.Errorf("find node by node id fail: %w", err)
	}

	if !nodeResult.IsCurrent {
		return fmt.Errorf("this node is not eligible to start the node")
	}

	node := model.Nodes{}
	utils.Mapper(nodeResult, &node)

	// Update Current Node Status To In Process
	node.Status = string(constants.NodeStatusInProgress)

	// Set actual start time
	now := time.Now().UTC().Add(7 * time.Hour)
	node.ActualStartTime = &now

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return fmt.Errorf("update node fail: %w", err)
	}

	// Sync with Jira
	if err := s.SyncJiraWhenStartNode(ctx, tx, node); err != nil {
		return fmt.Errorf("sync jira when start node fail: %w", err)
	}

	// Notify
	if err := s.NotificationService.NotifyTaskStarted(ctx, tx, node); err != nil {
		return err
	}

	// History
	oldStatus := string(constants.NodeStatusTodo)
	err = s.HistoryService.HistoryChangeNodeStatus(ctx, tx, userId, node.RequestID, nodeId, &oldStatus, string(constants.NodeStatusInProgress))
	if err != nil {
		return err
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

	if err := s.CompleteNodeLogic(ctx, tx, nodeId, userId); err != nil {
		return err
	}

	nodeResult, err := s.NodeRepo.FindOneNodeByNodeIdTx(ctx, tx, nodeId)
	if err != nil {
		return fmt.Errorf("find node by node id fail: %w", err)
	}

	node := model.Nodes{}
	utils.Mapper(nodeResult, &node)

	// Sync with Jira
	if err := s.SyncJiraWhenCompleteNode(ctx, tx, node); err != nil {
		return fmt.Errorf("sync jira when complete node fail: %w", err)
	}

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

		fieldsResponse := [][]responses.FormTemplateFieldsFindAll{}
		for _, formformTemplateField := range nodeForm.FormTemplateFields {

			colIndex := formformTemplateField.ColNum

			for len(fieldsResponse) <= int(colIndex) {
				fieldsResponse = append(fieldsResponse, []responses.FormTemplateFieldsFindAll{})
			}

			fieldResponse := responses.FormTemplateFieldsFindAll{}

			//Mapping AdvancedOptions
			var advancedOptions map[string]interface{}
			if err := json.Unmarshal([]byte(*formformTemplateField.AdvancedOptions), &advancedOptions); err != nil {
				return response, err
			}
			fieldResponse.AdvancedOptions = advancedOptions
			if err := utils.Mapper(formformTemplateField, &fieldResponse); err != nil {
				return response, err
			}

			fieldsResponse[colIndex] = append(fieldsResponse[colIndex], fieldResponse)
		}

		formDatas := []responses.NodeFormDataResponse{}
		for _, formData := range nodeForm.FormFieldData {
			formDatas = append(formDatas, responses.NodeFormDataResponse{
				FieldId: fieldMap[formData.FormTemplateFieldID],
				Value:   formData.Value,
			})
		}

		nodeFormRes := responses.NodeFormDetailResponse{
			Template:    formTemplate,
			Fields:      fieldsResponse,
			Data:        formDatas,
			IsSubmitted: nodeForm.NodeForms.IsSubmitted,
			IsApproved:  nodeForm.NodeForms.IsApproved,
			IsRejected:  nodeForm.NodeForms.IsRejected,
		}
		if nodeForm.NodeForms.DataID != nil {
			nodeFormRes.DataId = *nodeForm.NodeForms.DataID
		}

		response = append(response, nodeFormRes)
	}

	return response, nil
}

func (s *NodeService) GetNodeJiraForm(ctx context.Context, nodeId string) (responses.JiraFormDetailResponse, error) {
	nodeJiraForm, err := s.NodeRepo.FindJiraFormByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return responses.JiraFormDetailResponse{}, err
	}

	formTemplateSystem, err := s.FormService.FindOneFormTemplateDetailByFormTemplateId(ctx, constants.FormTemplateIDJiraSystemForm)
	if err != nil {
		return responses.JiraFormDetailResponse{}, err
	}

	fieldMap := map[int32]string{}
	for _, formTemplateField := range nodeJiraForm.FormTemplateFields {
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
		Template: formTemplateSystem.Template,
		Fields:   formTemplateSystem.Fields,
		Data:     formDatas,
	}

	return response, nil
}

func (s *NodeService) ReassignNode(ctx context.Context, nodeId string, userId int32, userIdReq int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	nodeResult, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	// Old Assignee
	oldAssigneeId := nodeResult.AssigneeID

	node := model.Nodes{}
	utils.Mapper(nodeResult, &node)

	node.AssigneeID = &userId

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return err
	}

	// Notify
	users, err := s.UserAPI.FindUsersByUserIds([]int32{userId, userIdReq})
	if err != nil {
		return err
	}

	userMap := make(map[int32]types.Assignee)
	for _, user := range users.Data {
		userMap[user.ID] = types.Assignee{
			Id:           user.ID,
			Name:         user.Name,
			Email:        user.Email,
			AvatarUrl:    user.AvatarUrl,
			IsSystemUser: user.IsSystemUser,
		}
	}

	notification := types.Notification{
		ToUserIds: []string{strconv.Itoa(int(userId))},
		Subject:   fmt.Sprintf("New task assigned: %s", node.Title),
		Body:      fmt.Sprintf("You have been assigned a new task by %s.", userMap[userIdReq].Name),
	}
	s.NotificationService.SendNotification(ctx, notification)

	// History
	err = s.HistoryService.HistoryChangeNodeAssignee(ctx, tx, userId, node.RequestID, nodeId, oldAssigneeId, userIdReq)
	if err != nil {
		slog.Error("Error when creating history change node assignee", "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Check if node has Jira Key, then send update to Jira
	if node.JiraKey != nil {
		err := s.NatsService.SyncJiraWhenReassignNode(node)
		if err != nil {
			slog.Error("Error when syncing Jira when reassign node", "error", err)
			return err
		}
	}

	return nil
}

func (s *NodeService) ApproveNode(ctx context.Context, userId int32, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	nodeResult, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	node := model.Nodes{}
	utils.Mapper(nodeResult, &node)

	node.IsApproved = true

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return fmt.Errorf("update node status to completed fail: %w", err)
	}

	//
	if err := s.CompleteNodeLogic(ctx, tx, nodeId, userId); err != nil {
		return err
	}

	// History
	if err := s.HistoryService.HistoryApproveNode(ctx, tx, userId, node.RequestID, nodeId); err != nil {
		return err
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeService) RejectNode(ctx context.Context, userId int32, nodeId string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	nodeResult, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	node := model.Nodes{}
	utils.Mapper(nodeResult, &node)

	node.IsRejected = true

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return fmt.Errorf("update node status to completed fail: %w", err)
	}

	//
	if err := s.CompleteNodeLogic(ctx, tx, nodeId, userId); err != nil {
		return err
	}

	// History
	if err := s.HistoryService.HistoryRejectNode(ctx, tx, userId, node.RequestID, nodeId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) SubmitNodeForm(ctx context.Context, userId int32, nodeId string, formDataId string, req *[]requests.SubmitNodeFormRequest) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get Node Form
	formData, err := s.FormRepo.FindFormDataById(ctx, s.DB, formDataId)
	if err != nil {
		return fmt.Errorf("find node form by node id and form id fail: %w", err)
	}

	fieldMap := map[string]int32{}
	for _, formField := range formData.FormTemplateFields {
		fieldMap[formField.FieldID] = formField.ID
	}

	formFieldData := []model.FormFieldData{}
	for _, formField := range *req {
		formFieldData = append(formFieldData, model.FormFieldData{
			FormTemplateFieldID: fieldMap[formField.FieldId],
			Value:               formField.Value,
			FormDataID:          formDataId,
		})
	}
	if err := s.FormRepo.CreateFormFieldDatas(ctx, tx, formFieldData); err != nil {
		return fmt.Errorf("create form field data fail: %w", err)
	}

	nodeForm, err := s.NodeRepo.FindOneNodeFormByNodeIdAndFormId(ctx, s.DB, nodeId, formDataId)
	if err != nil {
		return fmt.Errorf("find node form by node id and form id fail: %w", err)
	}

	// Update Node Form Is Submitted
	nodeForm.SubmittedByUserID = &userId
	now := time.Now().UTC().Add(7 * time.Hour)
	nodeForm.SubmittedAt = &now
	nodeForm.LastUpdateUserID = &userId
	nodeForm.IsSubmitted = true
	if err := s.NodeRepo.UpdateNodeForm(ctx, tx, nodeForm); err != nil {
		return fmt.Errorf("update node form is submitted fail: %w", err)
	}

	node, err := s.NodeRepo.FindOneNodeByNodeIdTx(ctx, tx, nodeId)
	if err != nil {
		return fmt.Errorf("find node by node id fail: %w", err)
	}

	isCompletedNode := true
	for _, nodeForm := range node.NodeForms {
		if !nodeForm.IsSubmitted && (nodeForm.Permission == string(constants.NodeFormPermissionInput) || nodeForm.Permission == string(constants.NodeFormPermissionEdit)) {
			isCompletedNode = false
			break
		}
	}

	// Update Node To Completed
	if isCompletedNode {
		if err := s.CompleteNodeLogic(ctx, tx, nodeId, userId); err != nil {
			return fmt.Errorf("complete node handler fail: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) EditNodeForm(ctx context.Context, userId int32, nodeId string, formDataId string, req *[]requests.SubmitNodeFormRequest) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove All Form Field Data
	if err := s.FormRepo.RemoveAllFormFieldDataByFormDataId(ctx, tx, formDataId); err != nil {
		return fmt.Errorf("remove all form field data by form data id fail: %w", err)
	}

	// Find Form Data
	formData, err := s.FormRepo.FindFormDataById(ctx, s.DB, formDataId)
	if err != nil {
		return fmt.Errorf("find node form by node id and form id fail: %w", err)
	}

	fieldMap := map[string]int32{}
	for _, formField := range formData.FormTemplateFields {
		fieldMap[formField.FieldID] = formField.ID
	}

	formFieldData := []model.FormFieldData{}
	for _, formField := range *req {
		formFieldData = append(formFieldData, model.FormFieldData{
			FormTemplateFieldID: fieldMap[formField.FieldId],
			Value:               formField.Value,
			FormDataID:          formDataId,
		})
	}
	if err := s.FormRepo.CreateFormFieldDatas(ctx, tx, formFieldData); err != nil {
		return fmt.Errorf("create form field data fail: %w", err)
	}

	nodeForm, err := s.NodeRepo.FindOneNodeFormByNodeIdAndFormId(ctx, s.DB, nodeId, formDataId)
	if err != nil {
		return fmt.Errorf("find node form by node id and form id fail: %w", err)
	}

	// Update Node Form Is Submitted
	nodeForm.SubmittedByUserID = &userId
	now := time.Now().UTC().Add(7 * time.Hour)
	nodeForm.SubmittedAt = &now
	nodeForm.LastUpdateUserID = &userId
	nodeForm.IsSubmitted = true
	if err := s.NodeRepo.UpdateNodeForm(ctx, tx, nodeForm); err != nil {
		return fmt.Errorf("update node form is submitted fail: %w", err)
	}

	node, err := s.NodeRepo.FindOneNodeByNodeIdTx(ctx, tx, nodeId)
	if err != nil {
		return fmt.Errorf("find node by node id fail: %w", err)
	}

	isCompletedNode := true
	for _, nodeForm := range node.NodeForms {
		if !nodeForm.IsSubmitted && (nodeForm.Permission == string(constants.NodeFormPermissionInput) || nodeForm.Permission == string(constants.NodeFormPermissionEdit)) {
			isCompletedNode = false
			break
		}
	}

	// Update Node To Completed
	if isCompletedNode {
		if err := s.CompleteNodeLogic(ctx, tx, nodeId, userId); err != nil {
			return fmt.Errorf("complete node handler fail: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) ApproveNodeForm(ctx context.Context, nodeId string, formId string, userId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update Node Form Is Approved
	nodeForm, err := s.NodeRepo.FindOneNodeFormByNodeIdAndFormId(ctx, s.DB, nodeId, formId)
	if err != nil {
		return err
	}
	nodeForm.IsApproved = true
	if err := s.NodeRepo.UpdateNodeForm(ctx, tx, nodeForm); err != nil {
		return err
	}

	approveOrRejectUser := model.NodeFormApproveOrRejectUsers{
		UserID:     userId,
		NodeFormID: nodeForm.ID,
		IsApproved: true,
	}
	if err := s.NodeRepo.CreateApproveOrRejectUser(ctx, tx, approveOrRejectUser); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeService) RejectNodeForm(ctx context.Context, nodeId string, formId string, userId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update Node Form Is Rejected
	nodeForm, err := s.NodeRepo.FindOneNodeFormByNodeIdAndFormId(ctx, s.DB, nodeId, formId)
	if err != nil {
		return err
	}
	nodeForm.IsRejected = true
	if err := s.NodeRepo.UpdateNodeForm(ctx, tx, nodeForm); err != nil {
		return err
	}

	approveOrRejectUser := model.NodeFormApproveOrRejectUsers{
		UserID:     userId,
		NodeFormID: nodeForm.ID,
		IsApproved: false,
	}
	if err := s.NodeRepo.CreateApproveOrRejectUser(ctx, tx, approveOrRejectUser); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *NodeService) GetNodeTaskDetail(ctx context.Context, nodeId string) (responses.TaskDetail, error) {
	node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return responses.TaskDetail{}, err
	}

	request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, node.RequestID)
	if err != nil {
		return responses.TaskDetail{}, err
	}

	nodeTaskRelateds := []model.Nodes{}
	if err := utils.Mapper(request.Nodes, &nodeTaskRelateds); err != nil {
		return responses.TaskDetail{}, err
	}

	userIds := []int32{}
	existingUserIds := map[int32]bool{}

	if !existingUserIds[request.UserID] {
		userIds = append(userIds, request.UserID)
		existingUserIds[request.UserID] = true
	}

	if node.AssigneeID != nil && !existingUserIds[*node.AssigneeID] {
		userIds = append(userIds, *node.AssigneeID)
		existingUserIds[*node.AssigneeID] = true
	}

	var parentNode *model.Nodes
	if node.ParentID != nil {
		nodeResult, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, *node.ParentID)
		if err != nil {
			return responses.TaskDetail{}, err
		}
		parentNode = &nodeResult.Nodes

		if parentNode.AssigneeID != nil && !existingUserIds[*parentNode.AssigneeID] {
			userIds = append(userIds, *parentNode.AssigneeID)
			existingUserIds[*parentNode.AssigneeID] = true
		}
	}

	for _, related := range nodeTaskRelateds {
		if related.AssigneeID != nil && !existingUserIds[*related.AssigneeID] {
			userIds = append(userIds, *related.AssigneeID)
			existingUserIds[*related.AssigneeID] = true
		}
	}

	userApiMap := map[int32]results.UserApiDataResult{}
	if len(userIds) > 0 {
		assigneeResult, err := s.UserAPI.FindUsersByUserIds(userIds)
		if err != nil {
			return responses.TaskDetail{}, err
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

	requestTaskRes := responses.RequestTaskResponse{
		Id:               nodeId,
		Title:            node.Title,
		Type:             node.Type,
		RequestID:        node.RequestID,
		RequestTitle:     request.Title,
		RequestProgress:  request.Progress,
		Assignee:         mapUser(node.AssigneeID),
		PlannedStartTime: node.PlannedStartTime,
		PlannedEndTime:   node.PlannedEndTime,
		ActualStartTime:  node.ActualStartTime,
		ActualEndTime:    node.ActualEndTime,
		EstimatePoint:    node.EstimatePoint,
		Status:           node.Status,
		IsCurrent:        node.IsCurrent,
		IsApproved:       node.IsApproved,
		IsRejected:       node.IsRejected,
		ProjectKey:       request.Workflow.ProjectKey,
		JiraLinkUrl:      node.JiraLinkURL,
		Description:      node.Description,
	}

	if node.AttachFile != nil {
		var attachFiles map[string]interface{}
		if err := json.Unmarshal([]byte(*node.AttachFile), &attachFiles); err != nil {
			return responses.TaskDetail{}, err
		}

		requestTaskRes.AttachFiles = &attachFiles
	}

	taskDetail := responses.TaskDetail{
		RequestTaskResponse: requestTaskRes,
		RequestRequestBy:    mapUser(&request.UserID),
		IsApproval:          node.IsApproved,
		UpdatedAt:           node.UpdatedAt,
		JiraLinkURL:         node.JiraLinkURL,
		ProjectKey:          request.Workflow.ProjectKey,
		SprintId:            request.SprintID,
	}

	if parentNode != nil {
		taskDetail.Parent = &responses.TaskRelated{
			Title:    parentNode.Title,
			Type:     parentNode.Type,
			Status:   parentNode.Status,
			Assignee: mapUser(parentNode.AssigneeID),
		}
		if parentNode.JiraKey != nil {
			taskDetail.Parent.Key = *parentNode.JiraKey
		} else {
			taskDetail.Parent.Key = strconv.Itoa(int(parentNode.Key))
		}
	}

	nodeTaskRelatedsResponse := []responses.TaskRelated{}
	for _, nodeTaskRelated := range nodeTaskRelateds {

		if nodeTaskRelated.Type == string(constants.NodeTypeStart) || nodeTaskRelated.Type == string(constants.NodeTypeEnd) {
			continue
		}
		if nodeTaskRelated.ID == node.ID {
			continue
		}

		assigneeResponse := mapUser(nodeTaskRelated.AssigneeID)

		nodeTaskRelatedResponse := responses.TaskRelated{
			Title:    nodeTaskRelated.Title,
			Type:     nodeTaskRelated.Type,
			Status:   nodeTaskRelated.Status,
			Assignee: assigneeResponse,
		}

		if nodeTaskRelated.JiraKey != nil {
			nodeTaskRelatedResponse.Key = *nodeTaskRelated.JiraKey
		} else {
			nodeTaskRelatedResponse.Key = strconv.Itoa(int(nodeTaskRelated.Key))
		}
		nodeTaskRelatedsResponse = append(nodeTaskRelatedsResponse, nodeTaskRelatedResponse)
	}
	taskDetail.Related = nodeTaskRelatedsResponse

	if node.JiraKey != nil {
		taskDetail.Key = *node.JiraKey
	} else {
		taskDetail.Key = strconv.Itoa(int(node.Key))
	}

	return taskDetail, nil
}

func (s *NodeService) GetNodeStoryByAssignee(ctx context.Context, userId int32) ([]responses.WorkflowResponse, error) {

	workflows := []responses.WorkflowResponse{}

	stories, err := s.NodeRepo.FindAllNodeStoryByAssigneeId(ctx, s.DB, userId)
	if err != nil {
		return nil, err
	}

	for _, story := range stories {
		workflowResponse := responses.WorkflowResponse{}
		if err := utils.Mapper(story.Workflows, &workflowResponse); err != nil {
			return nil, err
		}
		if story.SubRequestID != nil {
			workflowResponse.RequestId = *story.SubRequestID
		}

		categoryResponse := responses.CategoryResponse{}
		if err := utils.Mapper(story.Category, &categoryResponse); err != nil {
			return nil, err
		}
		workflowResponse.Category = categoryResponse
		workflows = append(workflows, workflowResponse)
	}

	return workflows, nil
}

// JIRA ==========================
func (s *NodeService) SyncJiraWhenCompleteNode(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	// Update node status to completed
	node.Status = string(constants.NodeStatusCompleted)

	// Get request info for Jira sync
	request, err := s.RequestRepo.FindOneRequestByRequestIdTx(ctx, tx, node.RequestID)
	if err != nil {
		return fmt.Errorf("get request info fail: %w", err)
	}

	// Sync with Jira
	if err := s.NatsService.SyncNodeStatusToJira(ctx, tx, node, request.Requests, request.Workflow); err != nil {
		slog.Error("Failed to sync with Jira", "error", err)
		// Continue execution even if Jira sync fails
	}

	return nil
}

func (s *NodeService) SyncJiraWhenStartNode(ctx context.Context, tx *sql.Tx, node model.Nodes) error {
	// Update node status to in progress
	node.Status = string(constants.NodeStatusInProgress)

	// Get request info for Jira sync
	request, err := s.RequestRepo.FindOneRequestByRequestIdTx(ctx, tx, node.RequestID)
	if err != nil {
		return fmt.Errorf("get request info fail: %w", err)
	}

	// Sync with Jira
	if err := s.NatsService.SyncNodeStatusToJira(ctx, tx, node, request.Requests, request.Workflow); err != nil {
		slog.Error("Failed to sync with Jira", "error", err)
		// Continue execution even if Jira sync fails
	}

	return nil
}

func (s *NodeService) GetNodeTaskCount(ctx context.Context, userId int32) (responses.NodeTaskCountResponse, error) {
	response := responses.NodeTaskCountResponse{}

	count, err := s.RequestRepo.CountActiveRequests(ctx, s.DB, userId)
	if err != nil {
		return responses.NodeTaskCountResponse{}, err
	}
	response.ActiveRequests = count

	approvalCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, "", queryparams.RequestTaskCount{
		Type:         string(constants.NodeTypeApproval),
		WorkflowType: string(constants.WorkflowTypeGeneral),
	})
	if err != nil {
		return responses.NodeTaskCountResponse{}, err
	}
	response.ApprovalTasks = int32(approvalCount)

	inputTaskCount, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, "", queryparams.RequestTaskCount{
		Type:         string(constants.NodeTypeInput),
		WorkflowType: string(constants.WorkflowTypeGeneral),
	})
	if err != nil {
		return responses.NodeTaskCountResponse{}, err
	}
	response.InputTasks = int32(inputTaskCount)

	projectTasks, err := s.RequestRepo.CountRequestTaskByStatusAndUserIdAndQueryParams(ctx, s.DB, userId, "", queryparams.RequestTaskCount{
		WorkflowType: string(constants.WorkflowTypeProject),
	})
	if err != nil {
		return responses.NodeTaskCountResponse{}, err
	}
	response.ProjectTasks = int32(projectTasks)

	return response, nil
}

func (s *NodeService) CreateComment(ctx context.Context, req *requests.CreateComment, nodeId string, userId int32) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	commentModel := model.Comments{}
	if err := utils.Mapper(req, &commentModel); err != nil {
		return fmt.Errorf("mapping comment failed: %w", err)
	}

	commentModel.UserID = userId
	commentModel.NodeID = nodeId

	if err := s.CommentRepository.CreateComment(ctx, tx, commentModel); err != nil {
		return fmt.Errorf("create comment failed: %w", err)
	}

	node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	//Notify
	if node.AssigneeID != &userId {
		if err := s.NotificationService.NotifyComment(ctx, node.Title, userId); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	return nil
}

func (s *NodeService) GetAllComments(ctx context.Context, nodeId string) ([]responses.CommentResponse, error) {
	commentResponse := []responses.CommentResponse{}

	comments, err := s.CommentRepository.FindAllCommentByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return commentResponse, err
	}

	userIds := []int32{}
	existedUserIds := make(map[int32]bool)
	for _, comment := range comments {
		if existedUserIds[comment.UserID] {
			continue
		}
		userIds = append(userIds, comment.UserID)
		existedUserIds[comment.UserID] = true
	}

	results, err := s.UserAPI.FindUsersByUserIds(userIds)
	if err != nil {
		return commentResponse, err
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

	for _, comment := range comments {
		commentRes := responses.CommentResponse{
			CreatedAt: comment.CreatedAt,
			Assignee:  userMap[comment.UserID],
			Content:   comment.Content,
		}

		commentResponse = append(commentResponse, commentRes)
	}

	return commentResponse, err

}
