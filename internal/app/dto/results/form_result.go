package results

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"

type FormTemplateResult struct {
	model.FormTemplates
	Version model.FormTemplateVersions
}

type FormSystemResult struct {
	model.FormTemplates
	Version model.FormTemplateVersions
	Fields  []model.FormTemplateFields
}
