package services

import (
	"encoding/json"

	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/models"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/repositories"
	"gorm.io/datatypes"
)

type IFormService interface {
	Create(form *models.FormCreateRequest) error
}

type formServiceImpl struct {
	formRepo repositories.IFormRepository
}

func NewFormService(formRepo repositories.IFormRepository) *formServiceImpl {
	return &formServiceImpl{
		formRepo: formRepo,
	}
}

func (s *formServiceImpl) Create(req *models.FormCreateRequest) error {
	for i := range req.Forms {
		optionsJSON, err := json.Marshal(req.Forms[i].AdvancedOptions)
		if err != nil {
			return err
		}

		req.Forms[i].AdvancedOptions = datatypes.JSON(optionsJSON)
	}

	return s.formRepo.Create(req)
}
