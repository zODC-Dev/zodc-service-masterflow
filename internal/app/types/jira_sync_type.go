package types

// JiraSyncIssue represents an issue to be created/updated in Jira
type JiraSyncIssue struct {
	NodeId     string `json:"node_id"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	AssigneeId *int32 `json:"assignee_id,omitempty"`
	JiraKey    string `json:"jira_key,omitempty"`
	Action     string `json:"action"`
}

// JiraSyncLink represents a link between two issues in Jira
type JiraSyncConnection struct {
	FromIssueKey string `json:"from_issue_key"`
	ToIssueKey   string `json:"to_issue_key"`
	Type         string `json:"type"`
}

// JiraSyncRequest represents the request sent to Jira via NATS
type JiraSyncRequest struct {
	TransactionId string               `json:"transaction_id"`
	ProjectKey    string               `json:"project_key"`
	Issues        []JiraSyncIssue      `json:"issues"`
	Connections   []JiraSyncConnection `json:"connections"`
}

// JiraSyncResponse represents the response received from Jira via NATS
type JiraSyncResponse struct {
	Issues []struct {
		NodeId  string `json:"node_id"`
		JiraKey string `json:"jira_key"`
	} `json:"issues"`
}
