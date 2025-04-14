package constants

type WorkflowType string

const (
	WorkflowTypeGeneral WorkflowType = "GENERAL"
	WorkflowTypeProject WorkflowType = "PROJECT"
)

type NodeType string

const (
	NodeTypeStart       NodeType = "START"
	NodeTypeEnd         NodeType = "END"
	NodeTypeSubWorkflow NodeType = "SUB_WORKFLOW"
	NodeTypeStory       NodeType = "STORY"
	NodeTypeCondition   NodeType = "CONDITION"
	NodeTypeTask        NodeType = "TASK"
	NodeTypeBug         NodeType = "BUG"
)

type NodeStatus string

const (
	NodeStatusCompleted  NodeStatus = "COMPLETED"
	NodeStatusTodo       NodeStatus = "TO_DO"
	NodeStatusInProgress NodeStatus = "IN_PROGRESS"
	NodeStatusOverDue    NodeStatus = "OVER_DUE"
)

type RequestStatus string

const (
	RequestStatusCompleted  RequestStatus = "COMPLETED"
	RequestStatusTodo       RequestStatus = "TO_DO"
	RequestStatusInProgress RequestStatus = "IN_PROGRESS"
	RequestStatusCanceled   RequestStatus = "CANCELED"
	RequestStatusTerminated RequestStatus = "TERMINATED"
)

type NodeEndType string

const (
	NodeEndTypeComplete  NodeEndType = "COMPLETE"
	NodeEndTypeTerminate NodeEndType = "TERMINATE"
)
