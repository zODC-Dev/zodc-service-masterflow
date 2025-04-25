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

func (s *NodeService) CheckIfAllNodeFormIsApprovedOrRejected(ctx context.Context, nodeId string) (bool, error) {
	nodeForms, err := s.NodeRepo.FindAllNodeFormByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return false, err
	}

	for _, nodeForm := range nodeForms {
		if !nodeForm.IsApproved && !nodeForm.IsRejected {
			return false, nil
		}
	}

	return true, nil
}

func (s *NodeService) LogicForConditionNode(ctx context.Context, tx *sql.Tx, nodeId string, isTrue bool, userId int32) error {
	// Update Node Form
	node, err := s.NodeRepo.FindOneNodeByNodeIdTx(ctx, tx, nodeId)
	if err != nil {
		return err
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

	// Next node will alway be Condition Node
	connections, err := s.ConnectionRepo.FindConnectionsWithToNodesByFromNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return err
	}

	for _, mainConnection := range connections {
		// Update Condition Node
		mainConnection.IsCompleted = true
		connectionModel := model.Connections{}
		if err := utils.Mapper(mainConnection, &connectionModel); err != nil {
			return fmt.Errorf("mapper connection fail: %w", err)
		}
		if err := s.ConnectionRepo.UpdateConnection(ctx, tx, connectionModel); err != nil {
			return fmt.Errorf("update connection fail: %w", err)
		}

		if mainConnection.Node.Type == string(constants.NodeTypeCondition) {

			// Update Condition Node
			conditionNode := model.Nodes{}
			utils.Mapper(mainConnection.Node, &conditionNode)
			conditionNode.Status = string(constants.NodeStatusCompleted)
			conditionNode.IsCurrent = true
			now := time.Now().UTC().Add(7 * time.Hour)
			conditionNode.ActualStartTime = &now
			conditionNode.ActualEndTime = &now
			if err := s.NodeRepo.UpdateNode(ctx, tx, conditionNode); err != nil {
				return fmt.Errorf("update condition node fail: %w", err)
			}

			// condition destination
			nodeConditionDestinations, err := s.NodeRepo.FindAllNodeConditionDestinationByNodeId(ctx, s.DB, mainConnection.Node.ID, isTrue)
			if err != nil {
				return fmt.Errorf("find all node condition destination by node id fail: %w", err)
			}

			// Update Connection Condition Destination Node To Completed
			for _, nodeConditionDestination := range nodeConditionDestinations {

				node, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeConditionDestination.DestinationNodeID)
				if err != nil {
					return fmt.Errorf("find one node by node id fail: %w", err)
				}

				// Update Connection Condition Destination Node To Completed
				connectionConditionDestination, err := s.ConnectionRepo.FindConnectionsByToNodeIdTx(ctx, tx, nodeConditionDestination.DestinationNodeID)
				if err != nil {
					return fmt.Errorf("find connections by to node id fail: %w", err)
				}
				for _, connectionDestination := range connectionConditionDestination {
					if mainConnection.Node.ID == connectionDestination.FromNodeID {
						connectionDestination.IsCompleted = true
						connectionModel := model.Connections{}
						if err := utils.Mapper(connectionDestination, &connectionModel); err != nil {
							return fmt.Errorf("mapper connection fail: %w", err)
						}
						if err := s.ConnectionRepo.UpdateConnection(ctx, tx, connectionModel); err != nil {
							return fmt.Errorf("update connection fail: %w", err)
						}
					}
				}

				if node.Type == string(constants.NodeTypeEnd) {
					now := time.Now().UTC().Add(7 * time.Hour)

					nodeModel := model.Nodes{}
					utils.Mapper(node, &nodeModel)
					nodeModel.IsCurrent = true
					nodeModel.Status = string(constants.NodeStatusCompleted)

					nodeModel.ActualStartTime = &now
					nodeModel.ActualEndTime = &now
					if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
						return err
					}

					// Update Request
					request, err := s.RequestRepo.FindOneRequestByRequestIdTx(ctx, tx, node.RequestID)
					if err != nil {
						return err
					}

					requestModel := model.Requests{}
					utils.Mapper(request, &requestModel)

					// Notify
					userIds := []string{}
					existingUserIds := map[string]bool{}
					for _, node := range request.Nodes {
						if node.AssigneeID != nil {
							if !existingUserIds[strconv.Itoa(int(*node.AssigneeID))] {
								userIds = append(userIds, strconv.Itoa(int(*node.AssigneeID)))
								existingUserIds[strconv.Itoa(int(*node.AssigneeID))] = true
							}
						}
					}

					if *node.EndType == string(constants.NodeEndTypeTerminate) {
						requestModel.Status = string(constants.RequestStatusTerminated)
						requestModel.TerminatedAt = &now

						// Notify
						s.NotificationService.NotifyRequestTerminated(ctx, request.Title, userIds)
					} else {
						requestModel.Status = string(constants.RequestStatusCompleted)
						requestModel.CompletedAt = &now

						s.NotificationService.NotifyRequestCompleted(ctx, request.Title, userIds)
					}

					requestModel.CompletedAt = &now
					requestModel.Progress = 100

					if err := s.RequestRepo.UpdateRequest(ctx, tx, requestModel); err != nil {
						return err
					}

					return nil
				} else if node.Type == string(constants.NodeTypeNotification) {
					// Send Notification
					var cc []string
					if node.CcEmails != nil {
						err := json.Unmarshal([]byte(*node.CcEmails), &cc)
						if err != nil {
							return err
						}
					}
					var to []string
					if node.ToEmails != nil {
						err := json.Unmarshal([]byte(*node.ToEmails), &to)
						if err != nil {
							return err
						}
					}
					var bcc []string
					if node.BccEmails != nil {
						err := json.Unmarshal([]byte(*node.BccEmails), &bcc)
						if err != nil {
							return err
						}
					}
					notification := types.Notification{
						ToEmails:    to,
						ToCcEmails:  cc,
						ToBccEmails: bcc,
						Subject:     *node.Subject,
					}
					if node.Body != nil {
						notification.Body = *node.Body
					}
					if node.Subject != nil {
						notification.Subject = *node.Subject
					}

					// Is Send Form
					if node.IsSendApprovedForm || node.IsSendRejectedForm {
						request, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, node.RequestID)
						if err != nil {
							return err
						}

						notification.Body += "\n\n\n"
						for _, node := range request.Nodes {
							for _, nodeForm := range node.NodeForms {
								if nodeForm.IsApproved && node.IsSendApprovedForm {
									notification.Body += "\n" + configs.Env.FE_HOST + "/form-management/review/" + *nodeForm.DataID
								}
								if nodeForm.IsRejected && node.IsSendRejectedForm {
									notification.Body += "\n" + configs.Env.FE_HOST + "/form-management/review/" + *nodeForm.DataID
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
					node.Status = string(constants.NodeStatusCompleted)
					node.IsCurrent = true
					now := time.Now().UTC().Add(7 * time.Hour)
					node.ActualStartTime = &now
					node.ActualEndTime = &now

					nodeModel := model.Nodes{}
					utils.Mapper(node, &nodeModel)
					if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
						return err
					}

					connectionToNode, err := s.ConnectionRepo.FindConnectionsWithToNodesByFromNodeId(ctx, s.DB, node.ID)
					if err != nil {
						return err
					}
					for _, connection := range connectionToNode {
						connection.IsCompleted = true
						connectionModel := model.Connections{}
						if err := utils.Mapper(connection, &connectionModel); err != nil {
							return err
						}
						if err := s.ConnectionRepo.UpdateConnection(ctx, tx, connectionModel); err != nil {
							return err
						}

						if connection.Node.Type == string(constants.NodeTypeEnd) {
							connection.Node.IsCurrent = true
							now := time.Now().UTC().Add(7 * time.Hour)
							connection.Node.ActualStartTime = &now
							connection.Node.ActualEndTime = &now
							connection.Node.Status = string(constants.NodeStatusCompleted)

							if err := s.NodeRepo.UpdateNode(ctx, tx, connection.Node); err != nil {
								return err
							}

							// Update Request
							request, err := s.RequestRepo.FindOneRequestByRequestIdTx(ctx, tx, node.RequestID)
							if err != nil {
								return err
							}

							requestModel := model.Requests{}
							utils.Mapper(request, &requestModel)

							if *node.EndType == string(constants.NodeEndTypeTerminate) {
								requestModel.Status = string(constants.RequestStatusTerminated)
								requestModel.TerminatedAt = &now
							} else {
								requestModel.Status = string(constants.RequestStatusCompleted)
								requestModel.CompletedAt = &now
							}

							requestModel.Progress = 100

							if err := s.RequestRepo.UpdateRequest(ctx, tx, requestModel); err != nil {
								return err
							}

							// History
							err = s.HistoryService.HistoryEndRequest(ctx, tx, request.ID, node.ID)
							if err != nil {
								return err
							}

							return nil
						}
					}
				}
				nodeModel := model.Nodes{}
				utils.Mapper(node, &nodeModel)
				nodeModel.IsCurrent = true
				nodeModel.Status = string(constants.NodeStatusInProgress)
				now := time.Now().UTC().Add(7 * time.Hour)
				nodeModel.ActualStartTime = &now
				if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
					return err
				}

				// Notify
				users, err := s.UserAPI.FindUsersByUserIds([]int32{*node.AssigneeID})
				if err != nil {
					return err
				}

				err = s.NotificationService.NotifyTaskAvailable(ctx, node.Title, users.Data[0].ID)
				if err != nil {
					return err
				}

				// History
				err = s.HistoryService.HistoryNewTask(ctx, tx, node.RequestID, node.ID, *node.AssigneeID)
				if err != nil {
					return err
				}
			}
		}
	}

	// Calculate Request Process
	if err := s.RequestService.UpdateCalculateRequestProgress(ctx, tx, node.RequestID); err != nil {
		return fmt.Errorf("update calculate request progress fail: %w", err)
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

	//Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	// Notify
	if err := s.NotificationService.NotifyTaskStarted(ctx, node); err != nil {
		return err
	}

	// History
	oldStatus := string(constants.NodeStatusTodo)
	err = s.HistoryService.HistoryChangeNodeStatus(ctx, tx, userId, node.RequestID, nodeId, &oldStatus, string(constants.NodeStatusInProgress))
	if err != nil {
		return err
	}

	return nil
}

