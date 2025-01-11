package services

import (
	"encoding/json"

	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/interfaces"
	"gorm.io/datatypes"
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
	for i := range req.FormFields {
		optionsJSON, err := json.Marshal(req.FormFields[i].AdvancedOptions)
		if err != nil {
			return err
		}

		req.FormFields[i].AdvancedOptions = datatypes.JSON(optionsJSON)
	}

	return s.formRepo.Create(req)
}

func (s *formServiceImpl) FindAll() (*[]entities.Form, error) {
	forms, err := s.formRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return forms, nil
}
