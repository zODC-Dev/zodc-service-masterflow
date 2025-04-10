package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/constants"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	natsModel "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types/nats"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/nats"
)

//////////////////// JIRA ////////////////////

type NatsService struct {
	NodeRepo    *repositories.NodeRepository
	NatsClient  *nats.NATSClient
	RequestRepo *repositories.RequestRepository
}

func NewNatsService(cfg NatsService) *NatsService {
	return &NatsService{
		NodeRepo:    cfg.NodeRepo,
		NatsClient:  cfg.NatsClient,
		RequestRepo: cfg.RequestRepo,
	}
}

// publishWorkflowToJira gửi dữ liệu workflow đến Jira và trả về phản hồi
func (s *NatsService) publishWorkflowToJira(ctx context.Context, tx *sql.Tx, nodes []requests.Node, stories []requests.Story, connections []requests.Connection, projectKey string, sprintId int32) (natsModel.WorkflowSyncResponse, error) {
	slog.Info("Starting Jira synchronization", "projectKey", projectKey, "sprintId", sprintId)
	slog.Info("Processing stories", "len stories", len(stories), "len nodes", len(nodes), "len connections", len(connections))
	slog.Info("Processing stories", "stories", stories)
	slog.Info("Processing nodes", "nodes", nodes)
	slog.Info("Processing connections", "connections", connections)

	// First, ensure story assignees have feature_leader role
	if err := s.assignFeatureLeaderRoles(stories, projectKey); err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to assign feature leader roles: %w", err)
	}

	syncRequest := natsModel.WorkflowSyncRequest{
		TransactionId: uuid.New().String(),
		ProjectKey:    projectKey,
		SprintId:      sprintId,
		Issues:        make([]natsModel.WorkflowSyncIssue, 0),
		Connections:   make([]natsModel.WorkflowSyncConnection, 0),
	}

	// Process Stories
	for _, story := range stories {
		slog.Info("Processing story",
			"id", story.Node.Id,
			"title", story.Title,
			"jiraKey", story.Node.JiraKey)

		issue := natsModel.WorkflowSyncIssue{
			NodeId:     story.Node.Id,
			Type:       "Story",
			Title:      story.Title,
			AssigneeId: &story.Node.Data.Assignee.Id,
			Action:     "create",
		}

		if story.Node.JiraKey != nil {
			issue.Action = "update"
			issue.JiraKey = *story.Node.JiraKey
		}

		syncRequest.Issues = append(syncRequest.Issues, issue)
	}

	// Process Tasks and Bugs
	for _, node := range nodes {
		slog.Info("Processing node details",
			"id", node.Id,
			"type", node.Type,
			"title", node.Data.Title,
			"jiraKey", node.JiraKey,
			"assigneeId", node.Data.Assignee.Id,
			"parentId", node.ParentId)

		if node.Type != string(constants.NodeTypeTask) && node.Type != string(constants.NodeTypeBug) {
			slog.Info("Skipping node - not a task or bug",
				"id", node.Id,
				"type", node.Type)
			continue
		}

		slog.Info("Processing node",
			"id", node.Id,
			"type", node.Type,
			"title", node.Data.Title,
			"jiraKey", node.JiraKey)

		issue := natsModel.WorkflowSyncIssue{
			NodeId:        node.Id,
			Type:          node.Type,
			Title:         node.Data.Title,
			AssigneeId:    &node.Data.Assignee.Id,
			EstimatePoint: node.Data.EstimatePoint,
			Action:        "create",
		}

		if node.JiraKey != nil {
			issue.Action = "update"
			issue.JiraKey = *node.JiraKey
		}

		syncRequest.Issues = append(syncRequest.Issues, issue)
	}

	// Process Connections
	processedConnections := make(map[string]bool) // Để tránh duplicate connections

	// 1. Xử lý connections từ connections hiện có
	for _, conn := range connections {
		fromNode := findNodeByIdFromRequest(nodes, stories, conn.From)
		toNode := findNodeByIdFromRequest(nodes, stories, conn.To)

		if fromNode == nil || toNode == nil {
			slog.Info("Skipping connection - node not found",
				"fromId", conn.From,
				"toId", conn.To)
			continue
		}

		// Skip START/END connections
		if fromNode.Type == string(constants.NodeTypeStart) ||
			toNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Sử dụng node ID để tạo connectionKey
		connectionKey := fmt.Sprintf("%s-%s", fromNode.Id, toNode.Id)
		if processedConnections[connectionKey] {
			continue
		}

		connection := natsModel.WorkflowSyncConnection{
			FromIssueKey: fromNode.Id, // Luôn sử dụng node ID
			ToIssueKey:   toNode.Id,   // Luôn sử dụng node ID
			Type:         "relates to",
		}

		syncRequest.Connections = append(syncRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 2. Xử lý connections từ parent-child relationships
	for _, node := range nodes {
		if node.ParentId == "" {
			continue
		}

		// Tìm parent node (có thể là story hoặc node khác)
		parentNode := findNodeByIdFromRequest(nodes, stories, node.ParentId)
		if parentNode == nil {
			continue
		}

		// Skip nếu parent là START/END
		if parentNode.Type == string(constants.NodeTypeStart) ||
			parentNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Sử dụng node ID để tạo connectionKey
		connectionKey := fmt.Sprintf("%s-%s", parentNode.Id, node.Id)
		if processedConnections[connectionKey] {
			continue
		}

		slog.Info("Creating parent-child connection",
			"parent", parentNode.Data.Title,
			"child", node.Data.Title,
			"parentType", parentNode.Type,
			"childType", node.Type)

		connection := natsModel.WorkflowSyncConnection{
			FromIssueKey: parentNode.Id, // Luôn sử dụng node ID
			ToIssueKey:   node.Id,       // Luôn sử dụng node ID
			Type:         "contains",
		}

		syncRequest.Connections = append(syncRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 3. Xử lý connections giữa stories và nodes
	for _, story := range stories {
		// Tìm tất cả nodes liên quan đến story này
		for _, node := range nodes {
			// Skip nếu node là START/END
			if node.Type == string(constants.NodeTypeStart) ||
				node.Type == string(constants.NodeTypeEnd) {
				continue
			}

			// Sử dụng node ID để tạo connectionKey
			connectionKey := fmt.Sprintf("%s-%s", story.Node.Id, node.Id)
			if processedConnections[connectionKey] {
				continue
			}

			// Kiểm tra mối quan hệ giữa story và node
			if node.ParentId == story.Node.Id {
				slog.Info("Creating story-node connection",
					"story", story.Title,
					"node", node.Data.Title)

				connection := natsModel.WorkflowSyncConnection{
					FromIssueKey: story.Node.Id, // Luôn sử dụng node ID
					ToIssueKey:   node.Id,       // Luôn sử dụng node ID
					Type:         "contains",
				}

				syncRequest.Connections = append(syncRequest.Connections, connection)
				processedConnections[connectionKey] = true
			}
		}
	}

	// Send to NATS
	slog.Info("Sending sync request to NATS", "request", syncRequest)
	requestBytes, err := json.Marshal(syncRequest)
	if err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to marshal sync request: %w", err)
	}

	if s.NatsClient == nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("NatsClient is nil")
	}

	response, err := s.NatsClient.Request(constants.NatsTopicWorkflowSyncRequest, requestBytes, 30*time.Second)
	if err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to sync with Jira: %w", err)
	}

	// Log raw response để debug
	slog.Info("Received raw response from Jira sync service", "response", string(response.Data))

	// Process response
	var syncResponse natsModel.WorkflowSyncResponse
	if err := json.Unmarshal(response.Data, &syncResponse); err != nil {
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to unmarshal Jira response: %w", err)
	}

	// Kiểm tra response success
	if !syncResponse.Success || !syncResponse.Data.Success {
		slog.Error("Jira synchronization failed",
			"outerSuccess", syncResponse.Success,
			"innerSuccess", syncResponse.Data.Success)
		return natsModel.WorkflowSyncResponse{}, fmt.Errorf("Jira synchronization failed")
	}

	// Update JiraKeys in database from nested structure
	for _, issue := range syncResponse.Data.Data.Issues {
		slog.Info("Updating JiraKey",
			"nodeId", issue.NodeId,
			"jiraKey", issue.JiraKey)

		if err := s.NodeRepo.UpdateJiraKey(ctx, tx, issue.NodeId, issue.JiraKey); err != nil {
			return natsModel.WorkflowSyncResponse{}, fmt.Errorf("failed to update JiraKey: %w", err)
		}
	}

	slog.Info("Completed Jira synchronization", "projectKey", projectKey)
	return syncResponse, nil
}

