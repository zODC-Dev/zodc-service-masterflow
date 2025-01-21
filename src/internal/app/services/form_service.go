package services

import (
	"encoding/json"

	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/interfaces"
	"github.com/zODC-Dev/zodc-service-masterflow/src/pkg/utils"
)

type formServiceImpl struct {
	formRepo interfaces.FormRepository
}

func NewFormService(formRepo interfaces.FormRepository) *formServiceImpl {
	return &formServiceImpl{
		formRepo: formRepo,
	}
}

func (s *formServiceImpl) FindAll() (*[]responses.FormFindAll, error) {
	forms, err := s.formRepo.FindAll()
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

func (s *formServiceImpl) Create(formCreate *requests.FormCreate) error {
	var formModel model.Forms

	if err := utils.Mapper(formCreate, &formModel); err != nil {
		return err
	}

	formFieldsModels := []model.FormFields{}

	dataSheet, err := json.Marshal(formCreate.DataSheet)
	if err != nil {
		return err
	}
	dataSheetPtr := string(dataSheet)
	formModel.DataSheet = &dataSheetPtr

	for i := range formCreate.FormFields {
		for j := range formCreate.FormFields[i] {

			var formFieldModel model.FormFields

			if err := utils.Mapper(formCreate.FormFields[i][j], &formFieldModel); err != nil {
				return err
			}

			advancedOptions, err := json.Marshal(formCreate.FormFields[i][j].AdvancedOptions)
			if err != nil {
				return err
			}

			formFieldModel.AdvancedOptions = string(advancedOptions)

			formFieldModel.ColNum = int32(i)

			formFieldsModels = append(formFieldsModels, formFieldModel)
		}
	}

	if err := s.formRepo.CreateForm(formModel, formFieldsModels); err != nil {
		return err
	}

	return err

}
