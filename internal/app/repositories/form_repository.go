package repositories

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/results"
)

type FormRepository struct{}

func NewFormRepository() *FormRepository {
	return &FormRepository{}
}

func (r *FormRepository) FindAllFormTemplate(ctx context.Context, db *sql.DB, queryParam queryparams.FormQueryParam) ([]results.FormTemplateResult, error) {
	FormTemplates := table.FormTemplates
	FormTemplateVersions := table.FormTemplateVersions
	Categories := table.Categories

	statement := postgres.SELECT(
		FormTemplates.AllColumns,
		FormTemplateVersions.AllColumns,
		Categories.AllColumns,
	).FROM(
		FormTemplates.
			LEFT_JOIN(FormTemplateVersions, FormTemplates.ID.EQ(FormTemplateVersions.FormTemplateID)).
			LEFT_JOIN(Categories, Categories.ID.EQ(FormTemplates.CategoryID)),
	)

	conditions := []postgres.BoolExpression{
		FormTemplates.Type.EQ(postgres.String("USER")),
	}

	if queryParam.CategoryID != "" {
		categoryIdInt, err := strconv.Atoi(queryParam.CategoryID)
		if err != nil {
			return []results.FormTemplateResult{}, err
		}

		conditions = append(conditions, Categories.ID.EQ(postgres.Int32(int32(categoryIdInt))))
	}

	if queryParam.Search != "" {
		conditions = append(conditions, postgres.LOWER(FormTemplates.Title).LIKE(postgres.LOWER(postgres.String("%"+queryParam.Search+"%"))))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

	results := []results.FormTemplateResult{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *FormRepository) FindAllFormTemplateFieldsByFormTemplateId(ctx context.Context, db *sql.DB, formTemplateId int32) ([]model.FormTemplateFields, error) {
	FormTemplateFields := table.FormTemplateFields
	FormTemplateVersions := table.FormTemplateVersions
	FormTemplates := table.FormTemplates

	statement := postgres.SELECT(
		FormTemplateFields.AllColumns,
	).FROM(
		FormTemplateFields.LEFT_JOIN(
			FormTemplateVersions,
			FormTemplateVersions.FormTemplateID.EQ(FormTemplateFields.FormTemplateVersionID),
		).LEFT_JOIN(
			FormTemplates,
			FormTemplates.ID.EQ(FormTemplateVersions.FormTemplateID),
		),
	).WHERE(
		FormTemplates.ID.EQ(postgres.Int32(formTemplateId)).AND(
			FormTemplateVersions.Version.EQ(FormTemplates.CurrentVersion),
		),
	)

	results := []model.FormTemplateFields{}

	err := statement.QueryContext(ctx, db, &results)

	return results, err
}

func (r *FormRepository) CreateFormTemplate(ctx context.Context, tx *sql.Tx, formTemplate model.FormTemplates) (model.FormTemplates, error) {
	FormTemplates := table.FormTemplates

	columns := FormTemplates.AllColumns.Except(FormTemplates.ID, FormTemplates.CreatedAt, FormTemplates.UpdatedAt, FormTemplates.DeletedAt, FormTemplates.CurrentVersion)

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
			LEFT_JOIN(FormTemplateVersions, FormTemplateVersions.FormTemplateID.EQ(FormTemplates.ID)).
			LEFT_JOIN(FormTemplateFields, FormTemplateFields.FormTemplateVersionID.EQ(FormTemplateVersions.ID)),
	).WHERE(
		FormTemplates.Type.EQ(postgres.String("SYSTEM")),
	)

	formSystemResults := []results.FormSystemResult{}

	err := statement.QueryContext(ctx, db, &formSystemResults)

	return formSystemResults, err
}

func (r *FormRepository) CreateFormData(ctx context.Context, tx *sql.Tx, formData model.FormData) (model.FormData, error) {
	FormData := table.FormData

	columns := FormData.AllColumns.Except(FormData.CreatedAt, FormData.UpdatedAt, FormData.DeletedAt)

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

func (r *FormRepository) FindOneFormTemplateByFormTemplateId(ctx context.Context, db *sql.DB, formTemplateId int32) (results.FormSystemResult, error) {
	FormTemplates := table.FormTemplates
	FormTemplateVersions := table.FormTemplateVersions
	FormTemplateFields := table.FormTemplateFields
	Categories := table.Categories

	statement := FormTemplates.SELECT(
		FormTemplates.AllColumns,
		FormTemplateVersions.AllColumns,
		FormTemplateFields.AllColumns,
		Categories.AllColumns,
	).FROM(
		FormTemplates.
			LEFT_JOIN(FormTemplateVersions, FormTemplateVersions.FormTemplateID.EQ(FormTemplates.ID)).
			LEFT_JOIN(FormTemplateFields, FormTemplateFields.FormTemplateVersionID.EQ(FormTemplateVersions.ID)).
			LEFT_JOIN(Categories, Categories.ID.EQ(FormTemplates.CategoryID)),
	).WHERE(
		FormTemplates.ID.EQ(postgres.Int32(formTemplateId)),
	)

	formTemplate := results.FormSystemResult{}

	err := statement.QueryContext(ctx, db, &formTemplate)

	return formTemplate, err
}

func (r *FormRepository) FindFormDataById(ctx context.Context, db *sql.DB, formDataId string) (results.FormDataResult, error) {
	FormData := table.FormData
	FormFieldData := table.FormFieldData
	FormTemplateVersions := table.FormTemplateVersions
	FormTemplateFields := table.FormTemplateFields

	statement := FormData.SELECT(
		FormData.AllColumns,
		FormTemplateFields.AllColumns,
		FormFieldData.AllColumns,
	).FROM(
		FormData.
			LEFT_JOIN(FormFieldData, FormData.ID.EQ(FormFieldData.FormDataID)).
			LEFT_JOIN(FormTemplateVersions, FormData.FormTemplateVersionID.EQ(FormTemplateVersions.ID)).
			LEFT_JOIN(FormTemplateFields, FormTemplateVersions.ID.EQ(FormTemplateFields.FormTemplateVersionID)),
	).WHERE(
		FormData.ID.EQ(postgres.String(formDataId)),
	)

	result := results.FormDataResult{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}
