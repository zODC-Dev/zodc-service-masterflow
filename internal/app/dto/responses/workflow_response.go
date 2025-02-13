package responses

import "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"

type NodeResponse struct {
	Id       string         `json:"id"`
	Type     string         `json:"type"`
	Summary  string         `json:"summary"`
	Position types.Position `json:"position"`
	Size     types.Size     `json:"size"`
	EndType  string         `json:"endType"`
	ParentId string         `json:"parentId"`
	Key      string         `json:"key"`
}

type GroupResponse struct {
	Id       string         `json:"id"`
	Summary  string         `json:"summary"`
	Position types.Position `json:"position"`
	Size     types.Size     `json:"size"`
	Key      string         `json:"key"`
	Type     string         `json:"type"`
}

type ConnectionResponse struct {
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type WorkflowResponse struct {
	Title       string               `json:"title"`
	Type        string               `json:"type"`
	CategoryId  int32                `json:"categoryId"`
	Version     int32                `json:"version"`
	Description string               `json:"description"`
	Decoration  string               `json:"decoration"`
	Nodes       []NodeResponse       `json:"nodes"`
	Groups      []GroupResponse      `json:"groups"`
	Connections []ConnectionResponse `json:"connections"`
}
