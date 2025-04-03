package responses

import (
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type CategoryResponse struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Key  string `json:"key"`
}

type NodeDataResponse struct {
	Type  string `json:"type"`
	Title string `json:"title"`

	Assignee types.Assignee `json:"assignee"`
	EndType  *string        `json:"endType"`

	SubRequestID *int32 `json:"subRequestId"`

	EstimatePoint int32 `json:"estimatePoint"`
}

type NodeResponse struct {
	Id   string `json:"id"`
	Type string `json:"type"`

	Status    string `json:"status"`
	IsCurrent bool   `json:"isCurrent"`

	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`

	Position types.Position `json:"position"`
	Size     types.Size     `json:"size"`

	Data NodeDataResponse `json:"data"`

	ParentId *string `json:"parentId"`
}

type ConnectionResponse struct {
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

type StoryResponse struct {
	Decoration  string `json:"decoration"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	CategoryKey string `json:"categoryKey"`

	Node NodeResponse `json:"node"`
}

type WorkflowDetailResponse struct {
	WorkflowResponse

	Nodes       []NodeResponse       `json:"nodes"`
	Stories     []StoryResponse      `json:"stories"`
	Connections []ConnectionResponse `json:"connections"`
}

type WorkflowResponse struct {
	Id             int32            `json:"id"`
	Title          string           `json:"title"`
	Type           string           `json:"type"`
	Category       CategoryResponse `json:"category"`
	CurrentVersion int32            `json:"currentVersion"`
	Description    string           `json:"description"`
	Decoration     string           `json:"decoration"`
	IsArchived     bool             `json:"isArchived"`
	ProjectKey     string           `json:"projectKey"`

	RequestId int32 `json:"requestId"`

	types.Metadata
}