// New function to assign feature_leader roles to story assignees
func (s *NatsService) assignFeatureLeaderRoles(stories []requests.Story, projectKey string) error {
	// Keep track of users we've already assigned the role to avoid duplicate requests
	assignedUsers := make(map[int32]bool)

	for _, story := range stories {
		// Skip if no assignee or already processed
		if story.Node.Data.Assignee.Id == 0 || assignedUsers[story.Node.Data.Assignee.Id] {
			continue
		}

		// Mark this user as processed
		assignedUsers[story.Node.Data.Assignee.Id] = true

		// Create role assignment request
		roleRequest := natsModel.RoleAssignmentRequest{
			UserID:     story.Node.Data.Assignee.Id,
			ProjectKey: projectKey,
			RoleName:   constants.RoleFeatureLeader,
		}

		slog.Info("Assigning feature_leader role",
			"userId", roleRequest.UserID,
			"projectKey", roleRequest.ProjectKey)

		// Convert to JSON
		requestBytes, err := json.Marshal(roleRequest)
		if err != nil {
			return fmt.Errorf("failed to marshal role assignment request: %w", err)
		}

		// Send request via NATS with a 10-second timeout
		response, err := s.NatsClient.Request(constants.NatsTopicRoleAssignmentRequest, requestBytes, 10*time.Second)
		if err != nil {
			return fmt.Errorf("failed to assign feature_leader role to user %d: %w",
				story.Node.Data.Assignee.Id, err)
		}

		// Process response
		var roleResponse natsModel.RoleAssignmentResponse
		if err := json.Unmarshal(response.Data, &roleResponse); err != nil {
			return fmt.Errorf("failed to unmarshal role assignment response: %w", err)
		}

		// Check if successful
		if !roleResponse.Success {
			return fmt.Errorf("role assignment failed for user %d: %s",
				story.Node.Data.Assignee.Id, roleResponse.Message)
		}

		slog.Info("Successfully assigned feature_leader role",
			"userId", roleRequest.UserID,
			"projectKey", roleRequest.ProjectKey)
	}

	return nil
}

