package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type RequestDetailFormFieldData struct {
	model.FormFieldData
	FormTemplateField model.FormTemplateFields
}

type RequestDetailFormData struct {
	model.FormData
	FormFieldData []RequestDetailFormFieldData
}

type RequestDetailNodeForm struct {
	model.NodeForms
	NodeFormApproveUsers []model.NodeFormApproveUsers
}

type RequestDetailNode struct {
	model.Nodes
	NodeForms []RequestDetailNodeForm
	FormData  RequestDetailFormData
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
	model.Requests
	Nodes []model.Nodes
}

type Count struct {
	Count int64 `sql:"count"`
}

type RequestSubRequest struct {
	model.Requests
	Nodes            []model.Nodes
	Workflows        model.Workflows
	WorkflowVersions model.WorkflowVersions
}
