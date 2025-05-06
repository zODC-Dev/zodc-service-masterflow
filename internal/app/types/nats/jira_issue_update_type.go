package nats

import "time"

type JiraIssueUpdateMessage struct {
	JiraKey       string     `json:"jiraKey"`
	Summary       string     `json:"summary"`
	Description   string     `json:"description"`
	AssigneeEmail string     `json:"assigneeEmail"`
	AssigneeId    *string    `json:"assigneeId"`
	EstimatePoint *float32   `json:"estimatePoint"`
	Status        string     `json:"status"`
	OldStatus     *string    `json:"oldStatus"`
	SprintId      *int32     `json:"sprintId"`
	LastSyncedAt  *time.Time `json:"lastSyncedAt"`
}
