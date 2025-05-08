package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	natslib "github.com/nats-io/nats.go"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	natsModel "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types/nats"
	nats "github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

// NatsSubscriberService handles listening for messages from NATS
// Updated to use go-jet SQL builder instead of raw SQL queries
// Improved error handling and transaction management
type NatsSubscriberService struct {
	NatsClient  *nats.NATSClient
	DB          *sql.DB
	NodeRepo    *repositories.NodeRepository
	RequestRepo *repositories.RequestRepository
	NodeService *NodeService
	// Store subscriptions to properly unsubscribe later
	subscriptions []*natslib.Subscription
}

// NewNatsSubscriberService creates a new instance of NatsSubscriberService
func NewNatsSubscriberService(natsClient *nats.NATSClient, db *sql.DB, nodeRepo *repositories.NodeRepository, requestRepo *repositories.RequestRepository, nodeService *NodeService) *NatsSubscriberService {
	return &NatsSubscriberService{
		NatsClient:    natsClient,
		DB:            db,
		NodeRepo:      nodeRepo,
		RequestRepo:   requestRepo,
		NodeService:   nodeService,
		subscriptions: make([]*natslib.Subscription, 0),
	}
}

// Start initializes and starts the subscriber service
func (s *NatsSubscriberService) Start(ctx context.Context) error {
	slog.Info("Starting NATS Subscriber Service")

	// Subscribe to topics
	if err := s.subscribeToTopics(ctx); err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	slog.Info("NATS Subscriber Service started successfully")
	return nil
}

// subscribeToTopics registers subscriptions for each topic
func (s *NatsSubscriberService) subscribeToTopics(ctx context.Context) error {
	// List of topics to subscribe to
	topics := []string{
		"jira.issue.update",
	}

	for _, topic := range topics {
		slog.Info("Subscribing to topic", "topic", topic)

		// Use the Subscribe method from pkg/nats/subscribe.go
		subscription, err := s.NatsClient.Subscribe(topic, s.handleMessage)
		if err != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
		}

		// Store subscription for later cleanup
		s.subscriptions = append(s.subscriptions, subscription)
	}

	return nil
}

// handleMessage processes received messages
func (s *NatsSubscriberService) handleMessage(msg *natslib.Msg) {
	slog.Info("Received message", "subject", msg.Subject, "data_length", len(msg.Data))

	// Biến để lưu thông tin cần xử lý sau khi commit
	var postCommitActions []func()

	// Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		slog.Error("Failed to start transaction", "error", err)
		return
	}

	// Use a flag to track if we've committed the transaction
	var committed bool
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
				slog.Error("Failed to rollback transaction", "error", err)
			}
		}
	}()

	// Process message based on subject
	var processErr error
	switch msg.Subject {
	case "jira.issue.update":
		postCommitActions, processErr = s.handleJiraIssueUpdate(tx, msg.Data)
	default:
		slog.Warn("No handler for subject", "subject", msg.Subject)
	}

	if processErr != nil {
		slog.Error("Failed to process message", "subject", msg.Subject, "error", processErr)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", "error", err)
		return
	}
	committed = true

	// Thực hiện các action sau khi commit thành công
	for _, action := range postCommitActions {
		action()
	}

	slog.Info("Successfully processed message", "subject", msg.Subject)
}

