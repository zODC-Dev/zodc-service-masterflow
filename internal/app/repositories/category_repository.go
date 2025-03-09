package repositories

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/model"
	"github.com/zODC-Dev/zodc-service-masterflow/database/generated/zodc_masterflow_dev/public/table"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/dto/queryparams"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) FindAll(ctx context.Context, db *sql.DB, queryParam queryparams.CategoryQueryParam) ([]model.Categories, error) {
	Categories := table.Categories

	categories := []model.Categories{}

	statement := postgres.SELECT(
		Categories.AllColumns,
	).FROM(
		Categories,
	)

	conditions := []postgres.BoolExpression{}

	if queryParam.Type != "" {
		conditions = append(conditions, postgres.LOWER(Categories.Type).EQ(postgres.LOWER(postgres.String(queryParam.Type))))
	}

	if queryParam.Search != "" {
		conditions = append(conditions, postgres.LOWER(Categories.Name).EQ(postgres.LOWER(postgres.String(queryParam.Search))))
	}

	if queryParam.IsActive != "" {
		isActiveBool, err := strconv.ParseBool(queryParam.IsActive)
		if err != nil {
			return categories, err
		}

		conditions = append(conditions, Categories.IsActive.EQ(postgres.Bool(isActiveBool)))
	}

	if len(conditions) > 0 {
		statement = statement.WHERE(postgres.AND(conditions...))
	}

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
