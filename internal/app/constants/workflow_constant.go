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
)

type NodeStatus string

const (
	NodeStatusCompleted  NodeStatus = "COMPLETED"
	NodeStatusTodo       NodeStatus = "TO_DO"
	NodeStatusInProccess NodeStatus = "IN_PROCESS"
	NodeStatusOverDue    NodeStatus = "OVER_DUE"
)

type RequestStatus string

const (
	RequestStatusCompleted RequestStatus = "COMPLETED"
	RequestStatusTodo      RequestStatus = "TO_DO"
	RequestStatusInProcess RequestStatus = "IN_PROCESS"
)