// handleJiraIssueUpdate processes messages from jira.issue.update
func (s *NatsSubscriberService) handleJiraIssueUpdate(tx *sql.Tx, data []byte) ([]func(), error) {
	// Parse message data
	var message natsModel.JiraIssueUpdateMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Jira issue update message: %w", err)
	}

	// Validation
	if message.JiraKey == "" {
		return nil, fmt.Errorf("missing required field: jiraKey")
	}

	// Map Jira status to system status
	var systemStatus string
	switch message.Status {
	case "To Do":
		systemStatus = string(constants.NodeStatusTodo)
	case "In Progress":
		systemStatus = string(constants.NodeStatusInProgress)
	case "Done":
		systemStatus = string(constants.NodeStatusCompleted)
	default:
		// If status isn't recognized, use the original value
		systemStatus = message.Status
	}

	// Map old status if present
	var oldSystemStatus *string
	if message.OldStatus != nil {
		var oldStatusValue string
		switch *message.OldStatus {
		case "To Do":
			oldStatusValue = string(constants.NodeStatusTodo)
		case "In Progress":
			oldStatusValue = string(constants.NodeStatusInProgress)
		case "Done":
			oldStatusValue = string(constants.NodeStatusCompleted)
		default:
			oldStatusValue = *message.OldStatus
		}
		oldSystemStatus = &oldStatusValue
	}

	ctx := context.Background()

	// Step 1: Find and update all nodes with this Jira key
	nodesToUpdate, err := s.findNodesByJiraKey(ctx, tx, message.JiraKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find nodes: %w", err)
	}

	// Update each node with new data
	for _, node := range nodesToUpdate {
		// Update node fields
		node.Title = message.Summary
		node.Status = systemStatus

		// Add 7 hours to LastSyncedAt due to Vietnamese timezone
		if message.LastSyncedAt != nil {
			adjustedTime := message.LastSyncedAt.Add(7 * time.Hour)
			node.LastSyncedAt = &adjustedTime
		}

		// Update estimate point if provided
		if message.EstimatePoint != nil {
			node.EstimatePoint = message.EstimatePoint
		}

		// Update assignee ID if provided
		if message.AssigneeId != nil {
			// Convert string assignee ID to int32
			assigneeId, err := strconv.ParseInt(*message.AssigneeId, 10, 32)
			if err != nil {
				slog.Warn("Failed to parse assignee ID", "assigneeId", *message.AssigneeId, "error", err)
			} else {
				assigneeId32 := int32(assigneeId)
				node.AssigneeID = &assigneeId32
			}
		}

		// Update the node using the repository
		if err := s.NodeRepo.UpdateNode(ctx, tx, node); err != nil {
			return nil, fmt.Errorf("failed to update node %s: %w", node.ID, err)
		}
	}

	// Step 2: Update form data for all corresponding forms
	formDataIDs, err := s.findFormDataByJiraKey(ctx, tx, message.JiraKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find form data IDs: %w", err)
	}

	// For each form data ID, update relevant fields
	for _, formDataID := range formDataIDs {
		// Update summary field
		if err := s.updateFormFieldValue(ctx, tx, formDataID, "summary", message.Summary); err != nil {
			slog.Warn("Failed to update summary in form field data", "error", err, "formDataId", formDataID)
		}

		// Update assignee email field
		if message.AssigneeEmail != "" {
			if err := s.updateFormFieldValue(ctx, tx, formDataID, "assignee_email", message.AssigneeEmail); err != nil {
				slog.Warn("Failed to update assignee_email in form field data", "error", err, "formDataId", formDataID)
			}
		}

		// Update status field
		if err := s.updateFormFieldValue(ctx, tx, formDataID, "status", systemStatus); err != nil {
			slog.Warn("Failed to update status in form field data", "error", err, "formDataId", formDataID)
		}

		// Update estimate point if provided
		if message.EstimatePoint != nil {
			if err := s.updateFormFieldValue(ctx, tx, formDataID, "estimate_point", fmt.Sprintf("%f", *message.EstimatePoint)); err != nil {
				slog.Warn("Failed to update estimate_point in form field data", "error", err, "formDataId", formDataID)
			}
		}
	}

	// Step 3: Handle workflow state transitions
	var postCommitActions []func()

	if message.SprintId != nil && message.OldStatus != nil && oldSystemStatus != nil {
		// Find active node for potential workflow transition
		activeNodeIds, err := s.findActiveNodesByJiraKey(ctx, tx, message.JiraKey, *message.SprintId)
		if err != nil {
			slog.Warn("Failed to find active nodes", "error", err)
		} else if len(activeNodeIds) > 0 {
			// Handle state transitions for each active node
			for _, nodeId := range activeNodeIds {
				// Determine if we need to apply workflow state transition
				if *oldSystemStatus == string(constants.NodeStatusTodo) && systemStatus == string(constants.NodeStatusInProgress) {
					// From Todo to In Progress -> Start Node
					var assigneeId int32 = 0
					if message.AssigneeId != nil {
						if parsedId, err := strconv.ParseInt(*message.AssigneeId, 10, 32); err == nil {
							assigneeId = int32(parsedId)
						}
					}

					// Add to post-commit actions
					postCommitActions = append(postCommitActions, func() {
						ctx := context.Background()
						if err := s.NodeService.StartNodeHandler(ctx, assigneeId, nodeId); err != nil {
							slog.Error("Failed to start node", "nodeId", nodeId, "error", err)
						}
					})
				} else if (*oldSystemStatus == string(constants.NodeStatusInProgress) || *oldSystemStatus == string(constants.NodeStatusTodo)) && systemStatus == string(constants.NodeStatusCompleted) {
					// From In Progress to Completed -> Complete Node
					var assigneeId int32 = 1 // default system user ID
					if message.AssigneeId != nil {
						if parsedId, err := strconv.ParseInt(*message.AssigneeId, 10, 32); err == nil {
							assigneeId = int32(parsedId)
						}
					}

					// Add to post-commit actions
					postCommitActions = append(postCommitActions, func() {
						ctx := context.Background()
						// Create a new transaction for post-commit action
						newTx, err := s.DB.BeginTx(ctx, nil)
						if err != nil {
							slog.Error("Failed to begin transaction for completing node", "error", err)
							return
						}
						defer func() {
							if err := newTx.Rollback(); err != nil && err != sql.ErrTxDone {
								slog.Error("Failed to rollback transaction", "error", err)
							}
						}()

						if err := s.NodeService.CompleteNodeLogic(ctx, newTx, nodeId, assigneeId); err != nil {
							slog.Error("Failed to complete node", "nodeId", nodeId, "error", err)
							return
						}

						if err := newTx.Commit(); err != nil {
							slog.Error("Failed to commit transaction for completing node", "error", err)
						}
					})
				}
			}
		}
	}

	return postCommitActions, nil
}

