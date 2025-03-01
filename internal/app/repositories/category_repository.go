package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) FindAll(ctx context.Context, db *sql.DB, typeQueryParam string) ([]model.Categories, error) {
	Categories := table.Categories

	statement := postgres.SELECT(
		Categories.AllColumns,
	).FROM(
		Categories,
	)

	if typeQueryParam != "" {
		statement.WHERE(
			postgres.LOWER(Categories.Type).EQ(postgres.LOWER(postgres.String(typeQueryParam))),
		)
	}

	categories := []model.Categories{}
	err := statement.QueryContext(ctx, db, &categories)

	return categories, err
}

func (r *CategoryRepository) FindOneCategoryByKey(ctx context.Context, db *sql.DB, key string) (model.Categories, error) {
	Categories := table.Categories

	statement := postgres.SELECT(
		Categories.AllColumns,
	).FROM(
		Categories,
	).WHERE(
		Categories.Key.EQ(postgres.String(key)),
	)

	category := model.Categories{}
	err := statement.QueryContext(ctx, db, &category)

	return category, err
}
