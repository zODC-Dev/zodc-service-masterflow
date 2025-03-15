package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type RequestDetail struct {
	model.Requests
	Workflow    model.Workflows
	Version     model.WorkflowVersions
	Nodes       []model.Nodes
	Connections []model.Connections
	Category    model.Categories
}

type WorkflowTemplate struct {
	model.Workflows
	Version  model.WorkflowVersions
	Category model.Categories
}

type ConnectionWithNode struct {
	model.Connections
	Node model.Nodes
}