// updateFormFieldValue updates a specific field in a form data record
func (s *NatsSubscriberService) updateFormFieldValue(ctx context.Context, tx *sql.Tx, formDataID, fieldID, value string) error {
	FormFieldData := table.FormFieldData
	FormTemplateFields := table.FormTemplateFields

	// Find the field ID for the given field
	findFieldQuery := postgres.SELECT(
		FormFieldData.ID,
	).FROM(
		FormFieldData.
			INNER_JOIN(FormTemplateFields, FormFieldData.FormTemplateFieldID.EQ(FormTemplateFields.ID)),
	).WHERE(
		FormFieldData.FormDataID.EQ(postgres.String(formDataID)).
			AND(FormTemplateFields.FieldID.EQ(postgres.String(fieldID))),
	)

	var fieldData struct {
		ID int32
	}

	if err := findFieldQuery.QueryContext(ctx, tx, &fieldData); err != nil {
		if err == sql.ErrNoRows {
			return nil // Field doesn't exist, skip update
		}
		return fmt.Errorf("failed to find field %s: %w", fieldID, err)
	}

	// Update the field value
	updateStmt := FormFieldData.UPDATE(FormFieldData.Value).
		SET(postgres.String(value)).
		WHERE(FormFieldData.ID.EQ(postgres.Int32(fieldData.ID)))

	_, err := updateStmt.ExecContext(ctx, tx)
	return err
}