// Helper function to find node by ID from request structs
func findNodeByIdFromRequest(nodes []requests.Node, stories []requests.Story, nodeId string) *requests.Node {
	// Tìm trong nodes
	for i := range nodes {
		if nodes[i].Id == nodeId {
			return &requests.Node{
				Id:       nodes[i].Id,
				Type:     nodes[i].Type,
				Data:     nodes[i].Data,
				ParentId: nodes[i].ParentId,
				JiraKey:  nodes[i].JiraKey,
			}
		}
	}

	// Tìm trong stories
	for i := range stories {
		if stories[i].Node.Id == nodeId {
			return &requests.Node{
				Id:       stories[i].Node.Id,
				Type:     stories[i].Node.Type,
				Data:     stories[i].Node.Data,
				ParentId: "",
				JiraKey:  stories[i].Node.JiraKey,
			}
		}
	}

	return nil
}

// publishWorkflowToGanttChart gửi dữ liệu workflow đến service tính toán Gantt Chart
func (s *NatsService) publishWorkflowToGanttChart(ctx context.Context, tx *sql.Tx, nodes []requests.Node, stories []requests.Story, connections []requests.Connection, projectKey string, sprintId int32, workflowId int32) error {
	slog.Info("Starting Gantt Chart calculation",
		"projectKey", projectKey,
		"sprintId", sprintId,
		"workflowId", workflowId,
		"nodes", len(nodes),
		"stories", len(stories))

	// Log chi tiết nodes IDs và JiraKeys từ request
	for i, node := range nodes {
		slog.Info(fmt.Sprintf("Request node #%d details", i),
			"nodeId", node.Id,
			"title", node.Data.Title,
			"type", node.Type,
			"jiraKey", node.JiraKey)
	}

	// Lấy thông tin từ database để so sánh
	nodeMap := make(map[string]string) // map nodeId -> jiraKey

	// Trước tiên sử dụng JiraKeys từ parameters
	for _, node := range nodes {
		if node.JiraKey != nil {
			nodeMap[node.Id] = *node.JiraKey
			slog.Info("Using JiraKey from request", "nodeId", node.Id, "jiraKey", node.JiraKey)
		}
	}

	for _, story := range stories {
		if story.Node.JiraKey != nil {
			nodeMap[story.Node.Id] = *story.Node.JiraKey
			slog.Info("Using JiraKey from request story", "nodeId", story.Node.Id, "jiraKey", story.Node.JiraKey)
		}
	}

	// Chuẩn bị request
	ganttRequest := natsModel.GanttChartCalculationRequest{
		WorkflowId:  workflowId,
		SprintId:    sprintId,
		ProjectKey:  projectKey,
		Issues:      []natsModel.GanttChartJiraIssue{},
		Connections: []natsModel.GanttChartConnection{},
	}

	// Processing Stories - Add to issues
	for _, story := range stories {
		jiraKey := story.Node.JiraKey
		// Ưu tiên JiraKey từ database
		if dbJiraKey, exists := nodeMap[story.Node.Id]; exists && dbJiraKey != "" {
			jiraKey = &dbJiraKey
		}

		slog.Info("Processing story for Gantt Chart",
			"id", story.Node.Id,
			"title", story.Title,
			"type", "STORY",
			"requestJiraKey", story.Node.JiraKey,
			"dbJiraKey", nodeMap[story.Node.Id],
			"finalJiraKey", jiraKey)

		issue := natsModel.GanttChartJiraIssue{
			NodeId:  story.Node.Id,
			Type:    "STORY",
			JiraKey: *jiraKey,
		}

		ganttRequest.Issues = append(ganttRequest.Issues, issue)
	}

	// Processing Tasks and Bugs - Add to issues
	for _, node := range nodes {
		if node.Type != string(constants.NodeTypeTask) &&
			node.Type != string(constants.NodeTypeBug) &&
			node.Type != string(constants.NodeTypeStory) {
			continue
		}

		jiraKey := node.JiraKey
		// Ưu tiên JiraKey từ database
		if dbJiraKey, exists := nodeMap[node.Id]; exists && dbJiraKey != "" {
			jiraKey = &dbJiraKey
		}

		slog.Info("Processing node for Gantt Chart",
			"id", node.Id,
			"title", node.Data.Title,
			"type", node.Type,
			"requestJiraKey", node.JiraKey,
			"dbJiraKey", nodeMap[node.Id],
			"finalJiraKey", jiraKey)

		// Cảnh báo nếu không có JiraKey
		if jiraKey == nil {
			slog.Warn("Node missing JiraKey for Gantt Chart calculation",
				"nodeId", node.Id,
				"type", node.Type,
				"title", node.Data.Title)
		}

		issue := natsModel.GanttChartJiraIssue{
			NodeId:  node.Id,
			Type:    node.Type,
			JiraKey: *jiraKey,
		}

		ganttRequest.Issues = append(ganttRequest.Issues, issue)
	}

	// Processing Connections
	processedConnections := make(map[string]bool)

	// 1. Connections từ connections hiện có
	for _, conn := range connections {
		fromNode := findNodeByIdFromRequest(nodes, stories, conn.From)
		toNode := findNodeByIdFromRequest(nodes, stories, conn.To)

		if fromNode == nil || toNode == nil {
			continue
		}

		// Skip START/END connections
		if fromNode.Type == string(constants.NodeTypeStart) ||
			toNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Skip if not story/task/bug nodes
		if (fromNode.Type != string(constants.NodeTypeStory) &&
			fromNode.Type != string(constants.NodeTypeTask) &&
			fromNode.Type != string(constants.NodeTypeBug)) ||
			(toNode.Type != string(constants.NodeTypeStory) &&
				toNode.Type != string(constants.NodeTypeTask) &&
				toNode.Type != string(constants.NodeTypeBug)) {
			continue
		}

		// Tạo connection key để tránh duplicate
		connectionKey := fmt.Sprintf("%s-%s", fromNode.Id, toNode.Id)
		if processedConnections[connectionKey] {
			continue
		}

		connection := natsModel.GanttChartConnection{
			FromNodeId: fromNode.Id,
			ToNodeId:   toNode.Id,
			Type:       "relates to",
		}

		ganttRequest.Connections = append(ganttRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 2. Mối quan hệ parent-child
	for _, node := range nodes {
		if node.ParentId == "" {
			continue
		}

		// Skip if not task/bug
		if node.Type != string(constants.NodeTypeTask) &&
			node.Type != string(constants.NodeTypeBug) {
			continue
		}

		// Tìm parent node
		parentNode := findNodeByIdFromRequest(nodes, stories, node.ParentId)
		if parentNode == nil {
			continue
		}

		// Skip nếu parent không phải story/task/bug
		if parentNode.Type != string(constants.NodeTypeStory) &&
			parentNode.Type != string(constants.NodeTypeTask) &&
			parentNode.Type != string(constants.NodeTypeBug) {
			continue
		}

		// Tạo connection key để tránh duplicate
		connectionKey := fmt.Sprintf("%s-%s", parentNode.Id, node.Id)
		if processedConnections[connectionKey] {
			continue
		}

		// Tạo connection type "contains" giữa parent và node
		connection := natsModel.GanttChartConnection{
			FromNodeId: parentNode.Id,
			ToNodeId:   node.Id,
			Type:       "contains",
		}

		ganttRequest.Connections = append(ganttRequest.Connections, connection)
		processedConnections[connectionKey] = true
	}

	// 3. Mối quan hệ story-tasks
	for _, story := range stories {
		for _, node := range nodes {
			// Skip if not task/bug
			if node.Type != string(constants.NodeTypeTask) &&
				node.Type != string(constants.NodeTypeBug) {
				continue
			}

			// Tạo connection key để tránh duplicate
			connectionKey := fmt.Sprintf("%s-%s", story.Node.Id, node.Id)
			if processedConnections[connectionKey] {
				continue
			}

			// Kiểm tra nếu node thuộc story
			if node.ParentId == story.Node.Id {
				// Tạo connection type "contains" giữa story và node
				connection := natsModel.GanttChartConnection{
					FromNodeId: story.Node.Id,
					ToNodeId:   node.Id,
					Type:       "contains",
				}

				ganttRequest.Connections = append(ganttRequest.Connections, connection)
				processedConnections[connectionKey] = true
			}
		}
	}

	// Gửi request đến NATS
	slog.Info("Sending Gantt Chart calculation request to NATS", "request", ganttRequest)
	requestBytes, err := json.Marshal(ganttRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal Gantt Chart request: %w", err)
	}

	response, err := s.NatsClient.Request(constants.NatsTopicGanttChartCalculationRequest, requestBytes, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to calculate Gantt Chart: %w", err)
	}

	// Log raw response để debug
	slog.Info("Received raw response from Gantt Chart service", "response", string(response.Data))

	// Xử lý response có cấu trúc lồng ghép
	var ganttResponse natsModel.GanttChartCalculationResponse
	if err := json.Unmarshal(response.Data, &ganttResponse); err != nil {
		slog.Error("Failed to unmarshal Gantt Chart response", "error", err)
		return fmt.Errorf("failed to unmarshal Gantt Chart response: %w", err)
	}

	// Kiểm tra response success
	if !ganttResponse.Success || !ganttResponse.Data.Success {
		slog.Error("Gantt Chart calculation failed",
			"outerSuccess", ganttResponse.Success,
			"innerSuccess", ganttResponse.Data.Success)
		return fmt.Errorf("Gantt Chart calculation failed")
	}

	// Log cụ thể các node được cập nhật
	slog.Info("Issues returned from Gantt Chart service", "count", len(ganttResponse.Data.Data.Issues))
	for i, issue := range ganttResponse.Data.Data.Issues {
		slog.Info(fmt.Sprintf("Issue #%d details", i),
			"nodeId", issue.NodeId,
			"plannedStart", issue.PlannedStartTime,
			"plannedEnd", issue.PlannedEndTime)
	}

	// Cập nhật PlannedStartTime và PlannedEndTime vào database
	nodeUpdates := make([]struct {
		NodeId           string
		PlannedStartTime time.Time
		PlannedEndTime   time.Time
	}, len(ganttResponse.Data.Data.Issues))

	for i, issue := range ganttResponse.Data.Data.Issues {
		// Cập nhật PlannedStartTime và PlannedEndTime
		nodeUpdates[i] = struct {
			NodeId           string
			PlannedStartTime time.Time
			PlannedEndTime   time.Time
		}{
			NodeId:           issue.NodeId,
			PlannedStartTime: issue.PlannedStartTime,
			PlannedEndTime:   issue.PlannedEndTime,
		}
	}

	if len(nodeUpdates) > 0 {
		if err := s.NodeRepo.UpdateNodePlannedTimes(ctx, tx, nodeUpdates); err != nil {
			slog.Error("Failed to update node planned times", "error", err)
			return fmt.Errorf("failed to update node planned times: %w", err)
		}

		// Log thành công sau khi đã lưu
		slog.Info("Successfully updated planned times for nodes", "count", len(nodeUpdates))
	} else {
		slog.Warn("No planned times received from Gantt Chart service")
	}

	slog.Info("Completed Gantt Chart calculation", "projectKey", projectKey)
	return nil
}

// SyncNodeStatusToJira syncs node status changes to Jira
func (s *NatsService) SyncNodeStatusToJira(ctx context.Context, tx *sql.Tx, node model.Nodes, request model.Requests, workflow model.Workflows) error {
	// Only sync if it's a project workflow with project key
	if workflow.Type != string(constants.WorkflowTypeProject) || workflow.ProjectKey == nil {
		return nil
	}

	// Skip if no Jira key
	if node.JiraKey == nil {
		return nil
	}

	// Create sync request with only node status update
	syncRequest := natsModel.NodeStatusSyncRequest{
		TransactionId: uuid.New().String(),
		ProjectKey:    *workflow.ProjectKey,
		JiraKey:       *node.JiraKey,
		NodeId:        node.ID,
		Status:        node.Status,
	}

	// Send to NATS
	slog.Info("Sending node status update to Jira", "nodeId", node.ID, "jiraKey", node.JiraKey, "status", node.Status)
	requestBytes, err := json.Marshal(syncRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal sync request: %w", err)
	}

	response, err := s.NatsClient.Request(constants.NatsTopicNodeStatusSyncRequest, requestBytes, 30*time.Second)
	if err != nil {
		return fmt.Errorf("failed to sync with Jira: %w", err)
	}

	// Process response
	var syncResponse natsModel.NodeStatusSyncResponse
	if err := json.Unmarshal(response.Data, &syncResponse); err != nil {
		return fmt.Errorf("failed to unmarshal Jira response: %w", err)
	}

	// Check response success
	if !syncResponse.Success || !syncResponse.Data.Success {
		slog.Error("Jira synchronization failed",
			"outerSuccess", syncResponse.Success,
			"innerSuccess", syncResponse.Data.Success)
		return fmt.Errorf("Jira synchronization failed")
	}

	slog.Info("Successfully synced node status to Jira", "nodeId", node.ID, "jiraKey", node.JiraKey)
	return nil
}

// publishWorkflowEditToJira sends workflow edit data to Jira and returns the response
func (s *NatsService) publishWorkflowEditToJira(ctx context.Context, tx *sql.Tx,
	nodes []requests.Node, origNodes []model.Nodes,
	stories []requests.Story,
	connections []requests.Connection, origConnections []model.Connections,
	projectKey string, sprintId *int32) (natsModel.WorkflowEditResponse, error) {

	slog.Info("Starting Jira edit synchronization", "projectKey", projectKey, "sprintId", sprintId)
	slog.Info("Processing edit", "len stories", len(stories), "len nodes", len(nodes), "len connections", len(connections))

	// First, ensure story assignees have feature_leader role
	if err := s.assignFeatureLeaderRoles(stories, projectKey); err != nil {
		return natsModel.WorkflowEditResponse{}, fmt.Errorf("failed to assign feature leader roles: %w", err)
	}

	syncRequest := natsModel.WorkflowEditRequest{
		TransactionId:       uuid.New().String(),
		ProjectKey:          projectKey,
		SprintId:            sprintId,
		Issues:              make([]natsModel.WorkflowEditIssue, 0),
		Connections:         make([]natsModel.WorkflowEditConnection, 0),
		ConnectionsToRemove: make([]natsModel.WorkflowEditConnection, 0),
		NodeMappings:        make([]natsModel.NodeJiraMapping, 0),
	}

	// Create maps for original nodes for quick lookup
	origNodeMap := make(map[string]model.Nodes)
	for _, origNode := range origNodes {
		origNodeMap[origNode.ID] = origNode

		// Add existing mappings to the request
		if origNode.JiraKey != nil {
			syncRequest.NodeMappings = append(syncRequest.NodeMappings, natsModel.NodeJiraMapping{
				NodeId:  origNode.ID,
				JiraKey: *origNode.JiraKey,
			})
		}
	}

	// Process Stories
	for _, story := range stories {
		slog.Info("Processing story for edit",
			"id", story.Node.Id,
			"title", story.Title,
			"jiraKey", story.Node.JiraKey)

		issue := natsModel.WorkflowEditIssue{
			NodeId:     story.Node.Id,
			Type:       "Story",
			Title:      story.Title,
			AssigneeId: &story.Node.Data.Assignee.Id,
			Action:     "create",
		}

		// Check if story existed before
		origNode, storyExisted := origNodeMap[story.Node.Id]
		if storyExisted && origNode.JiraKey != nil {
			issue.Action = "update"
			issue.JiraKey = *origNode.JiraKey
		} else if story.Node.JiraKey != nil {
			issue.Action = "update"
			issue.JiraKey = *story.Node.JiraKey

			// Also add to mapping
			syncRequest.NodeMappings = append(syncRequest.NodeMappings, natsModel.NodeJiraMapping{
				NodeId:  story.Node.Id,
				JiraKey: *story.Node.JiraKey,
			})
		}

		syncRequest.Issues = append(syncRequest.Issues, issue)
	}

	// Process Tasks and Bugs
	for _, node := range nodes {
		slog.Info("Processing node for edit",
			"id", node.Id,
			"type", node.Type,
			"title", node.Data.Title,
			"jiraKey", node.JiraKey)

		if node.Type != string(constants.NodeTypeTask) && node.Type != string(constants.NodeTypeBug) {
			slog.Info("Skipping node - not a task or bug",
				"id", node.Id,
				"type", node.Type)
			continue
		}

		issue := natsModel.WorkflowEditIssue{
			NodeId:     node.Id,
			Type:       node.Type,
			Title:      node.Data.Title,
			AssigneeId: &node.Data.Assignee.Id,
			Action:     "create",
		}

		// Check if node existed before
		origNode, nodeExisted := origNodeMap[node.Id]
		if nodeExisted && origNode.JiraKey != nil {
			issue.Action = "update"
			issue.JiraKey = *origNode.JiraKey
		} else if node.JiraKey != nil {
			issue.Action = "update"
			issue.JiraKey = *node.JiraKey

			// Also add to mapping
			syncRequest.NodeMappings = append(syncRequest.NodeMappings, natsModel.NodeJiraMapping{
				NodeId:  node.Id,
				JiraKey: *node.JiraKey,
			})
		}

		syncRequest.Issues = append(syncRequest.Issues, issue)
	}

	// Create maps to track processed connections to avoid duplicates
	newConnectionsMap := make(map[string]bool)
	oldConnectionsMap := make(map[string]bool)

	// Process new connections
	for _, conn := range connections {
		fromNode := findNodeByIdFromRequest(nodes, stories, conn.From)
		toNode := findNodeByIdFromRequest(nodes, stories, conn.To)

		if fromNode == nil || toNode == nil {
			slog.Info("Skipping connection - node not found",
				"fromId", conn.From,
				"toId", conn.To)
			continue
		}

		// Skip START/END connections
		if fromNode.Type == string(constants.NodeTypeStart) ||
			toNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Create a unique key for this connection
		connectionKey := fmt.Sprintf("%s-%s-%s", fromNode.Id, toNode.Id, "relates to")
		if newConnectionsMap[connectionKey] {
			continue
		}

		connection := natsModel.WorkflowEditConnection{
			FromIssueKey: fromNode.Id,
			ToIssueKey:   toNode.Id,
			Type:         "relates to",
		}

		syncRequest.Connections = append(syncRequest.Connections, connection)
		newConnectionsMap[connectionKey] = true
	}

	// Process parent-child relationships for new connections
	for _, node := range nodes {
		if node.ParentId == "" {
			continue
		}

		// Find parent node
		parentNode := findNodeByIdFromRequest(nodes, stories, node.ParentId)
		if parentNode == nil {
			continue
		}

		// Skip if parent is START/END
		if parentNode.Type == string(constants.NodeTypeStart) ||
			parentNode.Type == string(constants.NodeTypeEnd) {
			continue
		}

		// Create a unique key for this connection
		connectionKey := fmt.Sprintf("%s-%s-%s", parentNode.Id, node.Id, "contains")
		if newConnectionsMap[connectionKey] {
			continue
		}

		connection := natsModel.WorkflowEditConnection{
			FromIssueKey: parentNode.Id,
			ToIssueKey:   node.Id,
			Type:         "contains",
		}

		syncRequest.Connections = append(syncRequest.Connections, connection)
		newConnectionsMap[connectionKey] = true
	}

	// Process original connections to identify those to be removed
	for _, origConn := range origConnections {
		// Create map of original connections
		// For connections we'll use "relates to" as default type since it's the common one
		connType := "relates to"

		// Create a unique key for this connection
		connectionKey := fmt.Sprintf("%s-%s-%s", origConn.FromNodeID, origConn.ToNodeID, connType)
		oldConnectionsMap[connectionKey] = true

		// If this connection doesn't exist in the new connections, it should be removed
		if !newConnectionsMap[connectionKey] {
			connection := natsModel.WorkflowEditConnection{
				FromIssueKey: origConn.FromNodeID,
				ToIssueKey:   origConn.ToNodeID,
				Type:         connType,
			}

			syncRequest.ConnectionsToRemove = append(syncRequest.ConnectionsToRemove, connection)
		}
	}

	// Sort and deduplicate node mappings
	nodeMapSet := make(map[string]string)
	for _, mapping := range syncRequest.NodeMappings {
		nodeMapSet[mapping.NodeId] = mapping.JiraKey
	}

	// Recreate the node mappings without duplicates
	syncRequest.NodeMappings = make([]natsModel.NodeJiraMapping, 0, len(nodeMapSet))
	for nodeId, jiraKey := range nodeMapSet {
		syncRequest.NodeMappings = append(syncRequest.NodeMappings, natsModel.NodeJiraMapping{
			NodeId:  nodeId,
			JiraKey: jiraKey,
		})
	}

	// Send to NATS
	slog.Info("Sending workflow edit request to NATS",
		"issues", len(syncRequest.Issues),
		"connections", len(syncRequest.Connections),
		"connectionsToRemove", len(syncRequest.ConnectionsToRemove),
		"nodeMappings", len(syncRequest.NodeMappings))

	requestBytes, err := json.Marshal(syncRequest)
	if err != nil {
		return natsModel.WorkflowEditResponse{}, fmt.Errorf("failed to marshal workflow edit request: %w", err)
	}

	response, err := s.NatsClient.Request(constants.NatsTopicWorkflowEditRequest, requestBytes, 30*time.Second)
	if err != nil {
		return natsModel.WorkflowEditResponse{}, fmt.Errorf("failed to sync workflow edit with Jira: %w", err)
	}

	// Process response
	var syncResponse natsModel.WorkflowEditResponse
	if err := json.Unmarshal(response.Data, &syncResponse); err != nil {
		return natsModel.WorkflowEditResponse{}, fmt.Errorf("failed to unmarshal Jira response: %w", err)
	}

	// Check response success
	if !syncResponse.Success || !syncResponse.Data.Success {
		slog.Error("Jira edit synchronization failed",
			"outerSuccess", syncResponse.Success,
			"innerSuccess", syncResponse.Data.Success)
		return natsModel.WorkflowEditResponse{}, fmt.Errorf("Jira edit synchronization failed")
	}

	// Update JiraKeys in database from response
	for _, issue := range syncResponse.Data.Data.Issues {
		slog.Info("Updating JiraKey from edit response",
			"nodeId", issue.NodeId,
			"jiraKey", issue.JiraKey)

		if err := s.NodeRepo.UpdateJiraKey(ctx, tx, issue.NodeId, issue.JiraKey); err != nil {
			return natsModel.WorkflowEditResponse{}, fmt.Errorf("failed to update JiraKey: %w", err)
		}
	}

	// Also process the node mappings from the response
	if len(syncResponse.Data.Data.NodeMappings) > 0 {
		slog.Info("Processing node mappings from response", "count", len(syncResponse.Data.Data.NodeMappings))
		for _, mapping := range syncResponse.Data.Data.NodeMappings {
			slog.Info("Mapping from response", "nodeId", mapping.NodeId, "jiraKey", mapping.JiraKey)
			if err := s.NodeRepo.UpdateJiraKey(ctx, tx, mapping.NodeId, mapping.JiraKey); err != nil {
				slog.Error("Failed to update node mapping", "error", err, "nodeId", mapping.NodeId)
				// Continue processing other mappings
			}
		}
	}

	slog.Info("Completed Jira edit synchronization", "projectKey", projectKey)
	return syncResponse, nil
}
