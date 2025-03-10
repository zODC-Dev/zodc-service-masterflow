package responses

import "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"

type CategoryResponse struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Key  string `json:"key"`
}

type NodeDataAssigneeResponse struct {
	Id           int32  `json:"id"`
	Email        string `json:"email"`
	AvatarUrl    string `json:"avatarUrl"`
	IsSystemUser bool   `json:"isSystemUser"`
}

type NodeDataResponse struct {
	Type     string                   `json:"type"`
	Title    string                   `json:"title"`
	DueIn    *int32                   `json:"dueIn"`
	Assignee NodeDataAssigneeResponse `json:"assignee"`
	EndType  *string                  `json:"endType"`

	SubWorkflowVersionId *int32 `json:"subWorkflowVersionId"`
}

type NodeResponse struct {
	Id   string `json:"id"`
	Type string `json:"type"`

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
	ProjectKey  string           `json:"procjectKey"`

	types.Metadata
}

type RequestResponse struct {
	Id           int32
	Key          int32
	Title        int32
	Parent_id    int32
	Task         string
	Participants []struct {
		avatar string
		name   string
	}
	UpdatedAt string
	CreatedAt string
}
