package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/types"
)

type FormRepository struct{}

func NewFormRepository() *FormRepository {
	return &FormRepository{}
}

func (r *FormRepository) FindAll(ctx context.Context, db *sql.DB) (*[]types.FormWithFields, error) {
	Forms := table.Forms
	FormFields := table.FormFields

	stmt := postgres.SELECT(
		Forms.AllColumns,
		FormFields.AllColumns,
	).FROM(
		Forms.
			LEFT_JOIN(FormFields, Forms.ID.EQ(FormFields.FormID)),
	)

	var forms []types.FormWithFields
	err := stmt.QueryContext(ctx, db, &forms)

	return &forms, err
}

func (r *FormRepository) Create(ctx context.Context, tx *sql.Tx, form model.Forms) (model.Forms, error) {
	Forms := table.Forms

	formInsertColumns := Forms.AllColumns.Except(Forms.ID, Forms.CreatedAt, Forms.UpdatedAt, Forms.DeletedAt)

	formStmt := Forms.INSERT(formInsertColumns).MODEL(form).RETURNING(Forms.ID)

	if err := formStmt.QueryContext(ctx, tx, &form); err != nil {
		return form, err
	}

	return form, nil
}
