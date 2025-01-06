package repositories

import (
	"github.com/zODC-Dev/zodc-service-masterflow/src/internal/app/models"
	"gorm.io/gorm"
)

type IFormRepository interface {
	Create(form *models.FormCreateRequest) error
}

type formRepositoryImlp struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) IFormRepository {
	return &formRepositoryImlp{
		db: db,
	}
}

func (r *formRepositoryImlp) Create(req *models.FormCreateRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		formExcel := models.FormExcel{
			FileName:    req.FileName,
			Title:       req.Title,
			Function:    req.Function,
			Template:    req.Template,
			DataSheet:   req.DataSheet,
			Description: req.Description,
		}

		if err := tx.Create(&formExcel).Error; err != nil {
			return err
		}

		for i := range req.Forms {
			req.Forms[i].FormExcelID = formExcel.ID
		}

		if len(req.Forms) > 0 {
			if err := tx.Create(&req.Forms).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
