package repositories

import (
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

func (r *formRepositoryImpl) Create(form *entities.Form) error {
	return r.db.Create(&form).Error
}

func (r *formRepositoryImpl) Delete(form *entities.Form) error {
	//Soft delete
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

func (r *formRepositoryImpl) FindById(id string) (*entities.Form, error) {
	var form entities.Form
	err := r.db.First(&form, id).Error
	if err != nil {
		return nil, err
	}

	return &form, nil
}
