package nats

// WorkflowSyncIssue represents an issue to be created/updated in Jira
type WorkflowSyncIssue struct {
	NodeId     string `json:"node_id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	AssigneeId *int32 `json:"assignee_id,omitempty"`
	JiraKey    string `json:"jira_key,omitempty"`
	Action     string `json:"action"`
}

// JiraSyncLink represents a link between two issues in Jira
type WorkflowSyncConnection struct {
	FromIssueKey string `json:"from_issue_key"`
	ToIssueKey   string `json:"to_issue_key"`
	Type         string `json:"type"`
}

// WorkflowSyncRequest represents the request sent to Jira via NATS
type WorkflowSyncRequest struct {
	TransactionId string                   `json:"transaction_id"`
	ProjectKey    string                   `json:"project_key"`
	SprintId      int32                    `json:"sprint_id"`
	Issues        []WorkflowSyncIssue      `json:"issues"`
	Connections   []WorkflowSyncConnection `json:"connections"`
}

// WorkflowSyncResponse represents the response received from Jira via NATS
type WorkflowSyncResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Success bool `json:"success"`
		Data    struct {
			Issues []struct {
				NodeId  string `json:"node_id"`
				JiraKey string `json:"jira_key"`
			} `json:"issues"`
		} `json:"data"`
	} `json:"data"`
}
