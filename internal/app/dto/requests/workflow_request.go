package requests

import "github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"

type NodeRequest struct {
	Id       string
	NodeType string
	Title    string
	Position types.Position
	Size     types.Size
	EndType  string
	ParentId string
	TicketID string
}

type GroupRequest struct {
	Id        string
	Title     string
	Position  types.Position
	Size      types.Size
	TicketID  string
	GroupType string
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
