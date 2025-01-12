package services

import (
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/interfaces"
)

type formServiceImpl struct {
	formRepo interfaces.IFormRepository
}

func NewFormService(formRepo interfaces.IFormRepository) *formServiceImpl {
	return &formServiceImpl{
		formRepo: formRepo,
	}
}

func (s *formServiceImpl) Create(req *requests.FormCreateRequest) error {
	var form = entities.Form{
		FileId:      req.FileId,
		FileName:    req.FileName,
		Title:       req.Title,
		Function:    req.Function,
		Template:    req.Template,
		DataSheet:   req.DataSheet,
		Description: req.Description,
	}

	var formFields []entities.FormField

	for i := range req.FormFields {
		for j := range req.FormFields[i] {
			var data = req.FormFields[i][j]
			formFields = append(formFields, entities.FormField{
				Icon:            data.Icon,
				Title:           data.Title,
				Category:        data.Category,
				FieldName:       data.FieldName,
				FieldType:       data.FieldType,
				Required:        data.Required,
				AdvancedOptions: data.AdvancedOptions,
				ColNum:          uint(i),
			})
		}
	}

	form.FormFields = formFields

	return s.formRepo.Create(&form)
}

func (s *formServiceImpl) FindAll() (*[]responses.FormResponse, error) {
	forms, err := s.formRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var formsResponses []responses.FormResponse

	for i := range *forms {
		data := (*forms)[i]

		formFieldsResponses := [][]entities.FormField{}

		for _, field := range data.FormFields {

			rowIndex := int(field.ColNum)

			for len(formFieldsResponses) <= rowIndex {
				formFieldsResponses = append(formFieldsResponses, []entities.FormField{})
			}

			formFieldsResponses[rowIndex] = append(formFieldsResponses[rowIndex], field)
		}

		formsResponses = append(formsResponses, responses.FormResponse{
			BaseModel: entities.BaseModel{
				ID:        data.ID,
				CreatedAt: data.CreatedAt,
				UpdatedAt: data.UpdatedAt,
			},
			FileId:      data.FileId,
			FileName:    data.FileName,
			Title:       data.Title,
			Function:    data.Function,
			Template:    data.Template,
			DataSheet:   data.DataSheet,
			Description: data.Description,
			FormFields:  formFieldsResponses,
		})
	}

	return &formsResponses, nil
}

func (s *formServiceImpl) Delete(form *entities.Form) error {
	return s.formRepo.Delete(form)
}

func (s *formServiceImpl) FindById(id string) (*entities.Form, error) {
	return s.formRepo.FindById(id)
}
