package services

import (
	"context"
	"database/sql"

	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/responses"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/pkg/utils"
)

type CategoryService struct {
	db           *sql.DB
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(db *sql.DB, categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		db:           db,
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) FindAll(ctx context.Context, queryParam queryparams.CategoryQueryParam) ([]responses.CategoryFindAll, error) {
	categoriesResponse := []responses.CategoryFindAll{}

	categories, err := s.categoryRepo.FindAll(ctx, s.db, queryParam)
	if err != nil {
		return categoriesResponse, err
	}

	if err := utils.Mapper(categories, &categoriesResponse); err != nil {
		return categoriesResponse, err
	}

	return categoriesResponse, nil

}
