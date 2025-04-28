package nats

// WorkflowEditIssue represents an issue to be created or updated during workflow edit
type WorkflowEditIssue struct {
	NodeId        string   `json:"node_id"`
	Type          string   `json:"type"` // "Story", "Task", "Bug", etc.
	Title         string   `json:"title"`
	AssigneeId    *int32   `json:"assignee_id,omitempty"`
	JiraKey       string   `json:"jira_key,omitempty"`
	EstimatePoint *float32 `json:"estimate_point,omitempty"`
	Action        string   `json:"action"` // "create" or "update"
}

// WorkflowEditConnection represents a connection between issues
type WorkflowEditConnection struct {
	FromIssueKey string `json:"from_issue_key"` // node_id or jira_key
	ToIssueKey   string `json:"to_issue_key"`   // node_id or jira_key
	Type         string `json:"type"`           // "relates to" or "contains"
}

// NodeJiraMapping represents a mapping between node ID and Jira key
type NodeJiraMapping struct {
	NodeId  string `json:"node_id"`
	JiraKey string `json:"jira_key"`
}

// WorkflowEditRequest represents the request to edit a workflow in Jira
type WorkflowEditRequest struct {
	TransactionId       string                   `json:"transaction_id"`
	ProjectKey          string                   `json:"project_key"`
	SprintId            *int32                   `json:"sprint_id"`
	Issues              []WorkflowEditIssue      `json:"issues"`                // New/updated issues
	Connections         []WorkflowEditConnection `json:"connections"`           // New connections
	ConnectionsToRemove []WorkflowEditConnection `json:"connections_to_remove"` // Connections to remove
	NodeMappings        []NodeJiraMapping        `json:"node_mappings"`         // Mappings between node IDs and Jira keys
}

// WorkflowEditResponse represents the response from Jira edit operation
type WorkflowEditResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Success bool `json:"success"`
		Data    struct {
			Issues []struct {
				NodeId  string `json:"node_id"`
				JiraKey string `json:"jira_key"`
			} `json:"issues"`
			NodeMappings []NodeJiraMapping `json:"node_mappings"` // Updated mappings between node IDs and Jira keys
		} `json:"data"`
	} `json:"data"`
}
