package repositories

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
)

type FormFieldRepository struct{}

func NewFormFieldRepository() *FormFieldRepository {
	return &FormFieldRepository{}
}

func (r *FormFieldRepository) Create(ctx context.Context, tx *sql.Tx, formFields []model.FormFields, formId int32) error {
	FormFields := table.FormFields

	for i := range formFields {
		formFields[i].FormID = formId
	}

	formFieldsInsertColumns := FormFields.AllColumns.Except(FormFields.ID, FormFields.CreatedAt, FormFields.UpdatedAt, FormFields.DeletedAt)
	formFieldsStmt := FormFields.INSERT(formFieldsInsertColumns).MODELS(formFields)

	if _, err := formFieldsStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
