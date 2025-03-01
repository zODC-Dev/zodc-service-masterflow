package routes

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/controllers"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/repositories"
	"github.com/zODC-Dev/zodc-service-masterflow/internal/app/services"
)

func CategoryRoute(group *echo.Group, db *sql.DB) {
	categoryRepo := repositories.NewCategoryRepository()
	categoryService := services.NewCategoryService(db, categoryRepo)
	categoryController := controllers.NewCategoryController(categoryService)

	categoryRoute := group.Group("/categories")
	{
		categoryRoute.GET("", categoryController.FindAll)
	}
}
