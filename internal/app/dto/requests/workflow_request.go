package requests

import (
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type Connection struct {
	Id   string
	From string
	To   string
	Text string
}

type NodeDataFormAttached struct {
	Key                      string
	FormTemplateId           int32
	DataId                   *string
	OptionId                 *string
	FromUserId               *int32
	FromFormAttachedPosition *int32
	Permission               string
	IsOriginal               bool
	ApproveUserIds           []int32
}

type NodeDataAssignee struct {
	Id int32
}

type NodeDataCondition struct {
	TrueDestinations  []string
	FalseDestinations []string
}

type NodeData struct {
	Type         string
	Title        string
	DueIn        int32
	Assignee     NodeDataAssignee
	EndType      string
	SubRequestID *int32
	JiraKey      string

	Condition NodeDataCondition

	FormAttached []NodeDataFormAttached
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

	EstimatePoint *int32
	JiraKey       *string

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
	Title      string
	RequestID  int32
	SprintID   *int32
	IsTemplate bool

	NodesConnectionsStories
}
