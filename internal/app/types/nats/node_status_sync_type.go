package nats

import "time"

// NodeStatusSyncRequest represents the request to sync node status with Jira
type NodeStatusSyncRequest struct {
	TransactionId string     `json:"transaction_id"`
	ProjectKey    string     `json:"project_key"`
	JiraKey       string     `json:"jira_key"`
	NodeId        string     `json:"node_id"`
	Status        string     `json:"status"`
	LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`
}

// NodeStatusSyncResponse represents the response from Jira sync service
type NodeStatusSyncResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Success bool `json:"success"`
		Data    struct {
			NodeId  string `json:"node_id"`
			JiraKey string `json:"jira_key"`
			Status  string `json:"status"`
		} `json:"data"`
	} `json:"data"`
}
