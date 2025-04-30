package requests

import (
	"time"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type Connection struct {
	Id          string
	From        string
	To          string
	Text        string
	IsCompleted bool
}

type NodeDataFormAttachedApprovalOrRejectUsers struct {
	IsApproved bool
	Assignee   types.Assignee
}

type NodeDataFormAttached struct {
	Key                           string
	FormTemplateId                int32
	FormTemplateVersionId         int32
	DataId                        string
	OptionKey                     *string
	FromUserId                    *int32
	FromFormAttachedPosition      *int32
	Permission                    string
	IsOriginal                    bool
	NodeFormApprovalOrRejectUsers []NodeDataFormAttachedApprovalOrRejectUsers
	Level                         *int32
}

type NodeDataAssignee struct {
	Id int32
}

type NodeDataCondition struct {
	TrueDestinations  []string
	FalseDestinations []string
}

type NodeDataEditorContent struct {
	Subject            *string
	Body               *string
	Cc                 *[]string
	To                 *[]string
	Bcc                *[]string
	IsSendApprovedForm bool
	IsSendRejectedForm bool
}

type TaskConfigNotification struct {
	Requester    bool
	Assignee     bool
	Participants bool
}

type NodeData struct {
	Type          string
	Title         string
	EndDate       *time.Time
	Assignee      NodeDataAssignee
	EndType       string
	SubRequestID  *int32
	JiraKey       string
	EstimatePoint *float32

	JiraLinkUrl *string

	Condition NodeDataCondition

	FormAttached []NodeDataFormAttached

	EditorContent NodeDataEditorContent

	TaskCompleted TaskConfigNotification
	TaskStarted   TaskConfigNotification

	// Add for update request
	PlannedStartTime *time.Time
	PlannedEndTime   *time.Time
	ActualStartTime  *time.Time
	ActualEndTime    *time.Time

	//
	Description *string
	AttachFile  *string
}

type NodeForm struct {
	FieldId string
	Value   string
}

type Node struct {
	Id           string
	Position     types.Position
	Size         types.Size
	ParentId     string
	JiraKey      *string
	LastSyncedAt *time.Time

	Level *int32

	Type string

	Data NodeData

	Status    string
	IsCurrent bool

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

	IsSystemLinked bool
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
	SprintId    *int32

	NodesConnectionsStories
}

type StartWorkflow struct {
	Title      string
	RequestID  int32
	SprintID   *int32
	IsTemplate bool
	Template   *NodesConnectionsStories

	NodesConnectionsStories
}

type UpdateWorkflow struct {
	Decoration  string
	Description string
	Title       string
	Type        string
	CategoryId  int32
	ProjectKey  string
}
