package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type FormTemplateVersionResult struct {
	model.FormTemplateVersions
	Fields []model.FormTemplateFields
}

type FormTemplateResult struct {
	model.FormTemplates
	Versions []FormTemplateVersionResult
}

type FormSystemResult struct {
	model.FormTemplates
	Version model.FormTemplateVersions
	Fields  []model.FormTemplateFields
}
