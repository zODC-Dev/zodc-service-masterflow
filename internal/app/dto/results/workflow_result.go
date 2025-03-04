package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type WorkflowDetailResult struct {
	model.Workflows
	Version     model.WorkflowVersions
	Nodes       []model.WorkflowNodes
	Connections []model.WorkflowConnections
	Category    model.Categories
}

type WorkflowTemplateResult struct {
	model.Workflows
	Version  model.WorkflowVersions
	Category model.Categories
}
