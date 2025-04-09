package constants

const (
	NatsTopicRoleAssignmentRequest = "role.assign.request"
	NatsTopicWorkflowSyncRequest   = "workflow.sync.request"
	NatsTopicNodeStatusSyncRequest = "node.status.sync.request"

	// Gantt Chart calculation
	NatsTopicGanttChartCalculationRequest  = "ganttchart.calculation.request"
	NatsTopicGanttChartCalculationResponse = "ganttchart.calculation.response"
)
