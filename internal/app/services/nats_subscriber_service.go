package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	natslib "github.com/nats-io/nats.go"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	natsModel "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types/nats"
	nats "github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

// NatsSubscriberService handles listening for messages from NATS
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

	// Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		slog.Error("Failed to start transaction", "error", err)
		return
	}
	defer tx.Rollback()

	// Process message based on subject
	var processErr error
	switch msg.Subject {
	case "jira.issue.update":
		processErr = s.handleJiraIssueUpdate(tx, msg.Data)
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

	slog.Info("Successfully processed message", "subject", msg.Subject)
}

// handleJiraIssueUpdate processes messages from jira.issue.update
func (s *NatsSubscriberService) handleJiraIssueUpdate(tx *sql.Tx, data []byte) error {
	// Parse message data
	var message natsModel.JiraIssueUpdateMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return fmt.Errorf("failed to unmarshal Jira issue update message: %w", err)
	}

	// Validation
	if message.JiraKey == "" {
		return fmt.Errorf("missing required field: jiraKey")
	}

	// Map Jira status to system status
	var systemStatus string
	switch message.Status {
	case "To Do":
		systemStatus = string(constants.NodeStatusTodo) // "TO_DO"
	case "In Progress":
		systemStatus = string(constants.NodeStatusInProgress) // "IN_PROGRESS"
	case "Done":
		systemStatus = string(constants.NodeStatusCompleted) // "COMPLETED"
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

	// Step 1: Update all nodes with the same Jira key to maintain data consistency
	updateNodeQuery := `
		UPDATE nodes
		SET title = $1,
			status = $2
		WHERE jira_key = $3
	`

	// Additional fields that may be null
	var updateParams []interface{}
	var additionalSets []string

	// Add estimate point if provided
	if message.EstimatePoint != nil {
		additionalSets = append(additionalSets, "estimate_point = $"+strconv.Itoa(4+len(updateParams)))
		updateParams = append(updateParams, *message.EstimatePoint)
	}

	// Add assignee ID if provided
	if message.AssigneeId != nil {
		// Convert string assignee ID to int32
		assigneeId, err := strconv.ParseInt(*message.AssigneeId, 10, 32)
		if err != nil {
			slog.Warn("Failed to parse assignee ID", "assigneeId", *message.AssigneeId, "error", err)
		} else {
			assigneeId32 := int32(assigneeId)
			additionalSets = append(additionalSets, "assignee_id = $"+strconv.Itoa(4+len(updateParams)))
			updateParams = append(updateParams, assigneeId32)
		}
	}

	// Build the final query with additional fields
	finalUpdateQuery := updateNodeQuery
	if len(additionalSets) > 0 {
		// Cắt bỏ phần WHERE từ câu query gốc để thêm các fields
		wherePos := len(finalUpdateQuery) - 21 // Độ dài của "\n\t\tWHERE jira_key = $3\n\t"
		finalUpdateQuery = finalUpdateQuery[:wherePos]

		// Thêm các fields bổ sung
		for _, set := range additionalSets {
			finalUpdateQuery += ", " + set
		}

		// Thêm lại mệnh đề WHERE
		finalUpdateQuery += "\n\t\tWHERE jira_key = $3\n\t"
	}

	// Add base parameters
	baseParams := []interface{}{
		message.Summary,
		systemStatus,
		message.JiraKey,
	}

	// Combine all parameters
	allParams := append(baseParams, updateParams...)

	// Execute the update for all nodes with this Jira key
	_, err := tx.ExecContext(context.Background(), finalUpdateQuery, allParams...)
	if err != nil {
		return fmt.Errorf("failed to update nodes: %w", err)
	}

	// Step 2: Update form data for all corresponding forms
	formUpdateQuery := `
		UPDATE form_field_data ffd
		SET value = CASE
			WHEN form_template_field_id = (
				SELECT ftf.id
				FROM form_template_fields ftf
				WHERE ftf.field_id = 'summary'
			) THEN $1
			WHEN form_template_field_id = (
				SELECT ftf.id
				FROM form_template_fields ftf
				WHERE ftf.field_id = 'assignee_email'
			) THEN $2
			WHEN form_template_field_id = (
				SELECT ftf.id
				FROM form_template_fields ftf
				WHERE ftf.field_id = 'status'
			) THEN $3
			ELSE value
		END
		WHERE form_data_id IN (
			SELECT fd.id
			FROM form_data fd
			INNER JOIN form_field_data ffd ON ffd.form_data_id = fd.id
			INNER JOIN form_template_fields ftf ON ftf.id = ffd.form_template_field_id
			INNER JOIN form_template_versions ftv ON ftv.id = fd.form_template_version_id
			INNER JOIN form_templates ft ON ft.id = ftv.form_template_id
			WHERE ft.tag = 'TASK' AND ftf.field_id = 'key' AND ffd.value = $4
		)
	`

	// Execute the update for form data
	_, err = tx.ExecContext(
		context.Background(),
		formUpdateQuery,
		message.Summary,       // $1 - Summary field
		message.AssigneeEmail, // $2 - Assignee email field
		systemStatus,          // $3 - Status field
		message.JiraKey,       // $4 - Jira key for finding matching forms
	)
	if err != nil {
		return fmt.Errorf("failed to update form field data: %w", err)
	}

	// Step 3: Find active node for workflow state transition
	if message.SprintId != nil && message.OldStatus != nil {
		findActiveNodeQuery := `
			SELECT n.id
			FROM requests r 
			INNER JOIN nodes n ON n.request_id = r.id
			WHERE r.status = 'IN_PROGRESS' 
			AND r.sprint_id = $1 
			AND n.jira_key = $2
		`

		var activeNodeId string
		err := tx.QueryRowContext(context.Background(), findActiveNodeQuery, *message.SprintId, message.JiraKey).Scan(&activeNodeId)

		if err != nil {
			slog.Error("Error finding active node", "error", err)
		} else {

			// Step 4: Apply workflow state transition based on status change
			if message.OldStatus != nil && message.Status != "" {
				// Use mapped system status values for state transition logic
				if oldSystemStatus != nil && systemStatus != "" {
					// From Todo to In Progress -> Start Node
					if *oldSystemStatus == string(constants.NodeStatusTodo) && systemStatus == string(constants.NodeStatusInProgress) {
						// Use NodeService to start the node - this will be handled in a separate transaction
						ctx := context.Background()

						// GET ASSIGNEE ID
						assigneeId := 0
						if message.AssigneeId != nil {
							assigneeId, _ = strconv.Atoi(*message.AssigneeId)
						}
						if err := s.NodeService.StartNodeHandler(ctx, int32(assigneeId), activeNodeId); err != nil {
							slog.Error("Failed to start node", "nodeId", activeNodeId, "error", err)
							// Continue execution even if node start fails
						}
					} else if *oldSystemStatus == string(constants.NodeStatusInProgress) && systemStatus == string(constants.NodeStatusCompleted) {
						// Use NodeService to complete the node - this will be handled in a separate transaction
						ctx := context.Background()

						// Use system user ID for completion or a default value
						var systemUserId int32 = 1 // default system user ID
						if message.AssigneeId != nil {
							if assigneeId, err := strconv.ParseInt(*message.AssigneeId, 10, 32); err == nil {
								systemUserId = int32(assigneeId)
							}
						}

						if err := s.NodeService.CompleteNodeHandler(ctx, activeNodeId, systemUserId, false); err != nil {
							slog.Error("Failed to complete node", "nodeId", activeNodeId, "error", err)
							// Continue execution even if node completion fails
						}
					}
				}
			}
		}
	}

	// Optionally, also update the estimate point in form data if provided
	if message.EstimatePoint != nil {
		estimatePointQuery := `
			UPDATE form_field_data ffd
			SET value = $1
			WHERE form_template_field_id = (
				SELECT ftf.id
				FROM form_template_fields ftf
				WHERE ftf.field_id = 'estimate_point'
			)
			AND form_data_id IN (
				SELECT fd.id
				FROM form_data fd
				INNER JOIN form_field_data ffd ON ffd.form_data_id = fd.id
				INNER JOIN form_template_fields ftf ON ftf.id = ffd.form_template_field_id
				INNER JOIN form_template_versions ftv ON ftv.id = fd.form_template_version_id
				INNER JOIN form_templates ft ON ft.id = ftv.form_template_id
				WHERE ft.tag = 'TASK' AND ftf.field_id = 'key' AND ffd.value = $2
			)
		`

		_, err = tx.ExecContext(
			context.Background(),
			estimatePointQuery,
			fmt.Sprintf("%f", *message.EstimatePoint), // Convert to string
			message.JiraKey,
		)
		if err != nil {
			slog.Warn("Failed to update estimate point", "error", err)
		}
	}

	return nil
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
