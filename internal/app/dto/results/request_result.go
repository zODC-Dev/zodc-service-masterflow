package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type RequestOverviewResult struct {
	MyRequests  int32 `json:"myrequest"`
	InProcess   int32 `json:"in_process"`
	Completed   int32 `json:"completed"`
	AllRequests int32 `json:"all_request"`
}

type Request struct {
	model.Requests
	Nodes    []model.Nodes
	Workflow model.Workflows
}

type RequestSubRequest struct {
	model.Requests
	Nodes            []model.Nodes
	Workflows        model.Workflows
	WorkflowVersions model.WorkflowVersions
}