func (s *NodeService) CompleteNodeHandler(ctx context.Context, nodeId string, userId int32, isChangeToInProgress bool) error {

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	nodeResult, err := s.NodeRepo.FindOneNodeByNodeId(ctx, s.DB, nodeId)
	if err != nil {
		return fmt.Errorf("find node by node id fail: %w", err)
	}

	node := model.Nodes{}
	utils.Mapper(nodeResult, &node)

	// Update Current Node Status To Completed
	node.Status = string(constants.NodeStatusCompleted)

	// Set actual finish time
	now := time.Now().UTC().Add(7 * time.Hour)
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

		// History
		oldStatus := string(constants.NodeStatusInProgress)
		err = s.HistoryService.HistoryChangeNodeStatus(ctx, tx, userId, node.RequestID, nodeId, &oldStatus, string(constants.NodeStatusCompleted))
		if err != nil {
			return err
		}

		// If Prevous Nodes not finish yet // If More than one node not completed then next node dont need to update status
		isUpdateNodeStatus := true

		connections, err := s.ConnectionRepo.FindConnectionsByToNodeIdTx(ctx, tx, connectionsToNode[i].Node.ID)
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

			// If Node is Approval Node
			if connectionsToNode[i].Node.Type == string(constants.NodeTypeApproval) {
				if connectionsToNode[i].Node.AssigneeID != nil {
					s.NotificationService.NotifyNodeApproveNeeded(ctx, nodeResult.Request.Title, *connectionsToNode[i].Node.AssigneeID)
				}
			}

			// If Node is End Node
			if connectionsToNode[i].Node.Type == string(constants.NodeTypeEnd) {
				// Update end node to completed
				connectionsToNode[i].Node.Status = string(constants.NodeStatusCompleted)
				connectionsToNode[i].Node.ActualEndTime = &now
				connectionsToNode[i].Node.ActualStartTime = &now
				if err := s.NodeRepo.UpdateNode(ctx, tx, connectionsToNode[i].Node); err != nil {
					return err
				}

				// Mark request completed
				request, err := s.RequestRepo.FindRequestByNodeId(ctx, s.DB, connectionsToNode[i].Node.ID)
				if err != nil {
					return err
				}

				requestDetail, err := s.RequestRepo.FindOneRequestByRequestId(ctx, s.DB, request.ID)
				if err != nil {
					return err
				}

				userIds := []string{}
				existingUserIds := map[string]bool{}
				for _, node := range requestDetail.Nodes {
					if node.AssigneeID != nil {
						if !existingUserIds[strconv.Itoa(int(*node.AssigneeID))] {
							userIds = append(userIds, strconv.Itoa(int(*node.AssigneeID)))
							existingUserIds[strconv.Itoa(int(*node.AssigneeID))] = true
						}
					}
				}

				if node.EndType != nil && *node.EndType == string(constants.NodeEndTypeTerminate) {
					request.Status = string(constants.RequestStatusTerminated)
					request.TerminatedAt = &now

					// Notify
					err = s.NotificationService.NotifyRequestTerminated(ctx, request.Title, userIds)
					if err != nil {
						return err
					}

					// History
					err = s.HistoryService.HistoryEndRequest(ctx, tx, request.ID, node.ID)
					if err != nil {
						return err
					}

				} else {
					request.Status = string(constants.RequestStatusCompleted)
					request.CompletedAt = &now

					// Notify
					err = s.NotificationService.NotifyRequestCompleted(ctx, request.Title, userIds)
					if err != nil {
						return err
					}

					// History
					err = s.HistoryService.HistoryEndRequest(ctx, tx, request.ID, node.ID)
					if err != nil {
						return err
					}
				}

				request.Progress = 100

				if err := s.RequestRepo.UpdateRequest(ctx, tx, request); err != nil {
					return err
				}

				//
				nodeSubRequest, err := s.NodeRepo.FindOneNodeBySubRequestID(ctx, tx, request.ID)
				if err != nil {
					errStr := err.Error()
					if errStr != "qrm: no rows in result set" {
						return err
					}
				} else {
					err = s.CompleteNodeHandler(ctx, nodeSubRequest.ID, userId, isChangeToInProgress)
					if err != nil {
						return err
					}
				}

			} else {
				connectionsToNode[i].Node.IsCurrent = true
				if isChangeToInProgress {
					connectionsToNode[i].Node.Status = string(constants.NodeStatusInProgress)

					now := time.Now().UTC().Add(7 * time.Hour)
					connectionsToNode[i].Node.ActualStartTime = &now
				}
				err := s.NodeRepo.UpdateNode(ctx, tx, connectionsToNode[i].Node)
				if err != nil {
					return err
				}

				users, err := s.UserAPI.FindUsersByUserIds([]int32{*connectionsToNode[i].Node.AssigneeID})
				if err != nil {
					return err
				}

				// Notify
				err = s.NotificationService.NotifyTaskAvailable(ctx, connectionsToNode[i].Node.Title, users.Data[0].ID)
				if err != nil {
					return err
				}

				// History
				err = s.HistoryService.HistoryNewTask(ctx, tx, connectionsToNode[i].Node.RequestID, connectionsToNode[i].Node.ID, *connectionsToNode[i].Node.AssigneeID)
				if err != nil {
					return err
				}

			}
		}
	}

	// send notification
	s.NotificationService.NotifyTaskCompleted(ctx, node)

	// Calculate Request Process
	if err := s.RequestService.UpdateCalculateRequestProgress(ctx, tx, node.RequestID); err != nil {
		return fmt.Errorf("update calculate request progress fail: %w", err)
	}

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

	if err := tx.Commit(); err != nil {
		return err
	}

	// Check if node has Jira Key, then send update to Jira
	if node.JiraKey != nil {
		err := s.NatsService.SyncJiraWhenReassignNode(ctx, tx, node)
		if err != nil {
			return err
		}
	}

	// History
	err = s.HistoryService.HistoryChangeNodeAssignee(ctx, tx, userId, node.RequestID, nodeId, oldAssigneeId, userIdReq)
	if err != nil {
		return err
	}

	return nil
}

