package services

import (
	"encoding/json"

	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
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
	for i := range req.Forms {
		optionsJSON, err := json.Marshal(req.Forms[i].AdvancedOptions)
		if err != nil {
			return err
		}

		req.Forms[i].AdvancedOptions = datatypes.JSON(optionsJSON)
	}

	return s.formRepo.Create(req)
}