// findNodesByJiraKey finds all nodes with a specific Jira key
func (s *NatsSubscriberService) findNodesByJiraKey(ctx context.Context, tx *sql.Tx, jiraKey string) ([]model.Nodes, error) {
	Nodes := table.Nodes
	findNodesStmt := postgres.SELECT(
		Nodes.AllColumns,
	).FROM(
		Nodes,
	).WHERE(
		Nodes.JiraKey.EQ(postgres.String(jiraKey)),
	)

	var nodes []model.Nodes
	if err := findNodesStmt.QueryContext(ctx, tx, &nodes); err != nil {
		return nil, fmt.Errorf("failed to find nodes with jiraKey %s: %w", jiraKey, err)
	}

	return nodes, nil
}

// findFormDataByJiraKey finds all form data IDs containing a specific Jira key
func (s *NatsSubscriberService) findFormDataByJiraKey(ctx context.Context, tx *sql.Tx, jiraKey string) ([]string, error) {
	FormData := table.FormData
	FormFieldData := table.FormFieldData
	FormTemplateFields := table.FormTemplateFields
	FormTemplateVersions := table.FormTemplateVersions
	FormTemplates := table.FormTemplates

	findFormDataQuery := postgres.SELECT(
		FormData.ID,
	).FROM(
		FormData.
			INNER_JOIN(FormFieldData, FormData.ID.EQ(FormFieldData.FormDataID)).
			INNER_JOIN(FormTemplateFields, FormFieldData.FormTemplateFieldID.EQ(FormTemplateFields.ID)).
			INNER_JOIN(FormTemplateVersions, FormData.FormTemplateVersionID.EQ(FormTemplateVersions.ID)).
			INNER_JOIN(FormTemplates, FormTemplateVersions.FormTemplateID.EQ(FormTemplates.ID)),
	).WHERE(
		FormTemplates.Tag.EQ(postgres.String("TASK")).
			AND(FormTemplateFields.FieldID.EQ(postgres.String("key"))).
			AND(FormFieldData.Value.EQ(postgres.String(jiraKey))),
	)

	var formDataIDs []struct {
		ID string
	}

	if err := findFormDataQuery.QueryContext(ctx, tx, &formDataIDs); err != nil {
		return nil, fmt.Errorf("failed to find form data with jiraKey %s: %w", jiraKey, err)
	}

	result := make([]string, len(formDataIDs))
	for i, formData := range formDataIDs {
		result[i] = formData.ID
	}

	return result, nil
}

// findActiveNodesByJiraKey finds all active nodes for a sprint with a specific Jira key
func (s *NatsSubscriberService) findActiveNodesByJiraKey(ctx context.Context, tx *sql.Tx, jiraKey string, sprintId int32) ([]string, error) {
	Nodes := table.Nodes

	// Tạo truy vấn với go-jet - chú ý không đặt alias cho Nodes.ID
	findActiveNodeQuery := postgres.SELECT(
		Nodes.AllColumns,
	).FROM(
		Nodes,
	).WHERE(
		Nodes.JiraKey.EQ(postgres.String(jiraKey)),
	)

	// In ra thông tin truy vấn để debug
	results := []model.Nodes{}

	// Thực thi truy vấn
	if err := findActiveNodeQuery.QueryContext(ctx, tx, &results); err != nil {
		slog.Error("Error finding active nodes", "error", err)
		return nil, fmt.Errorf("failed to find active nodes: %w", err)
	}

	// Chuyển đổi kết quả
	nodeIds := make([]string, len(results))
	for i, node := range results {
		nodeIds[i] = node.ID
	}

	return nodeIds, nil
}

// Shutdown stops the subscriber service
func (s *NatsSubscriberService) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down NATS Subscriber Service")

	// Unsubscribe from all subscriptions
	for _, sub := range s.subscriptions {
		if sub != nil {
			if err := sub.Unsubscribe(); err != nil {
				slog.Error("Failed to unsubscribe", "error", err)
			}
		}
	}

	// Close connections
	if s.NatsClient != nil {
		s.NatsClient.Close()
	}

	return nil
}
