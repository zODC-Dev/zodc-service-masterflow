package services

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type FormService struct {
	db            *sql.DB
	formRepo      *repositories.FormRepository
	formFieldRepo *repositories.FormFieldRepository
}

func NewFormService(db *sql.DB, formRepo *repositories.FormRepository, formFieldRepo *repositories.FormFieldRepository) *FormService {
	return &FormService{
		db:            db,
		formRepo:      formRepo,
		formFieldRepo: formFieldRepo,
	}
}

func (s *FormService) FindAll(ctx context.Context) (*[]responses.FormFindAll, error) {
	forms, err := s.formRepo.FindAll(ctx, s.db)
	if err != nil {
		return nil, err
	}

	var formsResponse = []responses.FormFindAll{}

	for _, form := range *forms {

		formFieldsList := [][]responses.FormFieldsFindAll{}

		for _, formField := range form.FormFields {
			colIndex := formField.ColNum

			for len(formFieldsList) <= int(colIndex) {
				formFieldsList = append(formFieldsList, []responses.FormFieldsFindAll{})
			}

			var formFieldFindAll responses.FormFieldsFindAll

			var advancedOptions map[string]interface{}
			if err := json.Unmarshal([]byte(formField.AdvancedOptions), &advancedOptions); err != nil {
				return nil, err
			}

			formFieldFindAll.AdvancedOptions = advancedOptions

			if err := utils.Mapper(formField, &formFieldFindAll); err != nil {
				return nil, err
			}

			formFieldsList[colIndex] = append(formFieldsList[colIndex], formFieldFindAll)
		}

		var formFindAll responses.FormFindAll
		err := utils.Mapper(form, &formFindAll)
		if err != nil {
			return nil, err
		}

		formFindAll.FormFields = formFieldsList

		var dataSheet map[string]interface{}
		if err := json.Unmarshal([]byte(*form.DataSheet), &dataSheet); err != nil {
			return nil, err
		}

		formFindAll.DataSheet = &dataSheet

		formsResponse = append(formsResponse, formFindAll)
	}

	return &formsResponse, nil

}

func (s *FormService) Create(ctx context.Context, req *requests.FormCreate) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var formModel model.Forms
	var formFieldsModels = []model.FormFields{}

	if err := utils.Mapper(req, &formModel); err != nil {
		return err
	}

	dataSheet, err := utils.MapToString(req.DataSheet)
	if err != nil {
		return err
	}

	formModel.DataSheet = &dataSheet

	for i := range req.FormFields {
		for j := range req.FormFields[i] {

			var formFieldModel model.FormFields

			if err := utils.Mapper(req.FormFields[i][j], &formFieldModel); err != nil {
				return err
			}

			advancedOptions, err := utils.MapToString(req.FormFields[i][j].AdvancedOptions)
			if err != nil {
				return err
			}

			formFieldModel.AdvancedOptions = advancedOptions

			formFieldModel.ColNum = int32(i)

			formFieldsModels = append(formFieldsModels, formFieldModel)
		}
	}

	form, err := s.formRepo.Create(ctx, tx, formModel)
	if err != nil {
		return err
	}

	if err := s.formFieldRepo.Create(ctx, tx, formFieldsModels, form.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
