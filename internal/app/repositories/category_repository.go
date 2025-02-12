package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow/public/table"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) FindAll(ctx context.Context, db *sql.DB) ([]model.Categories, error) {
	Categories := table.Categories

	stmt := postgres.SELECT(
		Categories.AllColumns,
	).FROM(Categories)

	categories := []model.Categories{}
	err := stmt.QueryContext(ctx, db, &categories)

	return categories, err
}
