package types

import "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"

type FormWithFields struct {
	model.Forms
	FormFields []model.FormFields
}
