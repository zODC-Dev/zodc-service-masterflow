package results

import (
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
)

type NodeResult struct {
	model.Nodes
	Request   model.Requests
	NodeForms []model.NodeForms
}

type NodeFormResult struct {
	model.NodeForms
	FormData             model.FormData
	FormFieldData        []model.FormFieldData
	FormTemplateFields   []model.FormTemplateFields
	FormTemplateVersions model.FormTemplateVersions
	FormTemplates        model.FormTemplates
}

type FormDataResult struct {
	model.FormData
	FormFieldData        []model.FormFieldData
	FormTemplateFields   []model.FormTemplateFields
	FormTemplateVersions model.FormTemplateVersions
	FormTemplates        model.FormTemplates
}
