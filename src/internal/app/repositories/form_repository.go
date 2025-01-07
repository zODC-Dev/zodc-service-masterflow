package repositories

import (
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/interfaces"
	"gorm.io/gorm"
)

type formRepositoryImpl struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) interfaces.IFormRepository {
	return &formRepositoryImpl{
		db: db,
	}
}

func (r *formRepositoryImpl) Create(req *requests.FormCreateRequest) error {
	formExcel := entities.FormExcel{
		FileName:    req.FileName,
		Title:       req.Title,
		Function:    req.Function,
		Template:    req.Template,
		DataSheet:   req.DataSheet,
		Description: req.Description,
	}

	if err := r.db.Create(&formExcel).Error; err != nil {
		return err
	}

	for i := range req.Forms {
		req.Forms[i].FormExcelID = formExcel.ID
	}

	if len(req.Forms) > 0 {
		if err := r.db.Create(&req.Forms).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *formRepositoryImpl) Delete(form *entities.Form) error {
	return r.db.Delete(&form).Error
}

func (r *formRepositoryImpl) FindAll() ([]entities.Form, error) {
	var forms []entities.Form
	err := r.db.Find(&forms).Error
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (r *formRepositoryImpl) Find(form *entities.Form) error {
	return r.db.Where(&entities.Form{}).Find(&form).Error
}
