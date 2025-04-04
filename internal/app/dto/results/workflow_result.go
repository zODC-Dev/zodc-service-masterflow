package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type RequestDetailNode struct {
	model.Nodes
	NodeForms []model.NodeForms
}

type RequestDetail struct {
	model.Requests
	Workflow    model.Workflows
	Version     model.WorkflowVersions
	Nodes       []RequestDetailNode
	Connections []model.Connections
	Category    model.Categories
}

type WorkflowTemplate struct {
	model.Workflows
	Version  model.WorkflowVersions
	Category model.Categories
	Request  model.Requests
}

type ConnectionWithNode struct {
	model.Connections
	Node model.Nodes
}

type Request struct {
	Count int64 `sql:"count"`
	model.Requests
	Nodes []model.Nodes
}

type Count struct {
	Count int64 `sql:"count"`
}
