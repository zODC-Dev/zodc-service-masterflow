package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
)

type FormRepository struct{}

func NewFormRepository() *FormRepository {
	return &FormRepository{}
}

func (r *FormRepository) FindAllFormTemplate(ctx context.Context, db *sql.DB) ([]results.FormTemplateResult, error) {
	FormTemplates := table.FormTemplates
	FormTemplateVersions := table.FormTemplateVersions

	statement := postgres.SELECT(
		FormTemplates.AllColumns,
		FormTemplateVersions.AllColumns,
	).FROM(
		FormTemplates.
			INNER_JOIN(FormTemplateVersions, FormTemplates.ID.EQ(FormTemplateVersions.FormTemplateID).
				AND(FormTemplateVersions.IsArchived.EQ(postgres.Bool(false)))),
	)

	results := []results.FormTemplateResult{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *FormRepository) FindAllFormTemplateFieldsByFormTemplateVersionId(ctx context.Context, db *sql.DB, formTemplateVersionId int32) ([]model.FormTemplateFields, error) {
	FormTemplateFields := table.FormTemplateFields

	statement := postgres.SELECT(
		FormTemplateFields.AllColumns,
	).FROM(
		FormTemplateFields,
	).WHERE(
		FormTemplateFields.FormTemplateVersionID.EQ(postgres.Int32(formTemplateVersionId)),
	)

	results := []model.FormTemplateFields{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *FormRepository) CreateFormTemplate(ctx context.Context, tx *sql.Tx, formTemplate model.FormTemplates) (model.FormTemplates, error) {
	FormTemplates := table.FormTemplates

	columns := FormTemplates.AllColumns.Except(FormTemplates.ID, FormTemplates.CreatedAt, FormTemplates.UpdatedAt, FormTemplates.DeletedAt)

	statement := FormTemplates.INSERT(columns).MODEL(formTemplate).RETURNING(FormTemplates.ID)

	err := statement.QueryContext(ctx, tx, &formTemplate)

	return formTemplate, err
}

func (r *FormRepository) CreateFormTemplateVersion(ctx context.Context, tx *sql.Tx, formTemplateVersionModel model.FormTemplateVersions) (model.FormTemplateVersions, error) {
	FormTemplateVersions := table.FormTemplateVersions

	columns := FormTemplateVersions.AllColumns.Except(FormTemplateVersions.ID, FormTemplateVersions.CreatedAt, FormTemplateVersions.UpdatedAt, FormTemplateVersions.DeletedAt)

	statement := FormTemplateVersions.INSERT(columns).MODEL(formTemplateVersionModel).RETURNING(FormTemplateVersions.ID)

	err := statement.QueryContext(ctx, tx, &formTemplateVersionModel)

	return formTemplateVersionModel, err
}

func (r *FormRepository) CreateFormTemplateFields(ctx context.Context, tx *sql.Tx, formTemplateFieldModels []model.FormTemplateFields) error {
	FormTemplateFields := table.FormTemplateFields

	columns := FormTemplateFields.AllColumns.Except(FormTemplateFields.ID, FormTemplateFields.CreatedAt, FormTemplateFields.UpdatedAt, FormTemplateFields.DeletedAt)

	statement := FormTemplateFields.INSERT(columns).MODELS(formTemplateFieldModels)

	err := statement.QueryContext(ctx, tx, &formTemplateFieldModels)

	return err
}

func (r *FormRepository) FindAllFormSystem(ctx context.Context, db *sql.DB) ([]results.FormSystemResult, error) {
	FormTemplates := table.FormTemplates
	FormTemplateVersions := table.FormTemplateVersions
	FormTemplateFields := table.FormTemplateFields

	statement := postgres.SELECT(
		FormTemplates.AllColumns,
		FormTemplateVersions.AllColumns,
		FormTemplateFields.AllColumns,
	).FROM(
		FormTemplates.
			INNER_JOIN(FormTemplateVersions, FormTemplateVersions.FormTemplateID.EQ(FormTemplates.ID)).
			INNER_JOIN(FormTemplateFields, FormTemplateFields.FormTemplateVersionID.EQ(FormTemplateVersions.ID)),
	).WHERE(
		FormTemplates.Type.EQ(postgres.String("SYSTEM")),
	)

	formSystemResults := []results.FormSystemResult{}

	err := statement.QueryContext(ctx, db, &formSystemResults)

	return formSystemResults, err
}

func (r *FormRepository) CreateFormData(ctx context.Context, tx *sql.Tx, formData model.FormData) (model.FormData, error) {
	FormData := table.FormData

	columns := FormData.AllColumns.Except(FormData.ID, FormData.CreatedAt, FormData.UpdatedAt, FormData.DeletedAt)

	statement := FormData.INSERT(columns).MODEL(formData).RETURNING(FormData.AllColumns)

	err := statement.QueryContext(ctx, tx, &formData)

	return formData, err

}

func (r *FormRepository) CreateFormFieldDatas(ctx context.Context, tx *sql.Tx, formFieldDatas []model.FormFieldData) error {
	FormFieldData := table.FormFieldData

	columns := FormFieldData.AllColumns.Except(FormFieldData.ID, FormFieldData.CreatedAt, FormFieldData.UpdatedAt, FormFieldData.DeletedAt)

	statement := FormFieldData.INSERT(columns).MODELS(formFieldDatas)

	err := statement.QueryContext(ctx, tx, &formFieldDatas)

	return err
}
