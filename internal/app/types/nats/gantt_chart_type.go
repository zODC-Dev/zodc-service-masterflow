package nats

import "time"

// GanttChartCalculationRequest - Yêu cầu tính toán Gantt Chart
type GanttChartCalculationRequest struct {
	WorkflowId  int32                  `json:"workflow_id"`
	SprintId    int32                  `json:"sprint_id"`
	ProjectKey  string                 `json:"project_key"`
	Issues      []GanttChartJiraIssue  `json:"issues"`
	Connections []GanttChartConnection `json:"connections"`
}

// GanttChartJiraIssue - Thông tin một node/issue
type GanttChartJiraIssue struct {
	NodeId  string `json:"node_id"`
	JiraKey string `json:"jira_key,omitempty"`
	Type    string `json:"type"` // TASK, BUG, STORY
}

// GanttChartConnection - Mối quan hệ giữa các node
type GanttChartConnection struct {
	FromNodeId string `json:"from_node_id"`
	ToNodeId   string `json:"to_node_id"`
	Type       string `json:"type"` // "contains", "relates to", etc.
}

// GanttChartCalculationResponse - Kết quả tính toán Gantt Chart
type GanttChartCalculationResponse struct {
	Issues []GanttChartJiraIssueResult `json:"issues"`
}

// GanttChartJiraIssueResult - Kết quả tính toán cho một issue
type GanttChartJiraIssueResult struct {
	NodeId           string    `json:"node_id"`
	PlannedStartTime time.Time `json:"planned_start_time"`
	PlannedEndTime   time.Time `json:"planned_end_time"`
}
