package repositories

import (
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/dto/requests"
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/entities"
	"gorm.io/gorm"
)

type formRepositoryImpl struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) *formRepositoryImpl {
	return &formRepositoryImpl{
		db: db,
	}
}

func (r *formRepositoryImpl) Create(req *requests.FormCreateRequest) error {
	form := entities.Form{
		FileName:    req.FileName,
		Title:       req.Title,
		Function:    req.Function,
		Template:    req.Template,
		DataSheet:   req.DataSheet,
		Description: req.Description,
	}

	if err := r.db.Create(&form).Error; err != nil {
		return err
	}

	for i := range req.FormFields {
		req.FormFields[i].FormID = form.ID
	}

	if err := r.db.Create(&req.FormFields).Error; err != nil {
		return err
	}

	return nil
}

func (r *formRepositoryImpl) Delete(form *entities.Form) error {
	return r.db.Delete(&form).Error
}

func (r *formRepositoryImpl) FindAll() (*[]entities.Form, error) {
	var forms []entities.Form
	err := r.db.Preload("FormFields").Find(&forms).Error
	if err != nil {
		return nil, err
	}

	return &forms, nil
}

func (r *formRepositoryImpl) Find(form *entities.Form) error {
	return r.db.Where(&entities.Form{}).Find(&form).Error
}