// Only For Condition Node
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
	node.Status = string(constants.NodeStatusCompleted)

	now := time.Now().UTC().Add(7 * time.Hour)
	node.ActualEndTime = &now

	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return fmt.Errorf("update node status to completed fail: %w", err)
	}

	//
	err = s.LogicForConditionNode(ctx, tx, nodeId, true, userId)
	if err != nil {
		return fmt.Errorf("logic for condition node fail: %w", err)
	}

	// Commit
	if err := tx.Commit(); err != nil {
		return err
	}

	// History
	err = s.HistoryService.HistoryApproveNode(ctx, tx, userId, node.RequestID, nodeId)
	if err != nil {
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
	node.Status = string(constants.NodeStatusCompleted)

	now := time.Now().UTC().Add(7 * time.Hour)
	node.ActualEndTime = &now
	if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
		return fmt.Errorf("update node status to completed fail: %w", err)
	}

	//
	err = s.LogicForConditionNode(ctx, tx, nodeId, false, userId)
	if err != nil {
		return fmt.Errorf("logic for condition node fail: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit fail: %w", err)
	}

	err = s.HistoryService.HistoryRejectNode(ctx, tx, userId, node.RequestID, nodeId)
	if err != nil {
		return err
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
		if err := s.CompleteNodeHandler(ctx, nodeId, userId, true); err != nil {
			return fmt.Errorf("complete node handler fail: %w", err)
		}

		nodeModel := model.Nodes{}
		utils.Mapper(node, &nodeModel)
		nodeModel.Status = string(constants.NodeStatusCompleted)

		actualEndTime := time.Now().UTC().Add(7 * time.Hour)
		nodeModel.ActualEndTime = &actualEndTime

		if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
			return fmt.Errorf("update node status to completed fail: %w", err)
		}

	}

	// Update Calculate Request Progress
	if err := s.RequestService.UpdateCalculateRequestProgress(ctx, tx, node.RequestID); err != nil {
		return fmt.Errorf("update calculate request progress fail: %w", err)
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
		if err := s.CompleteNodeHandler(ctx, nodeId, userId, true); err != nil {
			return fmt.Errorf("complete node handler fail: %w", err)
		}

		nodeModel := model.Nodes{}
		utils.Mapper(node, &nodeModel)
		nodeModel.Status = string(constants.NodeStatusCompleted)

		actualEndTime := time.Now().UTC().Add(7 * time.Hour)
		nodeModel.ActualEndTime = &actualEndTime

		if err := s.NodeRepo.UpdateNode(ctx, tx, nodeModel); err != nil {
			return fmt.Errorf("update node status to completed fail: %w", err)
		}

	}

	// Update Calculate Request Progress
	if err := s.RequestService.UpdateCalculateRequestProgress(ctx, tx, node.RequestID); err != nil {
		return fmt.Errorf("update calculate request progress fail: %w", err)
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

	taskDetail := responses.TaskDetail{
		RequestTaskResponse: responses.RequestTaskResponse{
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
		},
		RequestRequestBy: mapUser(&request.UserID),
		IsApproval:       node.IsApproved,
		UpdatedAt:        node.UpdatedAt,
		JiraLinkURL:      node.JiraLinkURL,
		ProjectKey:       request.Workflow.ProjectKey,
		SprintId:         request.SprintID,
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
