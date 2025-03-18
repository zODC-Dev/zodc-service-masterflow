package requests

import (
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type Connection struct {
	Id   string
	From string
	To   string
	Type string
}

type NodeDataAssignee struct {
	Id int32
}

type NodeData struct {
	Type                 string
	Title                string
	DueIn                int32
	Assignee             NodeDataAssignee
	EndType              string
	SubWorkflowVersionID *int32
}

type NodeForm struct {
	FieldId string
	Value   string
}

type Node struct {
	Id       string
	Position types.Position
	Size     types.Size
	ParentId string

	Type string

	Data NodeData

	Form []NodeForm
}

type Story struct {
	Decoration  string
	Description string
	Title       string
	Type        string
	CategoryKey string
	CategoryId  int32

	Node Node
}

type NodesConnectionsStories struct {
	Nodes       []Node /*ko chứa story*/
	Connections []Connection
	Stories     []Story
}

type CreateWorkflow struct {
	Decoration  string
	Description string
	Title       string
	Type        string
	CategoryId  int32
	ProjectKey  string

	NodesConnectionsStories
}

type StartWorkflow struct {
	Title             string
	WorkflowVersionId int32

	NodesConnectionsStories
}
