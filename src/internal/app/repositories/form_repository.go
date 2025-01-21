package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	. "github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/types"
)

type formRepositoryImpl struct {
	db *sql.DB
}

func NewFormRepository(db *sql.DB) *formRepositoryImpl {
	return &formRepositoryImpl{
		db: db,
	}
}

func (r *formRepositoryImpl) FindAll() (*[]types.FormWithFields, error) {
	stmt := postgres.SELECT(
		Forms.AllColumns,
		FormFields.AllColumns,
	).FROM(
		Forms.
			LEFT_JOIN(FormFields, Forms.ID.EQ(FormFields.FormID)),
	)

	var forms []types.FormWithFields
	err := stmt.Query(r.db, &forms)

	return &forms, err
}

func (r *formRepositoryImpl) CreateForm(form model.Forms, formFields []model.FormFields) error {

	var err error

	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	formInsertColumns := Forms.AllColumns.Except(Forms.ID, Forms.CreatedAt, Forms.UpdatedAt, Forms.DeletedAt)

	formStmt := Forms.INSERT(formInsertColumns).MODEL(form).RETURNING(Forms.ID)

	formModel := model.Forms{}

	if err = formStmt.Query(tx, &formModel); err != nil {
		return err
	}

	for i := range formFields {
		formFields[i].FormID = formModel.ID
	}

	formFieldsInsertColumns := FormFields.AllColumns.Except(FormFields.ID, FormFields.CreatedAt, FormFields.UpdatedAt, FormFields.DeletedAt)
	formFieldsStmt := FormFields.INSERT(formFieldsInsertColumns).MODELS(formFields)

	if _, err = formFieldsStmt.Exec(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
