package requests

import "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"

type NodeRequest struct {
	Id       string
	Summary  string
	Type     string
	Position types.Position
	Size     types.Size
	EndType  string
	ParentId string
	Key      string
}

type GroupRequest struct {
	Id       string
	Summary  string
	Position types.Position
	Size     types.Size
	Key      string
	Type     string
}

type ConnectionRequest struct {
	Id   string
	From string
	To   string
	Type string
}

type WorkflowRequest struct {
	Title       string
	Type        string
	CategoryId  int32
	Version     int32
	Description string
	Decoration  string
	Nodes       []NodeRequest
	Groups      []GroupRequest
	Connections []ConnectionRequest
}
