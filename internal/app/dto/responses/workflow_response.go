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

type NodeDataAssigneeResponse struct {
	Id           int32  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AvatarUrl    string `json:"avatarUrl"`
	IsSystemUser bool   `json:"isSystemUser"`
}

type NodeDataResponse struct {
	Type  string `json:"type"`
	Title string `json:"title"`

	Assignee NodeDataAssigneeResponse `json:"assignee"`
	EndType  *string                  `json:"endType"`

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
	Type string `json:"type"`
}

type StoryResponse struct {
	Decoration  string `json:"decoration"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	// CategoryId  int32 /* Tự xử lý */

	Node NodeResponse `json:"node"`
}

type WorkflowDetailResponse struct {
	WorkflowResponse

	Nodes       []NodeResponse       `json:"nodes"`
	Stories     []StoryResponse      `json:"stories"`
	Connections []ConnectionResponse `json:"connections"`
}

type WorkflowResponse struct {
	Id          int32            `json:"id"`
	Title       string           `json:"title"`
	Type        string           `json:"type"`
	Category    CategoryResponse `json:"category"`
	Version     int32            `json:"version"`
	Description string           `json:"description"`
	Decoration  string           `json:"decoration"`
	IsArchived  bool             `json:"isArchived"`
	ProjectKey  string           `json:"projectKey"`

	RequestId         int32 `json:"requestId"`
	WorkflowVersionId int32 `json:"workflowVersionId"`

	types.Metadata
}

type ParticipantResponse struct {
	Avatar string
	Name   string
}

type TaskResponse struct {
	Name      string
	UpdatedAt time.Time
	Status    string
}

type RequestResponse struct {
	Id           int32
	Key          int32
	Title        int32
	ParentId     int32
	Tasks        []TaskResponse
	Participants []ParticipantResponse
	Process      int32
	CreatedAt    time.Time
	CompletedAt  *time.Time
}
