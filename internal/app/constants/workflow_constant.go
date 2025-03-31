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
	NodeStatusCompleted     NodeStatus = "COMPLETED"
	NodeStatusTodo          NodeStatus = "TO_DO"
	NodeStatusInProccessing NodeStatus = "IN_PROCESSING"
	NodeStatusOverDue       NodeStatus = "OVER_DUE"
)

type RequestStatus string

const (
	RequestStatusCompleted    RequestStatus = "COMPLETED"
	RequestStatusTodo         RequestStatus = "TO_DO"
	RequestStatusInProcessing RequestStatus = "IN_PROCESSING"
)
