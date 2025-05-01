package results

import (
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
)

type NodeResult struct {
	model.Nodes
	Request   model.Requests
	NodeForms []model.NodeForms
	Workflows model.Workflows
	Category  model.Categories
}

type NodeFormResult struct {
	model.NodeForms
	FormData             model.FormData
	FormFieldData        []model.FormFieldData
	FormTemplateFields   []model.FormTemplateFields
	FormTemplateVersions model.FormTemplateVersions
	FormTemplates        model.FormTemplates
	ApproveOrRejectUsers []model.NodeFormApproveOrRejectUsers
}

type FormDataResult struct {
	model.FormData
	FormFieldData        []model.FormFieldData
	FormTemplateFields   []model.FormTemplateFields
	FormTemplateVersions model.FormTemplateVersions
	FormTemplates        model.FormTemplates
}

type NodeFormCompletedFormFieldDataResult struct {
	model.FormFieldData
	FormTemplateFields model.FormTemplateFields
}

type NodeFormCompletedFormDataResult struct {
	model.FormData
	FormFieldData []NodeFormCompletedFormFieldDataResult
	FormTemplate  model.FormTemplates
}

type NodeFormCompletedResult struct {
	model.NodeForms
	Node     model.Nodes
	FormData NodeFormCompletedFormDataResult
}

type NodeRetrospectiveReportResult struct {
	model.Nodes
	Request            model.Requests
	FormFieldData      []model.FormFieldData
	FormTemplateFields []model.FormTemplateFields
}
