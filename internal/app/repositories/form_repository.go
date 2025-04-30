package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

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
		FormTemplates.CurrentVersion.EQ(FormTemplateVersions.Version),
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

func (r *FormRepository) FindAllFormSystem(ctx context.Context, db *sql.DB) ([]results.FormResult, error) {
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

	formSystemResults := []results.FormResult{}

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

func (r *FormRepository) RemoveAllFormFieldDataByFormDataId(ctx context.Context, tx *sql.Tx, formDataId string) error {
	FormFieldData := table.FormFieldData

	statement := FormFieldData.DELETE().WHERE(
		FormFieldData.FormDataID.EQ(postgres.String(formDataId)),
	)

	_, err := statement.ExecContext(ctx, tx)

	return err
}

func (r *FormRepository) FindOneFormTemplateByFormTemplateId(ctx context.Context, db *sql.DB, formTemplateId int32) (results.FormResult, error) {
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
		FormTemplates.ID.EQ(postgres.Int32(formTemplateId)).AND(
			FormTemplateVersions.Version.EQ(FormTemplates.CurrentVersion),
		),
	)

	formTemplate := results.FormResult{}

	err := statement.QueryContext(ctx, db, &formTemplate)

	return formTemplate, err
}

func (r *FormRepository) FindOneFormTemplateByFormTemplateVersionId(ctx context.Context, db *sql.DB, formTemplateVersionId int32) (results.FormResult, error) {
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
		FormTemplateVersions.ID.EQ(postgres.Int32(formTemplateVersionId)),
	)

	formTemplate := results.FormResult{}

	err := statement.QueryContext(ctx, db, &formTemplate)

	return formTemplate, err
}

func (r *FormRepository) FindFormDataById(ctx context.Context, db *sql.DB, formDataId string) (results.FormDataResult, error) {
	FormData := table.FormData
	FormFieldData := table.FormFieldData
	FormTemplateVersions := table.FormTemplateVersions
	FormTemplateFields := table.FormTemplateFields
	FormTemplates := table.FormTemplates

	statement := FormData.SELECT(
		FormData.AllColumns,
		FormTemplateFields.AllColumns,
		FormFieldData.AllColumns,
		FormTemplates.AllColumns,
	).FROM(
		FormData.
			LEFT_JOIN(FormFieldData, FormData.ID.EQ(FormFieldData.FormDataID)).
			LEFT_JOIN(FormTemplateVersions, FormData.FormTemplateVersionID.EQ(FormTemplateVersions.ID)).
			LEFT_JOIN(FormTemplateFields, FormTemplateVersions.ID.EQ(FormTemplateFields.FormTemplateVersionID)).
			LEFT_JOIN(FormTemplates, FormTemplateVersions.FormTemplateID.EQ(FormTemplates.ID)),
	).WHERE(
		FormData.ID.EQ(postgres.String(formDataId)),
	)

	result := results.FormDataResult{}
	err := statement.QueryContext(ctx, db, &result)

	return result, err
}

func (r *FormRepository) UpdateFormFieldJiraKey(ctx context.Context, tx *sql.Tx, nodeId string, jiraKey string) error {
	Nodes := table.Nodes
	FormFieldData := table.FormFieldData
	FormTemplateFields := table.FormTemplateFields

	// 2. Lấy formDataId từ nodeId
	formDataQuery := postgres.SELECT(
		Nodes.FormDataID,
	).FROM(
		Nodes,
	).WHERE(
		Nodes.ID.EQ(postgres.String(nodeId)),
	)

	var nodeData struct {
		FormDataID *string
	}
	if err := formDataQuery.QueryContext(ctx, tx, &nodeData); err != nil {
		slog.Error("Error when querying formDataId", "nodeId", nodeId, "error", err)
		return fmt.Errorf("error when querying formDataId from nodeId: %w", err)
	}

	// Nếu node không có formDataId thì return
	if nodeData.FormDataID == nil || *nodeData.FormDataID == "" {
		return nil
	}

	formDataId := *nodeData.FormDataID

	// 3. Tìm và cập nhật trường có fieldId là "key" trong form_field_data
	// Trước tiên, tìm xem có field nào có fieldId="key" không
	findKeyFieldQuery := postgres.SELECT(
		FormFieldData.ID,
		FormTemplateFields.FieldID,
	).FROM(
		FormFieldData.
			INNER_JOIN(FormTemplateFields, FormFieldData.FormTemplateFieldID.EQ(FormTemplateFields.ID)),
	).WHERE(
		FormFieldData.FormDataID.EQ(postgres.String(formDataId)).
			AND(FormTemplateFields.FieldID.EQ(postgres.String("key"))),
	)

	var keyFields []struct {
		ID      int32
		FieldID string
	}

	if err := findKeyFieldQuery.QueryContext(ctx, tx, &keyFields); err != nil {
		slog.Error("Error when finding field with fieldId='key'", "formDataId", formDataId, "error", err)
		return fmt.Errorf("error when finding field with fieldId='key': %w", err)
	}

	if len(keyFields) == 0 {
		return nil
	}

	// 4. Thực hiện update
	subQuery := postgres.SELECT(
		FormFieldData.ID,
	).FROM(
		FormFieldData.
			INNER_JOIN(FormTemplateFields, FormFieldData.FormTemplateFieldID.EQ(FormTemplateFields.ID)),
	).WHERE(
		FormFieldData.FormDataID.EQ(postgres.String(formDataId)).
			AND(FormTemplateFields.FieldID.EQ(postgres.String("key"))),
	)

	statement := FormFieldData.UPDATE(FormFieldData.Value).
		SET(postgres.String(jiraKey)).
		WHERE(FormFieldData.ID.IN(subQuery))

	_, err := statement.ExecContext(ctx, tx)
	if err != nil {
		slog.Error("Error when updating jiraKey", "formDataId", formDataId, "error", err)
		return fmt.Errorf("error when updating jiraKey in form field data: %w", err)
	}

	return nil
}

func (r *FormRepository) UpdateFormTemplate(ctx context.Context, tx *sql.Tx, formTemplate model.FormTemplates) error {
	FormTemplates := table.FormTemplates

	formTemplate.UpdatedAt = time.Now()

	columns := FormTemplates.AllColumns.Except(FormTemplates.ID, FormTemplates.CreatedAt, FormTemplates.DeletedAt)

	statement := FormTemplates.UPDATE(columns).MODEL(formTemplate).WHERE(FormTemplates.ID.EQ(postgres.Int32(formTemplate.ID)))

	_, err := statement.ExecContext(ctx, tx)

	return err
}
