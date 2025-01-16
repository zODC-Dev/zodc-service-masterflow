package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	database "github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database/generated"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/database/models"
)

type formServiceImpl struct {
	db      *pgx.Conn
	queries *database.Queries
}

func NewFormService(db *pgx.Conn, queries *database.Queries) *formServiceImpl {
	return &formServiceImpl{
		db:      db,
		queries: queries,
	}
}

func (s *formServiceImpl) Create(ctx context.Context, createFormRequest *models.CreateFormRequest) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	/**
	 * Create form
	 */
	form, err := s.queries.CreateForm(ctx, &database.CreateFormParams{
		FileName:    createFormRequest.FileName,
		Title:       createFormRequest.Title,
		Function:    createFormRequest.Function,
		Version:     createFormRequest.Version,
		Template:    createFormRequest.Template,
		Datasheet:   createFormRequest.Datasheet,
		Description: createFormRequest.Description,
		Decoration:  createFormRequest.Decoration,
	})
	if err != nil {
		return err
	}

	/**
	 * Create form fields
	 */
	for i := range createFormRequest.FormFields {
		for j := range createFormRequest.FormFields[i] {
			var formFieldData = createFormRequest.FormFields[i][j]
			_, err = qtx.CreateFormField(ctx, &database.CreateFormFieldParams{
				FieldID:         formFieldData.FieldID,
				Icon:            formFieldData.Icon,
				Title:           formFieldData.Title,
				Category:        formFieldData.Category,
				FieldName:       formFieldData.FieldName,
				FieldType:       formFieldData.FieldType,
				Required:        formFieldData.Required,
				AdvancedOptions: formFieldData.AdvancedOptions,
				ColNum:          int32(i),
				FormID:          form.ID,
			})
		}
	}
	if err != nil {
		return err
	}

	/**
	 * Commit transaction
	 */
	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *formServiceImpl) FindAll(ctx context.Context) ([]*database.FindAllFormsRow, error) {
	forms, err := s.queries.FindAllForms(ctx)
	return forms, err
}

func (s *formServiceImpl) FindAllView(ctx context.Context) ([]*database.FormView, error) {
	forms, err := s.queries.FormView(ctx)
	return forms, err
}

// func (s *formServiceImpl) Delete(form *entities.Form) error {
// 	return s.formRepo.Delete(form)
// }

// func (s *formServiceImpl) FindById(id string) (*entities.Form, error) {
// 	return s.formRepo.FindById(id)
// }
