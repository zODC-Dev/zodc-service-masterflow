package requests

import (
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type ConnectionRequest struct {
	Id   string
	From string
	To   string
	Type string
}

type NodeDataRequest struct {
	Type       string
	Title      string
	DueIn      int32
	AssigneeId int32
	EndType    string
}

type NodeFormRequest struct {
	FieldId string
	Value   string
}

type NodeRequest struct {
	Id       string
	Position types.Position
	Size     types.Size
	ParentId string

	Type string

	Data NodeDataRequest

	Form []NodeFormRequest
}

type StoriesRequest struct {
	Decoration  string
	Description string
	Title       string
	Type        string
	// CategoryId  int32 /* Tự xử lý */

	Node NodeRequest
}

type WorkflowRequest struct {
	Decoration  string
	Description string
	Title       string
	Type        string
	CategoryId  int32

	Nodes       []NodeRequest /*ko chứa story*/
	Connections []ConnectionRequest

	Stories []StoriesRequest
}
